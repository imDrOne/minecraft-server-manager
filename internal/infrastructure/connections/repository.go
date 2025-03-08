package connections

import (
	"context"
	"fmt"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/query"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate go tool mockgen -destination mock_test.go -package connections . ConnectionQueries
type ConnectionQueries interface {
	FindConnectionsById(context.Context, int64) ([]query.Connection, error)
	SaveConnection(context.Context, query.SaveConnectionParams) (query.SaveConnectionRow, error)
	UpdateConnectionById(context.Context, query.UpdateConnectionByIdParams) error
	CheckExistsConnection(ctx context.Context, checksum string) (bool, error)
}

type ConnectionRepository struct {
	q ConnectionQueries
}

func (r ConnectionRepository) Save(ctx context.Context, createConn func() (*domain.Connection, error)) (*domain.Connection, error) {
	conn, err := createConn()
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	checksum, err := conn.ChecksumStr()
	idExists, err := r.q.CheckExistsConnection(ctx, checksum)
	if err != nil {
		return nil, fmt.Errorf("failed to check connection exist: %w", err)
	}
	if idExists {
		return nil, domain.ErrConnectionAlreadyExists
	}

	data, err := r.q.SaveConnection(ctx, query.SaveConnectionParams{
		NodeID:   conn.Id(),
		Key:      conn.Key(),
		User:     pgtype.Text{},
		Checksum: checksum,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert connection: %w", err)
	}

	return conn.WithDBGeneratedValues(data), nil
}

func NewConnectionRepository(p *pgxpool.Pool) *ConnectionRepository {
	return &ConnectionRepository{q: query.New(p)}
}
