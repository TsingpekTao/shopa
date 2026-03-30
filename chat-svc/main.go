package main

import (
	_ "github.com/TsingpekTao/shopa/chat-svc/internal/logic"
	_ "github.com/TsingpekTao/shopa/chat-svc/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/TsingpekTao/shopa/chat-svc/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
