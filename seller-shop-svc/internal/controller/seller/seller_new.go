package seller

import (
	sellerapi "github.com/TsingpekTao/shopa/seller-shop-svc/api/seller"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/service"
)

type ControllerV1 struct {
	svc service.ISellerShop
}

// NewV1 创建 seller-shop HTTP 控制器。
func NewV1() sellerapi.ISellerV1 {
	return &ControllerV1{svc: service.SellerShop()}
}
