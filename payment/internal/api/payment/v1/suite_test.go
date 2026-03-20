package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	serviceMocks "github.com/zhenklchhh/KozProject/payment/internal/service/mocks"
)

type ApiSuite struct {
	suite.Suite

	ctx     context.Context
	service *serviceMocks.PaymentService
	handler *api
}

func (s *ApiSuite) SetupTest() {
	s.ctx = context.Background()
	s.service = serviceMocks.NewPaymentService(s.T())
	s.handler = NewApi(s.service)
}

func (s *ApiSuite) TearDownTest() {
}

func TestApiIntegration(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
