package aftersale

import (
	aftersaleapi "github.com/TsingpekTao/shopa/aftersale-svc/api/aftersale"
	"github.com/TsingpekTao/shopa/aftersale-svc/internal/service"
)

type ControllerV1 struct {
	afterSale service.IAfterSale
}

func NewV1() aftersaleapi.IAftersaleV1 {
	return &ControllerV1{afterSale: service.AfterSale()}
}
