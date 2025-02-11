package nodes

import (
	repository "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/nodes"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func NewRouter(pool *pgxpool.Pool) *http.ServeMux {
	router := http.NewServeMux()

	repo := repository.NewPostgresRepo(pool)
	handler := NewHandlers(repo)

	router.HandleFunc("POST /", handler.Create)
	router.HandleFunc("POST /pageable", handler.GetPaginated)

	return router
}
