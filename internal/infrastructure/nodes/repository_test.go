package nodes

import (
	"context"
	"errors"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
)

type ControllerFinishFinish func()

type NodeRepositoryTestSuite struct {
	suite.Suite
	ctx                     context.Context
	mockNodeQueriesProvider func(*testing.T) (*MockNodeQueries, ControllerFinishFinish)
}

var errInternalSql = errors.New("DB error")

func (suite *NodeRepositoryTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockNodeQueriesProvider = func(t *testing.T) (*MockNodeQueries, ControllerFinishFinish) {
		ctrl := gomock.NewController(t)
		mockQueries := NewMockNodeQueries(ctrl)
		return mockQueries, ctrl.Finish
	}
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Save_ErrorOnCheckExists() {
	mockQueries, finish := suite.mockNodeQueriesProvider(suite.T())
	defer finish()

	repo := NodeRepository{q: mockQueries}
	mockNode := &domain.Node{}
	createNode := func() (*domain.Node, error) { return mockNode, nil }

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
	mockNode := &domain.Node{}
	createNode := func() (*domain.Node, error) { return mockNode, nil }

	mockQueries.EXPECT().
		CheckExistsNode(suite.ctx, gomock.Any()).
		Return(false, nil)

	mockQueries.EXPECT().
		SaveNode(suite.ctx, gomock.Any()).
		Return(int64(0), errors.New("DB error"))

	_, err := repo.Save(suite.ctx, createNode)
	require.Error(suite.T(), err)
	require.Contains(suite.T(), errInternalSql, "failed to insert node")
}

func TestRunNodeRepositorySuite(t *testing.T) {
	suite.Run(t, new(NodeRepositoryTestSuite))
}
