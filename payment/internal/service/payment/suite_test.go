package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	repositoryMocks "github.com/zhenklchhh/KozProject/payment/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx         context.Context
	paymentRepo *repositoryMocks.PaymentRepository
	service     *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.paymentRepo = repositoryMocks.NewPaymentRepository(s.T())
	s.service = NewService(s.paymentRepo)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
