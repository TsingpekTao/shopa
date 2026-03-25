package main

import (
	_ "github.com/TsingpekTao/shopa/iam-svc/internal/logic"
	_ "github.com/TsingpekTao/shopa/iam-svc/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/TsingpekTao/shopa/iam-svc/internal/cmd"
)

// main 是 iam-svc 的启动入口，初始化命令并运行服务。
func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}

