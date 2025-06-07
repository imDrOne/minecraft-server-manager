package service

import (
	"context"
	"github.com/hashicorp/vault-client-go"
	"github.com/imDrOne/minecraft-server-manager/config"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	conndb "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/db"
	connvt "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/vault"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/nodes"
	"github.com/imDrOne/minecraft-server-manager/internal/service/connections"
	sshservice "github.com/imDrOne/minecraft-server-manager/internal/service/ssh"
	"github.com/imDrOne/minecraft-server-manager/pkg/db"
	vaultpkg "github.com/imDrOne/minecraft-server-manager/pkg/vault"
	seeds "github.com/imDrOne/minecraft-server-manager/test/generated/seeds"
	"github.com/imDrOne/minecraft-server-manager/test/lib"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConnectionServiceTestSuite struct {
	suite.Suite
	cfg           *config.Config
	postgres      *db.Postgres
	vaultConnRepo *connvt.ConnectionSshKeyRepository
	service       *connections.ConnectionService
	seedQuery     *seeds.Queries
	ctx           context.Context
}

func (suite *ConnectionServiceTestSuite) SetupSuite() {
	var err error

	suite.cfg = config.NewWithEnvironment("test")
	suite.ctx = context.Background()
	pgContainer := lib.GetPgContainer()
	suite.postgres, err = db.NewWithConnectionString(pgContainer.ConnectionString)
	suite.Require().NoError(err)

	connDbRepo := conndb.NewConnectionRepository(suite.postgres.Pool)
	nodeDbRepo := nodes.NewNodeRepository(suite.postgres.Pool)

	suite.seedQuery = seeds.New(suite.postgres.Pool)

	vaultContainer := lib.GetVaultContainer()
	vaultClient, err := vaultpkg.New(vaultContainer.HostAddress, func(cl *vault.Client) error {
		return cl.SetToken("root")
	})
	suite.Require().NoError(err)

	suite.vaultConnRepo = connvt.NewConnSshKeyRepository(vaultClient, suite.cfg.Vault)

	sshKeygenService := sshservice.NewKeygenService(suite.cfg.SSHKeygen)

	suite.service = connections.NewConnectionService(connections.Dependencies{
		NodeRepo:         nodeDbRepo,
		ConnRepo:         connDbRepo,
		ConnSshKeyRepo:   suite.vaultConnRepo,
		SshKeygenService: sshKeygenService,
	})
}

func (suite *ConnectionServiceTestSuite) BeforeTest() {
	_ = suite.seedQuery.InsertNodeSeed(context.Background())
}

func (suite *ConnectionServiceTestSuite) TestConnectionService_SaveConnection() {
	expectedNodeId := int64(100)
	expectedUserName := "test"
	actual, err := suite.service.Create(suite.ctx, expectedNodeId, func() (*domain.Connection, error) {
		return domain.CreateConnection(expectedNodeId, expectedUserName)
	})
	suite.NoError(err)
	suite.Equal(expectedNodeId, actual.NodeId())
	suite.Equal(expectedUserName, actual.User())
	suite.NotEmpty(actual.PublicKey)

	keys, err := suite.vaultConnRepo.Get(suite.ctx, actual.Id())
	suite.NoError(err)
	suite.Equal(keys.Public(), actual.PublicKey)
}

func TestServiceConnectionSuite(t *testing.T) {
	suite.Run(t, new(ConnectionServiceTestSuite))
}
