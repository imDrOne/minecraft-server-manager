package connections

import (
	"encoding/json"
	"errors"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	"net/http"
	"strconv"
)

type ConnectionController struct {
	service Service
}

func NewController(service Service) ConnectionController {
	return ConnectionController{service: service}
}

func (c ConnectionController) Create(w http.ResponseWriter, r *http.Request) {
	var connDto CreateConnectionRequestDto
	if err := json.NewDecoder(r.Body).Decode(&connDto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	conn, err := c.service.Create(r.Context(), connDto.NodeId, func() (*domain.Connection, error) {
		return domain.CreateConnection(connDto.NodeId, connDto.User)
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
		Id: conn.Id(),
		//Key:       conn.Key(),
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

func (c ConnectionController) Update(w http.ResponseWriter, r *http.Request) {
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

	var connDto UpdateConnectionRequestDto
	if err := json.NewDecoder(r.Body).Decode(&connDto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = c.service.Update(r.Context(), int64(id), func(connection domain.Connection) (*domain.Connection, error) {
		return connection.Update(connDto.User)
	})

	if err != nil {
		statusCode := http.StatusInternalServerError

		switch {
		case errors.Is(err, domain.ErrValidationConnection):
			statusCode = http.StatusBadRequest
		case errors.Is(err, domain.ErrConnectionNotFound):
			statusCode = http.StatusNotFound
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (c ConnectionController) FindById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("node-id")
	if idStr == "" {
		http.Error(w, "expected id - got empty string", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "error during parsing id", http.StatusBadRequest)
		return
	}

	connections, err := c.service.FindByNodeId(r.Context(), int64(id))

	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domain.ErrConnectionNotFound) {
			statusCode = http.StatusNotFound
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	connectionsDto := make([]ConnectionResponseDto, 0, len(connections))
	for _, val := range connections {
		connectionsDto = append(connectionsDto, ConnectionResponseDto{
			Id: val.Id(),
			//Key:       val.Key(),
			User:      val.User(),
			CreatedAt: val.CreatedAt(),
		})
	}

	j, err := json.Marshal(connectionsDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
