package promotion

import (
	promotionapi "github.com/TsingpekTao/shopa/promotion-svc/api/promotion"
	"github.com/TsingpekTao/shopa/promotion-svc/internal/service"
)

type ControllerV1 struct {
	promotion service.IPromotion
}

func NewV1() promotionapi.IPromotionV1 {
	return &ControllerV1{promotion: service.Promotion()}
}
