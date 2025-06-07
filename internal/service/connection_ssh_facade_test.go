package service

import (
	"github.com/imDrOne/minecraft-server-manager/internal/app/remotes"
	"github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/domain/nodes"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/vault"
	testutils "github.com/imDrOne/minecraft-server-manager/internal/pkg/test"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

type Mocks struct {
	*MockNodeRepository
	*MockConnectionRepository
	*MockConnectionSshKeyRepository
	*MockSshService
}

type ConnectionSshFacadeTestSuite struct {
	testutils.Suite[Mocks, *ConnectionSshFacade]
}

func (suite *ConnectionSshFacadeTestSuite) SetupTest() {
	suite.Suite.SetupTest(func(ctrl *gomock.Controller) Mocks {
		return Mocks{
			MockNodeRepository:             NewMockNodeRepository(ctrl),
			MockConnectionRepository:       NewMockConnectionRepository(ctrl),
			MockConnectionSshKeyRepository: NewMockConnectionSshKeyRepository(ctrl),
			MockSshService:                 NewMockSshService(ctrl),
		}
	}, func(mocks Mocks) *ConnectionSshFacade {
		return NewConnectionSshFacade(Dependencies{
			NodeRepo:       mocks.MockNodeRepository,
			ConnRepo:       mocks.MockConnectionRepository,
			ConnSshKeyRepo: mocks.MockConnectionSshKeyRepository,
			SshService:     mocks.MockSshService,
		})
	})
}

func (suite *ConnectionSshFacadeTestSuite) TestConnectionSshFacade_InjectPublicKey_ErrorOnFetchingConnection() {
	mocks := suite.MockSupplier()
	mocks.MockConnectionRepository.EXPECT().
		FindById(suite.Ctx, gomock.Any()).
		Return(nil, connections.ErrConnectionNotFound)

	service := suite.TargetSupplier()
	err := service.InjectPublicKey(suite.Ctx, 101, remotes.ForwardPublicKeyDto{})
	suite.ErrorIs(err, connections.ErrConnectionNotFound)
	suite.ErrorContains(err, "error on supplying node-ssh-connection obj by conn-id=101")
}

func (suite *ConnectionSshFacadeTestSuite) TestConnectionSshFacade_InjectPublicKey_ErrorOnFetchingNode() {
	conn, _ := connections.NewConnection(1, 101, "test", time.Now())
	mocks := suite.MockSupplier()
	mocks.MockConnectionRepository.EXPECT().
		FindById(suite.Ctx, gomock.Any()).
		Return(conn, nil)

	mocks.MockNodeRepository.EXPECT().
		FindById(suite.Ctx, int64(101)).
		Return(nil, nodes.ErrNodeNotFound)

	service := suite.TargetSupplier()
	err := service.InjectPublicKey(suite.Ctx, 101, remotes.ForwardPublicKeyDto{})
	suite.ErrorIs(err, nodes.ErrNodeNotFound)
	suite.ErrorContains(err, "error on supplying node-ssh-connection obj by conn-id=101")
}

func (suite *ConnectionSshFacadeTestSuite) TestConnectionSshFacade_InjectPublicKey_ErrorOnFetchingSshKeys() {
	conn, _ := connections.NewConnection(1, 101, "test", time.Now())
	node, _ := nodes.NewNode(101, "test-host", 1080, time.Now())

	mocks := suite.MockSupplier()
	mocks.MockConnectionRepository.EXPECT().
		FindById(suite.Ctx, gomock.Any()).
		Return(conn, nil)

	mocks.MockNodeRepository.EXPECT().
		FindById(suite.Ctx, gomock.Any()).
		Return(node, nil)

	mocks.MockConnectionSshKeyRepository.EXPECT().
		Get(suite.Ctx, gomock.Any()).
		Return(nil, vault.GetPrivateKeyError)

	service := suite.TargetSupplier()
	err := service.InjectPublicKey(suite.Ctx, 101, remotes.ForwardPublicKeyDto{})
	suite.ErrorIs(err, vault.GetPrivateKeyError)
	suite.ErrorContains(err, "error on fetching keys by conn-id=101")
}

func (suite *ConnectionSshFacadeTestSuite) TestConnectionSshFacade_InjectPublicKey_ErrorOnInjectingKeys() {
	conn, _ := connections.NewConnection(1, 101, "test", time.Now())
	node, _ := nodes.NewNode(101, "test-host", 1080, time.Now())

	mocks := suite.MockSupplier()
	mocks.MockConnectionRepository.EXPECT().
		FindById(suite.Ctx, gomock.Any()).
		Return(conn, nil)

	mocks.MockNodeRepository.EXPECT().
		FindById(suite.Ctx, gomock.Any()).
		Return(node, nil)

	mocks.MockConnectionSshKeyRepository.EXPECT().
		Get(suite.Ctx, gomock.Any()).
		Return(nil, vault.GetPrivateKeyError)

	service := suite.TargetSupplier()
	err := service.InjectPublicKey(suite.Ctx, 101, remotes.ForwardPublicKeyDto{})
	suite.ErrorIs(err, vault.GetPrivateKeyError)
	suite.ErrorContains(err, "error on fetching keys by conn-id=101")
}

func Test(t *testing.T) {
	suite.Run(t, new(ConnectionSshFacadeTestSuite))
}
