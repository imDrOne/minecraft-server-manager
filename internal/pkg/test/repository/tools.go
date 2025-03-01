package repository

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var ErrInternalSql = errors.New("DB error")

type MockQueriesProvider[T any] func() T

type RepoTestSuite[T, R any] struct {
	suite.Suite
	Ctx               context.Context
	MockQuerySupplier MockQueriesProvider[T]
	RepoSupplier      func() R
}

func (suite *RepoTestSuite[T, R]) SetupTest(newMockQueries func(ctrl *gomock.Controller) T, newRepo func(T) R) {
	var mockQueries T

	suite.Ctx = context.Background()
	suite.MockQuerySupplier = func() T {
		ctrl := gomock.NewController(suite.T())
		mockQueries = newMockQueries(ctrl)
		return mockQueries
	}

	suite.RepoSupplier = func() R {
		return newRepo(mockQueries)
	}
}
