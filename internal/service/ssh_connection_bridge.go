package service

import (
	"context"
	"fmt"
	domainconn "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	domainnode "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/nodes"
	"golang.org/x/crypto/ssh"
	"log"
)

type SSHConnectionBridge struct {
	nodeRepo       nodes.NodeRepository
	connectionRepo connections.ConnectionRepository
}

func NewSSHBridge(n nodes.NodeRepository, c connections.ConnectionRepository) *SSHConnectionBridge {
	return &SSHConnectionBridge{
		nodeRepo:       n,
		connectionRepo: c,
	}
}

func (s *SSHConnectionBridge) ConnectByConnectionID(ctx context.Context, connId int64, command string) (string, error) {
	conn, err := s.connectionRepo.FindById(ctx, connId)
	if err != nil {
		return "", fmt.Errorf("failed to get connection: %w", err)
	}

	node, err := s.nodeRepo.FindById(ctx, conn.NodeId())
	if err != nil {
		return "", fmt.Errorf("failed to get node: %w", err)
	}

	return s.connectAndRun(*node, *conn, command)
}

func (s *SSHConnectionBridge) connectAndRun(node domainnode.Node, conn domainconn.Connection, command string) (string, error) {
	config := &ssh.ClientConfig{
		User: conn.User(),
		Auth: []ssh.AuthMethod{
			s.resolveAuth(conn),
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
func (s *SSHConnectionBridge) resolveAuth(conn domainconn.Connection) ssh.AuthMethod {
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
