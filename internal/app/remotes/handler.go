package remotes

import "net/http"

func NewHandler(handler RemoteController) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc(
		"POST connections/{connection-id}/forward-public-key",
		handler.ForwardPublicKey)

	return router
}
