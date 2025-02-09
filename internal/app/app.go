package app

import (
	"fmt"
	"github.com/imDrOne/minecraft-server-manager/config"
	node "github.com/imDrOne/minecraft-server-manager/internal/node/controller/http"
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
		slog.Error("%w", err)
		os.Exit(1)
	}

	nodeRouter := node.NewRouter(db.Pool)
	http.ListenAndServe(fmt.Sprintf(":%s", config.HTTPServer.Port), nodeRouter)
}
