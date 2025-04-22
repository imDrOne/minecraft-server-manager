package connections

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

type KeyPair struct {
	Private string
	Public  string
}

type Keystore interface {
	Save(context context.Context, create func() (int, KeyPair)) error
	Get(context context.Context, id int) (KeyPair, error)
}

type KeyStoreClient struct {
	vault *vault.Client
	config.ConnectionsVault
}

func NewKeyStoreClient(vault *vault.Client, cfg config.ConnectionsVault) KeyStoreClient {
	return KeyStoreClient{
		vault,
		cfg,
	}
}

func (cl *KeyStoreClient) Save(ctx context.Context, create func() (int, KeyPair)) error {

	id, keypair := create()
	_, err := cl.vault.Secrets.KvV2Write(ctx, cl.setPath(id), schema.KvV2WriteRequest{
		Data: map[string]any{
			privateKey: keypair.Private,
			publicKey:  keypair.Public,
		}}, vault.WithMountPath(cl.MountPath))

	if err != nil {
		return fmt.Errorf("%w: %s", SaveKeyPairError, err)
	}
	return nil
}

func (cl *KeyStoreClient) Get(ctx context.Context, id int) (KeyPair, error) {

	result, err := cl.vault.Secrets.KvV2Read(ctx, cl.setPath(id), vault.WithMountPath(cl.MountPath))
	if err != nil {
		return KeyPair{}, fmt.Errorf("%w: %s", GetKeyPairError, err)
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

func (cl *KeyStoreClient) setPath(id int) string {
	return filepath.Join(cl.Path, strconv.Itoa(id))
}
