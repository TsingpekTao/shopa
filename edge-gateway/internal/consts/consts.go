package consts

import "github.com/gogf/gf/v2/errors/gcode"

var (
	// CodeForbidden 表示请求已认证但无权限访问目标资源。
	CodeForbidden = gcode.New(403, "Forbidden", nil)
)
