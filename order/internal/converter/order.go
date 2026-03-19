package converter

import (
	"github.com/zhenklchhh/KozProject/order/internal/model"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

func OrderServiceToRepo(order *orderV1.Order) *model.Order {
	if order == nil {
		return nil
	}

	var transactionUUID *string
	if order.TransactionUUID.Set {
		if order.TransactionUUID.Null {
			transactionUUID = nil
		} else {
			transactionUUID = &order.TransactionUUID.Value
		}
	}

	var paymentMethod *model.PaymentMethod
	if order.PaymentMethod.Set {
		paymentMethodValue := PaymentMethodServiceToRepo(order.PaymentMethod.Value)
		paymentMethod = &paymentMethodValue
	}

	return &model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          OrderStatusServiceToRepo(order.Status),
	}
}

func OrderRepoToService(order *model.Order) *orderV1.Order {
	if order == nil {
		return nil
	}

	var transactionUUID orderV1.OptNilString
	if order.TransactionUUID != nil {
		transactionUUID = orderV1.OptNilString{
			Value: *order.TransactionUUID,
			Set:   true,
		}
	} else {
		transactionUUID = orderV1.OptNilString{
			Null: true,
			Set:  true,
		}
	}

	var paymentMethod orderV1.OptPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderV1.NewOptPaymentMethod(PaymentMethodRepoToService(*order.PaymentMethod))
	}

	return &orderV1.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          OrderStatusRepoToService(order.Status),
	}
}

func CreateOrderRequestServiceToRepo(req *orderV1.CreateOrderRequest) *model.CreateOrderRequest {
	if req == nil {
		return nil
	}

	return &model.CreateOrderRequest{
		UserUUID:  req.UserUUID,
		PartUuids: req.PartUuids,
	}
}

func CreateOrderRequestRepoToService(req *model.CreateOrderRequest) *orderV1.CreateOrderRequest {
	if req == nil {
		return nil
	}

	return &orderV1.CreateOrderRequest{
		UserUUID:  req.UserUUID,
		PartUuids: req.PartUuids,
	}
}

func CreateOrderResponseServiceToRepo(resp *orderV1.CreateOrderResponse) *model.CreateOrderResponse {
	if resp == nil {
		return nil
	}

	return &model.CreateOrderResponse{
		OrderUUID:  resp.OrderUUID,
		TotalPrice: resp.TotalPrice,
	}
}

func CreateOrderResponseRepoToService(resp *model.CreateOrderResponse) *orderV1.CreateOrderResponse {
	if resp == nil {
		return nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  resp.OrderUUID,
		TotalPrice: resp.TotalPrice,
	}
}

func PayOrderRequestServiceToRepo(req *orderV1.PayOrderRequest) *model.PayOrderRequest {
	if req == nil {
		return nil
	}

	return &model.PayOrderRequest{
		PaymentMethod: PaymentMethodServiceToRepo(req.PaymentMethod),
	}
}

func PayOrderRequestRepoToService(req *model.PayOrderRequest) *orderV1.PayOrderRequest {
	if req == nil {
		return nil
	}

	return &orderV1.PayOrderRequest{
		PaymentMethod: PaymentMethodRepoToService(req.PaymentMethod),
	}
}

func PayOrderResponseServiceToRepo(resp *orderV1.PayOrderResponse) *model.PayOrderResponse {
	if resp == nil {
		return nil
	}

	return &model.PayOrderResponse{
		TransactionUUID: resp.TransactionUUID,
	}
}

func PayOrderResponseRepoToService(resp *model.PayOrderResponse) *orderV1.PayOrderResponse {
	if resp == nil {
		return nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: resp.TransactionUUID,
	}
}

func OrderStatusServiceToRepo(status orderV1.OrderStatus) model.OrderStatus {
	switch status {
	case orderV1.OrderStatusPAID:
		return model.OrderStatusPaid
	case orderV1.OrderStatusCANCELLED:
		return model.OrderStatusCancelled
	default:
		return model.OrderStatusPendingPayment
	}
}

func OrderStatusRepoToService(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusPaid:
		return orderV1.OrderStatusPAID
	case model.OrderStatusCancelled:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}

func PaymentMethodServiceToRepo(method orderV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodPAYMENTMETHODCARD:
		return model.PaymentMethodCard
	case orderV1.PaymentMethodPAYMENTMETHODSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodPAYMENTMETHODCREDITCARD:
		return model.PaymentMethodCreditCard
	case orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodCard
	}
}

func PaymentMethodRepoToService(method model.PaymentMethod) orderV1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return orderV1.PaymentMethodPAYMENTMETHODCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodPAYMENTMETHODSBP
	case model.PaymentMethodCreditCard:
		return orderV1.PaymentMethodPAYMENTMETHODCREDITCARD
	case model.PaymentMethodInvestorMoney:
		return orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY
	default:
		return orderV1.PaymentMethodPAYMENTMETHODCARD
	}
}

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
