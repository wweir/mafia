package zip

import (
	"archive/zip"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/wweir/mafia/drivers"
	"github.com/wweir/mafia/pkg/fsutil"
)

var _ drivers.FTPAdaptor = new(FTP)

type FTP struct {
	*common
}

func NewFTP(path string) (*FTP, error) {
	c, err := newCommon(path)
	if err != nil {
		return nil, err
	}

	return &FTP{
		common: c,
	}, nil
}

func (ftp *FTP) GetFile(path string, offset int64) (int64, io.ReadCloser, error) {
	zrc, err := zip.OpenReader(ftp.path)
	if err != nil {
		return 0, nil, errors.WithStack(err)
	}

	for _, zf := range zrc.File {
		hdr := zf.FileHeader
		if hdr.Flags == 0 {
			if str, err := scDecoder.String(hdr.Name); err == nil {
				hdr.Name = str
			}
		}

		if fsutil.SumPathRelation(path, hdr.Name) == fsutil.PathSelf {
			rc, err := zf.Open()
			if err != nil {
				zrc.Close()
				return 0, nil, errors.WithStack(err)
			}

			return hdr.FileInfo().Size(), &fsutil.MockReaderCloser{
				ReadFn: rc.Read,
				CloseFn: func() {
					rc.Close()
					zrc.Close()
				},
			}, nil
		}
	}

	return 0, nil, os.ErrNotExist
}
