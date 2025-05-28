package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

type ClientConfig struct {
	Auth    Auth
	Host    string
	Port    int64
	User    string
	Timeout time.Duration
}

func ProvideSshClient(cfg ClientConfig) (*ssh.Client, error) {
	authMethod, err := cfg.Auth.ToSSHAuthMethod()
	if err != nil {
		return nil, fmt.Errorf("error on creating ssh client for node-connection{%s@%s:%d}: %w", cfg.User, cfg.Host, cfg.Port, err)
	}
	config := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         cfg.Timeout,
	}
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	return ssh.Dial("tcp", addr, config)
}
