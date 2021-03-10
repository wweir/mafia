package drivers

import (
	"context"
	"os"

	"github.com/wweir/mafia/pkg/fspath"
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
	path = fspath.Relative(path)
	defer func() {
		Defer.Err(err).
			Str("path", path).
			Msg("Mkdir")
	}()
	return dav.Adaptor.Mkdir(path, perm)
}
func (dav *WebdavDriver) OpenFile(ctx context.Context, path string, flag int, perm os.FileMode) (_ webdav.File, err error) {
	path = fspath.Relative(path)
	defer func() {
		Defer.Err(err).
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
	path = fspath.Relative(path)
	defer func() {
		Defer.Err(err).
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
	oldpath = fspath.Relative(oldpath)
	newpath = fspath.Relative(newpath)
	defer func() {
		Defer.Err(err).
			Str("from", oldpath).
			Str("to", newpath).
			Msg("DeleteFile")
	}()
	return dav.Adaptor.Rename(oldpath, newpath)
}
func (dav *WebdavDriver) Stat(ctx context.Context, path string) (fi os.FileInfo, err error) {
	path = fspath.Relative(path)
	defer func() {
		Defer.Err(err).
			Str("path", path).
			Interface("fi", fi).
			Msg("Stat")
	}()

	return dav.Adaptor.Stat(path)
}
