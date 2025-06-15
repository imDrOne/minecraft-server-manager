package service

import (
	"context"
	"fmt"
	"github.com/imDrOne/minecraft-server-manager/internal/app/remotes"
	domainconn "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	domainnode "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/imDrOne/minecraft-server-manager/internal/service/ssh/model"
	sshpkg "github.com/imDrOne/minecraft-server-manager/pkg/ssh"
)

//go:generate go tool mockgen -destination mock_node_repo_test.go -package service . NodeRepository
type NodeRepository interface {
	FindById(context.Context, int64) (*domainnode.Node, error)
}

//go:generate go tool mockgen -destination mock_conn_repo_test.go -package service . ConnectionRepository
type ConnectionRepository interface {
	FindById(context.Context, int64) (*domainconn.Connection, error)
}

//go:generate go tool mockgen -destination mock_sshkey_repo_test.go -package service . ConnectionSshKeyRepository
type ConnectionSshKeyRepository interface {
	Get(context context.Context, id int64) (*domainconn.ConnectionSshKeyPair, error)
}

//go:generate go tool mockgen -destination mock_ssh_service_test.go -package service . SshService
type SshService interface {
	InjectPublicKey(cfg model.NodeSSHConnectionTO, publicKey string) error
	Ping(cfg model.NodeSSHConnectionTO) error
}

type Dependencies struct {
	NodeRepo       NodeRepository
	ConnRepo       ConnectionRepository
	ConnSshKeyRepo ConnectionSshKeyRepository
	SshService     SshService
}

type ConnectionSshFacade struct {
	nodeRepo       NodeRepository
	connRepo       ConnectionRepository
	connSshKeyRepo ConnectionSshKeyRepository
	sshService     SshService
}

func NewConnectionSshFacade(deps Dependencies) *ConnectionSshFacade {
	return &ConnectionSshFacade{
		nodeRepo:       deps.NodeRepo,
		connRepo:       deps.ConnRepo,
		connSshKeyRepo: deps.ConnSshKeyRepo,
		sshService:     deps.SshService,
	}
}

func (r *ConnectionSshFacade) InjectPublicKey(ctx context.Context, id int64, dto remotes.ForwardPublicKeyDto) error {
	nodeSshConnectionTO, err := r.supplyNodeSshConnection(ctx, id)
	if err != nil {
		return fmt.Errorf("error on supplying node-ssh-connection obj by conn-id=%d: %w", id, err)
	}

	keys, err := r.connSshKeyRepo.Get(ctx, nodeSshConnectionTO.NodeId)
	if err != nil {
		return fmt.Errorf("error on fetching keys by conn-id=%d: %w", id, err)
	}

	return r.sshService.InjectPublicKey(nodeSshConnectionTO.WithAuth(sshpkg.Auth{
		Type:     sshpkg.AuthPassword,
		Password: dto.Password,
	}), keys.Public())
}

func (r *ConnectionSshFacade) Ping(ctx context.Context, id int64) error {
	nodeSshConnectionTO, err := r.supplyNodeSshConnection(ctx, id)
	if err != nil {
		return fmt.Errorf("error on supplying node-ssh-connection obj by conn-id=%d: %w", id, err)
	}

	keys, err := r.connSshKeyRepo.Get(ctx, nodeSshConnectionTO.NodeId)
	if err != nil {
		return fmt.Errorf("error on fetching keys by conn-id=%d: %w", id, err)
	}

	return r.sshService.Ping(nodeSshConnectionTO.WithAuth(sshpkg.Auth{
		Type:       sshpkg.AuthPrivateKey,
		PrivateKey: keys.PrivatePem(),
	}))
}

func (r *ConnectionSshFacade) supplyNodeSshConnection(ctx context.Context, id int64) (val model.NodeSSHConnectionTO, err error) {
	connection, err := r.connRepo.FindById(ctx, id)
	if err != nil {
		return val, fmt.Errorf("error on fetching connection by id %d: %w", id, err)
	}

	node, err := r.nodeRepo.FindById(ctx, connection.NodeId())
	if err != nil {
		return val, fmt.Errorf("error on fetching node by node-id %d: %w", connection.NodeId(), err)
	}

	return model.NodeSSHConnectionTO{
		NodeId: connection.NodeId(),
		Host:   node.Host(),
		Port:   node.Port(),
		User:   connection.User(),
	}, err
}
