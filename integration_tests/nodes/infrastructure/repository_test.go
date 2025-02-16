package infrastructure

import (
	"context"
	"github.com/imDrOne/minecraft-server-manager/integration_tests/lib"
	"github.com/imDrOne/minecraft-server-manager/internal"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/nodes"
	"github.com/imDrOne/minecraft-server-manager/pkg/db"
	"github.com/stretchr/testify/suite"
	"testing"
)

type NodeRepositoryTestSuite struct {
	suite.Suite
	DB          *db.Postgres
	Repo        *nodes.NodeRepository
	ctx         context.Context
	pgContainer *lib.PgContainer
}

func (suite *NodeRepositoryTestSuite) SetupSuite() {
	var err error
	suite.ctx = context.Background()
	suite.pgContainer, err = lib.StartPostgresContainer(suite.ctx)
	suite.Require().NoError(err)

	suite.DB, err = db.NewWithConnectionString(suite.pgContainer.ConnectionString)
	suite.Require().NoError(err)

	suite.Repo = nodes.NewNodeRepository(suite.DB.Pool)
	err = internal.MigrateUpWithConnectionString(suite.pgContainer.ConnectionString)
	suite.Require().NoError(err)
}

func (suite *NodeRepositoryTestSuite) TearDownSuite() {
	suite.DB.Close()
	lib.StopPostgresContainer(suite.ctx)
}

func (suite *NodeRepositoryTestSuite) TestCreateNode() {
	expectedHost := "test.com"
	expectedPort := uint(49158)
	actual, err := suite.Repo.Save(context.Background(), func() (*domain.Node, error) {
		return domain.CreateNode(expectedHost, expectedPort)
	})
	suite.Require().NoError(err)
	suite.EqualValues(expectedPort, actual.Port())
	suite.EqualValues(expectedHost, actual.Host())

	_, err = suite.Repo.Save(context.Background(), func() (*domain.Node, error) {
		return domain.CreateNode(expectedHost, expectedPort)
	})
	suite.Require().ErrorIs(err, domain.ErrNodeAlreadyExist)
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(NodeRepositoryTestSuite))
}
