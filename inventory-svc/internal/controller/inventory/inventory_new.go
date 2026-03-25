package inventory

import (
	inventoryapi "github.com/TsingpekTao/shopa/inventory-svc/api/inventory"
	"github.com/TsingpekTao/shopa/inventory-svc/internal/service"
)

type ControllerV1 struct {
	inventory service.IInventory
}

func NewV1() inventoryapi.IInventoryV1 {
	return &ControllerV1{inventory: service.Inventory()}
}
