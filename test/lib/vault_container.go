package lib

import (
	"context"
	"github.com/testcontainers/testcontainers-go/modules/vault"
	"log/slog"
)

var (
	VaultHostAddress string
)

func StartVaultContainer(ctx context.Context) error {
	vaultContainer, err := vault.Run(ctx, "hashicorp/vault:1.13.0", vault.WithToken("root"))

	if err != nil {
		slog.Error("failed to start container: %s", err)
		return err
	}
	VaultHostAddress, err = vaultContainer.HttpHostAddress(ctx)
	if err != nil {
		slog.Error("failed to get host address: %s", err)
		return err
	}

	return nil
}
