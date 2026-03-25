package hello

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/hello/v1"
	"github.com/gogf/gf/v2/frame/g"
)

// Hello 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (c *ControllerV1) Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error) {
	g.RequestFromCtx(ctx).Response.Writeln("Hello World!")
	return
}
