package db

import (
	"context"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/nodes"
	"github.com/imDrOne/minecraft-server-manager/pkg/db"
	"github.com/imDrOne/minecraft-server-manager/pkg/pagination"
	seeds "github.com/imDrOne/minecraft-server-manager/test/generated/seeds"
	"github.com/imDrOne/minecraft-server-manager/test/lib"
	"github.com/stretchr/testify/suite"
	"testing"
)

type NodeRepositoryTestSuite struct {
	suite.Suite
	db        *db.Postgres
	repo      *nodes.NodeRepository
	seedQuery *seeds.Queries
	ctx       context.Context
}

func (suite *NodeRepositoryTestSuite) SetupSuite() {
	var err error
	suite.ctx = context.Background()

	pgContainer := lib.GetPgContainer()
	suite.db, err = db.NewWithConnectionString(pgContainer.ConnectionString)
	suite.Require().NoError(err)

	suite.repo = nodes.NewNodeRepository(suite.db.Pool)
	suite.seedQuery = seeds.New(suite.db.Pool)
}

func (suite *NodeRepositoryTestSuite) AfterTest(_, _ string) {
	_, _ = suite.db.Pool.Exec(context.Background(), "TRUNCATE node RESTART IDENTITY CASCADE")
}

func (suite *NodeRepositoryTestSuite) BeforeTest(_, _ string) {
	_ = suite.seedQuery.InsertNodeSeeds(context.Background())
}

func (suite *NodeRepositoryTestSuite) TearDownSuite() {
	defer func() {
		if err := suite.db.Close(); err != nil {
			panic("error during down node-repo-test suite")
		}
	}()
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_Save_CreatedNode() {
	expectedHost := "test.com"
	expectedPort := uint(49158)
	actual, err := suite.repo.Save(suite.ctx, func() (*domain.Node, error) {
		return domain.CreateNode(expectedHost, expectedPort)
	})
	suite.Require().NoError(err)
	suite.Require().NotEmpty(actual.Id())
	suite.EqualValues(expectedPort, actual.Port())
	suite.EqualValues(expectedHost, actual.Host())
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_FindById_ErrNotFound() {
	_, err := suite.repo.FindById(suite.ctx, 999)
	suite.Require().ErrorIs(err, domain.ErrNodeNotFound)
}

func (suite *NodeRepositoryTestSuite) TestNodeRepository_FindPaginated() {
	getPageRequest := func(page, size uint64) pagination.PageRequest {
		p, _ := pagination.NewPageRequest(page, size)
		return p
	}

	tests := []struct {
		name    string
		pageReq pagination.PageRequest
		count   int
		total   uint64
		pages   uint64
	}{
		{
			name:    "valid page request",
			pageReq: getPageRequest(1, 30),
			count:   15,
			total:   15,
			pages:   1,
		},
		{
			name:    "exact fetch",
			pageReq: getPageRequest(1, 15),
			count:   15,
			total:   15,
			pages:   1,
		},
		{
			name:    "page = 1; page size = 4;",
			pageReq: getPageRequest(1, 4),
			count:   4,
			total:   15,
			pages:   4,
		},
		{
			name:    "page = 2; page size = 4;",
			pageReq: getPageRequest(2, 4),
			count:   4,
			total:   15,
			pages:   4,
		},
		{
			name:    "page = 4; page size = 4;",
			pageReq: getPageRequest(4, 4),
			count:   3,
			total:   15,
			pages:   4,
		},
		{
			name:    "empty page -> page = 4; page size = 4;",
			pageReq: getPageRequest(5, 4),
			count:   0,
			total:   15,
			pages:   4,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			data, _ := suite.repo.FindPaginated(context.Background(), tt.pageReq)
			suite.Require().Equal(tt.pages, data.Meta.Pages())
			suite.Require().Equal(tt.total, data.Meta.Total())
			suite.Require().Equal(tt.count, len(data.Data))
		})
	}
}

func TestNodeRepositorySuite(t *testing.T) {
	suite.Run(t, new(NodeRepositoryTestSuite))
}
