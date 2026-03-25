// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package hello

import (
	"github.com/TsingpekTao/shopa/iam-svc/api/hello"
)

type ControllerV1 struct{}

// NewV1 创建 V1 版本的 Hello 控制器实例。
func NewV1() hello.IHelloV1 {
	return &ControllerV1{}
}
