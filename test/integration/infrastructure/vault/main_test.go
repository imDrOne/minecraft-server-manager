package vault

import (
	"context"
	"github.com/imDrOne/minecraft-server-manager/test/lib"
	"log/slog"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	if err := lib.StartVaultContainer(ctx); err != nil {
		slog.Error(err.Error())
		panic("error during starting docker container")
	}

	code := m.Run()

	os.Exit(code)
}
