package drivers

import (
	"context"
	"io/fs"
	"os"

	"golang.org/x/net/webdav"
)

type WebdavAdaptor interface {
	FSAdaptor

	OpenFile(path string, flag int, perm os.FileMode) (FileAdaptor, error)
}

type WebdavDriver struct {
	Adaptor WebdavAdaptor
}

func (dav *WebdavDriver) Mkdir(ctx context.Context, path string, perm os.FileMode) (err error) {
	defer func() {
		DeferLog.Err(err).
			Str("path", path).
			Msg("Mkdir")
	}()
	return dav.Adaptor.Mkdir(path, perm)
}
func (dav *WebdavDriver) OpenFile(ctx context.Context, path string, flag int, perm os.FileMode) (_ webdav.File, err error) {
	defer func() {
		DeferLog.Err(err).
			Str("path", path).
			Msg("OpenFile")
	}()

	file, err := dav.Adaptor.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}
	return &webdavFile{
		Adaptor: file,
		path:    path,
		stat:    dav.Adaptor.Stat,
		readdir: dav.Adaptor.ReadDir,
	}, nil
}
func (dav *WebdavDriver) RemoveAll(ctx context.Context, path string) (err error) {
	defer func() {
		DeferLog.Err(err).
			Str("path", path).
			Msg("RemoveAll")
	}()

	stat, err := dav.Adaptor.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if stat.IsDir() {
		return dav.Adaptor.DeleteDir(path)
	}
	return dav.Adaptor.DeleteFile(path)
}
func (dav *WebdavDriver) Rename(ctx context.Context, oldpath, newpath string) (err error) {
	defer func() {
		DeferLog.Err(err).
			Str("from", oldpath).
			Str("to", newpath).
			Msg("DeleteFile")
	}()
	return dav.Adaptor.Rename(oldpath, newpath)
}
func (dav *WebdavDriver) Stat(ctx context.Context, path string) (fi os.FileInfo, err error) {
	defer func() {
		DeferLog.Err(err).
			Str("path", path).
			Msg("Stat")
	}()

	return dav.Adaptor.Stat(path)
}

type webdavFile struct {
	Adaptor FileAdaptor

	path    string
	readdir func(string) ([]fs.FileInfo, error)
	stat    func(string) (os.FileInfo, error)
}

func (f *webdavFile) Readdir(count int) ([]fs.FileInfo, error) {
	fss, err := f.readdir(f.path)
	if err != nil {
		return nil, err
	}
	if len(fss) > count {
		return fss[:count], nil
	}
	return fss, nil
}
func (f *webdavFile) Stat() (fs.FileInfo, error) {
	return f.stat(f.path)
}
func (f *webdavFile) Read(p []byte) (n int, err error) {
	return f.Adaptor.Read(p)
}
func (f *webdavFile) Write(p []byte) (n int, err error) {
	return f.Adaptor.Write(p)
}
func (f *webdavFile) Seek(offset int64, whence int) (int64, error) {
	return f.Adaptor.Seek(offset, whence)
}
func (f *webdavFile) Close() (err error) {
	return f.Adaptor.Close()
}
