package connections

import (
	"bytes"
	_ "embed"
	"encoding/json"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	"github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections"
	testutils "github.com/imDrOne/minecraft-server-manager/internal/pkg/test"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

//go:embed test_key.pub
var validSSHKey string

type ConnectionHandlerTestSuite struct {
	testutils.Suite[*MockRepository, *ConnectionController]
}

func (suite *ConnectionHandlerTestSuite) SetupTest() {
	suite.Suite.SetupTest(
		func(ctrl *gomock.Controller) *MockRepository {
			return NewMockRepository(ctrl)
		},
		func(repository *MockRepository) *ConnectionController {
			return &ConnectionController{repo: repository}
		},
	)
}

func (suite *ConnectionHandlerTestSuite) TestConnectionHandler_Create_EmptyBody() {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	suite.TargetSupplier().Create(w, req)
	res := w.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			panic("err on writing result")
		}
	}()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Fail("err on read response body")
	}
	suite.EqualValues(http.StatusBadRequest, res.StatusCode)
	suite.Contains(string(data), "invalid json")
}

func (suite *ConnectionHandlerTestSuite) TestConnectionHandler_Create_DuplicateErr() {
	mockRepo := suite.MockSupplier()
	mockRepo.EXPECT().
		Save(suite.Ctx, gomock.Any(), gomock.Any()).
		Return(nil, domain.ErrConnectionAlreadyExists)

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(CreateConnectionRequestDto{})
	if err != nil {
		suite.T().Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/", &b)
	w := httptest.NewRecorder()

	suite.TargetSupplier().Create(w, req)
	res := w.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			panic("err on writing result")
		}
	}()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Fail("err on read response body")
	}
	suite.EqualValues(http.StatusConflict, res.StatusCode)
	suite.Contains(string(data), domain.ErrConnectionAlreadyExists.Error())
}

func TestConnectionController_Create_InvalidDomain(t *testing.T) {
	tests := []struct {
		name    string
		payload CreateConnectionRequestDto
	}{
		{
			name:    "empty dto",
			payload: CreateConnectionRequestDto{},
		},
		{
			name: "invalid user",
			payload: CreateConnectionRequestDto{
				Key: validSSHKey,
			},
		},
		{
			name: "invalid key",
			payload: CreateConnectionRequestDto{
				User: "happy-user",
			},
		},
	}

	repo := connections.NewConnectionRepository(nil)
	handler := ConnectionController{repo}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var b bytes.Buffer
			err := json.NewEncoder(&b).Encode(test.payload)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/", &b)
			w := httptest.NewRecorder()
			handler.Create(w, req)
			res := w.Result()
			defer func() {
				if err := res.Body.Close(); err != nil {
					panic("err on writing result")
				}
			}()

			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Error("err on read response body")
			}
			require.EqualValues(t, http.StatusBadRequest, res.StatusCode)
			require.Contains(t, string(data), domain.ErrValidationConnection.Error())
		})
	}
}

func Test(t *testing.T) {
	suite.Run(t, new(ConnectionHandlerTestSuite))
}
