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

//go:embed test_key.pub
var validSSHKey string

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

func (suite *ConnectionRepositoryTestSuite) BeforeTest(_, _ string) {
	_ = suite.seedQuery.InsertNodeSeeds(context.Background())
}

func (suite *ConnectionRepositoryTestSuite) AfterTest(_, _ string) {
	_, _ = suite.db.Pool.Exec(context.Background(), "TRUNCATE connection RESTART IDENTITY CASCADE")
}

func (suite *ConnectionRepositoryTestSuite) TearDownSuite() {
	defer func() {
		if err := suite.db.Close(); err != nil {
			panic("error during down node-repo-test suite")
		}
	}()
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Save_CreatedConnection() {
	expectedUser := "user"
	expectedNodeId := int64(100)
	actual, err := suite.repo.Save(suite.ctx, expectedNodeId, func() (*domain.Connection, error) {
		return domain.CreateConnection(validSSHKey, expectedUser)
	})
	suite.Require().NoError(err)
	suite.Require().NotEmpty(actual.Id())
	suite.EqualValues(expectedUser, actual.User())
	suite.EqualValues(validSSHKey, actual.Key())
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Update_UpdatedConnection() {
	nodeId := int64(100)
	conn, err := suite.repo.Save(suite.ctx, nodeId, func() (*domain.Connection, error) {
		return domain.CreateConnection(validSSHKey, "superuser")
	})
	suite.Require().NoError(err)
	suite.EqualValues(conn.User(), "superuser")

	expectedConn, _ := domain.CreateRootConnection(validSSHKey)
	err = suite.repo.Update(suite.ctx, conn.Id(), func(c domain.Connection) (*domain.Connection, error) {
		return expectedConn, nil
	})
	suite.Require().NoError(err)

	conn, err = suite.repo.FindById(suite.ctx, conn.Id())
	suite.Require().NoError(err)

	conns, err := suite.repo.FindByNodeId(suite.ctx, nodeId)
	suite.Require().NoError(err)
	suite.Len(conns, 1)

	suite.EqualValues(expectedConn.User(), conn.User())
	suite.EqualValues(expectedConn.Key(), conn.Key())
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Save_ErrNotFound() {
	_, err := suite.repo.FindById(suite.ctx, 999)
	suite.Require().ErrorIs(err, domain.ErrConnectionNotFound)
}

func TestConnectionRepositorySuite(t *testing.T) {
	suite.Run(t, new(ConnectionRepositoryTestSuite))
}
