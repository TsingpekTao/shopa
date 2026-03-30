package main

import (
	_ "github.com/TsingpekTao/shopa/order-svc/internal/packed"

	_ "github.com/TsingpekTao/shopa/order-svc/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/TsingpekTao/shopa/order-svc/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
