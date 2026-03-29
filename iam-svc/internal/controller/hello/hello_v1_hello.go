package hello

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/hello/v1"
	"github.com/gogf/gf/v2/frame/g"
)

// Hello 简单示例接口，向客户端输出 Hello World。
func (c *ControllerV1) Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error) {
	g.RequestFromCtx(ctx).Response.Writeln(g.I18n().T(ctx, "hello_world"))
	return
}
