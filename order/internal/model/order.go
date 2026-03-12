package model

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
}

type CreateOrderRequest struct {
	UserUUID  string
	PartUuids []string
}

type CreateOrderResponse struct {
	OrderUUID  string
	TotalPrice float64
}

type PayOrderRequest struct {
	PaymentMethod PaymentMethod
}

type PayOrderResponse struct {
	TransactionUUID string
}

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodCard           PaymentMethod = "PAYMENT_METHOD_CARD"
	PaymentMethodSBP            PaymentMethod = "PAYMENT_METHOD_SBP"
	PaymentMethodCreditCard     PaymentMethod = "PAYMENT_METHOD_CREDIT_CARD"
	PaymentMethodInvestorMoney  PaymentMethod = "PAYMENT_METHOD_INVESTOR_MONEY"
)

