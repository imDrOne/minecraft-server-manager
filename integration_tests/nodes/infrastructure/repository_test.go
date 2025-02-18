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
	db   *db.Postgres
	repo *nodes.NodeRepository
	ctx  context.Context
}

func (suite *NodeRepositoryTestSuite) SetupSuite() {
	var err error
	suite.ctx = context.Background()

	suite.db, err = db.NewWithConnectionString(lib.GetPgConnectionString())
	suite.Require().NoError(err)

	suite.repo = nodes.NewNodeRepository(suite.db.Pool)
}

func (suite *NodeRepositoryTestSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *NodeRepositoryTestSuite) TestCreateNode() {
	expectedHost := "test.com"
	expectedPort := uint(49158)
	actual, err := suite.repo.Save(suite.ctx, func() (*domain.Node, error) {
		return domain.CreateNode(expectedHost, expectedPort)
	})
	suite.Require().NoError(err)
	suite.Require().NotEmpty(actual.Id())
	suite.EqualValues(expectedPort, actual.Port())
	suite.EqualValues(expectedHost, actual.Host())

	_, err = suite.repo.Save(suite.ctx, func() (*domain.Node, error) {
		return domain.CreateNode(expectedHost, expectedPort)
	})
	suite.Require().ErrorIs(err, domain.ErrNodeAlreadyExist)

	if _, err = suite.db.Pool.Exec(context.Background(), "DELETE FROM node WHERE id = $1", actual.Id()); err != nil {
		suite.Fail("fail during clearing data")
	}
}

func (suite *NodeRepositoryTestSuite) TestFindById_Error() {
	_, err := suite.repo.FindById(suite.ctx, 999)
	suite.Require().ErrorIs(err, domain.ErrNodeNotFound)
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(NodeRepositoryTestSuite))
}
