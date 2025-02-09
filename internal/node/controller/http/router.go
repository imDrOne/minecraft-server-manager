package http

import (
	"github.com/imDrOne/minecraft-server-manager/internal/node"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func NewRouter(pool *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	repo := node.NewRepositoryImpl(pool)
	service := node.NewService(repo)
	handler := NewHandler(service)

	mux.HandleFunc("POST /nodes/", handler.save)
	mux.HandleFunc("POST /nodes/pageable", handler.find)

	return mux
}
