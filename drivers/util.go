package drivers

import (
	"io"
	"io/fs"
	"os"

	"goftp.io/server/v2"
)

var Drivers = map[string]server.Driver{}

func SeekRead(f interface {
	io.ReadSeekCloser
	Stat() (fs.FileInfo, error)
}, offset int64) (int64, io.ReadCloser, error) {

	info, err := f.Stat()
	if err != nil {
		return 0, nil, err
	}

	_, err = f.Seek(offset, io.SeekStart)
	if err != nil {
		return 0, nil, err
	}

	return info.Size() - offset, f, nil
}

type base interface {
	DeleteDir(path string) error // RemoveAll(name string) error
	DeleteFile(path string) error
	GetFile(path string, offset int64) (int64, io.ReadCloser, error)
	ListDir(path string, callback func(os.FileInfo) error) error
	PutFile(destPath string, data io.Reader, offset int64) (int64, error)
	Rename(fromPath string, toPath string) error
	Stat(path string) (os.FileInfo, error)

	Mkdir(name string, perm os.FileMode) error // MakeDir(path string) error
	OpenFile(name string, flag int, perm os.FileMode) (WebdavFileAdaptor, error)
}

type file interface {
	Readdir(count int) ([]fs.FileInfo, error)
	Stat() (fs.FileInfo, error)
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
	Close() error
}
