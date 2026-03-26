package main

import (
	_ "github.com/TsingpekTao/shopa/seller-shop-svc/internal/packed"

	_ "github.com/TsingpekTao/shopa/seller-shop-svc/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/cmd"
)

// main 是 seller-shop-svc 的程序入口。
func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
