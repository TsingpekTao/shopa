// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package hello

import (
	helloapi "github.com/TsingpekTao/shopa/agent-svc/api/hello"
)

type ControllerV1 struct{}

func NewV1() helloapi.IHelloV1 {
	return &ControllerV1{}
}
