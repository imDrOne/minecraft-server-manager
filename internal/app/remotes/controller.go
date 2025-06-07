package remotes

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type RemoteController struct {
	sshConnectionService SshConnectionService
}

func NewController(sshConnectionService SshConnectionService) RemoteController {
	return RemoteController{sshConnectionService: sshConnectionService}
}

func (c RemoteController) ForwardPublicKey(w http.ResponseWriter, r *http.Request) {
	connIdStr := r.PathValue("connectionId")
	if connIdStr == "" {
		http.Error(w, "expected id - got empty string", http.StatusBadRequest)
		return
	}

	connId, err := strconv.Atoi(connIdStr)
	if err != nil {
		http.Error(w, "error during parsing id", http.StatusBadRequest)
		return
	}

	var forwardedPublicKeyDto ForwardPublicKeyDto
	if err = json.NewDecoder(r.Body).Decode(&forwardedPublicKeyDto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = c.sshConnectionService.InjectPublicKey(r.Context(), int64(connId), forwardedPublicKeyDto)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error during injecting public key", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
