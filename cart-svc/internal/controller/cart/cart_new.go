package cart

import (
	cartapi "github.com/TsingpekTao/shopa/cart-svc/api/cart"
	"github.com/TsingpekTao/shopa/cart-svc/internal/service"
)

type ControllerV1 struct {
	cart service.ICart
}

func NewV1() cartapi.ICartV1 {
	return &ControllerV1{cart: service.Cart()}
}
