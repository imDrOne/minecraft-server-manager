package nodes

import (
	"encoding/json"
	"errors"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"net/http"
)

func (h NodeHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var nodeDto CreateNodeRequestDto
	if err := json.NewDecoder(r.Body).Decode(&nodeDto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	node, err := h.repo.Save(r.Context(), func() (*domain.Node, error) {
		return domain.CreateNode(nodeDto.Host, uint(nodeDto.Port))
	})
	if err != nil {
		msg := err.Error()
		if errors.Is(err, domain.ErrValidationNode) {
			http.Error(w, msg, http.StatusBadRequest)
		}
		if errors.Is(err, domain.ErrNodeAlreadyExist) {
			http.Error(w, msg, http.StatusConflict)
		} else {
			http.Error(w, msg, http.StatusInternalServerError)
		}
		return
	}

	j, err := json.Marshal(NodeResponseDto{
		Id:        node.Id(),
		Host:      node.Host(),
		Port:      int32(node.Port()),
		CreatedAt: node.CreatedAt(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
