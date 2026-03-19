package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	serviceMock "github.com/zhenklchhh/KozProject/inventory/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx     context.Context
	service *serviceMock.InventoryService
	handler *api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()
	s.service = serviceMock.NewInventoryService(s.T())
	s.handler = NewAPI(s.service)
}

func (s *APISuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
