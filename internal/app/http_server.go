package app

import (
	"github.com/imDrOne/minecraft-server-manager/internal/app/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/app/nodes"
	"net/http"
)

func SetupHttpServer(nodesRepo nodes.Repository, connRepo connections.Repository) *http.ServeMux {
	root := http.NewServeMux()

	nodeController := nodes.NewController(nodesRepo)
	connController := connections.NewController(connRepo)

	nodeHandler := nodes.NewHandler(nodeController)
	connHandler := connections.NewHandler(connController)

	root.Handle("/nodes/", http.StripPrefix("/nodes", nodeHandler))
	root.Handle("/connections/", http.StripPrefix("/connections", connHandler))

	return root
}
