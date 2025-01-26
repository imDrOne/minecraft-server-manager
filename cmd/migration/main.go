package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/imDrOne/minecraft-server-manager/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
)

func main() {
	dbConfig := config.New().DB
	connString := dbConfig.BuildConnectionString("disable", map[string]string{})

	db, err := sql.Open("pgx", connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during PG migration preparing: %v\n", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		dbConfig.Name, driver)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during preparing Migration: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	err = m.Up()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during UP migration: %v\n", err)
		os.Exit(1)
	}
}
