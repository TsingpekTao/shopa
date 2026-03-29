package admin

import "github.com/TsingpekTao/shopa/edge-gateway/api/admin"

type ControllerV1 struct{}

func NewV1() admin.IAdminV1 {
	return &ControllerV1{}
}
