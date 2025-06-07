package connections

import (
	"github.com/imDrOne/minecraft-server-manager/config"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/vault"
	testutils "github.com/imDrOne/minecraft-server-manager/internal/pkg/test"
	sshservice "github.com/imDrOne/minecraft-server-manager/internal/service/ssh"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

type Mocks struct {
	*MockConnectionRepository
	*MockConnectionSshKeyRepository
}

type ConnectionServiceTestSuite struct {
	testutils.Suite[Mocks, *ConnectionService]
	cfg *config.Config
}

func (suite *ConnectionServiceTestSuite) SetupTest() {
	suite.cfg = config.NewWithEnvironment("test")
	suite.Suite.SetupTest(func(ctrl *gomock.Controller) Mocks {
		return Mocks{
			MockConnectionRepository:       NewMockConnectionRepository(ctrl),
			MockConnectionSshKeyRepository: NewMockConnectionSshKeyRepository(ctrl),
		}
	}, func(mocks Mocks) *ConnectionService {
		return NewConnectionService(Dependencies{
			ConnRepo:       mocks.MockConnectionRepository,
			ConnSshKeyRepo: mocks.MockConnectionSshKeyRepository,
		})
	})
}

var createConn = func() (*domain.Connection, error) { return &domain.Connection{}, nil }

func (suite *ConnectionServiceTestSuite) TestConnectionService_Create_ErrOnSavingConn() {
	mocks := suite.MockSupplier()

	mocks.MockConnectionRepository.EXPECT().
		Save(suite.Ctx, gomock.Any(), gomock.Any()).
		Return(nil, testutils.ErrInternalSql)

	_, err := suite.TargetSupplier().
		Create(suite.Ctx, 1, createConn)

	suite.Error(err)
	suite.ErrorIs(err, testutils.ErrInternalSql)
}

func (suite *ConnectionServiceTestSuite) TestConnectionService_Create_ErrOnSavingSshKeys() {
	mocks := suite.MockSupplier()
	connection, _ := domain.NewConnection(1, 1, "test", time.Now())

	mocks.MockConnectionRepository.EXPECT().
		Save(suite.Ctx, gomock.Any(), gomock.Any()).
		Return(connection, nil)

	mocks.MockConnectionSshKeyRepository.EXPECT().
		Save(suite.Ctx, gomock.Any(), gomock.Any()).
		Return(nil, testutils.ErrInternalConnection)

	_, err := suite.TargetSupplier().
		Create(suite.Ctx, 1, createConn)

	suite.Error(err)
	suite.ErrorIs(err, testutils.ErrInternalConnection)
}

func (suite *ConnectionServiceTestSuite) TestConnectionService_Crete_ErrOnGeneratingSshPair() {
	sshCfg := suite.cfg.SSHKeygen
	sshCfg.Bits = 0
	mocks := suite.MockSupplier()
	connection, _ := domain.NewConnection(1, 1, "test", time.Now())

	service := suite.TargetSupplier()
	service.sshKeygenService = sshservice.NewKeygenService(sshCfg)

	service.connSshKeyRepo = vault.NewConnSshKeyRepository(nil, config.Vault{})
	mocks.MockConnectionRepository.EXPECT().
		Save(suite.Ctx, gomock.Any(), gomock.Any()).
		Return(connection, nil)

	_, err := service.Create(suite.Ctx, 1, createConn)
	suite.Error(err)
	suite.Contains(err.Error(), "crypto/rsa: 0-bit keys are insecure (see https://go.dev/pkg/crypto/rsa#hdr-Minimum_key_size)")
}

func Test(t *testing.T) {
	suite.Run(t, new(ConnectionServiceTestSuite))
}
