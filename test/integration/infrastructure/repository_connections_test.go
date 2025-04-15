package infrastructure

import (
	"context"
	_ "embed"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections"
	"github.com/imDrOne/minecraft-server-manager/pkg/db"
	seeds "github.com/imDrOne/minecraft-server-manager/test/generated/seeds"
	"github.com/imDrOne/minecraft-server-manager/test/lib"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConnectionRepositoryTestSuite struct {
	suite.Suite
	db        *db.Postgres
	repo      *connections.ConnectionRepository
	seedQuery *seeds.Queries
	ctx       context.Context
}

func (suite *ConnectionRepositoryTestSuite) SetupSuite() {
	var err error
	suite.ctx = context.Background()

	suite.db, err = db.NewWithConnectionString(lib.GetPgConnectionString())
	suite.Require().NoError(err)

	suite.repo = connections.NewConnectionRepository(suite.db.Pool)
	suite.seedQuery = seeds.New(suite.db.Pool)
}

func (suite *ConnectionRepositoryTestSuite) BeforeTest(_, testName string) {
	_ = suite.seedQuery.InsertNodeSeed(context.Background())

	switch testName {
	case "TestConnectionRepository_Update_UpdatedConnection":
		_ = suite.seedQuery.InsertConnectionSeed(context.Background())
	case "TestConnectionRepository_FindById_NodesConnections":
		_ = suite.seedQuery.InsertConnectionsSeed(context.Background())
	}
}

func (suite *ConnectionRepositoryTestSuite) AfterTest(_, _ string) {
	_, _ = suite.db.Pool.Exec(context.Background(), "TRUNCATE connection RESTART IDENTITY CASCADE")
	_, _ = suite.db.Pool.Exec(context.Background(), "TRUNCATE node RESTART IDENTITY CASCADE")
}

func (suite *ConnectionRepositoryTestSuite) TearDownSuite() {
	defer func() {
		if err := suite.db.Close(); err != nil {
			panic("error during down node-repo-test suite")
		}
	}()
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Save_CreatedConnection() {
	expectedNodeId := int64(100)
	expectedConn, _ := domain.CreateConnection(1, "user")
	actual, _ := suite.repo.Save(suite.ctx, expectedNodeId, func() (*domain.Connection, error) {
		return expectedConn, nil
	})
	suite.Require().NotEmpty(actual.Id())
	suite.EqualValues(expectedConn.User(), actual.User())
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Update_UpdatedConnection() {
	connId := int64(100)
	expectedConn, _ := domain.CreateRootConnection(1)
	err := suite.repo.Update(suite.ctx, connId, func(c domain.Connection) (*domain.Connection, error) {
		return expectedConn, nil
	})
	suite.Require().NoError(err)

	actualConn, err := suite.repo.FindById(suite.ctx, connId)
	suite.Require().NoError(err)

	suite.EqualValues(expectedConn.User(), actualConn.User())
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_FindById_NodesConnections() {
	actualConns, err := suite.repo.FindByNodeId(suite.ctx, 100)
	suite.Require().NoError(err)
	suite.Len(actualConns, 3)

	actualConns, err = suite.repo.FindByNodeId(suite.ctx, 1000)
	suite.Require().NoError(err)
	suite.Len(actualConns, 0)
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Save_ErrNotFound() {
	_, err := suite.repo.FindById(suite.ctx, 999)
	suite.Require().ErrorIs(err, domain.ErrConnectionNotFound)
}

func TestConnectionRepositorySuite(t *testing.T) {
	suite.Run(t, new(ConnectionRepositoryTestSuite))
}
