package sftp

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"github.com/wweir/mafia/drivers"
	"goftp.io/server/v2"
)

type Sftp struct {
	*sftp.Client

	*server.MultiDriver
}

func (s *Sftp) Stat(ctx *server.Context, path string) (os.FileInfo, error) {
	return s.Client.Stat(path)
}
func (s *Sftp) ListDir(ctx *server.Context, path string, callback func(os.FileInfo) error) error {
	fis, err := s.Client.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		if err := callback(fi); err != nil {
			return err
		}
	}
	return nil
}
func (s *Sftp) DeleteDir(ctx *server.Context, path string) error {
	return s.Client.RemoveDirectory(path)
}
func (s *Sftp) DeleteFile(ctx *server.Context, path string) error {
	return s.Client.Remove(path)
}
func (s *Sftp) Rename(ctx *server.Context, fromPath string, toPath string) error {
	return s.Client.Rename(fromPath, toPath)
}
func (s *Sftp) MakeDir(ctx *server.Context, path string) error {
	return s.Client.Mkdir(path)
}
func (s *Sftp) GetFile(ctx *server.Context, path string, offset int64) (int64, io.ReadCloser, error) {
	f, err := s.Client.Open(path)
	if err != nil {
		return 0, nil, err
	}

	return drivers.SeekRead(f, offset)
}
func (s *Sftp) PutFile(ctx *server.Context, destPath string, data io.Reader, offset int64) (int64, error) {
	return 0, os.ErrPermission

	var isExist bool
	f, err := s.Client.Lstat(destPath)
	if err == nil {
		isExist = true
		if f.IsDir() {
			return 0, errors.New("A dir has the same name")
		}
	} else {
		if os.IsNotExist(err) {
			isExist = false
		} else {
			return 0, errors.New(fmt.Sprintln("Put File error:", err))
		}
	}

	if offset > -1 && !isExist {
		offset = -1
	}

	if offset == -1 {
		if isExist {
			err = s.Client.Remove(destPath)
			if err != nil {
				return 0, err
			}
		}
		f, err := s.Client.Create(destPath)
		if err != nil {
			return 0, err
		}
		defer f.Close()
		bytes, err := io.Copy(f, data)
		if err != nil {
			return 0, err
		}
		return bytes, nil
	}

	of, err := os.OpenFile(destPath, os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return 0, err
	}
	defer of.Close()

	info, err := of.Stat()
	if err != nil {
		return 0, err
	}
	if offset > info.Size() {
		return 0, fmt.Errorf("Offset %d is beyond file size %d", offset, info.Size())
	}

	_, err = of.Seek(offset, os.SEEK_END)
	if err != nil {
		return 0, err
	}

	bytes, err := io.Copy(of, data)
	if err != nil {
		return 0, err
	}

	return bytes, nil
}
