package connections

import (
	"encoding/json"
	"errors"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	"net/http"
)

type ConnectionController struct {
	repo Repository
}

func NewController(repo Repository) ConnectionController {
	return ConnectionController{repo: repo}
}

func (c ConnectionController) Create(w http.ResponseWriter, r *http.Request) {
	var connDto CreateConnectionRequestDto
	if err := json.NewDecoder(r.Body).Decode(&connDto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	conn, err := c.repo.Save(r.Context(), connDto.NodeId, func() (*domain.Connection, error) {
		return domain.CreateConnection(connDto.Key, connDto.User)
	})
	if err != nil {
		msg := err.Error()
		statusCode := http.StatusInternalServerError

		switch {
		case errors.Is(err, domain.ErrValidationConnection):
			statusCode = http.StatusBadRequest
		case errors.Is(err, domain.ErrConnectionAlreadyExists):
			statusCode = http.StatusConflict
		}

		http.Error(w, msg, statusCode)
		return
	}

	j, err := json.Marshal(ConnectionResponseDto{
		Id:        conn.Id(),
		Key:       conn.Key(),
		User:      conn.User(),
		CreatedAt: conn.CreatedAt(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
