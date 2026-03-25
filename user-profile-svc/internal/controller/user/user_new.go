// =================================================================================
// 该控制器文件借鉴 GoFrame CLI 的生成风格以保持结构一致。
// =================================================================================

package user

import (
	"github.com/TsingpekTao/shopa/user-profile-svc/api/user"
)

// ControllerV1 实现 api/user/v1 的 HTTP 端点。
type ControllerV1 struct{}

// NewV1 构造 User API V1 控制器。
func NewV1() user.IUserV1 {
	return &ControllerV1{}
}
