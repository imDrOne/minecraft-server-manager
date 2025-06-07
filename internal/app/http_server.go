package app

import (
	"github.com/imDrOne/minecraft-server-manager/internal/app/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/app/nodes"
	"github.com/imDrOne/minecraft-server-manager/internal/app/remotes"
	connsshfacade "github.com/imDrOne/minecraft-server-manager/internal/service"
	connservice "github.com/imDrOne/minecraft-server-manager/internal/service/connections"
	sshservice "github.com/imDrOne/minecraft-server-manager/internal/service/ssh"
	"net/http"
	"time"
)

type Dependencies struct {
	NodesRepo      nodes.Repository
	ConnRepo       connections.Repository
	ConnSshKeyRepo connections.ConnectionSshKeyRepository
	KeygenService  *sshservice.KeygenService
}

func SetupHttpServer(deps Dependencies) *http.ServeMux {
	root := http.NewServeMux()

	connSshService := sshservice.NewSshService(time.Minute * 5)
	connSshFacade := connsshfacade.NewConnectionSshFacade(connsshfacade.Dependencies{
		NodeRepo:       deps.NodesRepo,
		ConnRepo:       deps.ConnRepo,
		ConnSshKeyRepo: deps.ConnSshKeyRepo,
		SshService:     connSshService,
	})
	connService := connservice.NewConnectionService(connservice.Dependencies{
		NodeRepo:         deps.NodesRepo,
		ConnRepo:         deps.ConnRepo,
		ConnSshKeyRepo:   deps.ConnSshKeyRepo,
		SshKeygenService: deps.KeygenService,
	})
	connController := connections.NewController(connService)
	nodeController := nodes.NewController(deps.NodesRepo)
	remoteController := remotes.NewController(connSshFacade)

	nodeHandler := nodes.NewHandler(nodeController)
	connHandler := connections.NewHandler(connController)
	remoteHandler := remotes.NewHandler(remoteController)

	root.Handle("/nodes/", http.StripPrefix("/nodes", nodeHandler))
	root.Handle("/connections/", http.StripPrefix("/connections", connHandler))
	root.Handle("/remote/", http.StripPrefix("/remote", remoteHandler))

	return root
}
