package app

import (
	"github.com/imDrOne/minecraft-server-manager/internal/app/nodes"
	"net/http"
)

func SetupHttpServer(nodesRepo nodes.Repository) *http.ServeMux {
	mux := http.NewServeMux()

	nodeRouter := nodes.NewRouter(nodesRepo)
	mux.Handle("/nodes/", http.StripPrefix("/nodes", nodeRouter))

	return mux
}
