package lib

import (
	"context"
	"github.com/testcontainers/testcontainers-go/modules/vault"
	"log/slog"
	"sync"
	"sync/atomic"
)

var (
	onceVault      sync.Once
	vaultContainer atomic.Value
)

type VaultContainer struct {
	*vault.VaultContainer
	HostAddress string
}

func StartVaultContainer(ctx context.Context) (err error) {
	onceVault.Do(func() {
		var container *vault.VaultContainer
		container, err = vault.Run(ctx, "hashicorp/vault:1.13.0", vault.WithToken("root"))
		if err != nil {
			slog.Error("failed to start container: %s", err)
			return
		}

		var addr string
		addr, err = container.HttpHostAddress(ctx)
		if err != nil {
			slog.Error("failed to get host address: %s", err)
			return
		}

		vaultContainer.Store(VaultContainer{
			VaultContainer: container,
			HostAddress:    addr,
		})
	})

	return err
}

func GetVaultContainer() VaultContainer {
	value := vaultContainer.Load()
	if value == nil {
		return VaultContainer{}
	}
	return value.(VaultContainer)
}
