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

type ClientSetup func(client *vault.Client) error

func New(addr string, setup ...ClientSetup) (*vault.Client, error) {
	client, err := vault.New(
		vault.WithAddress(addr),
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

func NewWithConfig(cfg config.Vault) (*vault.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Address, cfg.Port)

	client, err := New(addr, func(client *vault.Client) error {
		return client.SetToken(cfg.Token)
	})
	if err != nil {
		return nil, err
	}

	return client, err
}
