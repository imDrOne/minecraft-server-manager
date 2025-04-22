package vaultkeystores

import (
	"context"
	client "github.com/hashicorp/vault-client-go"
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections"
	"github.com/imDrOne/minecraft-server-manager/pkg/vault"
	"github.com/imDrOne/minecraft-server-manager/test/lib"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConnectionsKeyStoreTestSuite struct {
	suite.Suite
	ctx            context.Context
	keyStoreClient connections.KeyStoreClient
}

func (suite *ConnectionsKeyStoreTestSuite) SetupSuite() {

	vaultClient, err := vault.New(lib.VaultHostAddress, func(cl *client.Client) error {
		return cl.SetToken("root")
	})
	if err != nil {
		panic(err)
	}

	suite.ctx = context.Background()

	suite.keyStoreClient = connections.NewKeyStoreClient(
		vaultClient,
		config.ConnectionsVault{
			Path:      "ssh-keys/",
			MountPath: "/secret",
		},
	)

}

func (suite *ConnectionsKeyStoreTestSuite) TestConnectionsKeyStore_SaveAndGetPairs() {
	pair := connections.KeyPair{
		Private: "1",
		Public:  "2",
	}
	err := suite.keyStoreClient.Save(suite.ctx, func() (int, connections.KeyPair) {
		return 1, pair
	})

	require.NoError(suite.T(), err)

	newPair, err := suite.keyStoreClient.Get(suite.ctx, 1)

	require.NoError(suite.T(), err)
	require.Equal(suite.T(), pair, newPair)
}

func TestRepositoryConnectionsKeystoreSuite(t *testing.T) {
	suite.Run(t, new(ConnectionsKeyStoreTestSuite))
}
