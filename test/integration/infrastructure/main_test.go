package infrastructure

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

	if _, err := lib.StartPostgresContainer(ctx); err != nil {
		slog.Error(err.Error())
		panic("error during starting container")
	}

	pgConnStr := lib.GetPgConnectionString()
	if err := internal.MigrateUpWithConnectionString(pgConnStr); err != nil {
		slog.Error(err.Error())
		panic("error during running migrations")
	}

	code := m.Run()

	os.Exit(code)
}
