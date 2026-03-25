package auth

import (
	authapi "github.com/TsingpekTao/shopa/iam-svc/api/auth"
	"github.com/TsingpekTao/shopa/iam-svc/internal/service"
)

type ControllerV1 struct {
	auth service.IAuth
}

// NewV1 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func NewV1() authapi.IAuthV1 {
	return &ControllerV1{auth: service.Auth()}
}
