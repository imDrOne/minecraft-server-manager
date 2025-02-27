package connections

import (
	"context"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/query"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionQueries interface {
	FindConnectionsById(ctx context.Context, nodeID int64) ([]query.Connection, error)
	SaveConnection(ctx context.Context, arg query.SaveConnectionParams) (query.SaveConnectionRow, error)
	UpdateConnectionById(ctx context.Context, arg query.UpdateConnectionByIdParams) error
	CheckExistsConnection(ctx context.Context, arg string) (bool, error)
}

type ConnectionRepository struct {
	q ConnectionQueries
}

func (r ConnectionRepository) Save(ctx context.Context) {

}

func NewConnectionRepository(p *pgxpool.Pool) *ConnectionRepository {
	return &ConnectionRepository{q: query.New(p)}
}
