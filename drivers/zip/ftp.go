package zip

import (
	"archive/zip"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/wweir/mafia/drivers"
	"github.com/wweir/mafia/pkg/fsmock"
	"github.com/wweir/mafia/pkg/fspath"
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
			if hdr.Name, err = scDecoder.String(hdr.Name); err != nil {
				continue
			}
			if hdr.Comment, err = scDecoder.String(hdr.Comment); err != nil {
				continue
			}
		}

		if fspath.SumPathRelation(path, hdr.Name) == fspath.PathSelf {
			rc, err := zf.Open()
			if err != nil {
				zrc.Close()
				return 0, nil, errors.WithStack(err)
			}

			return hdr.FileInfo().Size(), &fsmock.MockReaderCloser{
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
