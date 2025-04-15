package connections

import (
	_ "embed"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/query"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewConnection(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		user      string
		createdAt time.Time
		wantErr   bool
	}{
		{"Valid connection", 1, "validuser", time.Now(), false},
		{"Invalid user (special chars)", 1, "Invalid-User!", time.Now(), true},
		{"Invalid user (empty)", 1, "", time.Now(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := NewConnection(tt.id, 1, tt.user, tt.createdAt)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewConnection() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				assert.NotNil(t, conn, "Connection should not be nil")
				assert.Equal(t, tt.id, conn.Id(), "ID should match")
				assert.Equal(t, tt.user, conn.User(), "User should match")
				assert.Equal(t, tt.createdAt, conn.CreatedAt(), "CreatedAt should match")
			}
		})
	}
}

func TestWithDBGeneratedValues(t *testing.T) {
	expectedCreatedAt := time.Now()
	expectedId := int64(123)

	dbRow := query.SaveConnectionRow{
		ID: 123,
		CreatedAt: pgtype.Timestamp{
			Time:             expectedCreatedAt,
			InfinityModifier: 0,
			Valid:            true,
		},
	}
	conn, err := NewConnection(0, 1, "validuser", time.Time{})
	assert.NoError(t, err, "Expected no error when creating a valid connection")

	updatedConn := conn.WithDBGeneratedValues(dbRow)

	t.Run("Check updated ID and CreatedAt", func(t *testing.T) {
		assert.Equal(t, expectedId, updatedConn.Id(), "ID should be updated correctly")
		assert.Equal(t, expectedCreatedAt, updatedConn.CreatedAt(), "CreatedAt should be updated correctly")
	})
}

func TestCreateConnection(t *testing.T) {
	tests := []struct {
		name    string
		user    string
		wantErr bool
	}{
		{
			name:    "Valid user",
			user:    "validuser",
			wantErr: false,
		},
		{
			name:    "Invalid user",
			user:    "Invalid!User",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := CreateConnection(1, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}
			if err == nil && conn == nil {
				t.Errorf("Expected a valid connection, got nil")
			}
		})
	}
}

func TestCreateRootConnection(t *testing.T) {
	conn, err := CreateRootConnection(1)
	if conn == nil || err != nil {
		t.Errorf("Not expected error on creating root connection")
		t.FailNow()
	}
	if conn.User() != RootUser {
		t.Errorf("Expected user: %s, got %s", RootUser, conn.User())
	}
}
