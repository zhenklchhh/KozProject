package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	serviceMock "github.com/zhenklchhh/KozProject/inventory/internal/service/mocks"
)

type ApiSuite struct {
	suite.Suite

	ctx     context.Context
	service *serviceMock.InventoryService
	handler *api
}

func (s *ApiSuite) SetupTest() {
	s.ctx = context.Background()
	s.service = serviceMock.NewInventoryService(s.T())
	s.handler = NewApi(s.service)
}

func (s *ApiSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
