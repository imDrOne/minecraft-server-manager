package connections

import "net/http"

func NewHandler(handler ConnectionController) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("POST /", handler.Create)
	return router
}
