package sftp

import (
	"os"

	"github.com/wweir/mafia/drivers"
)

var _ drivers.WebdavAdaptor = new(Webdav)

type Webdav struct {
	*common
}

func NewWebdav(conf *SSHConfig) (*Webdav, error) {
	c, err := newCommon(conf)
	if err != nil {
		return nil, err
	}

	return &Webdav{
		common: c,
	}, nil
}

func (dav *Webdav) OpenFile(name string, flag int, perm os.FileMode) (drivers.FileAdaptor, error) {
	return dav.Client.OpenFile(name, flag)
}
