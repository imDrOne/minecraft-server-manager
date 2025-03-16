package test

import (
	"context"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type MockProvider[T any] func() T

type Suite[M, T any] struct {
	suite.Suite
	Ctx            context.Context
	MockSupplier   MockProvider[M]
	TargetSupplier func() T
}

func (suite *Suite[T, R]) SetupTest(newMock func(ctrl *gomock.Controller) T, newTarget func(T) R) {
	var mockQueries T

	suite.Ctx = context.Background()
	suite.MockSupplier = func() T {
		ctrl := gomock.NewController(suite.T())
		mockQueries = newMock(ctrl)
		return mockQueries
	}

	suite.TargetSupplier = func() R {
		return newTarget(mockQueries)
	}
}
