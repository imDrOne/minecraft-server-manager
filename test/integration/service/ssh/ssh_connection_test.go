package ssh

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"github.com/imDrOne/minecraft-server-manager/config"
	"github.com/imDrOne/minecraft-server-manager/internal/service/ssh"
	"github.com/imDrOne/minecraft-server-manager/internal/service/ssh/model"
	sshpkg "github.com/imDrOne/minecraft-server-manager/pkg/ssh"
	"github.com/imDrOne/minecraft-server-manager/test/lib"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"time"
)

//go:embed test_ssh_key.pub
var testSshKeyPub string

//go:embed test_ssh_key_x1.pub
var testSshKeyPubX1 string

//go:embed test_ssh_key
var testSshKey string

func parseAuthorizedKeys(data []byte) []string {
	var keys []string
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		keys = append(keys, line)
	}

	return keys
}

func provideSshConnectionConfig() model.NodeSSHConnectionTO {
	sshdContainer := lib.GetSshdContainer()
	return model.NodeSSHConnectionTO{
		Host: sshdContainer.Host,
		Port: int64(sshdContainer.Port),
		User: "root",
		Auth: sshpkg.Auth{
			Type:     sshpkg.AuthPassword,
			Password: "test",
		},
	}
}

type ConnectionServiceTestSuite struct {
	suite.Suite
	cfg     *config.Config
	service *ssh.Service
	ctx     context.Context
}

func (suite *ConnectionServiceTestSuite) SetupSuite() {
	suite.cfg = config.NewWithEnvironment("test")
	suite.ctx = context.Background()
	suite.service = ssh.NewSshService(time.Second * 10)
}

func (suite *ConnectionServiceTestSuite) TestConnectionService_InjectPublicKey_SuccessfullyForwarded() {
	connectionConfig := provideSshConnectionConfig()

	err := suite.service.InjectPublicKey(connectionConfig, testSshKeyPub)
	suite.NoError(err)

	connectionConfig.Auth = sshpkg.Auth{
		Type:       sshpkg.AuthPrivateKey,
		PrivateKey: []byte(testSshKey),
	}
	pingResult := suite.service.Ping(connectionConfig)
	suite.NoError(pingResult)
}

func (suite *ConnectionServiceTestSuite) TestConnectionService_InjectPublicKey_ManyInserts() {
	tests := []struct {
		name             string
		keys             []string
		expectedKeyCount int
	}{
		{
			name:             "Trying to insert 2 identical keys",
			keys:             []string{testSshKeyPub, testSshKeyPub},
			expectedKeyCount: 1,
		},
		{
			name:             "Trying to insert 2 different keys",
			keys:             []string{testSshKeyPub, testSshKeyPubX1},
			expectedKeyCount: 2,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			connectionConfig := provideSshConnectionConfig()

			err := suite.service.InjectPublicKey(connectionConfig, tt.keys[0])
			suite.NoError(err)

			err = suite.service.InjectPublicKey(connectionConfig, tt.keys[1])
			suite.NoError(err)

			client, err := sshpkg.ProvideSshClient(sshpkg.ClientConfig{
				Auth:    connectionConfig.Auth,
				Host:    connectionConfig.Host,
				Port:    connectionConfig.Port,
				User:    connectionConfig.User,
				Timeout: time.Second * 10,
			})
			suite.NoError(err)
			session, err := client.NewSession()
			suite.NoError(err)
			defer session.Close()

			data, err := session.Output("cat ~/.ssh/authorized_keys")
			suite.NoError(err)
			keys := parseAuthorizedKeys(data)
			suite.Equal(tt.expectedKeyCount, len(keys))

			connectionConfig.Auth = sshpkg.Auth{
				Type:       sshpkg.AuthPrivateKey,
				PrivateKey: []byte(testSshKey),
			}
			pingResult := suite.service.Ping(connectionConfig)
			suite.NoError(pingResult)
		})
	}
}

func TestServiceConnectionSuite(t *testing.T) {
	suite.Run(t, new(ConnectionServiceTestSuite))
}
