package connections

import (
	"context"
	"fmt"
	domainconn "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	domainnode "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	conndb "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/db"
	sshservice "github.com/imDrOne/minecraft-server-manager/internal/service/ssh"
	"golang.org/x/crypto/ssh"
	"log"
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
	Get(context context.Context, id int64) (*domainconn.ConnectionSshKeyPair, error)
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

func (r *ConnectionService) ConnectByConnectionID(ctx context.Context, connId int64, command string) (string, error) {
	conn, err := r.connRepo.FindById(ctx, connId)
	if err != nil {
		return "", fmt.Errorf("failed to get connection: %w", err)
	}

	node, err := r.nodeRepo.FindById(ctx, conn.NodeId())
	if err != nil {
		return "", fmt.Errorf("failed to get node: %w", err)
	}

	return r.connectAndRun(*node, *conn, command)
}

func (r *ConnectionService) connectAndRun(node domainnode.Node, conn domainconn.Connection, command string) (string, error) {
	config := &ssh.ClientConfig{
		User: conn.User(),
		Auth: []ssh.AuthMethod{
			r.resolveAuth(conn),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", node.Host(), node.Port())
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %w", err)
	}

	return string(output), nil
}

// todo: Key must be private - now that it is public
func (r *ConnectionService) resolveAuth(conn domainconn.Connection) ssh.AuthMethod {
	//key, err := os.ReadFile(conn.Key())
	//if err != nil {
	//	log.Fatalf("unable to read private key: %v", err)
	//}
	signer, err := ssh.ParsePrivateKey([]byte{})
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	return ssh.PublicKeys(signer)
}
