package connections

import (
	"github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
)

type ConnectionDto struct {
	*connections.Connection
	*connections.ConnectionSshKeyPair
}
