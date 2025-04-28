package internal

import (
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal/app"
	conndb "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/db"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/nodes"
	"github.com/imDrOne/minecraft-server-manager/pkg/db"
	"log/slog"
	"net/http"
	"os"
)

func Run(config *config.Config) {
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

	pg, err := db.New(connData, db.MaxPoolSize(config.DB.MaxPoolSiz), db.ConnAttempts(config.DB.ConnAttempts))
	defer pg.Close()

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	nodeRepo := nodes.NewNodeRepository(pg.Pool)
	connRepo := conndb.NewConnectionRepository(pg.Pool)
	httpServer := app.SetupHttpServer(nodeRepo, connRepo)

	server := http.Server{
		Addr:    "0.0.0.0:" + config.HTTPServer.Port,
		Handler: httpServer,
	}
	server.ListenAndServe()
}
