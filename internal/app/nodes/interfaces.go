package nodes

import (
	"context"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/imDrOne/minecraft-server-manager/pkg/pagination"
)

type Repository interface {
	Save(context.Context, func() (*domain.Node, error)) (*domain.Node, error)
	FindPaginated(context.Context, pagination.PageRequest) (*domain.PagePaginatedNodes, error)
	FindById(context.Context, int64) (*domain.Node, error)
}
