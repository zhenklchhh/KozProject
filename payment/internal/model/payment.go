package model

type PayOrderRequest struct {
	OrderUuid     string
	UserUuid      string
	PaymentMethod PaymentMethod
}

type PayOrderResponse struct {
	TransactionUuid string
}

type PaymentMethod int32

const (
	PaymentMethodUnspecified PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

var (
	PaymentMethod_name = map[PaymentMethod]string{
		PaymentMethodUnspecified:   "PAYMENT_METHOD_UNSPECIFIED",
		PaymentMethodCard:          "PAYMENT_METHOD_CARD",
		PaymentMethodSBP:           "PAYMENT_METHOD_SBP",
		PaymentMethodCreditCard:    "PAYMENT_METHOD_CREDIT_CARD",
		PaymentMethodInvestorMoney: "PAYMENT_METHOD_INVESTOR_MONEY",
	}
	PaymentMethod_value = map[string]PaymentMethod{
		"PAYMENT_METHOD_UNSPECIFIED":    PaymentMethodUnspecified,
		"PAYMENT_METHOD_CARD":           PaymentMethodCard,
		"PAYMENT_METHOD_SBP":            PaymentMethodSBP,
		"PAYMENT_METHOD_CREDIT_CARD":    PaymentMethodCreditCard,
		"PAYMENT_METHOD_INVESTOR_MONEY": PaymentMethodInvestorMoney,
	}
)
