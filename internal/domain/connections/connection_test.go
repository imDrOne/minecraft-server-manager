package connections

import (
	_ "embed"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/query"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//go:embed test_key.pub
var validSSHKey string

func TestNewConnection(t *testing.T) {
	tests := []struct {
		name      string
		id        int64
		key       string
		user      string
		createdAt time.Time
		wantErr   bool
	}{
		{"Valid connection", 1, validSSHKey, "validuser", time.Now(), false},
		{"Invalid user (special chars)", 1, validSSHKey, "Invalid-User!", time.Now(), true},
		{"Invalid user (empty)", 1, validSSHKey, "", time.Now(), true},
		{"Invalid key (empty)", 1, "", "validuser", time.Now(), true},
		{"Invalid key (random string)", 1, "invalid-key", "validuser", time.Now(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := NewConnection(tt.id, tt.key, tt.user, tt.createdAt)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewConnection() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				assert.NotNil(t, conn, "Connection should not be nil")
				assert.Equal(t, tt.id, conn.Id(), "ID should match")
				assert.Equal(t, tt.key, conn.Key(), "Key should match")
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
	conn, err := NewConnection(0, validSSHKey, "validuser", time.Time{})
	assert.NoError(t, err, "Expected no error when creating a valid connection")

	updatedConn := conn.WithDBGeneratedValues(dbRow)

	t.Run("Check updated ID and CreatedAt", func(t *testing.T) {
		assert.Equal(t, expectedId, updatedConn.Id(), "ID should be updated correctly")
		assert.Equal(t, expectedCreatedAt, updatedConn.CreatedAt(), "CreatedAt should be updated correctly")
	})
}
