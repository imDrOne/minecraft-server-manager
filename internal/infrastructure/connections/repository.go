package connections

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/query"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate go tool mockgen -destination mock_test.go -package connections . ConnectionQueries
type ConnectionQueries interface {
	FindConnectionById(context.Context, int64) (query.Connection, error)
	FindConnectionsByNodeId(context.Context, int64) ([]query.Connection, error)
	SaveConnection(context.Context, query.SaveConnectionParams) (query.SaveConnectionRow, error)
	UpdateConnectionById(context.Context, query.UpdateConnectionByIdParams) error
	CheckExistsConnection(context.Context, string) (bool, error)
}

type (
	CreateConn func() (*domain.Connection, error)
	UpdateConn func(domain.Connection) (*domain.Connection, error)
)
type ConnectionRepository struct {
	q ConnectionQueries
}

func (r ConnectionRepository) Save(ctx context.Context, nodeId int64, createConn CreateConn) (*domain.Connection, error) {
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
		NodeID:   nodeId,
		Key:      conn.Key(),
		User:     "test",
		Checksum: checksum,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert connection: %w", err)
	}

	return conn.WithDBGeneratedValues(data), nil
}

func (r ConnectionRepository) FindByNodeId(ctx context.Context, id int64) ([]domain.Connection, error) {
	data, err := r.q.FindConnectionsByNodeId(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrConnectionNotFound
		}
		return nil, fmt.Errorf("failed to select connections by node-id %d: %w", id, err)
	}

	connections := make([]domain.Connection, 0, len(data))
	for i, c := range data {
		mapped, err := domain.FromDbModel(c)
		if err != nil {
			return nil, fmt.Errorf("failed to map connection by id %d: %w", data[i].ID, err)
		}
		connections = append(connections, *mapped)
	}

	return connections, nil
}

func (r ConnectionRepository) FindById(ctx context.Context, id int64) (*domain.Connection, error) {
	data, err := r.q.FindConnectionById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrConnectionNotFound
		}
		return nil, fmt.Errorf("failed to select connection by id %d: %w", id, err)
	}
	conn, err := domain.FromDbModel(data)
	if err != nil {
		return nil, fmt.Errorf("failed to map conn by id %d: %w", id, err)
	}

	return conn, nil
}

func (r ConnectionRepository) Update(ctx context.Context, id int64, updateConn UpdateConn) error {
	conn, err := r.FindById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get connection by id=%d: %w", id, err)
	}

	conn, err = updateConn(*conn)
	if err != nil {
		return fmt.Errorf("failed to update connection by id=%d %w", id, err)
	}

	err = r.q.UpdateConnectionById(ctx, query.UpdateConnectionByIdParams{
		ID:   id,
		Key:  conn.Key(),
		User: conn.User(),
	})
	if err != nil {
		return fmt.Errorf("failed to update connection by id=%d query: %w", id, err)
	}
	return nil
}

func NewConnectionRepository(p *pgxpool.Pool) *ConnectionRepository {
	return &ConnectionRepository{q: query.New(p)}
}
