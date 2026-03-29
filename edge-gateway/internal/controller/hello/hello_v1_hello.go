package hello

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/edge-gateway/api/hello/v1"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error) {
	g.RequestFromCtx(ctx).Response.Writeln(g.I18n().T(ctx, "hello_world"))
	return
}
