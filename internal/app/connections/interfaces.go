package connections

import (
	"context"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	conndb "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/db"
)

//go:generate go tool mockgen -destination mock_test.go -package connections . Repository
type Repository interface {
	Save(context.Context, int64, conndb.CreateConn) (*domain.Connection, error)
	FindByNodeId(context.Context, int64) ([]domain.Connection, error)
	FindById(context.Context, int64) (*domain.Connection, error)
	Update(ctx context.Context, id int64, updateConn conndb.UpdateConn) error
}
