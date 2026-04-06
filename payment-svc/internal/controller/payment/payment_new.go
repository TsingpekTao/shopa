package payment

import (
	paymentapi "github.com/TsingpekTao/shopa/payment-svc/api/payment"
	"github.com/TsingpekTao/shopa/payment-svc/internal/service"
)

type ControllerV1 struct {
	payment service.IPayment
}

func NewV1() paymentapi.IPaymentV1 {
	return &ControllerV1{payment: service.Payment()}
}
