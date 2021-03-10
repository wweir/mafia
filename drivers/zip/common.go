package zip

import (
	"archive/zip"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/wweir/mafia/drivers"
	"github.com/wweir/mafia/pkg/fsmock"
	"github.com/wweir/mafia/pkg/fspath"
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
		if hdr.Flags == 0 {
			if hdr.Name, err = scDecoder.String(hdr.Name); err != nil {
				continue
			}
			if hdr.Comment, err = scDecoder.String(hdr.Comment); err != nil {
				continue
			}
		}
		log.Info().Interface("hdr", hdr).Msg("")

		c.fis[hdr.Name] = hdr.FileInfo()
	}

	return c, nil
}

func (c *common) Stat(name string) (os.FileInfo, error) {
	if fi, ok := c.fis[name]; ok {
		return fi, nil
	}

	switch name {
	case "":
		return fsmock.MockFileInfo(".", true, 0, nil), nil
	case ".", "./", "/":
		return fsmock.MockFileInfo(name, true, 0, nil), nil
	}

	if name[:len(name)-1] != "/" {
		name += "/"
	}

	for path := range c.fis {
		if strings.HasPrefix(path, name) {
			return fsmock.MockFileInfo(name, true, 0, nil), nil
		}
	}
	return nil, os.ErrNotExist
}

func (c *common) Mkdir(name string, perm os.FileMode) error { return nil }

func (c *common) ReadDir(dir string) ([]os.FileInfo, error) {
	fis := make([]os.FileInfo, 0, len(c.fis))
	for file, fi := range c.fis {
		switch fspath.SumPathRelation(file, dir) {
		case fspath.PathParrent:
			fis = append(fis, fi)
		case fspath.PathSup:
			fis = append(fis, fsmock.MockFileInfo(file, true, 0, nil))
		}
	}
	return fis, nil
}
