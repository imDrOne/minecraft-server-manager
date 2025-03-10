package connections

import (
	"database/sql"
	_ "embed"
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
	updateConn = func(domain.Connection) (*domain.Connection, error) { return &domain.Connection{}, nil }
)

var (
	//go:embed test_key.pub
	validSSHKey string
	connUser    = pgtype.Text{
		String: "test",
		Valid:  true,
	}
	validConn = query.Connection{
		ID:        1,
		NodeID:    1,
		Key:       validSSHKey,
		Checksum:  "checksum",
		User:      connUser,
		CreatedAt: pgtype.Timestamp{},
		UpdatedAt: pgtype.Timestamp{},
	}
	invalidConn = query.Connection{
		ID:        2,
		NodeID:    1,
		Key:       "invalid-key",
		Checksum:  "checksum",
		User:      connUser,
		CreatedAt: pgtype.Timestamp{},
		UpdatedAt: pgtype.Timestamp{},
	}
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

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Save_ErrorOnCreateNode() {
	nodeCreateErr := errors.New("error node create")

	createNode := func() (*domain.Connection, error) { return nil, nodeCreateErr }
	_, err := suite.RepoSupplier().Save(suite.Ctx, createNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), nodeCreateErr.Error())
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_FindByNodeId_ErrOnSelect() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindConnectionsByNodeId(suite.Ctx, gomock.Any()).
		Return(nil, testutils.ErrInternalSql)

	_, err := suite.RepoSupplier().FindByNodeId(suite.Ctx, 11)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to select connections by node-id 11")
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_FindByNodeId_ErrNotFound() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindConnectionsByNodeId(suite.Ctx, gomock.Any()).
		Return(nil, sql.ErrNoRows)

	_, err := suite.RepoSupplier().FindByNodeId(suite.Ctx, 11)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), domain.ErrConnectionNotFound.Error())
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_FindByNodeId_ErrOnMapping() {
	mockQueries := suite.MockQuerySupplier()
	mockQueries.EXPECT().
		FindConnectionsByNodeId(suite.Ctx, gomock.Any()).
		Return([]query.Connection{validConn, invalidConn}, nil)

	_, err := suite.RepoSupplier().FindByNodeId(suite.Ctx, 11)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to map connection by id 2")
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_FindById_ErrOnSelect() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindConnectionById(suite.Ctx, gomock.Any()).
		Return(query.Connection{}, testutils.ErrInternalSql)

	_, err := suite.RepoSupplier().FindById(suite.Ctx, 11)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to select connection by id 11")
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_FindById_ErrNotFound() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindConnectionById(suite.Ctx, gomock.Any()).
		Return(query.Connection{}, sql.ErrNoRows)

	_, err := suite.RepoSupplier().FindById(suite.Ctx, 11)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), domain.ErrConnectionNotFound.Error())
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_FindById_ErrOnMapping() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindConnectionById(suite.Ctx, gomock.Any()).
		Return(invalidConn, nil)

	_, err := suite.RepoSupplier().FindById(suite.Ctx, 11)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to map conn by id")
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Update_ErrOnFind() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindConnectionById(suite.Ctx, gomock.Any()).
		Return(query.Connection{}, sql.ErrNoRows)

	err := suite.RepoSupplier().
		Update(suite.Ctx, 11, updateConn)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to get connection")
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Update_ErrOnUpdateFn() {
	connUpdateErr := errors.New("error connection update")
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindConnectionById(suite.Ctx, gomock.Any()).
		Return(validConn, nil)

	updateConn := func(conn domain.Connection) (*domain.Connection, error) {
		return nil, connUpdateErr
	}

	err := suite.RepoSupplier().
		Update(suite.Ctx, 11, updateConn)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to update connection by id=11")
}

func (suite *ConnectionRepositoryTestSuite) TestConnectionRepository_Update_ErrOnUpdateQuery() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindConnectionById(suite.Ctx, gomock.Any()).
		Return(validConn, nil)

	mockQueries.EXPECT().
		UpdateConnectionById(suite.Ctx, gomock.Any()).
		Return(testutils.ErrInternalSql)

	err := suite.RepoSupplier().Update(suite.Ctx, 11, updateConn)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to update connection by id=11")
}

func Test(t *testing.T) {
	suite.Run(t, new(ConnectionRepositoryTestSuite))
}
