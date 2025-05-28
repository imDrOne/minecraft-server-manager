package connections

import (
	"context"
	"fmt"
	domainconn "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	domainnode "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	conndb "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/db"
	sshservice "github.com/imDrOne/minecraft-server-manager/internal/service/ssh"
)

type ConnectionTO struct {
	*domainconn.Connection
	PublicKey string
}

//go:generate go tool mockgen -destination mock_node_repo_test.go -package connections . NodeRepository
type NodeRepository interface {
	FindById(context.Context, int64) (*domainnode.Node, error)
}

//go:generate go tool mockgen -destination mock_conn_repo_test.go -package connections . ConnectionRepository
type ConnectionRepository interface {
	Save(context.Context, int64, conndb.CreateConn) (*domainconn.Connection, error)
	FindById(context.Context, int64) (*domainconn.Connection, error)
	Update(ctx context.Context, id int64, updateConn conndb.UpdateConn) error
	FindByNodeId(ctx context.Context, id int64) ([]domainconn.Connection, error)
}

//go:generate go tool mockgen -destination mock_sshkey_repo_test.go -package connections . ConnectionSshKeyRepository
type ConnectionSshKeyRepository interface {
	Save(context context.Context, connId int64, create func() (*domainconn.ConnectionSshKeyPair, error)) (*domainconn.ConnectionSshKeyPair, error)
}

type Dependencies struct {
	NodeRepo         NodeRepository
	ConnRepo         ConnectionRepository
	ConnSshKeyRepo   ConnectionSshKeyRepository
	SshKeygenService *sshservice.KeygenService
}

type ConnectionService struct {
	nodeRepo         NodeRepository
	connRepo         ConnectionRepository
	connSshKeyRepo   ConnectionSshKeyRepository
	sshKeygenService *sshservice.KeygenService
}

func NewConnectionService(deps Dependencies) *ConnectionService {
	return &ConnectionService{
		nodeRepo:         deps.NodeRepo,
		connRepo:         deps.ConnRepo,
		sshKeygenService: deps.SshKeygenService,
		connSshKeyRepo:   deps.ConnSshKeyRepo,
	}
}

// Create todo: MSM-24
func (r *ConnectionService) Create(
	ctx context.Context,
	nodeId int64,
	createConn conndb.CreateConn,
) (*ConnectionTO, error) {
	conn, err := r.connRepo.Save(ctx, nodeId, createConn)
	if err != nil {
		return nil, fmt.Errorf("error on saving connection for node=%d: %w", nodeId, err)
	}

	pair, err := r.connSshKeyRepo.Save(ctx, conn.NodeId(), func() (*domainconn.ConnectionSshKeyPair, error) {
		pair, err := r.sshKeygenService.GeneratePair()
		return domainconn.ConnSshKeysFromPair(pair), err
	})
	if err != nil {
		return nil, fmt.Errorf("error on saving connection ssh-key pair: %w", err)
	}

	return &ConnectionTO{
		Connection: conn,
		PublicKey:  pair.Public(),
	}, nil
}

func (r *ConnectionService) Update(ctx context.Context, id int64, updateConn conndb.UpdateConn) error {
	return r.connRepo.Update(ctx, id, updateConn)
}

func (r *ConnectionService) FindByNodeId(ctx context.Context, id int64) ([]domainconn.Connection, error) {
	return r.connRepo.FindByNodeId(ctx, id)
}
