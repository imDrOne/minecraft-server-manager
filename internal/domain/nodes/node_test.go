package nodes

import (
	"fmt"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"testing"
	"time"
)

func TestNewNode(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		host    string
		port    uint
		wantErr bool
	}{
		{"Valid node", 1, "localhost", 50000, false},
		{"Invalid host (empty)", 1, "", 50000, true},
		{"Invalid port (too low)", 1, "localhost", 40000, true},
		{"Invalid port (too high)", 1, "localhost", 70000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := NewNode(tt.id, tt.host, tt.port, time.Time{})
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (node.id != tt.id || node.host != tt.host || node.port != tt.port) {
				t.Errorf("NewNode() = %+v, expected id=%d, host=%s, port=%d", node, tt.id, tt.host, tt.port)
			}
		})
	}
}

func TestCreateNode(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    uint
		wantErr bool
	}{
		{"Valid node", "localhost", 50000, false},
		{"Invalid host (empty)", "", 50000, true},
		{"Invalid port (too low)", "localhost", 40000, true},
		{"Invalid port (too high)", "localhost", 70000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := CreateNode(tt.host, tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (node.id != 0 || node.host != tt.host || node.port != tt.port) {
				t.Errorf("CreateNode() = %+v, expected id=0, host=%s, port=%d", node, tt.host, tt.port)
			}
		})
	}
}

func TestWithId(t *testing.T) {
	node, _ := NewNode(0, "localhost", 50000, time.Time{})
	newID := int64(42)
	newCreatedAt := time.Now()
	newNode := node.WithDBGeneratedValues(repository.SaveNodeRow{
		ID: newID,
		CreatedAt: pgtype.Timestamp{
			Time:  newCreatedAt,
			Valid: true,
		},
	})

	if newNode.id != newID {
		t.Errorf("WithDBGeneratedValues() = %+v, expected id=%d", newNode, newID)
	}

	if !newNode.createdAt.Equal(newCreatedAt) {
		t.Errorf("WithDBGeneratedValues() = %+v, expected id=%d", newNode, newID)
	}
}

func TestFromDbModel(t *testing.T) {
	dbNode := repository.Node{ID: 10, Host: "dbhost", Port: 50001}
	node, err := FromDbModel(dbNode)

	if err != nil {
		t.Fatalf("FromDbModel() error = %v, expected no error", err)
	}
	if node.id != dbNode.ID || node.host != dbNode.Host || node.port != uint(dbNode.Port) {
		t.Errorf("FromDbModel() = %+v, expected id=%d, host=%s, port=%d", node, dbNode.ID, dbNode.Host, dbNode.Port)
	}
}

func TestGetters(t *testing.T) {
	createdAt := time.Now()
	node := Node{id: 5, host: "testhost", port: 50002, createdAt: createdAt}

	if node.Id() != 5 {
		t.Errorf("Id() = %d, expected 5", node.Id())
	}
	if node.Host() != "testhost" {
		t.Errorf("Host() = %s, expected testhost", node.Host())
	}
	if node.Port() != 50002 {
		t.Errorf("Port() = %d, expected 50002", node.Port())
	}
	if !node.CreatedAt().Equal(createdAt) {
		t.Errorf("CreatedAt() = %v, expected %v", node.CreatedAt(), createdAt)
	}
}

func TestValidateHost(t *testing.T) {
	tests := []struct {
		host    string
		wantErr bool
	}{
		{"localhost", false},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Host=%s", tt.host), func(t *testing.T) {
			err := validateHost(tt.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHost(%q) error = %v, wantErr %v", tt.host, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		port    uint
		wantErr bool
	}{
		{50000, false},
		{49151, true},
		{65536, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Port=%d", tt.port), func(t *testing.T) {
			err := validatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePort(%d) error = %v, wantErr %v", tt.port, err, tt.wantErr)
			}
		})
	}
}
