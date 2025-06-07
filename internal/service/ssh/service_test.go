package ssh

import (
	_ "embed"
	"github.com/imDrOne/minecraft-server-manager/internal/service/ssh/model"
	sshpkg "github.com/imDrOne/minecraft-server-manager/pkg/ssh"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

//go:embed test_ssh_key.pub
var testSshKeyPub string

type SshServiceTestSuite struct {
	suite.Suite
	service          *Service
	sshConnectionCfg model.NodeSSHConnectionTO
}

func (s *SshServiceTestSuite) SetupTest() {
	s.sshConnectionCfg = model.NodeSSHConnectionTO{
		Auth: sshpkg.Auth{
			Type:     sshpkg.AuthPassword,
			Password: "test",
		},
		Host: "localhost",
		Port: 99,
		User: "test",
	}

	s.service = &Service{
		sshClientTimeout: time.Second * 10,
	}
}

func (s *SshServiceTestSuite) TestService_InjectPublicKey_ClientProvidingError() {
	err := s.service.InjectPublicKey(s.sshConnectionCfg, testSshKeyPub)
	s.NotNil(err)
	s.Contains(err.Error(), "ssh connection failed")
}

func (s *SshServiceTestSuite) TestService_Ping_ClientProvidingError() {
	err := s.service.Ping(s.sshConnectionCfg)
	s.NotNil(err)
	s.Contains(err.Error(), "ssh connection failed")
}

func Test(t *testing.T) {
	suite.Run(t, new(SshServiceTestSuite))
}
