package main

import (
	_ "github.com/TsingpekTao/shopa/notification-svc/internal/packed"

	_ "github.com/TsingpekTao/shopa/notification-svc/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/TsingpekTao/shopa/notification-svc/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
