package connections

import (
	"github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/pkg/ssh"
)

type ConnectionDto struct {
	RawConnection connections.Connection
	SshKeyPair    ssh.KeyPair
}
