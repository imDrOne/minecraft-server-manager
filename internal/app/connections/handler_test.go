package connections

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
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

func originalHandler() *ConnectionController {
	repo := connections.NewConnectionRepository(nil)
	return &ConnectionController{repo}
}

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

func TestConnectionHandler_Create_EmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	handler := originalHandler()
	handler.Create(w, req)
	res := w.Result()
	defer func() {
		if err := res.Body.Close(); err != nil {
			panic("err on writing result")
		}
	}()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		require.Fail(t, "err on read response body")
	}
	require.EqualValues(t, http.StatusBadRequest, res.StatusCode)
	require.Contains(t, string(data), "invalid json")
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

	handler := originalHandler()

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

func TestConnectionController_Create_InvalidPathId(t *testing.T) {
	tests := []struct {
		name  string
		value string
		error string
	}{
		{
			name:  "empty id",
			value: "",
			error: "expected id - got empty string",
		},
		{
			name:  "not numeric id",
			value: "test",
			error: "error during parsing id",
		},
	}

	handler := originalHandler()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/", nil)
			req.SetPathValue("id", test.value)

			w := httptest.NewRecorder()
			handler.Update(w, req)
			res := w.Result()
			defer func() {
				if err := res.Body.Close(); err != nil {
					panic("err on writing result")
				}
			}()

			data, err := io.ReadAll(res.Body)
			if err != nil {
				require.Fail(t, "err on read response body")
			}
			require.EqualValues(t, http.StatusBadRequest, res.StatusCode)
			require.Contains(t, string(data), test.error)
		})
	}
}

func (suite *ConnectionHandlerTestSuite) TestConnectionController_Update_BusinessLogicResult() {
	tests := []struct {
		name           string
		expectedError  error
		responseStatus int
		wantErr        bool
	}{
		{
			name:           "validation error",
			wantErr:        true,
			expectedError:  domain.ErrValidationConnection,
			responseStatus: http.StatusBadRequest,
		},
		{
			name:           "not found error",
			wantErr:        true,
			expectedError:  domain.ErrConnectionNotFound,
			responseStatus: http.StatusNotFound,
		},
		{
			name:           "internal error",
			wantErr:        true,
			expectedError:  errors.New("internal error"),
			responseStatus: http.StatusInternalServerError,
		},
		{
			name:           "successfully update",
			expectedError:  nil,
			responseStatus: http.StatusNoContent,
			wantErr:        false,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {
			mockRepo := suite.MockSupplier()
			mockRepo.EXPECT().
				Update(suite.Ctx, gomock.Any(), gomock.Any()).
				Return(test.expectedError)

			var b bytes.Buffer
			err := json.NewEncoder(&b).Encode(UpdateConnectionRequestDto{})
			if err != nil {
				suite.T().Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPut, "/", &b)
			req.SetPathValue("id", "101")

			w := httptest.NewRecorder()

			suite.TargetSupplier().Update(w, req)
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
			suite.EqualValues(test.responseStatus, res.StatusCode)
			if test.wantErr {
				suite.Contains(string(data), test.expectedError.Error())
			} else {
				suite.Empty(string(data))
			}
		})
	}
}

func TestConnectionController_FindById_InvalidPathId(t *testing.T) {
	tests := []struct {
		name  string
		value string
		error string
	}{
		{
			name:  "empty id",
			value: "",
			error: "expected id - got empty string",
		},
		{
			name:  "not numeric id",
			value: "test",
			error: "error during parsing id",
		},
	}

	handler := originalHandler()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.SetPathValue("node-id", test.value)

			w := httptest.NewRecorder()
			handler.FindById(w, req)
			res := w.Result()
			defer func() {
				if err := res.Body.Close(); err != nil {
					panic("err on writing result")
				}
			}()

			data, err := io.ReadAll(res.Body)
			if err != nil {
				require.Fail(t, "err on read response body")
			}
			require.EqualValues(t, http.StatusBadRequest, res.StatusCode)
			require.Contains(t, string(data), test.error)
		})
	}
}

func Test(t *testing.T) {
	suite.Run(t, new(ConnectionHandlerTestSuite))
}
