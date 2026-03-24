package api

import (
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zhenklchhh/KozProject/inventory/internal/converter"
	"github.com/zhenklchhh/KozProject/inventory/internal/model"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetPartSuccess() {
	req := &inventoryV1.GetPartRequest{
		Uuid: gofakeit.UUID(),
	}
	modelUuid := uuid.MustParse(req.Uuid)
	expectedResp := &inventoryV1.GetPartResponse{
		Part: &inventoryV1.Part{
			Uuid:          req.GetUuid(),
			Name:          "gas",
			Description:   "fuel for space rocket",
			Price:         20.3,
			StockQuantity: 2,
		},
	}
	modelPart := &model.Part{
		UUID:          modelUuid,
		Name:          "gas",
		Description:   "fuel for space rocket",
		Price:         20.3,
		StockQuantity: 2,
	}
	s.service.On("GetPart", s.ctx, req.GetUuid()).Return(modelPart, nil)
	resp, err := s.handler.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().Equal(resp, expectedResp)
	s.service.AssertCalled(s.T(), "GetPart", s.ctx, req.GetUuid())
}

func (s *APISuite) TestGetPartValidationErr() {
	req := &inventoryV1.GetPartRequest{
		Uuid: "invalid uuid",
	}
	expectedErr := status.Errorf(codes.InvalidArgument, "payment service: validation error")
	resp, err := s.handler.GetPart(s.ctx, req)
	s.Require().Error(err)
	s.Require().Equal(expectedErr, err)
	s.Require().Empty(resp)
	s.service.AssertNotCalled(s.T(), "GetPart", s.ctx, req.GetUuid())
}

func (s *APISuite) TestGetPartServiceError() {
	req := &inventoryV1.GetPartRequest{
		Uuid: gofakeit.UUID(),
	}
	expectedErr := errors.New("service err")
	s.service.On("GetPart", s.ctx, req.GetUuid()).Return(nil, expectedErr)
	resp, err := s.handler.GetPart(s.ctx, req)
	s.Require().Error(err)
	s.Require().Equal(expectedErr, err)
	s.Require().Empty(resp)
	s.service.AssertCalled(s.T(), "GetPart", s.ctx, req.GetUuid())
}

func (s *APISuite) TestListPartsSuccess() {
	partUuidsServer := []string{gofakeit.UUID(), gofakeit.UUID()}
	partUuidsRepo := make([]uuid.UUID, 0, len(partUuidsServer))
	for _, id := range partUuidsServer {
		partUuidsRepo = append(partUuidsRepo, uuid.MustParse(id))
	}
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartFilter{Uuids: partUuidsServer},
	}
	expectedResp := &inventoryV1.ListPartsResponse{
		Parts: []*inventoryV1.Part{
			{
				Uuid:          partUuidsServer[0],
				Name:          "gas",
				Description:   "fuel for space rocket",
				Price:         20.3,
				StockQuantity: 2,
			},
			{
				Uuid:          partUuidsServer[1],
				Name:          "engine v8",
				Description:   "rocket engine",
				Price:         1000.53,
				StockQuantity: 1,
			},
		},
	}
	modelParts := []*model.Part{
		{
			UUID:          partUuidsRepo[0],
			Name:          "gas",
			Description:   "fuel for space rocket",
			Price:         20.3,
			StockQuantity: 2,
		},
		{
			UUID:          partUuidsRepo[1],
			Name:          "engine v8",
			Description:   "rocket engine",
			Price:         1000.53,
			StockQuantity: 1,
		},
	}
	s.service.On("ListParts", s.ctx, converter.PartFilterServiceToRepo(req.GetFilter())).Return(modelParts, nil)
	resp, err := s.handler.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().Equal(resp, expectedResp)
	s.service.AssertCalled(s.T(), "ListParts", s.ctx, converter.PartFilterServiceToRepo(req.GetFilter()))
}

func (s *APISuite) TestListPartsValidationError() {
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartFilter{
			Uuids: []string{
				"invalid",
				"invalid",
			},
		},
	}
	expectedErr := status.Errorf(codes.InvalidArgument, "payment service: validation error")
	resp, err := s.handler.ListParts(s.ctx, req)
	s.Require().Error(err)
	s.Require().Equal(expectedErr, err)
	s.Require().Empty(resp)
	s.service.AssertNotCalled(s.T(), "ListParts", s.ctx, mock.AnythingOfType("*model.PartFilter"))
}

func (s *APISuite) TestListPartsServiceError() {
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartFilter{Uuids: partUuids},
	}
	s.service.On("ListParts", s.ctx, converter.PartFilterServiceToRepo(req.GetFilter())).Return(nil, errors.New("service err"))
	resp, err := s.handler.ListParts(s.ctx, req)
	s.Require().Error(err)
	s.Require().Equal(err, errors.New("service err"))
	s.Require().Empty(resp)
}
