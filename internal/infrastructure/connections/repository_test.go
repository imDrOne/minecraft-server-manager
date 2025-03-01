package connections

import (
	"errors"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/query"
	testutils "github.com/imDrOne/minecraft-server-manager/internal/pkg/test/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
)

type ConnectionRepositoryTestSuite struct {
	testutils.RepoTestSuite[*MockConnectionQueries, *ConnectionRepository]
}

func (suite *ConnectionRepositoryTestSuite) SetupTest() {
	suite.RepoTestSuite.SetupTest(
		func(ctrl *gomock.Controller) *MockConnectionQueries {
			return NewMockConnectionQueries(ctrl)
		},
		func(mockQueries *MockConnectionQueries) *ConnectionRepository {
			return &ConnectionRepository{q: mockQueries}
		},
	)
}

var (
	createConn = func() (*domain.Connection, error) { return &domain.Connection{}, nil }
	updateConn = func(*domain.Connection) (*domain.Connection, error) { return &domain.Connection{}, nil }
)

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Save_ErrorOnCheckExists() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		CheckExistsConnection(suite.Ctx, gomock.Any()).
		Return(false, testutils.ErrInternalSql)

	_, err := suite.RepoSupplier().Save(suite.Ctx, createConn)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to check connection exist")
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Save_AlreadyExists() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		CheckExistsConnection(suite.Ctx, gomock.Any()).
		Return(true, nil)

	_, err := suite.RepoSupplier().Save(suite.Ctx, createConn)
	require.Error(suite.T(), err)
	require.EqualError(suite.T(), err, domain.ErrConnectionAlreadyExists.Error())
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Save_ErrorOnSaveNode() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		CheckExistsConnection(suite.Ctx, gomock.Any()).
		Return(false, nil)

	mockQueries.EXPECT().
		SaveConnection(suite.Ctx, gomock.Any()).
		Return(query.SaveConnectionRow{
			ID:        0,
			CreatedAt: pgtype.Timestamp{},
		}, testutils.ErrInternalSql)

	_, err := suite.RepoSupplier().Save(suite.Ctx, createConn)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to insert connection")
}

func (suite *ConnectionRepositoryTestSuite) TestNodeRepository_Save_ErrorOnCreateNode() {
	nodeCreateErr := errors.New("error node create")

	createNode := func() (*domain.Connection, error) { return nil, nodeCreateErr }
	_, err := suite.RepoSupplier().Save(suite.Ctx, createNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), nodeCreateErr.Error())
}

func Test(t *testing.T) {
	suite.Run(t, new(ConnectionRepositoryTestSuite))
}
