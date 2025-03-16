package nodes

import (
	"net/http"
)

func NewHandler(handler NodeController) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("POST /", handler.Create)
	router.HandleFunc("GET /{id}", handler.GetById)
	router.HandleFunc("POST /pageable", handler.GetPaginated)
	return router
}
