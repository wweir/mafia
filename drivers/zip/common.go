package zip

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/wweir/mafia/drivers"
	"github.com/wweir/mafia/pkg/fsutil"
)

var _ drivers.FSAdaptor = new(common)

type common struct {
	path string
	fis  map[string]os.FileInfo

	drivers.MockFSFull
}

func newCommon(path string) (*common, error) {
	c := &common{
		path: path,
		fis:  map[string]os.FileInfo{},
	}

	zrc, err := zip.OpenReader(c.path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer zrc.Close()

	for _, zf := range zrc.File {
		hdr := zf.FileHeader
		if hdr.NonUTF8 {
			if str, err := scDecoder.String(hdr.Name); err == nil {
				hdr.Name = str
			}
			if str, err := scDecoder.String(hdr.Comment); err == nil {
				hdr.Comment = str
			}
		}

		c.fis[hdr.Name] = hdr.FileInfo()
	}

	return c, nil
}

func (c *common) Stat(name string) (os.FileInfo, error) {
	if fi, ok := c.fis[name]; ok {
		return fi, nil
	}

	if name == "/" {
		return &fsutil.MockFileInfo{
			Path:  name,
			Isdir: true,
		}, nil
	}

	if name[:len(name)-1] != "/" {
		name += "/"
	}

	for path := range c.fis {
		if strings.HasPrefix(path, name) {
			return &fsutil.MockFileInfo{
				Path:  name,
				Isdir: true,
			}, nil
		}
	}
	return nil, os.ErrNotExist
}

func (c *common) Mkdir(name string, perm os.FileMode) error { return nil }

func (c *common) ReadDir(dir string) ([]os.FileInfo, error) {
	fis := make([]os.FileInfo, 0, len(c.fis))
	for file, fi := range c.fis {
		switch fsutil.SumPathRelation(file, dir) {
		case fsutil.PathParrent:
			fis = append(fis, fi)
		case fsutil.PathSup:
			fis = append(fis, &fsutil.MockFileInfo{
				Path:  filepath.Base(file),
				Isdir: true,
			})
		}
	}
	return fis, nil
}
