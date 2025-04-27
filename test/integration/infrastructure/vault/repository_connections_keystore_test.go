package vault

import (
	"context"
	client "github.com/hashicorp/vault-client-go"
	"github.com/imDrOne/minecraft-server-manager/config"
	connssh "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/vault"
	"github.com/imDrOne/minecraft-server-manager/pkg/vault"
	"github.com/imDrOne/minecraft-server-manager/test/lib"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConnectionsKeyStoreTestSuite struct {
	suite.Suite
	ctx            context.Context
	keyStoreClient connssh.KeyStoreClient
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

	suite.keyStoreClient = connssh.NewKeyStoreClient(
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
	expected := connssh.KeyPair{
		Private: "0db52c7b-f398-479d-b52c-7bf398479d84",
		Public:  "bf5d0b43-6b80-479a-9d0b-436b80a79ac8",
	}
	err := suite.keyStoreClient.Save(suite.ctx, 1, func() (connssh.KeyPair, error) {
		return expected, nil
	})
	suite.NoError(err)

	actual, err := suite.keyStoreClient.Get(suite.ctx, 1)

	suite.NoError(err)
	suite.Equal(expected, actual)
}

func TestRepositoryConnectionsKeystoreSuite(t *testing.T) {
	suite.Run(t, new(ConnectionsKeyStoreTestSuite))
}
