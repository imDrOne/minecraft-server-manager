package nodes

import (
	"encoding/json"
	"errors"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/imDrOne/minecraft-server-manager/pkg/pagination"
	"net/http"
	"strconv"
)

func (h NodeHandlers) GetPaginated(w http.ResponseWriter, r *http.Request) {
	var p FindNodeRequestDto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	pagedNodes, err := h.repo.FindPaginated(r.Context(), p.ToValue())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nodesDto := make([]NodeResponseDto, 0, pagedNodes.Meta.Size())
	for _, v := range pagedNodes.Data {
		nodesDto = append(nodesDto, NodeResponseDto{
			Id:        v.Id(),
			Host:      v.Host(),
			Port:      int32(v.Port()),
			CreatedAt: v.CreatedAt(),
		})
	}

	j, err := json.Marshal(pagination.PagePaginationResponseWrapDto{
		Data: nodesDto,
		Meta: pagedNodes.Meta.ToDTO(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (h NodeHandlers) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "expected id - got empty string", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "error during parsing id", http.StatusBadRequest)
		return
	}

	node, err := h.repo.FindById(r.Context(), int64(id))
	if err != nil {
		msg := err.Error()
		if errors.Is(err, domain.ErrNodeNotFound) {
			http.Error(w, msg, http.StatusNotFound)
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
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
