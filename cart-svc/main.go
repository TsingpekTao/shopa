package main

import (
	_ "github.com/TsingpekTao/shopa/cart-svc/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/TsingpekTao/shopa/cart-svc/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
