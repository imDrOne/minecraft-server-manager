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
	cfg            *config.Config
	keyStoreClient *connssh.ConnectionSshKeyRepository
}

// todo: Decompose setup (look at service_connection_test.go)
func (suite *ConnectionsKeyStoreTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	suite.cfg = config.NewWithEnvironment("test")
	vaultContainer := lib.GetVaultContainer()
	vaultClient, err := vault.New(vaultContainer.HostAddress, func(cl *client.Client) error {
		return cl.SetToken("root")
	})
	suite.Require().NoError(err)
	suite.keyStoreClient = connssh.NewConnSshKeyRepository(vaultClient, suite.cfg.Vault)

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
