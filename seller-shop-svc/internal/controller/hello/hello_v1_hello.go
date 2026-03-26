package hello

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/hello/v1"
	"github.com/gogf/gf/v2/frame/g"
)

// Hello 返回基础连通性测试响应。
func (c *ControllerV1) Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error) {
	g.RequestFromCtx(ctx).Response.Writeln("Hello World!")
	return
}
