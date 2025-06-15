package ssh

import (
	"fmt"
	"github.com/imDrOne/minecraft-server-manager/internal/service/ssh/model"
	"github.com/imDrOne/minecraft-server-manager/internal/service/ssh/scripts"
	sshpkg "github.com/imDrOne/minecraft-server-manager/pkg/ssh"
	"github.com/melbahja/goph"
	"strings"
	"time"
)

type Service struct {
	sshClientTimeout time.Duration
}

func NewSshService(timeout time.Duration) *Service {
	return &Service{sshClientTimeout: timeout}
}

func (s *Service) newClient(cfg model.NodeSSHConnectionTO) (*goph.Client, error) {
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

func (s *Service) InjectPublicKey(cfg model.NodeSSHConnectionTO, publicKey string) error {
	client, err := s.newClient(cfg)
	if err != nil {
		return fmt.Errorf("ssh connection failed: %w", err)
	}
	defer client.Close()

	command, err := client.Command(scripts.InstallKeyScript)
	if err != nil {
		return fmt.Errorf("error during preparing injecting key command: %w", err)
	}

	command.Stdin = strings.NewReader(publicKey)
	result, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("script failed with details %v: %w", string(result), err)
	}

	return nil
}

func (s *Service) Ping(cfg model.NodeSSHConnectionTO) error {
	client, err := s.newClient(cfg)
	if err != nil {
		return fmt.Errorf("error on trying ping connection=%v: %w", cfg, err)
	}
	defer client.Close()

	return nil
}
