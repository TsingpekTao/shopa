package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/TsingpekTao/shopa/edge-gateway/internal/dao"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/model/do"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/model/entity"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
)

// sProxy 负责网关请求透传：
// - 命中白名单路由后执行鉴权、上下文注入、反向代理。
// - 请求结束后异步写审计日志。
// - 对大 body 采用流式转发，不读入网关内存。
type sProxy struct{}

// New 创建透传服务实例。
func New() *sProxy {
	return &sProxy{}
}

// init 启动时注册透传实现到 service 门面。
func init() {
	service.RegisterProxy(New())
}

// HandleProxyRequest 处理代理请求。
// 返回 handled=true 表示该请求已被 proxy 接管（无论成功或失败都不应再走后续 handler）。
func (s *sProxy) HandleProxyRequest(ctx context.Context, r *ghttp.Request) (bool, error) {
	// 空请求直接忽略。
	if r == nil {
		return false, nil
	}

	// 统一规范化 Method，避免大小写差异。
	method := strings.ToUpper(strings.TrimSpace(r.Method))
	// 只代理白名单方法；其他方法继续走常规链路。
	if !isProxyMethod(method) {
		return false, nil
	}

	// 先查数据库路由；数据库不可用时回退到内置路由，保证核心链路可用。
	route, err := s.matchRoute(ctx, method, r.URL.Path)
	if err != nil {
		g.Log().Warningf(ctx, "[edge-gateway] query proxy route failed, fallback builtin route, err=%+v", err)
	}
	if route == nil {
		route = s.matchBuiltinRoute(method, r.URL.Path)
	}
	if route == nil {
		return false, nil
	}

	// 读取或生成 request_id，保证全链路可追踪。
	requestID := strings.TrimSpace(r.Header.Get("X-Request-Id"))
	if requestID == "" {
		requestID = strings.ReplaceAll(guid.S(), "-", "")
		r.Header.Set("X-Request-Id", requestID)
	}

	// 记录起始时间用于审计耗时统计。
	startAt := time.Now()
	var (
		// userID 用于审计落库；未鉴权路由保持 0。
		userID uint64
		// statusCode 默认 502，若代理正常返回会被 recorder 覆盖。
		statusCode = http.StatusBadGateway
		// errCode 记录代理错误码，方便运维检索。
		errCode string
	)

	// 按路由配置决定是否要求鉴权。
	if route.AuthRequired == 1 {
		verified, verifyErr := service.Auth().VerifyAccessToken(ctx, service.Auth().ExtractAccessToken(r))
		if verifyErr != nil {
			return true, verifyErr
		}
		userID = verified.UserID

		// 需要注入用户上下文时，透传身份 Header 给下游。
		if route.InjectUserContext == 1 {
			r.Header.Set("X-User-Id", fmt.Sprintf("%d", verified.UserID))
			r.Header.Set("X-Account-Status-Code", verified.AccountStatusCode)
		}
	}

	// 解析上游服务地址。
	targetBase, err := s.resolveUpstream(ctx, route.UpstreamService)
	if err != nil {
		return true, err
	}
	// 规范化为 ReverseProxy 可用 URL。
	targetURL, err := parseUpstreamURL(targetBase)
	if err != nil {
		return true, err
	}

	// 支持配置化改写目标路径；未配置则沿用原始路径。
	upstreamPath := strings.TrimSpace(route.UpstreamPathTemplate)
	if upstreamPath == "" {
		upstreamPath = r.URL.Path
	}

	// 创建标准单上游反向代理。
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	defaultDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		// 先执行默认 Director，保留标准转发行为。
		defaultDirector(req)
		// 再按路由配置覆盖路径与 Host。
		req.URL.Path = upstreamPath
		req.Host = targetURL.Host
		// 强制透传 request_id，便于上下游日志对齐。
		req.Header.Set("X-Request-Id", requestID)
	}
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, proxyErr error) {
		// 代理失败时写 502，并打上统一错误码供审计。
		statusCode = http.StatusBadGateway
		errCode = "UPSTREAM_PROXY_ERROR"
		http.Error(rw, proxyErr.Error(), http.StatusBadGateway)
	}

	// 包装响应写入器，捕获最终 HTTP 状态码。
	recorder := &statusRecorder{ResponseWriter: r.Response.Writer, statusCode: http.StatusOK}
	// 大 body 场景必须流式透传：不读取 body，不做 io.ReadAll。
	proxy.ServeHTTP(recorder, r.Request)
	statusCode = recorder.StatusCode()

	// 异步写审计，避免阻塞主请求返回。
	go s.writeAudit(gctx.New(), &do.EdgeRequestAudit{
		RequestId:       requestID,
		RouteCode:       route.RouteCode,
		UserId:          userID,
		Method:          method,
		Path:            r.URL.Path,
		UpstreamService: route.UpstreamService,
		StatusCode:      statusCode,
		LatencyMs:       uint(time.Since(startAt).Milliseconds()),
		Partial:         0,
		ErrorCode:       errCode,
		ClientIp:        r.GetClientIp(),
		UserAgent:       r.Header.Get("User-Agent"),
	})
	return true, nil
}

// matchRoute 按 method+path 从数据库路由白名单中匹配可用路由。
func (s *sProxy) matchRoute(ctx context.Context, method, path string) (*entity.EdgeProxyRoute, error) {
	var routes []*entity.EdgeProxyRoute
	err := dao.EdgeProxyRoute.Ctx(ctx).
		Where(dao.EdgeProxyRoute.Columns().Method, method).
		Where(dao.EdgeProxyRoute.Columns().Status, 1).
		WhereNull(dao.EdgeProxyRoute.Columns().DeletedAt).
		OrderAsc(dao.EdgeProxyRoute.Columns().Id).
		Scan(&routes)
	if err != nil {
		return nil, gerror.Wrap(err, "query proxy routes failed")
	}

	// 按 ID 升序遍历，保证匹配顺序稳定可预期。
	for _, route := range routes {
		if route == nil {
			continue
		}
		if matchPathPattern(route.PathPattern, path) {
			return route, nil
		}
	}
	return nil, nil
}

// resolveUpstream 将逻辑服务名映射到具体配置项。
func (s *sProxy) resolveUpstream(ctx context.Context, upstreamService string) (string, error) {
	key := ""
	switch strings.ToLower(strings.TrimSpace(upstreamService)) {
	case "catalog":
		key = "upstream.catalogHttp"
	case "inventory":
		key = "upstream.inventoryHttp"
	case "media":
		key = "upstream.mediaHttp"
	case "iam":
		key = "upstream.iamHttp"
	case "user_profile", "userprofile":
		key = "upstream.userProfileHttp"
	case "seller_shop", "seller-shop", "sellershop":
		key = "upstream.sellerShopHttp"
	default:
		return "", gerror.NewCodef(gcode.CodeInvalidParameter, "unknown upstream service: %s", upstreamService)
	}

	value := strings.TrimSpace(g.Cfg().MustGet(ctx, key).String())
	if value == "" {
		return "", gerror.NewCodef(gcode.CodeInvalidParameter, "missing upstream config: %s", key)
	}
	return value, nil
}

// writeAudit 写入网关请求审计。
// 审计策略：仅记元数据，不保存请求 body，避免敏感信息泄露。
func (s *sProxy) writeAudit(ctx context.Context, data *do.EdgeRequestAudit) {
	if data == nil {
		return
	}
	// 兼容 degraded_fields 以 []string 传入时的序列化。
	if fields, ok := data.DegradedFieldsJson.([]string); ok {
		b, _ := json.Marshal(fields)
		data.DegradedFieldsJson = string(b)
	}
	_, _ = dao.EdgeRequestAudit.Ctx(ctx).Data(data).Insert()
}

// isProxyMethod 判断请求方法是否纳入透传处理。
func isProxyMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

// matchBuiltinRoute 数据库路由缺失或查询失败时，回退到内置路由。
func (s *sProxy) matchBuiltinRoute(method, path string) *entity.EdgeProxyRoute {
	for _, route := range builtinProxyRoutes() {
		if !strings.EqualFold(route.Method, method) {
			continue
		}
		if matchPathPattern(route.PathPattern, path) {
			route := route
			return &route
		}
	}
	return nil
}

// builtinProxyRoutes 定义网关核心路由的最小可用集合。
func builtinProxyRoutes() []entity.EdgeProxyRoute {
	return []entity.EdgeProxyRoute{
		{RouteCode: "BUILTIN_IAM_SMS_SEND", Method: http.MethodPost, PathPattern: "/v1/auth/sms/send", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/sms/send", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_IAM_REGISTER_PASSWORD", Method: http.MethodPost, PathPattern: "/v1/auth/register/password", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/register/password", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_IAM_LOGIN_PASSWORD", Method: http.MethodPost, PathPattern: "/v1/auth/login/password", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/login/password", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_IAM_LOGIN_SMS", Method: http.MethodPost, PathPattern: "/v1/auth/login/sms", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/login/sms", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_IAM_REFRESH_TOKEN", Method: http.MethodPost, PathPattern: "/v1/auth/token/refresh", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/token/refresh", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_IAM_LOGOUT", Method: http.MethodPost, PathPattern: "/v1/auth/logout", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/logout", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_PROFILE_GET", Method: http.MethodGet, PathPattern: "/v1/me/profile", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/profile", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_PROFILE_PATCH", Method: http.MethodPatch, PathPattern: "/v1/me/profile", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/profile", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ADDRESS_LIST", Method: http.MethodGet, PathPattern: "/v1/me/addresses", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/addresses", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ADDRESS_CREATE", Method: http.MethodPost, PathPattern: "/v1/me/addresses", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/addresses", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ADDRESS_PATCH", Method: http.MethodPatch, PathPattern: "/v1/me/addresses/{addressId}", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/addresses/{addressId}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ADDRESS_DELETE", Method: http.MethodDelete, PathPattern: "/v1/me/addresses/{addressId}", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/addresses/{addressId}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ADDRESS_SET_DEFAULT", Method: http.MethodPost, PathPattern: "/v1/me/addresses/default", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/addresses/default", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_LIST_MY_PRODUCTS", Method: http.MethodGet, PathPattern: "/v1/catalog/seller/products", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_GET_MY_PRODUCT", Method: http.MethodGet, PathPattern: "/v1/catalog/seller/products/{spu_no}", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products/{spu_no}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_CREATE_DRAFT", Method: http.MethodPost, PathPattern: "/v1/catalog/seller/products/draft", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products/draft", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_UPDATE_DRAFT", Method: http.MethodPut, PathPattern: "/v1/catalog/seller/products/draft", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products/draft", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_UPSERT_SKU", Method: http.MethodPost, PathPattern: "/v1/catalog/seller/products/skus:upsert", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products/skus:upsert", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_INVENTORY_ADJUST", Method: http.MethodPost, PathPattern: "/v1/inventory/seller/stock:adjust", UpstreamService: "inventory", UpstreamPathTemplate: "/v1/inventory/seller/stock:adjust", AuthRequired: 1, InjectUserContext: 1},
	}
}

// parseUpstreamURL 解析上游地址；缺失协议时默认补 http://。
func parseUpstreamURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "upstream url is empty")
	}
	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, gerror.Wrap(err, "parse upstream url failed")
	}
	return u, nil
}

// matchPathPattern 支持三种匹配：
// 1) 完全匹配。
// 2) 路由参数匹配（如 /v1/shops/{shop_no}）。
// 3) 通配符匹配（*）。
func matchPathPattern(pattern, actual string) bool {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return false
	}
	if pattern == actual {
		return true
	}
	if strings.Contains(pattern, "{") && strings.Contains(pattern, "}") {
		re := regexp.MustCompile(`\{[^/]+\}`)
		rePattern := "^" + re.ReplaceAllString(pattern, `[^/]+`) + "$"
		matched, _ := regexp.MatchString(rePattern, actual)
		return matched
	}
	if strings.Contains(pattern, "*") {
		rePattern := "^" + regexp.QuoteMeta(pattern) + "$"
		rePattern = strings.ReplaceAll(rePattern, `\*`, `.*`)
		matched, _ := regexp.MatchString(rePattern, actual)
		return matched
	}
	return false
}

// statusRecorder 用于捕获反向代理返回的最终状态码。
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader 写响应头时记录状态码。
func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// StatusCode 返回最终状态码；未显式写入时回退 200。
func (r *statusRecorder) StatusCode() int {
	if r.statusCode <= 0 {
		return http.StatusOK
	}
	return r.statusCode
}
