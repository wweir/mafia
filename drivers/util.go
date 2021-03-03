package drivers

import (
	"io"
	"io/fs"

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
