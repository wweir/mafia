package fsmock

import (
	"io/fs"
	"path/filepath"
	"time"
)

type MockReaderCloser struct {
	ReadFn  func(p []byte) (n int, err error)
	CloseFn func()
}

func (rc *MockReaderCloser) Read(p []byte) (n int, err error) {
	return rc.ReadFn(p)
}
func (rc *MockReaderCloser) Close() error {
	rc.CloseFn()
	return nil
}

type mockFileInfo struct {
	name    string
	isDir   bool
	size    int64
	modTime *time.Time
}

func MockFileInfo(name string, isDir bool, size int64, modTime *time.Time) fs.FileInfo {
	return &mockFileInfo{
		name:    filepath.Base(name),
		isDir:   isDir,
		size:    size,
		modTime: modTime,
	}
}

func (fi *mockFileInfo) Name() string { return fi.name }
func (fi *mockFileInfo) Size() int64  { return fi.size }
func (fi *mockFileInfo) Mode() fs.FileMode {
	if fi.isDir {
		return 0755
	}
	return 0644
}
func (fi *mockFileInfo) ModTime() time.Time {
	if fi.modTime != nil {
		return *fi.modTime
	}
	return time.Now()
}
func (fi *mockFileInfo) IsDir() bool      { return fi.isDir }
func (fi *mockFileInfo) Sys() interface{} { return nil }
