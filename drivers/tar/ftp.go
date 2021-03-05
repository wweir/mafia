package tar

import (
	"io"

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
	tr, close, err := ftp.openReader()
	if err != nil {
		return 0, nil, err
	}

	for {
		hdr, err := tr.Next()
		if err != nil {
			close()
			return 0, nil, errors.WithStack(err)
		}

		if fsutil.SumPathRelation(hdr.Name, path) == fsutil.PathSelf {
			return hdr.Size, &fsutil.MockReaderCloser{
				ReadFn:  tr.Read,
				CloseFn: close,
			}, nil
		}
	}
}
