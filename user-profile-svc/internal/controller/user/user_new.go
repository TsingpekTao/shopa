// =================================================================================
// This controller file is modeled after GoFrame CLI generated controller style.
// =================================================================================

package user

import (
	"github.com/TsingpekTao/shopa/user-profile-svc/api/user"
)

// ControllerV1 implements api/user/v1 endpoints.
type ControllerV1 struct{}

// NewV1 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func NewV1() user.IUserV1 {
	return &ControllerV1{}
}
