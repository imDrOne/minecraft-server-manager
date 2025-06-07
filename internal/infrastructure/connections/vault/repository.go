package vault

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	"path/filepath"
	"strconv"
)

const (
	privateKey = "private"
	publicKey  = "public"
)

var (
	SaveKeyPairError = errors.New("failed to save key pair")
	GetKeyPairError  = errors.New("failed to get key pair")

	GetPrivateKeyError = errors.New("failed to get private key")
	GetPublicKeyError  = errors.New("failed to get public key")
)

type ConnectionSshKeyRepository struct {
	vault *vault.Client
	cgf   config.Vault
}

func NewConnSshKeyRepository(vault *vault.Client, cfg config.Vault) *ConnectionSshKeyRepository {
	return &ConnectionSshKeyRepository{
		vault,
		cfg,
	}
}

func (r *ConnectionSshKeyRepository) Save(ctx context.Context, connId int64, create func() (*connections.ConnectionSshKeyPair, error)) (*connections.ConnectionSshKeyPair, error) {
	keypair, err := create()
	if err != nil {
		return nil, fmt.Errorf("error on creating keypair: %w", err)
	}

	_, err = r.vault.Secrets.KvV2Write(ctx, r.supplySavePath(connId), schema.KvV2WriteRequest{
		Data: map[string]any{
			privateKey: keypair.PrivatePemStr(),
			publicKey:  keypair.Public(),
		},
	}, vault.WithMountPath(r.cgf.MountPath))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", SaveKeyPairError, err)
	}

	return keypair, nil
}

func (r *ConnectionSshKeyRepository) Get(ctx context.Context, id int64) (*connections.ConnectionSshKeyPair, error) {
	result, err := r.vault.Secrets.KvV2Read(ctx, r.supplySavePath(id), vault.WithMountPath(r.cgf.MountPath))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", GetKeyPairError, err)
	}

	privateSshPem, ok := result.Data.Data[privateKey].(string)
	if !ok {
		return nil, GetPrivateKeyError
	}

	publicSsh, ok := result.Data.Data[publicKey].(string)
	if !ok {
		return nil, GetPublicKeyError
	}

	return connections.NewConnSshKeyPair([]byte(privateSshPem), publicSsh), nil
}

func (r *ConnectionSshKeyRepository) supplySavePath(id int64) string {
	return filepath.Join(r.cgf.Connections.Path, strconv.FormatInt(id, 10))
}
