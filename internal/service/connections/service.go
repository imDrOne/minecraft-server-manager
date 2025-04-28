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

//go:generate go tool mockgen -destination mock_test.go -package connections . NodeRepository
type NodeRepository interface {
	FindById(context.Context, int64) (*domainnode.Node, error)
}

//go:generate go tool mockgen -destination mock_test.go -package connections . ConnectionRepository
type ConnectionRepository interface {
	Save(context.Context, int64, conndb.CreateConn) (*domainconn.Connection, error)
	FindById(context.Context, int64) (*domainconn.Connection, error)
	Update(ctx context.Context, id int64, updateConn conndb.UpdateConn) error
}

type ConnectionService struct {
	nodeRepo         NodeRepository
	connRepo         ConnectionRepository
	sshKeygenService sshservice.KeygenService
}

func NewConnectionService(n NodeRepository, c ConnectionRepository) *ConnectionService {
	return &ConnectionService{
		nodeRepo: n,
		connRepo: c,
	}
}

func (r *ConnectionService) Create(
	ctx context.Context,
	nodeId int64,
	createConn conndb.CreateConn,
) (*ConnectionDto, error) {
	conn, err := r.connRepo.Save(ctx, nodeId, createConn)
	if err != nil {
		return nil, err
	}

	pair, err := r.sshKeygenService.GeneratePair()
	if err != nil {
		return nil, err
	}

	return &ConnectionDto{
		RawConnection: *conn,
		SshKeyPair:    pair,
	}, err
}

func (r *ConnectionService) Update(ctx context.Context, id int64, updateConn conndb.UpdateConn) error {
	return r.connRepo.Update(ctx, id, updateConn)
}

func (r *ConnectionService) FindById(ctx context.Context, id int64) (*domainconn.Connection, error) {
	return r.connRepo.FindById(ctx, id)
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
