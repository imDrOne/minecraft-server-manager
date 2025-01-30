package main

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/imDrOne/minecraft-server-manager/config"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func main() {
	dbConfig := config.New().DB
	connString := dbConfig.BuildConnectionString("disable", map[string]string{})

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", connString)
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

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Migrate: up error: %s", err)
		os.Exit(1)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		slog.Warn("Migrate: no change")
		return
	}

	log.Printf("Migrate: up success")
}
