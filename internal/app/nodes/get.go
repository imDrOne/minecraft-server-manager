package nodes

import (
	"encoding/json"
	"github.com/imDrOne/minecraft-server-manager/pkg/pagination"
	"net/http"
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
