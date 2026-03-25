package main

import (
	_ "github.com/TsingpekTao/shopa/user-profile-svc/internal/packed"

	_ "github.com/TsingpekTao/shopa/user-profile-svc/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/TsingpekTao/shopa/user-profile-svc/internal/cmd"
)

// main 绋嬪簭鍏ュ彛锛氬惎鍔ㄥ懡浠ゅ苟杩涘叆鏈嶅姟鐢熷懡鍛ㄦ湡銆
func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
