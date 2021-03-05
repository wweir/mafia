package sftp

import (
	"os"

	"github.com/pkg/sftp"
	"github.com/wweir/mafia/drivers"
)

var _ drivers.FSAdaptor = new(common)

type common struct {
	*sftp.Client
}

func (c *common) DeleteDir(path string) error {
	return c.Client.RemoveDirectory(path)
}

func (c *common) DeleteFile(path string) error {
	return c.Client.Remove(path)
}

func (c *common) Mkdir(name string, perm os.FileMode) error {
	return c.Client.Mkdir(name)
}
