package lib

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

var (
	oncePg      sync.Once
	pgContainer atomic.Value
)

const waitLogStr = "database system is ready to accept connections"

type PgContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func StartPostgresContainer(ctx context.Context) (err error) {
	oncePg.Do(func() {
		var container *postgres.PostgresContainer
		container, err = postgres.Run(
			ctx,
			"postgres:17.3-alpine",
			testcontainers.WithWaitStrategy(
				wait.ForLog(waitLogStr).
					WithOccurrence(2).
					WithStartupTimeout(5*time.Second)),
		)
		if err != nil {
			slog.Error("failed to start container", slog.String("error", err.Error()))
			return
		}

		var connectionStr string
		connectionStr, err = container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			return
		}

		pgContainer.Store(PgContainer{
			PostgresContainer: container,
			ConnectionString:  connectionStr,
		})
	})

	return err
}

func GetPgContainer() PgContainer {
	value := pgContainer.Load()
	if value == nil {
		return PgContainer{}
	}
	return value.(PgContainer)
}
