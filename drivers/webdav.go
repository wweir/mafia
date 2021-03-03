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
	WebdavAdaptor
}

func (dav *WebdavDriver) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return dav.WebdavAdaptor.Mkdir(name, perm)
}
func (dav *WebdavDriver) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	file, err := dav.WebdavAdaptor.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return &webdavFile{WebdavFileAdaptor: file}, nil
}
func (dav *WebdavDriver) RemoveAll(ctx context.Context, name string) error {
	return dav.WebdavAdaptor.RemoveAll(name)
}
func (dav *WebdavDriver) Rename(ctx context.Context, oldName, newName string) error {
	return dav.WebdavAdaptor.Rename(oldName, newName)
}
func (dav *WebdavDriver) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return dav.WebdavAdaptor.Stat(name)
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
	WebdavFileAdaptor
}

func (f *webdavFile) Readdir(count int) ([]fs.FileInfo, error) {
	return f.WebdavFileAdaptor.Readdir(count)
}
func (f *webdavFile) Stat() (fs.FileInfo, error) {
	return f.WebdavFileAdaptor.Stat()
}
func (f *webdavFile) Read(p []byte) (n int, err error) {
	return f.WebdavFileAdaptor.Read(p)
}
func (f *webdavFile) Write(p []byte) (n int, err error) {
	return f.WebdavFileAdaptor.Write(p)
}
func (f *webdavFile) Seek(offset int64, whence int) (int64, error) {
	return f.WebdavFileAdaptor.Seek(offset, whence)
}
func (f *webdavFile) Close() error {
	return f.WebdavFileAdaptor.Close()
}
