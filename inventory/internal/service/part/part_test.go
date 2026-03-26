package part

import (
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/zhenklchhh/KozProject/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	tc := &struct {
		part *model.Part
	}{
		part: &model.Part{
			Name:          "shovel",
			Price:         22.3,
			StockQuantity: 3,
		},
	}
	partUUID := gofakeit.UUID()
	tc.part.UUID = uuid.MustParse(partUUID)
	s.invRepo.On("GetPart", s.ctx, tc.part.UUID).Return(tc.part, nil)
	resp, err := s.service.GetPart(s.ctx, partUUID)
	s.Require().Equal(resp, tc.part)
	s.Require().NoError(err)
	s.invRepo.AssertCalled(s.T(), "GetPart", s.ctx, tc.part.UUID)
}

func (s *ServiceSuite) TestGetPartRepoError() {
	partUUID := gofakeit.UUID()
	id := uuid.MustParse(partUUID)
	s.invRepo.On("GetPart", s.ctx, id).Return(nil, errors.New("repo error"))
	resp, err := s.service.GetPart(s.ctx, partUUID)
	s.Require().Empty(resp)
	s.Require().Error(err)
	s.invRepo.AssertCalled(s.T(), "GetPart", s.ctx, id)
}

func (s *ServiceSuite) TestListPartsSuccess() {
	partUuids := []uuid.UUID{uuid.MustParse(gofakeit.UUID()), uuid.MustParse(gofakeit.UUID())}
	pf := &model.PartFilter{Uuids: partUuids}
	parts := []*model.Part{
		{
			UUID:     partUuids[0],
			Name:     "shovel",
			Price:    22.3,
			Category: model.CategoryUnspecified,
		},
		{
			UUID:     partUuids[1],
			Name:     "gas",
			Price:    100.8,
			Category: model.CategoryFuel,
		},
	}

	s.invRepo.On("ListParts", s.ctx, pf).Return(parts, nil)
	resp, err := s.service.ListParts(s.ctx, pf)
	s.Require().Equal(resp, parts)
	s.Require().NoError(err)
	s.invRepo.AssertCalled(s.T(), "ListParts", s.ctx, pf)
}

func (s *ServiceSuite) TestListPartsErr() {
	partUuids := []uuid.UUID{uuid.MustParse(gofakeit.UUID()), uuid.MustParse(gofakeit.UUID())}
	pf := &model.PartFilter{Uuids: partUuids}

	s.invRepo.On("ListParts", s.ctx, pf).Return(nil, errors.New("repo err"))
	resp, err := s.service.ListParts(s.ctx, pf)
	s.Require().Empty(resp)
	s.Require().Error(err)
	s.invRepo.AssertCalled(s.T(), "ListParts", s.ctx, pf)
}
