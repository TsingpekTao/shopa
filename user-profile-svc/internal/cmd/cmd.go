package cmd

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"google.golang.org/grpc"

	"github.com/TsingpekTao/shopa/user-profile-svc/internal/controller/api"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/router"
	"github.com/TsingpekTao/shopa/user-profile-svc/internal/worker"
)

var (
	// Main 鏄繘绋嬪叆鍙ｅ懡浠ゃ€?	// GoFrame 鍚姩鏃朵細鎵ц杩欎釜鍛戒护锛屼粠杩欓噷缁熶竴鎷夎捣 gRPC 涓?HTTP 涓ゅ鏈嶅姟銆
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start grpc server and http server",
		Func:  mainFunc,
	}
)

// mainFunc 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func mainFunc(ctx context.Context, parser *gcmd.Parser) (err error) {
	// 鍚姩鈥滄敞鍐屽悗鍒濆鍖栬祫鏂欌€濇秷鎭秷璐硅€呫€
	worker.StartRegisterInitConsumer(ctx)

	// gRPC 鐢ㄥ崗绋嬪惎鍔紝閬垮厤闃诲褰撳墠涓诲崗绋嬶紝纭繚 HTTP 涔熻兘鍦ㄥ悓杩涚▼鍐呭悓鏃跺惎鍔ㄣ€
	go func() {
		// 鍒涘缓 gRPC 鏈嶅姟鍣ㄩ厤缃紝鐢ㄤ簬鎸傝浇鎷︽埅鍣ㄩ摼銆
		c := grpcx.Server.NewConfig()
		// 鍦?Unary 鎷︽埅鍣ㄩ摼涓帴鍏ュ弬鏁版牎楠岋紝灏芥棭鎷︽埅闈炴硶璇锋眰銆
		c.Options = append(c.Options, []grpc.ServerOption{
			grpcx.Server.ChainUnary(
				grpcx.Server.UnaryValidate,
			),
		}...)
		// 鍩轰簬閰嶇疆鏋勯€?gRPC 鏈嶅姟瀹炰緥銆
		s := grpcx.Server.New(c)
		// 娉ㄥ唽 user-profile 鐨?RPC 澶勭悊鍣ㄣ€
		api.Register(s)
		// 鍚姩骞堕樆濉炲湪 gRPC 鏈嶅姟鍗忕▼涓€
		s.Run()
	}()

	// 鏋勫缓 HTTP 鏈嶅姟锛圫wagger/OpenAPI + 瀵瑰 HTTP 鎺ュ彛锛夈€
	httpServer := g.Server()
	// 浣跨敤缁熶竴鍝嶅簲涓棿浠讹紝淇濊瘉杩斿洖缁撴瀯涓€鑷淬€
	httpServer.Use(ghttp.MiddlewareHandlerResponse)
	// 鎸夋ā鍧楁敞鍐岃矾鐢卞垎缁勩€
	router.RegisterHTTP(httpServer)
	// 鑷畾涔?OpenAPI 鏂囨。淇℃伅銆
	router.EnhanceOpenAPIDoc(httpServer)
	// 鍚姩 HTTP 骞堕樆濉炰富鍗忕▼銆
	httpServer.Run()
	return nil
}
