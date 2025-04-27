package vault

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/imDrOne/minecraft-server-manager/config"
	"path/filepath"
	"strconv"
)

var (
	SaveKeyPairError = errors.New("failed to save key pair")
	GetKeyPairError  = errors.New("failed to get key pair")

	GetPrivateKeyError = errors.New("failed to get private key")
	GetPublicKeyError  = errors.New("failed to get public key")
)

type KeyStore interface {
	Save(context context.Context, create func() KeyPair) error
	Get(context context.Context, id int) (KeyPair, error)
}

type KeyStoreClient struct {
	vault *vault.Client
	cgf   config.Vault
}

func NewKeyStoreClient(vault *vault.Client, cfg config.Vault) KeyStoreClient {
	return KeyStoreClient{
		vault,
		cfg,
	}
}

func (r *KeyStoreClient) Save(ctx context.Context, create func() (int, KeyPair)) error {
	id, keypair := create()
	_, err := r.vault.Secrets.KvV2Write(ctx, r.supplySavePath(id), schema.KvV2WriteRequest{
		Data: map[string]any{
			privateKey: keypair.Private,
			publicKey:  keypair.Public,
		},
	}, vault.WithMountPath(r.cgf.MountPath))
	if err != nil {
		return fmt.Errorf("%w: %w", SaveKeyPairError, err)
	}

	return nil
}

func (r *KeyStoreClient) Get(ctx context.Context, id int) (KeyPair, error) {
	result, err := r.vault.Secrets.KvV2Read(ctx, r.supplySavePath(id), vault.WithMountPath(r.cgf.MountPath))
	if err != nil {
		return KeyPair{}, fmt.Errorf("%w: %w", GetKeyPairError, err)
	}

	privateSsh, ok := result.Data.Data[privateKey].(string)
	if !ok {
		return KeyPair{}, GetPrivateKeyError
	}

	publicSsh, ok := result.Data.Data[publicKey].(string)
	if !ok {
		return KeyPair{}, GetPublicKeyError
	}

	return KeyPair{
		privateSsh,
		publicSsh,
	}, nil
}

func (r *KeyStoreClient) supplySavePath(id int) string {
	return filepath.Join(r.cgf.Connections.Path, strconv.Itoa(id))
}
