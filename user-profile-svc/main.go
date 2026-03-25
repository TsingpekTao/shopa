package main

import (
	_ "github.com/TsingpekTao/shopa/user-profile-svc/internal/packed"

	_ "github.com/TsingpekTao/shopa/user-profile-svc/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/TsingpekTao/shopa/user-profile-svc/internal/cmd"
)

// main 是用户画像服务的入口，负责启动全局命令并执行主流程。
func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
