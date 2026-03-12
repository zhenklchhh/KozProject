package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
	payModel "github.com/zhenklchhh/KozProject/payment/internal/model"
)

type repository struct{}

func NewRepository() *repository {
	return &repository{}
}

func (s *repository) PayOrder(_ context.Context, req *payModel.PayOrderRequest) (string, error) {
	transactionUUID := uuid.New()
	log.Printf("Оплата прошла успешно, transaction_uuid: %v\n", transactionUUID)
	return transactionUUID.String(), nil
}
