package node

import (
	"context"
	"errors"
	rawrepo "github.com/imDrOne/minecraft-server-manager/internal/generated/repository"
	"github.com/imDrOne/minecraft-server-manager/pkg/pagination"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
)

var (
	ErrNodeAlreadyExists = errors.New("node with this host and port already exists")
)

type UseCaseImpl struct {
	repository Repository
}

type Repository interface {
	Save(ctx context.Context, arg rawrepo.SaveNodeParams) (rawrepo.Node, error)
	//Update(ctx context.Context, arg rawrepo.UpdateNodeByIdParams) error
	Find(ctx context.Context, arg rawrepo.FindNodesParams) ([]rawrepo.Node, error)
	Count(ctx context.Context) (int64, error)
	//FindById(ctx context.Context, id int64) (rawrepo.Node, error)
}

type PagedNodes struct {
	Nodes []rawrepo.Node
	Meta  pagination.PageMetadata
}

func NewUseCase(r Repository) *UseCaseImpl {
	return &UseCaseImpl{repository: r}
}

func (s UseCaseImpl) Save(ctx context.Context, args rawrepo.SaveNodeParams) (rawrepo.Node, error) {
	node, err := s.repository.Save(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return rawrepo.Node{}, ErrNodeAlreadyExists
		}
		return rawrepo.Node{}, err
	}
	return node, nil
}

func (s UseCaseImpl) Find(ctx context.Context, payload pagination.PageRequest) (PagedNodes, error) {
	nodes, err := s.repository.Find(ctx, rawrepo.FindNodesParams{
		Limit:  int32(payload.Size),
		Offset: int32(payload.Offset()),
	})
	if err != nil {
		slog.Error("error during fetching nodes %w", err)
		return PagedNodes{}, err
	}

	total, err := s.repository.Count(ctx)
	if err != nil {
		slog.Error("error during counting %w", err)
		return PagedNodes{}, err
	}

	meta := pagination.NewPageMetadata(payload.Page, payload.Size, uint64(total))

	return PagedNodes{
		Nodes: nodes,
		Meta:  *meta,
	}, nil
}
