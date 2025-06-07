package remotes

import (
	"context"
)

//go:generate go tool mockgen -destination mock_ssh_connectionservice_test.go -package remotes . SshConnectionService
type SshConnectionService interface {
	InjectPublicKey(ctx context.Context, id int64, dto ForwardPublicKeyDto) error
}
