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
	client, err := vault.New(
		vault.WithAddress(fmt.Sprintf("%s:%s", cfg.Address, cfg.Port)),
		vault.WithRequestTimeout(requestTimeout),
	)

	if err != nil {
		return Options{}, err
	}

	if err := client.SetToken(cfg.Token); err != nil {
		return Options{}, err
	}

	return Options{
		cfg,
		client,
	}, err
}
