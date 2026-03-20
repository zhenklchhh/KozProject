package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	serviceMock "github.com/zhenklchhh/KozProject/order/internal/service/mocks"
)

type ApiSuite struct {
	suite.Suite

	ctx     context.Context
	service *serviceMock.OrderService
	handler *api
}

func (s *ApiSuite) SetupTest() {
	s.ctx = context.Background()
	s.service = serviceMock.NewOrderService(s.T())
	s.handler = NewApi(s.service)
}

func (s *ApiSuite) TearDownTest() {
}

func TestApiIntegration(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
