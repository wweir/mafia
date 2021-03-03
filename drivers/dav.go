package drivers

import (
	"context"
	"io/fs"
	"os"

	"golang.org/x/net/webdav"
)

type WebdavAdaptor interface {
	Mkdir(name string, perm os.FileMode) error
	OpenFile(name string, flag int, perm os.FileMode) (WebdavFileAdaptor, error)
	RemoveAll(name string) error
	Rename(oldName, newName string) error
	Stat(name string) (os.FileInfo, error)
}

type WebdavDriver struct {
	Adaptor WebdavAdaptor
}

func (dav *WebdavDriver) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return dav.Adaptor.Mkdir(name, perm)
}
func (dav *WebdavDriver) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	file, err := dav.Adaptor.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return &webdavFile{Adaptor: file}, nil
}
func (dav *WebdavDriver) RemoveAll(ctx context.Context, name string) error {
	return dav.Adaptor.RemoveAll(name)
}
func (dav *WebdavDriver) Rename(ctx context.Context, oldName, newName string) error {
	return dav.Adaptor.Rename(oldName, newName)
}
func (dav *WebdavDriver) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return dav.Adaptor.Stat(name)
}

type WebdavFileAdaptor interface {
	Readdir(count int) ([]fs.FileInfo, error)
	Stat() (fs.FileInfo, error)
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
	Close() error
}
type webdavFile struct {
	Adaptor WebdavFileAdaptor
}

func (f *webdavFile) Readdir(count int) ([]fs.FileInfo, error) {
	return f.Adaptor.Readdir(count)
}
func (f *webdavFile) Stat() (fs.FileInfo, error) {
	return f.Adaptor.Stat()
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
func (f *webdavFile) Close() error {
	return f.Adaptor.Close()
}
