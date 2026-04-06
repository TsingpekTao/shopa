package search

import (
	searchapi "github.com/TsingpekTao/shopa/search-svc/api/search"
	"github.com/TsingpekTao/shopa/search-svc/internal/service"
)

type ControllerV1 struct {
	search service.ISearch
}

func NewV1() searchapi.ISearchV1 {
	return &ControllerV1{search: service.Search()}
}
