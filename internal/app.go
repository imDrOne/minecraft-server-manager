package internal

import (
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal/app"
	conndb "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/db"
	connvt "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/vault"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/nodes"
	"github.com/imDrOne/minecraft-server-manager/internal/service/ssh"
	"github.com/imDrOne/minecraft-server-manager/pkg/db"
	"github.com/imDrOne/minecraft-server-manager/pkg/vault"
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
		slog.Error("error on db connection preparing: " + err.Error())
		os.Exit(1)
	}

	pg, err := db.New(connData, db.MaxPoolSize(config.DB.MaxPoolSiz), db.ConnAttempts(config.DB.ConnAttempts))
	defer pg.Close()
	if err != nil {
		slog.Error("error on db connecting: " + err.Error())
		os.Exit(1)
	}

	vaultClient, err := vault.NewWithConfig(config.Vault)
	if err != nil {
		panic(err)
	}

	nodeRepo := nodes.NewNodeRepository(pg.Pool)
	connRepo := conndb.NewConnectionRepository(pg.Pool)
	connSshRepo := connvt.NewConnSshKeyRepository(vaultClient, config.Vault)
	keygenService := ssh.NewKeygenService(config.SSHKeygen)

	httpServer := app.SetupHttpServer(app.Dependencies{
		NodesRepo:      nodeRepo,
		ConnRepo:       connRepo,
		KeygenService:  keygenService,
		ConnSshKeyRepo: connSshRepo,
	})

	server := http.Server{
		Addr:    "0.0.0.0:" + config.HTTPServer.Port,
		Handler: httpServer,
	}
	server.ListenAndServe()
}
