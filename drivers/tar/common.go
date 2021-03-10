package tar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
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

	tr, close, err := c.openReader()
	if err != nil {
		return nil, err
	}
	defer close()

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return c, nil
		} else if err != nil {
			return nil, err
		}

		c.fis[hdr.Name] = hdr.FileInfo()
	}
}

func (c *common) openReader() (tr *tar.Reader, close func(), err error) {
	f, err := os.Open(c.path)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	var rc io.ReadCloser
	switch filepath.Ext(c.path) {
	case ".tar":
		rc = f
	case ".gz", ".tgz":
		if rc, err = gzip.NewReader(f); err != nil {
			f.Close()
			return nil, nil, errors.WithStack(err)
		}
	default:
		return nil, nil, errors.Errorf("unsupported file type: %s", c.path)
	}

	return tar.NewReader(rc), func() {
		rc.Close()
		f.Close()
	}, nil
}
func (c *common) openWriter() (tr *tar.Writer, close func(), err error) {
	f, err := os.Open(c.path)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	var wc io.WriteCloser
	switch filepath.Ext(c.path) {
	case ".tar":
		wc = f
	case ".gz", ".tgz":
		wc = gzip.NewWriter(f)
	default:
		return nil, nil, errors.Errorf("unsupported file type: %s", c.path)
	}

	tw := tar.NewWriter(wc)
	return tw, func() {
		tw.Close()
		f.Close()
	}, nil
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
