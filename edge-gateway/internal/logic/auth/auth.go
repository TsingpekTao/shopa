package auth

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	iamv1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// sAuth 负责网关侧统一鉴权能力：
// 1) 从 HTTP 请求中提取 AccessToken；
// 2) 调用 iam-svc 校验 token；
// 3) 将 IAM 的枚举状态转换为前端可读的字符串码。
type sAuth struct {
	// once 用于保证 gRPC 客户端仅初始化一次，避免并发重复建连。
	once sync.Once

	// iamConn 为到 iam-svc 的长连接。
	iamConn *grpc.ClientConn
	// iamClient 为 IAM 内部鉴权 RPC 客户端。
	iamClient iamv1.InternalServiceClient
	// initErr 记录首次初始化错误，后续调用直接复用该错误结果。
	initErr error
}

// New 创建鉴权服务实例。
func New() *sAuth {
	return &sAuth{}
}

// init 在程序启动时将当前实现注册到 service 门面。
func init() {
	service.RegisterAuth(New())
}

// ExtractAccessToken 从请求头中提取 Bearer Token。
// 兼容两种格式：
// 1) Authorization: Bearer <token>
// 2) Authorization: <token>
func (s *sAuth) ExtractAccessToken(r *ghttp.Request) string {
	// 请求为空时直接返回空串，调用方按未登录处理。
	if r == nil {
		return ""
	}
	// 去掉首尾空白，避免网关层大小写/空格差异影响。
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if auth == "" {
		return ""
	}
	// 标准 Bearer 前缀（大小写不敏感）。
	const prefix = "Bearer "
	// 命中 Bearer 前缀时，仅返回真实 token 部分。
	if len(auth) >= len(prefix) && strings.EqualFold(auth[:len(prefix)], prefix) {
		return strings.TrimSpace(auth[len(prefix):])
	}
	// 未使用 Bearer 前缀时，保留原值作为 token。
	return auth
}

// VerifyAccessToken 调 IAM 校验 token，并返回网关统一鉴权结果。
func (s *sAuth) VerifyAccessToken(ctx context.Context, accessToken string) (*service.VerifyResult, error) {
	// 空 token 直接返回 401，避免无效下游调用。
	if strings.TrimSpace(accessToken) == "" {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "missing access token")
	}
	// 确保 IAM gRPC 客户端就绪。
	if err := s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 调用 IAM 内部接口验证 access token 合法性。
	res, err := s.iamClient.VerifyAccessToken(ctx, &iamv1.VerifyAccessTokenReq{AccessToken: accessToken})
	if err != nil {
		return nil, gerror.Wrap(err, "verify access token failed")
	}
	// token 无效或无 user_id 均视为未授权。
	if !res.GetValid() || res.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "invalid access token")
	}
	// 将 IAM 返回结构转换为网关内部统一结果。
	out := &service.VerifyResult{
		UserID: res.GetUserId(),
		// 统一输出语义字符串，例如 ACTIVE / DISABLED。
		AccountStatusCode: enumCode(res.GetAccountStatus().String(), "ACCOUNT_STATUS_"),
	}
	// 逐条拷贝角色信息，供 BFF 侧做权限/归属判断。
	for _, role := range res.GetRoles() {
		// scope_no 对外使用字符串，避免直接暴露下游数字语义。
		scopeNo := ""
		if role.GetScopeId() > 0 {
			scopeNo = fmt.Sprintf("%d", role.GetScopeId())
		}
		out.Roles = append(out.Roles, service.AuthRole{
			// 将 IAM 枚举名裁剪前缀后输出，减少前端对“魔术数字”依赖。
			RoleCode:      enumCode(role.GetRoleCode().String(), "ROLE_CODE_"),
			ScopeTypeCode: enumCode(role.GetScopeType().String(), "SCOPE_TYPE_"),
			ScopeNo:       scopeNo,
			// ScopeID 保留给内部逻辑使用（如精确鉴权/审计）。
			ScopeID: role.GetScopeId(),
		})
	}
	return out, nil
}

// ensureClients 负责惰性初始化 IAM 客户端。
// 使用 sync.Once 的目的是避免高并发请求下重复创建连接。
func (s *sAuth) ensureClients(ctx context.Context) error {
	s.once.Do(func() {
		// 从配置读取 IAM gRPC 地址，默认回落本地开发地址。
		addr := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.iamGrpc", "127.0.0.1:9001").String())
		// 建连超时控制，避免请求线程被长期阻塞。
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 使用 insecure 是因为当前链路处于内网信任域。
		s.iamConn, s.initErr = grpc.DialContext(timeoutCtx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if s.initErr != nil {
			// 初始化失败时保留错误，后续调用直接返回该错误。
			return
		}
		// 连接成功后创建 IAM RPC 客户端。
		s.iamClient = iamv1.NewInternalServiceClient(s.iamConn)
	})
	return s.initErr
}

// enumCode 将枚举名从带前缀形式转换为纯业务码。
// 例如 ACCOUNT_STATUS_ACTIVE -> ACTIVE。
func enumCode(v, prefix string) string {
	if strings.HasPrefix(v, prefix) {
		return strings.TrimPrefix(v, prefix)
	}
	return v
}
