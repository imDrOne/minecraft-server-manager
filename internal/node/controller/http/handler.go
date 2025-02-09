package http

import (
	"context"
	"encoding/json"
	"errors"
	rawrepo "github.com/imDrOne/minecraft-server-manager/internal/generated/repository"
	"github.com/imDrOne/minecraft-server-manager/internal/node"
	"github.com/imDrOne/minecraft-server-manager/internal/node/controller/http/dto"
	"github.com/imDrOne/minecraft-server-manager/pkg/pagination"
	pdto "github.com/imDrOne/minecraft-server-manager/pkg/pagination/dto"
	"log/slog"
	"net/http"
)

type Handler struct {
	service Service
}

type Service interface {
	Save(ctx context.Context, arg rawrepo.SaveNodeParams) (rawrepo.Node, error)
	Find(ctx context.Context, payload pagination.PageRequest) (node.PagedNodes, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h Handler) save(w http.ResponseWriter, r *http.Request) {
	var nodeDto dto.CreateNodeDto
	if err := json.NewDecoder(r.Body).Decode(&nodeDto); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	val, err := h.service.Save(r.Context(), rawrepo.SaveNodeParams{
		Host: nodeDto.Host,
		Port: int32(nodeDto.Port),
	})

	if err != nil {
		if errors.Is(err, node.ErrNodeAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	resp, err := json.Marshal(dto.NodeDto{
		Id:        val.ID,
		Host:      val.Host,
		Port:      val.Port,
		CreatedAt: val.CreatedAt.Time,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func (h Handler) find(w http.ResponseWriter, r *http.Request) {
	var p dto.FindNodeDto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	pagedNode, err := h.service.Find(r.Context(), pagination.PageRequest{
		Page: p.Page,
		Size: p.Size,
	})

	if err != nil {
		slog.Error("Error during fetching data: %w", err)
		http.Error(w, "Error during fetching data", http.StatusInternalServerError)
		return
	}

	nodesDto := make([]dto.NodeDto, 0, pagedNode.Meta.Size)
	for _, v := range pagedNode.Nodes {
		nodesDto = append(nodesDto, dto.NodeDto{
			Id:        v.ID,
			Host:      v.Host,
			Port:      v.Port,
			CreatedAt: v.CreatedAt.Time,
		})
	}

	response := pdto.PageResponseWrapDto{
		Data: nodesDto,
		Meta: pdto.PagePaginationMetaDto(pagedNode.Meta),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error during fetching data", http.StatusInternalServerError)
	}
}
