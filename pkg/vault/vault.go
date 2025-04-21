package vault

import (
	"fmt"
	"github.com/hashicorp/vault-client-go"
	"github.com/imDrOne/minecraft-server-manager/config"
	"time"
)

const (
	requestTimeout = 30 * time.Second
)

type Options struct {
	config config.Vault
	client *vault.Client
}

func New(cfg config.Vault) (Options, error) {
	host := fmt.Sprintf("%s:%s", cfg.Address, cfg.Port)

	client, err := create(host, cfg.Token)

	if err != nil {
		return Options{}, err
	}

	return Options{
		cfg,
		client,
	}, err
}

func create(host string, token string) (*vault.Client, error) {
	client, err := vault.New(
		vault.WithAddress(host),
		vault.WithRequestTimeout(requestTimeout),
	)

	if err != nil {
		return nil, err
	}

	if err := client.SetToken(token); err != nil {
		return nil, err
	}

	return client, nil
}
