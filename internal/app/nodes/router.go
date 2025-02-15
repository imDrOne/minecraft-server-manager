package nodes

import (
	"net/http"
)

type NodeHandlers struct {
	repo Repository
}

func NewHandler(repo Repository) NodeHandlers {
	return NodeHandlers{repo}
}

func NewRouter(nodeRepository Repository) *http.ServeMux {
	router := http.NewServeMux()
	handler := NewHandler(nodeRepository)
	router.HandleFunc("POST /", handler.Create)
	router.HandleFunc("POST /pageable", handler.GetPaginated)

	return router
}
