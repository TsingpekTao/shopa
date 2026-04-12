package main

import (
	_ "github.com/TsingpekTao/shopa/agent-svc/internal/logic"
	_ "github.com/TsingpekTao/shopa/agent-svc/internal/packed"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/TsingpekTao/shopa/agent-svc/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
