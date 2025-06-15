package ssh

import (
	"fmt"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"time"
)

type ClientConfig struct {
	Auth    Auth
	Host    string
	Port    uint
	User    string
	Timeout time.Duration
}

func ProvideSshClient(cfg ClientConfig) (*goph.Client, error) {
	authMethod, err := cfg.Auth.ToSSHAuthMethod()
	if err != nil {
		return nil, fmt.Errorf("error on creating ssh client for node-connection{%s@%s:%d}: %w", cfg.User, cfg.Host, cfg.Port, err)
	}

	return goph.NewConn(&goph.Config{
		Auth:     authMethod,
		User:     cfg.User,
		Addr:     cfg.Host,
		Port:     cfg.Port,
		Timeout:  cfg.Timeout,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
}
