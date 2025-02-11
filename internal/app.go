package internal

import (
	"fmt"
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal/app/nodes"
	"github.com/imDrOne/minecraft-server-manager/pkg/db"
	"log/slog"
	"net/http"
	"os"
)

func Run(config *config.Config) {
	slog.With("app.Run - postgres.New")
	connStr := config.DB.BuildConnectionString("disable", map[string]string{})
	db, err := db.New(connStr, db.MaxPoolSize(config.DB.MaxPoolSiz), db.ConnAttempts(config.DB.ConnAttempts))
	defer db.Close()

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	mux := http.NewServeMux()
	nodeRouter := nodes.NewRouter(db.Pool)
	mux.Handle("/nodes/", http.StripPrefix("/nodes", nodeRouter))

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", config.HTTPServer.Port),
		Handler: mux,
	}
	server.ListenAndServe()
}
