package internal

import (
	"errors"
	"github.com/imDrOne/minecraft-server-manager/pkg/db"
	"log/slog"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/imDrOne/minecraft-server-manager/config"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func MigrateUp(config *config.Config) {
	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	connData, err := db.NewConnectionData(
		config.DB.Host,
		config.DB.Name,
		config.DB.User,
		config.DB.Password,
		config.DB.Port,
		false,
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	for attempts > 0 {
		m, err = migrate.New("file://db/migrations", connData.String())
		if err == nil {
			break
		}

		slog.Warn("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if m == nil || err != nil {
		slog.Error("Migrate: postgres connect error: %s", err)
		os.Exit(1)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Migrate: up error: %s", err)
		os.Exit(1)
	}
	defer m.Close()

	slog.Info("Migrate: up success")
}
