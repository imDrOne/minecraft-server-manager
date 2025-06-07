package connections

import (
	"bytes"
	"encoding/json"
	"errors"
	domain "github.com/imDrOne/minecraft-server-manager/internal/domain/connections"
	conndb "github.com/imDrOne/minecraft-server-manager/internal/infrastructure/connections/db"
	testutils "github.com/imDrOne/minecraft-server-manager/internal/pkg/test"
	connservice "github.com/imDrOne/minecraft-server-manager/internal/service/connections"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ConnectionHandlerTestSuite struct {
	testutils.Suite[*MockService, *ConnectionController]
}

func (suite *ConnectionHandlerTestSuite) SetupTest() {
	suite.Suite.SetupTest(
		func(ctrl *gomock.Controller) *MockService {
			return NewMockService(ctrl)
		},
		func(service *MockService) *ConnectionController {
			return &ConnectionController{service: service}
		},
	)
}

func (suite *ConnectionHandlerTestSuite) TestConnectionHandler_Create_EmptyBody() {
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	handler := suite.TargetSupplier()
	handler.Create(w, req)
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
		Create(suite.Ctx, gomock.Any(), gomock.Any()).
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

func (suite *ConnectionHandlerTestSuite) TestConnectionController_Create_InvalidDomain() {
	tests := []struct {
		name        string
		errorText   string
		isEmptyBody bool
		payload     CreateConnectionRequestDto
	}{
		{
			name:        "empty dto",
			isEmptyBody: true,
		},
		{
			name: "invalid user (try 1)",
			payload: CreateConnectionRequestDto{
				User: "invalid!!User&",
			},
		},
		{
			name:    "invalid user (try 2)",
			payload: CreateConnectionRequestDto{},
		},
	}

	handler := suite.TargetSupplier()
	connRepo := conndb.NewConnectionRepository(nil)
	handler.service = connservice.NewConnectionService(connservice.Dependencies{
		ConnRepo: connRepo,
	})

	for _, test := range tests {
		suite.Run(test.name, func() {
			var b bytes.Buffer
			if !test.isEmptyBody {
				err := json.NewEncoder(&b).Encode(test.payload)
				if err != nil {
					suite.FailNow(err.Error())
				}
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
				suite.Error(err, "err on read response body")
			}
			suite.EqualValues(http.StatusBadRequest, res.StatusCode)
			if test.isEmptyBody {
				suite.Contains(string(data), "invalid json")
			} else {
				suite.Contains(string(data), domain.ErrValidationConnection.Error())
			}
		})
	}
}

func (suite *ConnectionHandlerTestSuite) TestConnectionController_Create_InvalidPathId() {
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

	handler := suite.TargetSupplier()

	for _, test := range tests {
		suite.Run(test.name, func() {
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
				suite.Fail("err on read response body")
			}
			suite.EqualValues(http.StatusBadRequest, res.StatusCode)
			suite.Contains(string(data), test.error)
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

func (suite *ConnectionHandlerTestSuite) TestConnectionController_FindById_InvalidPathId() {
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

	handler := suite.TargetSupplier()

	for _, test := range tests {
		suite.Run(test.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.SetPathValue("nodeId", test.value)

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
				suite.Fail("err on read response body")
			}
			suite.EqualValues(http.StatusBadRequest, res.StatusCode)
			suite.Contains(string(data), test.error)
		})
	}
}

func Test(t *testing.T) {
	suite.Run(t, new(ConnectionHandlerTestSuite))
}
