package fsutil

import (
	"io/fs"
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

type MockFileInfo struct {
	Path  string
	Isdir bool
}

func (fi *MockFileInfo) Name() string { return "" }
func (fi *MockFileInfo) Size() int64  { return 0 }
func (fi *MockFileInfo) Mode() fs.FileMode {
	if fi.Isdir {
		return 0755
	}
	return 0644
}
func (fi *MockFileInfo) ModTime() time.Time { return time.Now() }
func (fi *MockFileInfo) IsDir() bool        { return fi.Isdir }
func (fi *MockFileInfo) Sys() interface{}   { return nil }
