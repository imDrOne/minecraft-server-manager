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

func WithConfig(cfg config.Vault) (Options, error) {
	host := fmt.Sprintf("%s:%s", cfg.Address, cfg.Port)

	client, err := New(host, func(client *vault.Client) error {
		return client.SetToken(cfg.Token)
	})

	if err != nil {
		return Options{}, err
	}

	return Options{
		cfg,
		client,
	}, err
}

func New(host string, setup ...ClientSetup) (*vault.Client, error) {
	client, err := vault.New(
		vault.WithAddress(host),
		vault.WithRequestTimeout(requestTimeout),
	)

	if err != nil {
		return nil, err
	}

	for _, clientSetup := range setup {
		err = clientSetup(client)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

type ClientSetup func(client *vault.Client) error
