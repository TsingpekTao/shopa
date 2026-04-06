package risk

import (
	riskapi "github.com/TsingpekTao/shopa/risk-svc/api/risk"
	"github.com/TsingpekTao/shopa/risk-svc/internal/service"
)

type ControllerV1 struct {
	risk service.IRisk
}

func NewV1() riskapi.IRiskV1 {
	return &ControllerV1{risk: service.Risk()}
}
