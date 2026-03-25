package seller

import "github.com/TsingpekTao/shopa/edge-gateway/api/seller"

type ControllerV1 struct{}

func NewV1() seller.ISellerV1 {
	return &ControllerV1{}
}
