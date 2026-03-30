package order

import (
	orderapi "github.com/TsingpekTao/shopa/order-svc/api/order"
	"github.com/TsingpekTao/shopa/order-svc/internal/service"
)

type ControllerV1 struct {
	order service.IOrder
}

func NewV1() orderapi.IOrderV1 {
	return &ControllerV1{order: service.Order()}
}
