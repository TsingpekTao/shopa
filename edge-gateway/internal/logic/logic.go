// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package logic

import (
	// 通过导入各子模块的空包，触发它们的 init() 并完成 service.RegisterXXX 注册。
	// 这样 controller 在调用 service 门面时，可以获得已注册的实现。
	_ "github.com/TsingpekTao/shopa/edge-gateway/internal/logic/auth"
	_ "github.com/TsingpekTao/shopa/edge-gateway/internal/logic/bff"
	_ "github.com/TsingpekTao/shopa/edge-gateway/internal/logic/proxy"
)
