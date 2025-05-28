package connections

import (
	"context"
	"github.com/imDrOne/minecraft-server-manager/internal"
	"github.com/imDrOne/minecraft-server-manager/test/lib"
	"log/slog"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	if err := lib.StartPostgresContainer(ctx); err != nil {
		slog.Error(err.Error())
		panic("error during starting pg-docker container")
	}

	pgContainer := lib.GetPgContainer()
	if err := internal.MigrateUpWithConnectionString(pgContainer.ConnectionString); err != nil {
		slog.Error(err.Error())
		panic("error during running migrations")
	}

	if err := lib.StartVaultContainer(ctx); err != nil {
		slog.Error(err.Error())
		panic("error during starting vault-docker container")
	}

	code := m.Run()

	os.Exit(code)
}
