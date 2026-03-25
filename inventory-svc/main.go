package main

import (
	_ "github.com/TsingpekTao/shopa/inventory-svc/internal/packed"

	_ "github.com/TsingpekTao/shopa/inventory-svc/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/TsingpekTao/shopa/inventory-svc/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
