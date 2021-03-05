package drivers

import (
	"fmt"
	"io"
	"os"
	"time"

	"goftp.io/server/v2"
)

type FTPAdaptor interface {
	FSAdaptor

	GetFile(path string, offset int64) (int64, io.ReadCloser, error)
	PutFile(destPath string, data io.Reader, offset int64) (int64, error)
}

type FTPDriver struct {
	fs  FSAdaptor
	ftp FTPAdaptor
	dav WebdavAdaptor
}

func NewFTPDriver(ftp FTPAdaptor, dav WebdavAdaptor) *FTPDriver {
	switch {
	case ftp != nil:
		return &FTPDriver{fs: ftp, ftp: ftp}
	case dav != nil:
		return &FTPDriver{fs: ftp, dav: dav}
	default:
		panic("ftp adaptor and webdav adaptor are both empty")
	}
}

func (ftp *FTPDriver) DeleteDir(ctx *server.Context, path string) (err error) {
	defer func() {
		DeferLog.Err(err).
			Str("path", path).
			Msg("DeleteFile")
	}()
	return ftp.fs.DeleteDir(path)
}
func (ftp *FTPDriver) DeleteFile(ctx *server.Context, path string) (err error) {
	defer func() {
		DeferLog.Err(err).
			Str("path", path).
			Msg("DeleteFile")
	}()
	return ftp.fs.DeleteFile(path)
}

func (ftp *FTPDriver) ListDir(ctx *server.Context, path string, callback func(os.FileInfo) error) (err error) {
	files := []string{}
	defer func() {
		DeferLog.Err(err).
			Str("path", path).
			Strs("files", files).
			Msg("ListDir")
	}()

	fis, err := ftp.fs.ReadDir(path)
	if err != nil {
		return err
	}

	for f := range fis {
		files = append(files, fis[f].Name())
		if err := callback(fis[f]); err != nil {
			return err
		}
	}
	return nil
}
func (ftp *FTPDriver) MakeDir(ctx *server.Context, path string) (err error) {
	defer func() {
		DeferLog.Err(err).
			Str("path", path).
			Msg("DeleteFile")
	}()
	return ftp.fs.Mkdir(path, 0755)
}

func (ftp *FTPDriver) Rename(ctx *server.Context, fromPath string, toPath string) (err error) {
	defer func() {
		DeferLog.Err(err).
			Str("from", fromPath).
			Str("to", toPath).
			Msg("DeleteFile")
	}()
	return ftp.fs.Rename(fromPath, toPath)
}
func (ftp *FTPDriver) Stat(ctx *server.Context, path string) (fi os.FileInfo, err error) {
	defer func() {
		DeferLog.Err(err).
			Str("path", path).
			Interface("stat", fi).
			Msg("Stat")
	}()
	return ftp.fs.Stat(path)
}

func (ftp *FTPDriver) GetFile(ctx *server.Context, path string, offset int64) (size int64, _ io.ReadCloser, err error) {
	start := time.Now()
	defer func() {
		DeferLog.Err(err).
			Str("path", path).
			Int64("offset", offset).
			Int64("size", size).
			Dur("duration", time.Since(start)).
			Msg("GetFile")
	}()

	if ftp.ftp != nil {
		return ftp.ftp.GetFile(path, offset)
	} else if ftp.dav == nil {
		return 0, nil, os.ErrPermission
	}

	f, err := ftp.dav.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return 0, nil, err
	}

	stat, err := ftp.fs.Stat(path)
	if err != nil {
		return 0, nil, err
	}
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return 0, nil, err
	}
	return stat.Size() - offset, f, nil
}

func (ftp *FTPDriver) PutFile(ctx *server.Context, path string, data io.Reader, offset int64) (size int64, err error) {
	start := time.Now()
	defer func() {
		DeferLog.Err(err).
			Str("path", path).
			Int64("offset", offset).
			Int64("size", size).
			Dur("duration", time.Since(start)).
			Msg("PutFile")
	}()

	if ftp.ftp != nil {
		return ftp.ftp.PutFile(path, data, offset)
	} else if ftp.dav == nil {
		return 0, os.ErrPermission
	}

	if _, err := ftp.dav.Stat(path); !os.IsNotExist(err) {
		return 0, fmt.Errorf("%s is already exists, err: %w", path, err)
	}

	f, err := ftp.dav.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return 0, err
	}

	return io.Copy(f, data)
}
