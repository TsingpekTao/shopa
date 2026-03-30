package fulfillment

import (
	fulfillmentapi "github.com/TsingpekTao/shopa/fulfillment-svc/api/fulfillment"
	"github.com/TsingpekTao/shopa/fulfillment-svc/internal/service"
)

type ControllerV1 struct {
	fulfillment service.IFulfillment
}

func NewV1() fulfillmentapi.IFulfillmentV1 {
	return &ControllerV1{fulfillment: service.Fulfillment()}
}
