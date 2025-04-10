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
	var mock T

	suite.Ctx = context.Background()
	suite.MockSupplier = func() T {
		ctrl := gomock.NewController(suite.T())
		mock = newMock(ctrl)
		return mock
	}

	suite.TargetSupplier = func() R {
		return newTarget(mock)
	}
}
