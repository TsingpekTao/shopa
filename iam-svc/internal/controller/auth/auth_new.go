package auth

import (
	authapi "github.com/TsingpekTao/shopa/iam-svc/api/auth"
	"github.com/TsingpekTao/shopa/iam-svc/internal/service"
)

type ControllerV1 struct {
	auth service.IAuth
}

// NewV1 创建 auth 控制器并与 service 层适配。
func NewV1() authapi.IAuthV1 {
	return &ControllerV1{auth: service.Auth()}
}
