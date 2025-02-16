package infrastructure

import (
	"context"
	"github.com/imDrOne/minecraft-server-manager/integration_tests/lib"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/nodes"
	"github.com/imDrOne/minecraft-server-manager/pkg/db"
	"github.com/stretchr/testify/suite"
	"testing"
)

type NodeRepositoryTestSuite struct {
	suite.Suite
	DB   *db.Postgres
	Repo *nodes.NodeRepository
	ctx  context.Context
}

func (suite *NodeRepositoryTestSuite) SetupSuite() {
	var err error
	suite.ctx = context.Background()

	suite.DB, err = db.NewWithConnectionString(lib.GetPgConnectionString())
	suite.Require().NoError(err)

	suite.Repo = nodes.NewNodeRepository(suite.DB.Pool)
}

func (suite *NodeRepositoryTestSuite) TearDownSuite() {
	suite.DB.Close()
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

func (suite *NodeRepositoryTestSuite) TestFindById_Error() {
	_, err := suite.Repo.FindById(context.Background(), 999)
	suite.Require().ErrorIs(err, domain.ErrNodeNotFound)
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(NodeRepositoryTestSuite))
}
