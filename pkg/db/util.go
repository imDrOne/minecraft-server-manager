package db

import (
	"errors"
	"fmt"
	"net"
)

type ConnectionData struct {
	Host     string
	Database string
	User     string
	Password string
	Port     string
	SSL      bool
}

func NewConnectionData(host, dbname, user, password, port string, ssl bool) (ConnectionData, error) {
	if host == "" {
		return ConnectionData{}, errors.New("no host found")
	}
	if port == "" {
		port = "5432"
	}
	return ConnectionData{
		Database: dbname,
		User:     user,
		Password: password,
		Port:     port,
		SSL:      ssl,
		Host:     host,
	}, nil
}

func (c ConnectionData) String() string {
	sslMode := "disable"
	if c.SSL {
		sslMode = "require"
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		c.User,
		c.Password,
		net.JoinHostPort(c.Host, c.Port),
		c.Database,
		sslMode,
	)
}
