package lib

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log/slog"
	"sync"
	"time"
)

var (
	once          sync.Once
	pgContainer   *postgres.PostgresContainer
	connectionStr string
)

const waitLogStr = "database system is ready to accept connections"

type PgContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

func StartPostgresContainer(ctx context.Context) (*PgContainer, error) {
	var err error
	once.Do(func() {
		pgContainer, err = postgres.Run(
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

		connectionStr, err = pgContainer.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			return
		}
	})

	return &PgContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connectionStr,
	}, nil
}

func StopPostgresContainer(ctx context.Context) {
	if pgContainer != nil {
		if err := pgContainer.Terminate(ctx); err != nil {
			slog.Error(err.Error())
			return
		}
		once = sync.Once{}
	}
}

func GetPgConnectionString() string {
	return connectionStr
}
