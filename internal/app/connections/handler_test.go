package connections

import (
	testutils "github.com/imDrOne/minecraft-server-manager/internal/pkg/test"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	suite.Contains("invalid json\n", string(data))
}

func Test(t *testing.T) {
	suite.Run(t, new(ConnectionHandlerTestSuite))
}
