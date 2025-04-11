package connections

import "net/http"

func NewHandler(handler ConnectionController) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("POST /", handler.Create)
	router.HandleFunc("PUT /{id}", handler.Update)
	return router
}
