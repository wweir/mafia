package drivers

import (
	"io/fs"
	"os"
)

type webdavFile struct {
	Adaptor FileAdaptor

	path    string
	readdir func(string) ([]fs.FileInfo, error)
	stat    func(string) (os.FileInfo, error)
}

func (f *webdavFile) Readdir(count int) (fis []fs.FileInfo, err error) {
	defer func() {
		Defer.Err(err).
			Str("file", f.path).
			Int("count", count).
			Int("fis", len(fis)).
			Msg("Readdir")
	}()

	fss, err := f.readdir(f.path)
	if err != nil {
		return nil, err
	}
	if count > 0 && len(fss) > count {
		return fss[:count], nil
	}
	return fss, nil
}
func (f *webdavFile) Stat() (fi fs.FileInfo, err error) {
	defer func() {
		Defer.Err(err).
			Str("file", f.path).
			Interface("fi", fi).
			Msg("Stat")
	}()
	return f.stat(f.path)
}
func (f *webdavFile) Read(p []byte) (n int, err error) {
	defer func() {
		Defer.Err(err).
			Str("file", f.path).
			Int("n", n).
			Msg("Read")
	}()
	return f.Adaptor.Read(p)
}
func (f *webdavFile) Write(p []byte) (n int, err error) {
	defer func() {
		Defer.Err(err).
			Str("file", f.path).
			Int("n", n).
			Msg("Write")
	}()
	return f.Adaptor.Write(p)
}
func (f *webdavFile) Seek(offset int64, whence int) (n int64, err error) {
	defer func() {
		Defer.Err(err).
			Str("file", f.path).
			Int64("n", n).
			Msg("Seek")
	}()
	return f.Adaptor.Seek(offset, whence)
}
func (f *webdavFile) Close() (err error) {
	defer func() {
		Defer.Err(err).
			Str("file", f.path).
			Msg("Close")
	}()
	return f.Adaptor.Close()
}
