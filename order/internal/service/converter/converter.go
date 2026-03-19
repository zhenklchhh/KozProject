package converter

import (
	"github.com/zhenklchhh/KozProject/order/internal/model"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

func ConvertPayOrderRequestApiToService(req *paymentV1.PayOrderRequest) *model.PayOrderServiceRequest {
	return &model.PayOrderServiceRequest{
		OrderUuid:     req.OrderUuid,
		UserUuid:      req.UserUuid,
		PaymentMethod: model.PaymentMethod(req.PaymentMethod.String()),
	}
}

func ConvertPayOrderRequestServiceToApi(req *model.PayOrderServiceRequest) *paymentV1.PayOrderRequest {
	return &paymentV1.PayOrderRequest{
		OrderUuid:     req.OrderUuid,
		UserUuid:      req.UserUuid,
		PaymentMethod: paymentV1.PaymentMethod(paymentV1.PaymentMethod_value[string(req.PaymentMethod)]),
	}
}
