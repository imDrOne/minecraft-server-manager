package nodes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/query"
	"github.com/imDrOne/minecraft-server-manager/pkg/pagination"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate go tool mockgen -destination mock_test.go -package nodes . NodeQueries
type NodeQueries interface {
	CheckExistsNode(context.Context, query.CheckExistsNodeParams) (bool, error)
	SaveNode(context.Context, query.SaveNodeParams) (query.SaveNodeRow, error)
	UpdateNodeById(context.Context, query.UpdateNodeByIdParams) error
	FindNodeById(context.Context, int64) (query.Node, error)
	FindNodes(context.Context, query.FindNodesParams) ([]query.Node, error)
	CountNode(context.Context) (int64, error)
}

type NodeRepository struct {
	q NodeQueries
}

func NewNodeRepository(p *pgxpool.Pool) *NodeRepository {
	return &NodeRepository{q: query.New(p)}
}

func (r NodeRepository) Save(ctx context.Context, createNode func() (*domain.Node, error)) (*domain.Node, error) {
	node, err := createNode()
	if err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	isExists, err := r.q.CheckExistsNode(ctx, query.CheckExistsNodeParams{
		Host: node.Host(),
		Port: int32(node.Port()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check node exist: %w", err)
	}
	if isExists {
		return nil, domain.ErrNodeAlreadyExist
	}

	data, err := r.q.SaveNode(ctx, query.SaveNodeParams{
		Host: node.Host(),
		Port: int32(node.Port()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert node: %w", err)
	}

	return node.WithDBGeneratedValues(data), nil
}

func (r NodeRepository) Update(ctx context.Context, id int64, updateNode func(*domain.Node) (*domain.Node, error)) error {
	node, err := r.FindById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get node: %w", err)
	}

	node, err = updateNode(node)
	if err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}

	err = r.q.UpdateNodeById(ctx, query.UpdateNodeByIdParams{
		ID:   node.Id(),
		Host: node.Host(),
		Port: int32(node.Port()),
	})
	if err != nil {
		return fmt.Errorf("failed to update node query: %w", err)
	}

	return nil
}

func (r NodeRepository) FindPaginated(ctx context.Context, req pagination.PageRequest) (*domain.PagePaginatedNodes, error) {
	nodes, err := r.Find(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to select pages nodes: %w", err)
	}

	total, err := r.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.PagePaginatedNodes{
		Data: *nodes,
		Meta: req.ToPageMeta(uint64(total)),
	}, nil
}

// Find todo: Change return data type, from *[]domain.Node to []domain.Node
func (r NodeRepository) Find(ctx context.Context, pagination pagination.PageRequest) (*[]domain.Node, error) {
	data, err := r.q.FindNodes(ctx, query.FindNodesParams{
		Limit:  int32(pagination.Size()),
		Offset: int32(pagination.Offset()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to select nodes: %w", err)
	}

	nodes := make([]domain.Node, 0, pagination.Size())
	for _, n := range data {
		mapped, err := domain.FromDbModel(n)
		if err != nil {
			return nil, fmt.Errorf("failed to map node by id %d: %w", n.ID, err)
		}
		nodes = append(nodes, *mapped)
	}

	return &nodes, nil
}

func (r NodeRepository) FindById(ctx context.Context, id int64) (*domain.Node, error) {
	data, err := r.q.FindNodeById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNodeNotFound
		}
		return nil, fmt.Errorf("failed to select node by id %d: %w", id, err)
	}

	node, err := domain.FromDbModel(data)
	if err != nil {
		return nil, fmt.Errorf("failed to map node by id %d: %w", id, err)
	}

	return node, nil
}

func (r NodeRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountNode(ctx)
}
