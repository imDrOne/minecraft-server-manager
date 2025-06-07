package ssh

import (
	"context"
	"github.com/imDrOne/minecraft-server-manager/test/lib"
	"log/slog"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	if err := lib.StartSshdContainer(ctx); err != nil {
		slog.Error(err.Error())
		panic("error during starting sshd-docker container")
	}

	code := m.Run()

	os.Exit(code)
}
