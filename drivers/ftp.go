package drivers

import (
	"io"
	"os"

	"goftp.io/server/v2"
)

type FTPAdaptor interface {
	DeleteDir(path string) error
	DeleteFile(path string) error
	GetFile(path string, offset int64) (int64, io.ReadCloser, error)
	ListDir(path string, callback func(os.FileInfo) error) error
	MakeDir(path string) error
	PutFile(destPath string, data io.Reader, offset int64) (int64, error)
	Rename(fromPath string, toPath string) error
	Stat(path string) (os.FileInfo, error)
}

type FTPDriver struct {
	FTPAdaptor
}

func (ftp *FTPDriver) DeleteDir(ctx *server.Context, path string) error {
	return ftp.FTPAdaptor.DeleteDir(path)
}
func (ftp *FTPDriver) DeleteFile(ctx *server.Context, path string) error {
	return ftp.FTPAdaptor.DeleteFile(path)
}
func (ftp *FTPDriver) GetFile(ctx *server.Context, path string, offset int64) (int64, io.ReadCloser, error) {
	return ftp.FTPAdaptor.GetFile(path, offset)
}
func (ftp *FTPDriver) ListDir(ctx *server.Context, path string, callback func(os.FileInfo) error) error {
	return ftp.FTPAdaptor.ListDir(path, callback)
}
func (ftp *FTPDriver) MakeDir(ctx *server.Context, path string) error {
	return ftp.FTPAdaptor.MakeDir(path)
}
func (ftp *FTPDriver) PutFile(ctx *server.Context, destPath string, data io.Reader, offset int64) (int64, error) {
	return ftp.FTPAdaptor.PutFile(destPath, data, offset)
}
func (ftp *FTPDriver) Rename(ctx *server.Context, fromPath string, toPath string) error {
	return ftp.FTPAdaptor.Rename(fromPath, toPath)
}
func (ftp *FTPDriver) Stat(ctx *server.Context, path string) (os.FileInfo, error) {
	return ftp.FTPAdaptor.Stat(path)
}