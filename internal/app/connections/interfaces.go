package connections

import (
	"context"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	conndb "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/db"
	connservice "github.com/imDrOne/minecraft-server-manager/internal/service/connections"
)

//go:generate go tool mockgen -destination mock_repo_test.go -package connections . Repository
type Repository interface {
	Save(context.Context, int64, conndb.CreateConn) (*domain.Connection, error)
	FindByNodeId(context.Context, int64) ([]domain.Connection, error)
	FindById(context.Context, int64) (*domain.Connection, error)
	Update(ctx context.Context, id int64, updateConn conndb.UpdateConn) error
}

//go:generate go tool mockgen -destination mock_sshkey_repo_test.go -package connections . ConnectionSshKeyRepository
type ConnectionSshKeyRepository interface {
	Save(context context.Context, connId int64, create func() (*domain.ConnectionSshKeyPair, error)) (*domain.ConnectionSshKeyPair, error)
	Get(context context.Context, id int64) (*domain.ConnectionSshKeyPair, error)
}

//go:generate go tool mockgen -destination mock_service_test.go -package connections . Service
type Service interface {
	Create(ctx context.Context, nodeId int64, createConn conndb.CreateConn) (*connservice.ConnectionTO, error)
	Update(ctx context.Context, id int64, updateConn conndb.UpdateConn) error
	FindByNodeId(ctx context.Context, id int64) ([]domain.Connection, error)
}
