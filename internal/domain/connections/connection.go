package connections

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/query"
	"golang.org/x/crypto/ssh"
	"io"
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
	key       string
	user      string
	createdAt time.Time
}

func NewConnection(id int64, key string, user string, createdAt time.Time) (*Connection, error) {
	if err := validateUser(user); err != nil {
		return nil, err
	}
	if err := validateKey(key); err != nil {
		return nil, err
	}

	return &Connection{id: id, key: key, user: user, createdAt: createdAt}, nil
}

func CreateConnection(key string, user string) (*Connection, error) {
	return NewConnection(0, key, user, time.Time{})
}

func CreateRootConnection(key string) (*Connection, error) {
	return NewConnection(0, key, RootUser, time.Time{})
}

func FromDbModel(c query.Connection) (*Connection, error) {
	return NewConnection(c.ID, c.Key, c.User, c.CreatedAt.Time)
}

func (c *Connection) Id() int64 {
	return c.id
}
func (c *Connection) Key() string {
	return c.key
}
func (c *Connection) User() string {
	return c.user
}
func (c *Connection) CreatedAt() time.Time {
	return c.createdAt
}
func (c *Connection) ChecksumStr() string {
	return hex.EncodeToString(c.Checksum())
}
func (c *Connection) Checksum() []byte {
	h := md5.New()
	_, _ = io.WriteString(h, c.key)
	_, _ = io.WriteString(h, c.user)
	return h.Sum(nil)
}

func (c *Connection) WithDBGeneratedValues(row query.SaveConnectionRow) *Connection {
	return &Connection{
		id:        row.ID,
		key:       c.key,
		user:      c.user,
		createdAt: row.CreatedAt.Time,
	}
}

func (c *Connection) Update(key, user string) (*Connection, error) {
	if err := validateUser(user); err != nil {
		return nil, err
	}
	if err := validateKey(key); err != nil {
		return nil, err
	}

	c.key = key
	c.user = user
	return c, nil
}

func validateKey(value string) error {
	if value == "" {
		return fmt.Errorf("%w: key is required", ErrValidationConnection)
	}

	if _, _, _, _, err := ssh.ParseAuthorizedKey([]byte(value)); err != nil {
		return fmt.Errorf("%w: invalid key", ErrValidationConnection)
	}
	return nil
}

func validateUser(value string) error {
	re := regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)
	if !re.MatchString(value) {
		return fmt.Errorf("%w: invalid user", ErrValidationConnection)
	}
	return nil
}
