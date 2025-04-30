package app

import (
	"github.com/imDrOne/minecraft-server-manager/internal/app/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/app/nodes"
	connservice "github.com/imDrOne/minecraft-server-manager/internal/service/connections"
	sshservice "github.com/imDrOne/minecraft-server-manager/internal/service/ssh"
	"net/http"
)

type Dependencies struct {
	NodesRepo      nodes.Repository
	ConnRepo       connections.Repository
	ConnSshKeyRepo connections.ConnectionSshKeyRepository
	KeygenService  *sshservice.KeygenService
}

func SetupHttpServer(deps Dependencies) *http.ServeMux {
	root := http.NewServeMux()

	nodeController := nodes.NewController(deps.NodesRepo)
	connService := connservice.NewConnectionService(connservice.Dependencies{
		NodeRepo:         deps.NodesRepo,
		ConnRepo:         deps.ConnRepo,
		ConnSshKeyRepo:   deps.ConnSshKeyRepo,
		SshKeygenService: deps.KeygenService,
	})
	connController := connections.NewController(connService)

	nodeHandler := nodes.NewHandler(nodeController)
	connHandler := connections.NewHandler(connController)

	root.Handle("/nodes/", http.StripPrefix("/nodes", nodeHandler))
	root.Handle("/connections/", http.StripPrefix("/connections", connHandler))

	return root
}
