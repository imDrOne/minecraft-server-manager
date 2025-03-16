package nodes

import (
	"database/sql"
	"errors"
	"fmt"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/query"
	testutils "github.com/imDrOne/minecraft-server-manager/internal/pkg/test/repository"
	"github.com/imDrOne/minecraft-server-manager/pkg/pagination"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
)

type NodeRepositoryTestSuite struct {
	testutils.RepoTestSuite[*MockNodeQueries, *NodeRepository]
}

func (suite *NodeRepositoryTestSuite) SetupTest() {
	suite.RepoTestSuite.SetupTest(
		func(ctrl *gomock.Controller) *MockNodeQueries {
			return NewMockNodeQueries(ctrl)
		},
		func(mockQueries *MockNodeQueries) *NodeRepository {
			return &NodeRepository{q: mockQueries}
		},
	)
}

var (
	createNode = func() (*domain.Node, error) { return &domain.Node{}, nil }
	updateNode = func(*domain.Node) (*domain.Node, error) { return &domain.Node{}, nil }
)

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Save_ErrorOnCheckExists() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		CheckExistsNode(suite.Ctx, gomock.Any()).
		Return(false, testutils.ErrInternalSql)

	_, err := suite.RepoSupplier().Save(suite.Ctx, createNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to check node exist")
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Save_AlreadyExists() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		CheckExistsNode(suite.Ctx, gomock.Any()).
		Return(true, nil)

	_, err := suite.RepoSupplier().Save(suite.Ctx, createNode)
	require.Error(suite.T(), err)
	require.EqualError(suite.T(), err, domain.ErrNodeAlreadyExist.Error())
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Save_ErrorOnSaveNode() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		CheckExistsNode(suite.Ctx, gomock.Any()).
		Return(false, nil)

	mockQueries.EXPECT().
		SaveNode(suite.Ctx, gomock.Any()).
		Return(query.SaveNodeRow{
			ID:        0,
			CreatedAt: pgtype.Timestamp{},
		}, testutils.ErrInternalSql)

	_, err := suite.RepoSupplier().Save(suite.Ctx, createNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to insert node")
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Save_ErrorOnCreateNode() {
	nodeCreateErr := errors.New("error node create")

	createNode := func() (*domain.Node, error) { return nil, nodeCreateErr }
	_, err := suite.RepoSupplier().Save(suite.Ctx, createNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), nodeCreateErr.Error())
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Update_ErrorOnFindById() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindNodeById(suite.Ctx, gomock.Any()).
		Return(query.Node{}, testutils.ErrInternalSql)

	err := suite.RepoSupplier().Update(suite.Ctx, 10, updateNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to get node")
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Update_ErrorOnUpdate() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindNodeById(suite.Ctx, gomock.Any()).
		Return(query.Node{Host: "test.t", Port: 64676}, nil)

	mockQueries.EXPECT().
		UpdateNodeById(suite.Ctx, gomock.Any()).
		Return(testutils.ErrInternalSql)

	err := suite.RepoSupplier().Update(suite.Ctx, 10, updateNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to update node query")
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Update_ErrorOnUpdateCallback() {
	nodeUpdateErr := errors.New("error node update")
	mockQueries := suite.MockQuerySupplier()

	updateNode := func(*domain.Node) (*domain.Node, error) { return nil, nodeUpdateErr }

	mockQueries.EXPECT().
		FindNodeById(suite.Ctx, gomock.Any()).
		Return(query.Node{Host: "test.t", Port: 64676}, nil)

	err := suite.RepoSupplier().Update(suite.Ctx, 10, updateNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "error node update")
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Find_ErrorOnFind() {
	nodesFindErr := errors.New("DB error")
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindNodes(suite.Ctx, gomock.Any()).
		Return(nil, nodesFindErr)

	_, err := suite.RepoSupplier().Find(suite.Ctx, pagination.PageRequest{})
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to select nodes")
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Find_ErrorOnMapping() {
	mockQueries := suite.MockQuerySupplier()

	failedId := int64(1740170404993)
	nodes := []query.Node{
		{
			ID:   1740170404992,
			Host: "test.t",
			Port: 49153,
		},
		{
			ID:   failedId,
			Host: "test.t",
			Port: 0,
		},
	}

	mockQueries.EXPECT().
		FindNodes(suite.Ctx, gomock.Any()).
		Return(nodes, nil)

	_, err := suite.RepoSupplier().Find(suite.Ctx, pagination.PageRequest{})
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), fmt.Sprintf("failed to map node by id %d", failedId))
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_FindById_ErrorOnFind() {
	failedId := int64(1740170936307)
	errOnFindById := errors.New("DB error")
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindNodeById(suite.Ctx, gomock.Any()).
		Return(query.Node{}, errOnFindById)

	_, err := suite.RepoSupplier().FindById(suite.Ctx, failedId)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), fmt.Sprintf("failed to select node by id %d", failedId))
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_FindById_NotFound() {
	mockQueries := suite.MockQuerySupplier()

	mockQueries.EXPECT().
		FindNodeById(suite.Ctx, gomock.Any()).
		Return(query.Node{}, sql.ErrNoRows)

	_, err := suite.RepoSupplier().FindById(suite.Ctx, 1740171226)
	require.Error(suite.T(), err)
	require.EqualError(suite.T(), err, domain.ErrNodeNotFound.Error())
}

func TestRunNodeRepositorySuite(t *testing.T) {
	suite.Run(t, new(NodeRepositoryTestSuite))
}
