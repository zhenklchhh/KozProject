package converter

import (
	"github.com/zhenklchhh/KozProject/payment/internal/model"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

// PayOrderRequestServiceToRepo конвертирует protobuf PayOrderRequest в repository model PayOrderRequest
func PayOrderRequestServiceToRepo(req *paymentV1.PayOrderRequest) *model.PayOrderRequest {
	if req == nil {
		return nil
	}

	return &model.PayOrderRequest{
		OrderUuid:     req.GetOrderUuid(),
		UserUuid:      req.GetUserUuid(),
		PaymentMethod: PaymentMethodServiceToRepo(req.GetPaymentMethod()),
	}
}

// PayOrderRequestRepoToService конвертирует repository model PayOrderRequest в protobuf PayOrderRequest
func PayOrderRequestRepoToService(req *model.PayOrderRequest) *paymentV1.PayOrderRequest {
	if req == nil {
		return nil
	}

	return &paymentV1.PayOrderRequest{
		OrderUuid:     req.OrderUuid,
		UserUuid:      req.UserUuid,
		PaymentMethod: PaymentMethodRepoToService(req.PaymentMethod),
	}
}

// PaymentMethodServiceToRepo конвертирует protobuf PaymentMethod в repository model PaymentMethod
func PaymentMethodServiceToRepo(method paymentV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

// PaymentMethodRepoToService конвертирует repository model PaymentMethod в protobuf PaymentMethod
func PaymentMethodRepoToService(method model.PaymentMethod) paymentV1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

// PayOrderResponseServiceToRepo конвертирует protobuf PayOrderResponse в repository model PayOrderResponse
func PayOrderResponseServiceToRepo(resp *paymentV1.PayOrderResponse) *model.PayOrderResponse {
	if resp == nil {
		return nil
	}

	return &model.PayOrderResponse{
		TransactionUuid: resp.GetTransactionUuid(),
	}
}

// PayOrderResponseRepoToService конвертирует repository model PayOrderResponse в protobuf PayOrderResponse
func PayOrderResponseRepoToService(resp *model.PayOrderResponse) *paymentV1.PayOrderResponse {
	if resp == nil {
		return nil
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: resp.TransactionUuid,
	}
}
