package vault

import (
	"context"
	client "github.com/hashicorp/vault-client-go"
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	connssh "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/vault"
	"github.com/imDrOne/minecraft-server-manager/internal/pkg/ssh"
	"github.com/imDrOne/minecraft-server-manager/pkg/vault"
	"github.com/imDrOne/minecraft-server-manager/test/lib"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConnectionsKeyStoreTestSuite struct {
	suite.Suite
	ctx            context.Context
	keyStoreClient *connssh.ConnectionSshKeyRepository
}

func (suite *ConnectionsKeyStoreTestSuite) SetupSuite() {

	vaultContainer := lib.GetVaultContainer()
	vaultClient, err := vault.New(vaultContainer.HostAddress, func(cl *client.Client) error {
		return cl.SetToken("root")
	})
	if err != nil {
		panic(err)
	}

	suite.ctx = context.Background()

	suite.keyStoreClient = connssh.NewConnSshKeyRepository(
		vaultClient,
		config.Vault{
			MountPath: "secret",
			Connections: config.ConnectionsVault{
				Path: "ssh-keys",
			},
		},
	)

}

func (suite *ConnectionsKeyStoreTestSuite) TestConnectionsKeyStore_SaveAndGetPairs() {
	expected, _ := ssh.GenerateKeyPair(2054, "", "test")
	actual, err := suite.keyStoreClient.Save(suite.ctx, 1, func() (*connections.ConnectionSshKeyPair, error) {
		return connections.ConnSshKeysFromPair(expected), nil
	})
	suite.NoError(err)
	suite.Equal(expected.Public, actual.Public())
	suite.Equal(expected.Private, actual.PrivatePem())

	actual, err = suite.keyStoreClient.Get(suite.ctx, 1)
	suite.NoError(err)
	suite.Equal(expected.Public, actual.Public())
	suite.Equal(expected.Private, actual.PrivatePem())
}

func TestRepositoryConnectionsKeystoreSuite(t *testing.T) {
	suite.Run(t, new(ConnectionsKeyStoreTestSuite))
}
