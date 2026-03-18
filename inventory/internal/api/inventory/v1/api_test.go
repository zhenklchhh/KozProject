package api

import (
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/zhenklchhh/KozProject/inventory/internal/converter"
	"github.com/zhenklchhh/KozProject/inventory/internal/model"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ApiSuite) TestGetPartSuccess() {
	req := &inventoryV1.GetPartRequest{
		Uuid: gofakeit.UUID(),
	}
	expectedResp := &inventoryV1.GetPartResponse{
		Part: &inventoryV1.Part{
			Uuid:          req.Uuid,
			Name:          "gas",
			Description:   "fuel for space rocket",
			Price:         20.3,
			StockQuantity: 2,
		},
	}
	modelPart := &model.Part{
		Uuid:          req.Uuid,
		Name:          "gas",
		Description:   "fuel for space rocket",
		Price:         20.3,
		StockQuantity: 2,
	}
	s.service.On("GetPart", s.ctx, req.Uuid).Return(modelPart, nil)
	resp, err := s.handler.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().Equal(resp, expectedResp)
	s.service.AssertCalled(s.T(), "GetPart", s.ctx, req.Uuid)
}

func (s *ApiSuite) TestGetPartValidationErr() {
	req := &inventoryV1.GetPartRequest{
		Uuid: "invalid uuid",
	}
	expectedErr := status.Errorf(codes.InvalidArgument, "payment service: validation error")
	resp, err := s.handler.GetPart(s.ctx, req)
	s.Require().Error(err)
	s.Require().Equal(expectedErr, err)
	s.Require().Empty(resp)
	s.service.AssertNotCalled(s.T(), "GetPart", s.ctx, req.Uuid)
}

func (s *ApiSuite) TestGetPartServiceError() {
	req := &inventoryV1.GetPartRequest{
		Uuid: gofakeit.UUID(),
	}
	expectedErr := errors.New("service err")
	s.service.On("GetPart", s.ctx, req.Uuid).Return(nil, expectedErr)
	resp, err := s.handler.GetPart(s.ctx, req)
	s.Require().Error(err)
	s.Require().Equal(expectedErr, err)
	s.Require().Empty(resp)
	s.service.AssertCalled(s.T(), "GetPart", s.ctx, req.Uuid)
}

func (s *ApiSuite) TestListPartsSuccess() {
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartFilter{Uuids: partUuids},
	}
	expectedResp := &inventoryV1.ListPartsResponse{
		Parts: []*inventoryV1.Part{
			{
				Uuid:          partUuids[0],
				Name:          "gas",
				Description:   "fuel for space rocket",
				Price:         20.3,
				StockQuantity: 2,
			},
			{
				Uuid:          partUuids[1],
				Name:          "engine v8",
				Description:   "rocket engine",
				Price:         1000.53,
				StockQuantity: 1,
			},
		},
	}
	modelParts := []*model.Part{
		{
			Uuid:          partUuids[0],
			Name:          "gas",
			Description:   "fuel for space rocket",
			Price:         20.3,
			StockQuantity: 2,
		},
		{
			Uuid:          partUuids[1],
			Name:          "engine v8",
			Description:   "rocket engine",
			Price:         1000.53,
			StockQuantity: 1,
		},
	}
	s.service.On("ListParts", s.ctx, converter.PartFilterServiceToRepo(req.Filter)).Return(modelParts, nil)
	resp, err := s.handler.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().Equal(resp, expectedResp)
	s.service.AssertCalled(s.T(), "ListParts", s.ctx, converter.PartFilterServiceToRepo(req.Filter))
}

func (s *ApiSuite) TestListPartsValidationError() {
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

func (s *ApiSuite) TestListPartsServiceError() {
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	req := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartFilter{Uuids: partUuids},
	}
	s.service.On("ListParts", s.ctx, converter.PartFilterServiceToRepo(req.Filter)).Return(nil, errors.New("service err"))
	resp, err := s.handler.ListParts(s.ctx, req)
	s.Require().Error(err)
	s.Require().Equal(err, errors.New("service err"))
	s.Require().Empty(resp)
}
