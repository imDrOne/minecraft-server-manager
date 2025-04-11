package connections

import (
	"context"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections"
)

//go:generate go tool mockgen -destination mock_test.go -package connections . Repository
type Repository interface {
	Save(context.Context, int64, connections.CreateConn) (*domain.Connection, error)
	FindByNodeId(context.Context, int64) ([]domain.Connection, error)
	FindById(context.Context, int64) (*domain.Connection, error)
	Update(ctx context.Context, id int64, updateConn connections.UpdateConn) error
}
