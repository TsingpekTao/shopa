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

// sAuth 负责网关统一鉴权：
// 1) 从 HTTP 请求提取 access token。
// 2) 调 IAM 校验 token 并拿到 user_id/角色信息。
// 3) 将下游枚举转换为网关统一可读字段，供 BFF/代理层复用。
type sAuth struct {
	// once 保证 gRPC 客户端只初始化一次，避免高并发重复建连。
	once sync.Once

	// iamConn 是到 IAM 的长连接。
	iamConn *grpc.ClientConn
	// iamClient 是 IAM 内部鉴权 RPC 客户端。
	iamClient iamv1.InternalServiceClient
	// initErr 记录首次初始化失败原因，后续请求直接返回同一错误。
	initErr error
}

// New 创建鉴权服务实现。
func New() *sAuth {
	return &sAuth{}
}

// init 启动时注册到 service 门面，供其他模块通过 service.Auth() 调用。
func init() {
	service.RegisterAuth(New())
}

// ExtractAccessToken 从请求头提取 token，兼容两种格式：
// 说明：1) Authorization: Bearer <token>
// 说明：2) Authorization: <token>
func (s *sAuth) ExtractAccessToken(r *ghttp.Request) string {
	// 请求对象为空时直接返回空 token，上层按未登录处理。
	if r == nil {
		return ""
	}
	// 去掉头尾空白，避免代理层空格差异影响解析。
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if auth == "" {
		return ""
	}

	// 标准 Bearer 前缀（大小写不敏感）。
	const prefix = "Bearer "
	if len(auth) >= len(prefix) && strings.EqualFold(auth[:len(prefix)], prefix) {
		// 命中 Bearer 时只返回 token 本体。
		return strings.TrimSpace(auth[len(prefix):])
	}
	// 未使用 Bearer 前缀时，兼容旧客户端，直接把整段当 token。
	return auth
}

// VerifyAccessToken 调 IAM 校验 token，并转换为网关统一鉴权结果。
func (s *sAuth) VerifyAccessToken(ctx context.Context, accessToken string) (*service.VerifyResult, error) {
	// 空 token 直接 401，避免无意义下游调用。
	if strings.TrimSpace(accessToken) == "" {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "missing access token")
	}
	// 确保 IAM 客户端就绪。
	if err := s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 向 IAM 发起 token 校验。
	res, err := s.iamClient.VerifyAccessToken(ctx, &iamv1.VerifyAccessTokenReq{AccessToken: accessToken})
	if err != nil {
		return nil, gerror.Wrap(err, "verify access token failed")
	}
	// valid=false 或 user_id=0 都视为未授权。
	if !res.GetValid() || res.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "invalid access token")
	}

	// 组装网关统一结果：保留 user_id，状态码转换为可读字符串。
	out := &service.VerifyResult{
		UserID:            res.GetUserId(),
		AccountStatusCode: enumCode(res.GetAccountStatus().String(), "ACCOUNT_STATUS_"),
	}
	permissionSet := make(map[string]struct{}, len(res.GetPermissions()))
	for _, permission := range res.GetPermissions() {
		permission = strings.TrimSpace(permission)
		if permission == "" {
			continue
		}
		if _, ok := permissionSet[permission]; ok {
			continue
		}
		permissionSet[permission] = struct{}{}
		out.Permissions = append(out.Permissions, permission)
	}

	// 逐条拷贝角色信息，BFF/代理层可基于角色做权限与归属判断。
	for _, role := range res.GetRoles() {
		// ScopeNo 对外使用字符串，减少前端对数值语义耦合。
		scopeNo := ""
		if role.GetScopeId() > 0 {
			scopeNo = fmt.Sprintf("%d", role.GetScopeId())
		}
		out.Roles = append(out.Roles, service.AuthRole{
			// 统一将 IAM 枚举裁剪前缀，输出业务码。
			RoleCode:      enumCode(role.GetRoleCode().String(), "ROLE_CODE_"),
			ScopeTypeCode: enumCode(role.GetScopeType().String(), "SCOPE_TYPE_"),
			ScopeNo:       scopeNo,
			// ScopeID 仍保留给网关内部逻辑（精细鉴权/审计）。
			ScopeID: role.GetScopeId(),
		})
	}
	return out, nil
}

// ensureClients 惰性初始化 IAM 客户端。
// 使用 sync.Once 避免高并发下重复创建连接。
func (s *sAuth) ensureClients(ctx context.Context) error {
	s.once.Do(func() {
		// 从配置读取 IAM gRPC 地址，未配置时回落本地默认值。
		addr := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.iamGrpc", "127.0.0.1:9001").String())

		// 建连超时保护，避免请求线程长期阻塞。
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 当前网关到 IAM 默认在内网链路，使用 insecure 传输。
		s.iamConn, s.initErr = grpc.DialContext(timeoutCtx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if s.initErr != nil {
			// 首次失败后保留错误，后续直接返回，避免反复拨号放大故障。
			return
		}
		s.iamClient = iamv1.NewInternalServiceClient(s.iamConn)
	})
	return s.initErr
}

// enumCode 将带前缀枚举名转换为业务码。
// 例如：ACCOUNT_STATUS_ACTIVE -> ACTIVE。
func enumCode(v, prefix string) string {
	if strings.HasPrefix(v, prefix) {
		return strings.TrimPrefix(v, prefix)
	}
	return v
}
