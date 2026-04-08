package main

import (
	_ "github.com/TsingpekTao/shopa/aftersale-svc/internal/packed"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/TsingpekTao/shopa/aftersale-svc/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
