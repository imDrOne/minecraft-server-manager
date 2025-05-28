package lib

import (
	"context"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

var (
	onceSshd      sync.Once
	sshdContainer atomic.Value
)

type SshdContainer struct {
	testcontainers.Container
	Host string
	Port int
}

func StartSshdContainer(ctx context.Context) (err error) {
	onceSshd.Do(func() {
		var container testcontainers.Container
		var host string
		var port nat.Port

		req := testcontainers.ContainerRequest{
			Image:        "testcontainers/sshd:1.2.0",
			ExposedPorts: []string{"22/tcp"},
			Env: map[string]string{
				"PASSWORD": "test",
			},
			WaitingFor: wait.ForListeningPort("22/tcp").WithStartupTimeout(10 * time.Second),
		}

		container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			slog.Error("failed to start container", slog.String("error", err.Error()))
			return
		}

		host, err = container.Host(ctx)
		if err != nil {
			slog.Error("failed to get container host", slog.String("error", err.Error()))
			return
		}

		port, err = container.MappedPort(ctx, "22")
		if err != nil {
			slog.Error("failed to get container port", slog.String("error", err.Error()))
			return
		}

		sshdContainer.Store(SshdContainer{
			Container: container,
			Host:      host,
			Port:      port.Int(),
		})
	})

	return err
}

func GetSshdContainer() SshdContainer {
	value := sshdContainer.Load()
	if value == nil {
		return SshdContainer{}
	}
	return value.(SshdContainer)
}
