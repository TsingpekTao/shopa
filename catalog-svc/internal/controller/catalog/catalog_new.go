package catalog

import (
	catalogapi "github.com/TsingpekTao/shopa/catalog-svc/api/catalog"
	"github.com/TsingpekTao/shopa/catalog-svc/internal/service"
)

type ControllerV1 struct {
	catalog service.ICatalog
}

func NewV1() catalogapi.ICatalogV1 {
	return &ControllerV1{catalog: service.Catalog()}
}
