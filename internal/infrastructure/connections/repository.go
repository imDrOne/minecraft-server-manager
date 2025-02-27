package connections

import (
	"context"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionQueries interface {
	FindConnectionsById(ctx context.Context, nodeID int64) ([]repository.Connection, error)
	SaveConnection(ctx context.Context, arg repository.SaveConnectionParams) (repository.SaveConnectionRow, error)
	UpdateConnectionById(ctx context.Context, arg repository.UpdateConnectionByIdParams) error
	CheckExistsConnection(ctx context.Context, arg string) (bool, error)
}

type ConnectionRepository struct {
	q ConnectionQueries
}

func (r ConnectionRepository) Save(ctx context.Context) {

}

func NewConnectionRepository(p *pgxpool.Pool) *ConnectionRepository {
	return &ConnectionRepository{q: repository.New(p)}
}
