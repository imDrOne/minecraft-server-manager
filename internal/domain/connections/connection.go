package connections

import (
	"errors"
	"fmt"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/query"
	"regexp"
	"time"
)

const RootUser = "root"

var (
	ErrConnectionNotFound      = errors.New("connection not found")
	ErrValidationConnection    = errors.New("invalid connection")
	ErrConnectionAlreadyExists = errors.New("invalid connection")
)

type Connection struct {
	id        int64
	nodeId    int64
	user      string
	createdAt time.Time
}

func NewConnection(id, nodeId int64, user string, createdAt time.Time) (*Connection, error) {
	if err := validateUser(user); err != nil {
		return nil, err
	}

	return &Connection{id: id, nodeId: nodeId, user: user, createdAt: createdAt}, nil
}

func CreateConnection(nodeId int64, user string) (*Connection, error) {
	return NewConnection(0, nodeId, user, time.Time{})
}

func CreateRootConnection(nodeId int64) (*Connection, error) {
	return NewConnection(0, nodeId, RootUser, time.Time{})
}

func FromDbModel(c query.Connection) (*Connection, error) {
	return NewConnection(c.ID, c.NodeID, c.User, c.CreatedAt.Time)
}

func (c *Connection) Id() int64 {
	return c.id
}
func (c *Connection) NodeId() int64 {
	return c.nodeId
}
func (c *Connection) User() string {
	return c.user
}
func (c *Connection) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Connection) WithDBGeneratedValues(row query.SaveConnectionRow) *Connection {
	return &Connection{
		id:        row.ID,
		nodeId:    c.nodeId,
		user:      c.user,
		createdAt: row.CreatedAt.Time,
	}
}

func (c *Connection) Update(user string) (*Connection, error) {
	if err := validateUser(user); err != nil {
		return nil, err
	}

	c.user = user
	return c, nil
}

func validateUser(value string) error {
	re := regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)
	if !re.MatchString(value) {
		return fmt.Errorf("invalid user: %w", ErrValidationConnection)
	}
	return nil
}
