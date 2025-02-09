package node

import (
	"context"
	rawrepo "github.com/imDrOne/minecraft-server-manager/internal/generated/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryImpl struct {
	q *rawrepo.Queries
}

func NewRepositoryImpl(pool *pgxpool.Pool) *RepositoryImpl {
	return &RepositoryImpl{q: rawrepo.New(pool)}
}

func (r RepositoryImpl) Save(ctx context.Context, arg rawrepo.SaveNodeParams) (rawrepo.Node, error) {
	return r.q.SaveNode(ctx, arg)
}

func (r RepositoryImpl) Update(ctx context.Context, arg rawrepo.UpdateNodeByIdParams) error {
	return r.q.UpdateNodeById(ctx, arg)
}

func (r RepositoryImpl) Find(ctx context.Context, arg rawrepo.FindNodesParams) ([]rawrepo.Node, error) {
	return r.q.FindNodes(ctx, arg)
}

func (r RepositoryImpl) FindById(ctx context.Context, id int64) (rawrepo.Node, error) {
	return r.q.FindNodeById(ctx, id)
}

func (r RepositoryImpl) Count(ctx context.Context) (int64, error) {
	return r.q.CountNode(ctx)
}
