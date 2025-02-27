package nodes

import (
	"errors"
	"fmt"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/query"
	"github.com/imDrOne/minecraft-server-manager/pkg/pagination"
	"time"
)

type PagePaginatedNodes = pagination.PaginatedResult[Node]

var (
	ErrNodeNotFound     = errors.New("node not found")
	ErrValidationNode   = errors.New("invalid node")
	ErrNodeAlreadyExist = errors.New("already exists")
)

type Node struct {
	id        int64
	host      string
	port      uint
	createdAt time.Time
}

func NewNode(id int64, host string, port uint, createdAt time.Time) (*Node, error) {
	if err := validateHost(host); err != nil {
		return nil, err
	}
	if err := validatePort(port); err != nil {
		return nil, err
	}

	return &Node{id: id, host: host, port: port, createdAt: createdAt}, nil
}
func CreateNode(host string, port uint) (*Node, error) {
	return NewNode(0, host, port, time.Time{})
}
func FromDbModel(n query.Node) (*Node, error) {
	return NewNode(n.ID, n.Host, uint(n.Port), n.CreatedAt.Time)
}

func (n *Node) WithDBGeneratedValues(row query.SaveNodeRow) *Node {
	return &Node{
		id:        row.ID,
		host:      n.host,
		port:      n.port,
		createdAt: row.CreatedAt.Time,
	}
}
func (n *Node) Id() int64 {
	return n.id
}
func (n *Node) Host() string {
	return n.host
}
func (n *Node) Port() uint {
	return n.port
}
func (n *Node) CreatedAt() time.Time {
	return n.createdAt
}

func validateHost(host string) error {
	if host == "" {
		return fmt.Errorf("%w: host is required", ErrValidationNode)
	}
	return nil
}
func validatePort(port uint) error {
	if port <= 49152 || port >= 65535 {
		return fmt.Errorf("%w: out of range 49152 - 65535", ErrValidationNode)
	}
	return nil
}
