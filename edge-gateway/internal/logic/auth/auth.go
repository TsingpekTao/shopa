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

// sAuth 负责网关统一鉴权。
// 1) 从 HTTP 请求里提取 access token。
// 2) 调 IAM 校验 token 并拿到 user_id、权限和角色信息。
// 3) 把下游枚举结果转换成网关统一可读字段，供 BFF 和代理层复用。
type sAuth struct {
	// once 保证 gRPC 客户端只初始化一次，避免高并发下重复建连。
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

// init 在启动时注册到 service 门面，供其他模块通过 service.Auth() 调用。
func init() {
	service.RegisterAuth(New())
}

// ExtractAccessToken 从请求头提取 token，兼容两种格式：
// 1) Authorization: Bearer <token>
// 2) Authorization: <token>
func (s *sAuth) ExtractAccessToken(r *ghttp.Request) string {
	// 请求对象为空时直接返回空 token，上层按未登录处理。
	if r == nil {
		return ""
	}
	// 去掉首尾空白，避免代理层空格差异影响解析。
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if auth == "" {
		return ""
	}

	// 标准 Bearer 前缀，大小写不敏感。
	const prefix = "Bearer "
	if len(auth) >= len(prefix) && strings.EqualFold(auth[:len(prefix)], prefix) {
		// 命中 Bearer 前缀时只返回真正的 token 本体。
		return strings.TrimSpace(auth[len(prefix):])
	}
	// 未使用 Bearer 前缀时兼容旧客户端，直接把整段内容视为 token。
	return auth
}

// VerifyAccessToken 调 IAM 校验 token，并转换成网关统一鉴权结果。
func (s *sAuth) VerifyAccessToken(ctx context.Context, accessToken string) (*service.VerifyResult, error) {
	// 空 token 直接按 401 处理，避免无意义地下游调用。
	if strings.TrimSpace(accessToken) == "" {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "missing access token")
	}
	// 先确保 IAM 客户端已经初始化完成。
	if err := s.ensureClients(ctx); err != nil {
		return nil, err
	}

	// 向 IAM 发起 token 校验。
	res, err := s.iamClient.VerifyAccessToken(ctx, &iamv1.VerifyAccessTokenReq{AccessToken: accessToken})
	if err != nil {
		return nil, gerror.Wrap(err, "verify access token failed")
	}
	// valid=false 或 user_id=0 都说明当前 token 不合法。
	if !res.GetValid() || res.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "invalid access token")
	}

	// 组装网关统一鉴权结果，保留 user_id 并把账户状态转成业务码。
	out := &service.VerifyResult{
		UserID:            res.GetUserId(),
		AccountStatusCode: enumCode(res.GetAccountStatus().String(), "ACCOUNT_STATUS_"),
	}
	// permissionSet 用来给权限列表去重，避免下游重复权限污染判断。
	permissionSet := make(map[string]struct{}, len(res.GetPermissions()))
	// 逐条收集 IAM 返回的权限集。
	for _, permission := range res.GetPermissions() {
		// 先清理空白，防止无效权限键混入结果。
		permission = strings.TrimSpace(permission)
		if permission == "" {
			continue
		}
		// 已经收录过的权限直接跳过，保证结果集唯一。
		if _, ok := permissionSet[permission]; ok {
			continue
		}
		// 把权限记进去重集合，供后续重复判断。
		permissionSet[permission] = struct{}{}
		// 把去重后的权限写回输出结果，供代理和 BFF 做细粒度授权判断。
		out.Permissions = append(out.Permissions, permission)
	}

	// 逐条拷贝角色信息，BFF 和代理层可基于角色做权限与归属判断。
	for _, role := range res.GetRoles() {
		// ScopeNo 对外使用字符串，减少前端和网关对数值语义的耦合。
		scopeNo := ""
		if role.GetScopeId() > 0 {
			// 有范围 ID 时转成字符串透出，保持响应结构统一。
			scopeNo = fmt.Sprintf("%d", role.GetScopeId())
		}
		// 把单个角色转换成网关统一角色结构。
		out.Roles = append(out.Roles, service.AuthRole{
			// 把 IAM 枚举名裁掉固定前缀，输出业务层可读码值。
			RoleCode:      enumCode(role.GetRoleCode().String(), "ROLE_CODE_"),
			ScopeTypeCode: enumCode(role.GetScopeType().String(), "SCOPE_TYPE_"),
			ScopeNo:       scopeNo,
			// ScopeID 继续保留给网关内部逻辑，例如更精细的鉴权或审计。
			ScopeID: role.GetScopeId(),
		})
	}
	return out, nil
}

// ensureClients 惰性初始化 IAM 客户端。
// 这里使用 sync.Once，避免高并发下重复创建连接放大故障。
func (s *sAuth) ensureClients(ctx context.Context) error {
	s.once.Do(func() {
		// 从配置读取 IAM gRPC 地址，未配置时回落到本地默认值。
		addr := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.iamGrpc", "127.0.0.1:9001").String())

		// 给建连动作套一个独立超时，避免初始化过程长期阻塞请求线程。
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 当前网关到 IAM 默认走内网链路，所以这里使用 insecure 传输。
		s.iamConn, s.initErr = grpc.DialContext(timeoutCtx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if s.initErr != nil {
			// 首次初始化失败后缓存错误，后续直接返回，避免反复拨号放大故障。
			return
		}
		// 建连成功后创建 IAM RPC 客户端，供后续校验调用复用。
		s.iamClient = iamv1.NewInternalServiceClient(s.iamConn)
	})
	return s.initErr
}

// enumCode 把带前缀的枚举名转换成业务码。
// 例如 ACCOUNT_STATUS_ACTIVE 会被裁剪成 ACTIVE。
func enumCode(v, prefix string) string {
	if strings.HasPrefix(v, prefix) {
		return strings.TrimPrefix(v, prefix)
	}
	return v
}
