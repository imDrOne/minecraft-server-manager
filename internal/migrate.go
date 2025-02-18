package internal

import (
	"errors"
	"fmt"
	"github.com/imDrOne/minecraft-server-manager/pkg/db"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/imDrOne/minecraft-server-manager/config"
)

const (
	_defaultAttempts = 5
	_defaultTimeout  = time.Second
)

var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

func MigrateUp(config *config.Config) {
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

	if err := MigrateUpWithConnectionString(connData.String()); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func MigrateUpWithConnectionString(connString string) error {
	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		migrationsPath := filepath.Join(basePath, "..", "db", "migrations")
		m, err = migrate.New(fmt.Sprintf("file://%s", migrationsPath), connString)
		if err == nil {
			break
		}

		slog.Warn(
			"migrate: postgres is trying to connect",
			slog.Int("attempts_left", attempts),
			slog.String("error", err.Error()))
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if m == nil || err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("migrate: up error", slog.String("error", err.Error()))
		return err
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			slog.Error("migrate: error during closing src", slog.String("error", err.Error()))
		}
		if dbErr != nil {
			slog.Error("migrate: error during closing db", slog.String("error", err.Error()))
		}
	}()

	slog.Info("Migrate: up success")
	return nil
}
