package ssh

import (
	"bytes"
	"fmt"
	"github.com/imDrOne/minecraft-server-manager/internal/service/ssh/model"
	"github.com/imDrOne/minecraft-server-manager/internal/service/ssh/scripts"
	sshpkg "github.com/imDrOne/minecraft-server-manager/pkg/ssh"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

type SshService struct {
	sshClientTimeout time.Duration
}

func NewSshService(timeout time.Duration) *SshService {
	return &SshService{sshClientTimeout: timeout}
}

func (s *SshService) newClient(cfg model.NodeSSHConnectionTO) (*ssh.Client, error) {
	client, err := sshpkg.ProvideSshClient(sshpkg.ClientConfig{
		Auth:    cfg.Auth,
		Host:    cfg.Host,
		Port:    cfg.Port,
		User:    cfg.User,
		Timeout: s.sshClientTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("ssh connection failed: %w", err)
	}
	return client, nil
}

func (s *SshService) InjectPublicKey(cfg model.NodeSSHConnectionTO, publicKey string) error {
	client, err := s.newClient(cfg)
	if err != nil {
		return fmt.Errorf("ssh connection failed: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("error on creating ssh session: %w", err)
	}
	defer session.Close()

	escapedKey := strings.ReplaceAll(publicKey, "'", "'\"'\"'")

	var out bytes.Buffer
	session.Stdout = &out
	session.Stderr = &out
	session.Stdin = strings.NewReader(escapedKey)

	err = session.Run(scripts.InstallKeyScript)
	if err != nil {
		return fmt.Errorf("script failed: %v\n%s", err, out.String())
	}

	return nil
}

func (s *SshService) Ping(cfg model.NodeSSHConnectionTO) error {
	client, err := s.newClient(cfg)
	defer client.Close()

	if err != nil {
		return fmt.Errorf("error on trying ping connection=%v: %w", cfg, err)
	}
	return nil
}
