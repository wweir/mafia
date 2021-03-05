package sftp

import (
	"os"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/wweir/mafia/drivers"
	"golang.org/x/crypto/ssh"
)

var _ drivers.FSAdaptor = new(common)

type common struct {
	*sftp.Client
}

func newCommon(conf *SSHConfig) (*common, error) {
	clientConf, err := conf.Validate()
	if err != nil {
		return nil, err
	}

	sshClient, err := ssh.Dial("tcp", conf.Addr, clientConf)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &common{
		Client: sftpClient,
	}, nil
}

func (c *common) Stat(path string) (os.FileInfo, error) {
	return c.Client.Lstat(path)
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
