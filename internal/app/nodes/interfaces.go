package nodes

import (
	"context"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/imDrOne/minecraft-server-manager/pkg/pagination"
)

type Repository interface {
	Save(ctx context.Context, createNode func() (*domain.Node, error)) (*domain.Node, error)
	FindPaginated(ctx context.Context, req pagination.PageRequest) (*domain.PagePaginatedNodes, error)
}
