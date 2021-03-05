package sftp

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Addr     string
	User     string
	Password string
	Key      string
}

func (conf *SSHConfig) Validate() (*ssh.ClientConfig, error) {
	if _, _, err := net.SplitHostPort(conf.Addr); err != nil {
		conf.Addr = net.JoinHostPort(conf.Addr, "22")
		if _, _, err = net.SplitHostPort(conf.Addr); err != nil {
			return nil, err
		}
	}

	if conf.User == "" {
		if conf.User = os.Getenv("USER"); conf.User == "" {
			return nil, errors.New("ssh user is not set")
		}
	}

	var auth ssh.AuthMethod
	if conf.Password != "" {
		auth = ssh.Password(conf.Password)

	} else {
		if conf.Key == "" {
			if conf.Key = GetUniqSSHKeyPath(); conf.Key == "" {
				return nil, errors.New("ssh key is not set")
			}
		}

		keyData, err := ioutil.ReadFile(conf.Key)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		privateKey, err := ssh.ParsePrivateKey(keyData)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		auth = ssh.PublicKeys(privateKey)
	}

	return &ssh.ClientConfig{
		User:            conf.User,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

// GetUniqSSHKeyPath the uniq private key path
func GetUniqSSHKeyPath() string {
	home, _ := os.UserHomeDir()
	fis, err := ioutil.ReadDir(filepath.Join(home, "/.ssh"))
	if err != nil {
		return ""
	}

	keys := make([]string, 0, 1)
	for f := range fis {
		filename := filepath.Base(fis[f].Name())
		if len(filename) >= 4 &&
			filename[:3] == "id_" &&
			filename[len(filename)-4:] != ".pub" {

			keys = append(keys, home+"/.ssh/"+filename)
		}
	}
	if len(keys) == 1 {
		return keys[0]
	}
	return ""
}
