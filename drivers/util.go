package drivers

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"goftp.io/server/v2"
)

var Defer zerolog.Logger
var Drivers = map[string]server.Driver{}

type FSAdaptor interface {
	DeleteDir(path string) error
	DeleteFile(path string) error
	Rename(fromPath string, toPath string) error
	Mkdir(dirpath string, perm os.FileMode) error
	Stat(path string) (os.FileInfo, error)
	ReadDir(dirpath string) ([]os.FileInfo, error)
}

type FileAdaptor interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

type MockFSFull struct{}

// common
func (f *MockFSFull) DeleteDir(path string) error                 { return os.ErrPermission }
func (f *MockFSFull) DeleteFile(path string) error                { return os.ErrPermission }
func (f *MockFSFull) Rename(fromPath string, toPath string) error { return os.ErrPermission }
func (f *MockFSFull) Mkdir(name string, perm os.FileMode) error   { return os.ErrPermission }
func (f *MockFSFull) Stat(name string) (os.FileInfo, error)       { return nil, os.ErrPermission }
func (f *MockFSFull) ReadDir(name string) ([]os.FileInfo, error)  { return nil, os.ErrPermission }

// webdav
func (f *MockFSFull) OpenFile(path string, flag int, perm os.FileMode) (FileAdaptor, error) {
	return nil, os.ErrPermission
}

// ftp
func (f *MockFSFull) GetFile(path string, offset int64) (int64, io.ReadCloser, error) {
	return 0, nil, os.ErrPermission
}
func (f *MockFSFull) PutFile(destPath string, data io.Reader, offset int64) (int64, error) {
	return 0, os.ErrPermission
}
