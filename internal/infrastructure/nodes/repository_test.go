package nodes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/imDrOne/minecraft-server-manager/internal/generated/repository"
	"github.com/imDrOne/minecraft-server-manager/pkg/pagination"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
)

type ControllerFinish func()
type MockNodeQueriesProvider func(*testing.T) (*MockNodeQueries, ControllerFinish)

type NodeRepositoryTestSuite struct {
	suite.Suite
	ctx                     context.Context
	mockNodeQueriesProvider MockNodeQueriesProvider
}

var errInternalSql = errors.New("DB error")

func (suite *NodeRepositoryTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockNodeQueriesProvider = func(t *testing.T) (*MockNodeQueries, ControllerFinish) {
		ctrl := gomock.NewController(t)
		mockQueries := NewMockNodeQueries(ctrl)
		return mockQueries, ctrl.Finish
	}
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Save_ErrorOnCheckExists() {
	mockQueries, finish := suite.mockNodeQueriesProvider(suite.T())
	defer finish()

	repo := NodeRepository{q: mockQueries}
	createNode := func() (*domain.Node, error) { return &domain.Node{}, nil }

	mockQueries.EXPECT().
		CheckExistsNode(suite.ctx, gomock.Any()).
		Return(false, errInternalSql)

	_, err := repo.Save(suite.ctx, createNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to check node exist")
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Save_ErrorOnSaveNode() {
	mockQueries, finish := suite.mockNodeQueriesProvider(suite.T())
	defer finish()

	repo := NodeRepository{q: mockQueries}
	createNode := func() (*domain.Node, error) { return &domain.Node{}, nil }

	mockQueries.EXPECT().
		CheckExistsNode(suite.ctx, gomock.Any()).
		Return(false, nil)

	mockQueries.EXPECT().
		SaveNode(suite.ctx, gomock.Any()).
		Return(int64(0), errInternalSql)

	_, err := repo.Save(suite.ctx, createNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to insert node")
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Save_ErrorOnCreateNode() {
	nodeCreateErr := errors.New("error node create")
	mockQueries, finish := suite.mockNodeQueriesProvider(suite.T())
	defer finish()

	repo := NodeRepository{q: mockQueries}
	createNode := func() (*domain.Node, error) { return nil, nodeCreateErr }
	_, err := repo.Save(suite.ctx, createNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), nodeCreateErr.Error())
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Update_ErrorOnFindById() {
	mockQueries, finish := suite.mockNodeQueriesProvider(suite.T())
	defer finish()

	repo := NodeRepository{q: mockQueries}
	updateNode := func(*domain.Node) (*domain.Node, error) { return &domain.Node{}, nil }

	mockQueries.EXPECT().
		FindNodeById(suite.ctx, gomock.Any()).
		Return(repository.Node{}, errInternalSql)

	err := repo.Update(suite.ctx, 10, updateNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to get node")
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Update_ErrorOnUpdate() {
	mockQueries, finish := suite.mockNodeQueriesProvider(suite.T())
	defer finish()

	repo := NodeRepository{q: mockQueries}
	updateNode := func(*domain.Node) (*domain.Node, error) { return &domain.Node{}, nil }

	mockQueries.EXPECT().
		FindNodeById(suite.ctx, gomock.Any()).
		Return(repository.Node{Host: "test.t", Port: 64676}, nil)

	mockQueries.EXPECT().
		UpdateNodeById(suite.ctx, gomock.Any()).
		Return(errInternalSql)

	err := repo.Update(suite.ctx, 10, updateNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to update node query")
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Update_ErrorOnUpdateCallback() {
	nodeUpdateErr := errors.New("error node update")
	mockQueries, finish := suite.mockNodeQueriesProvider(suite.T())
	defer finish()

	repo := NodeRepository{q: mockQueries}
	updateNode := func(*domain.Node) (*domain.Node, error) { return nil, nodeUpdateErr }

	mockQueries.EXPECT().
		FindNodeById(suite.ctx, gomock.Any()).
		Return(repository.Node{Host: "test.t", Port: 64676}, nil)

	err := repo.Update(suite.ctx, 10, updateNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "error node update")
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Find_ErrorOnFind() {
	nodesFindErr := errors.New("DB error")
	mockQueries, finish := suite.mockNodeQueriesProvider(suite.T())
	defer finish()

	repo := NodeRepository{q: mockQueries}
	mockQueries.EXPECT().
		FindNodes(suite.ctx, gomock.Any()).
		Return(nil, nodesFindErr)

	_, err := repo.Find(suite.ctx, pagination.PageRequest{})
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), "failed to select nodes")
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Find_ErrorOnMapping() {
	mockQueries, finish := suite.mockNodeQueriesProvider(suite.T())
	defer finish()

	failedId := int64(1740170404993)
	repo := NodeRepository{q: mockQueries}
	nodes := []repository.Node{
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
		FindNodes(suite.ctx, gomock.Any()).
		Return(nodes, nil)

	_, err := repo.Find(suite.ctx, pagination.PageRequest{})
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), fmt.Sprintf("failed to map node by id %d", failedId))
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_FindById_ErrorOnFind() {
	failedId := int64(1740170936307)
	errOnFindById := errors.New("DB error")
	mockQueries, finish := suite.mockNodeQueriesProvider(suite.T())
	defer finish()

	repo := NodeRepository{q: mockQueries}
	mockQueries.EXPECT().
		FindNodeById(suite.ctx, gomock.Any()).
		Return(repository.Node{}, errOnFindById)

	_, err := repo.FindById(suite.ctx, failedId)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), err.Error(), fmt.Sprintf("failed to select node by id %d", failedId))
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_FindById_NotFound() {
	mockQueries, finish := suite.mockNodeQueriesProvider(suite.T())
	defer finish()

	repo := NodeRepository{q: mockQueries}
	mockQueries.EXPECT().
		FindNodeById(suite.ctx, gomock.Any()).
		Return(repository.Node{}, sql.ErrNoRows)

	_, err := repo.FindById(suite.ctx, 1740171226)
	require.Error(suite.T(), err)
	require.EqualError(suite.T(), err, domain.ErrNodeNotFound.Error())
}

func TestRunNodeRepositorySuite(t *testing.T) {
	suite.Run(t, new(NodeRepositoryTestSuite))
}
