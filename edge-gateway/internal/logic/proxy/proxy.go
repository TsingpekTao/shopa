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

// sProxy 负责网关写请求透传：
// - 命中白名单路由后，执行鉴权、注入上下文、反向代理；
// - 请求结束后异步记录审计日志；
// - 对大 Body 保持流式转发，不读入内存。
type sProxy struct{}

// New 创建透传服务实例。
func New() *sProxy {
	return &sProxy{}
}

// init 注册透传服务到 service 门面。
func init() {
	service.RegisterProxy(New())
}

// HandleProxyRequest 处理写请求透传。
// 返回值 handled=true 表示请求已由 proxy 接管（成功或失败都不应再继续后续 handler）。
func (s *sProxy) HandleProxyRequest(ctx context.Context, r *ghttp.Request) (bool, error) {
	// 空请求直接忽略。
	if r == nil {
		return false, nil
	}
	// 统一标准化 Method，避免大小写差异。
	method := strings.ToUpper(strings.TrimSpace(r.Method))
	// 仅代理“写操作”方法，读操作走 BFF 聚合链路。
	if !isMutatingMethod(method) {
		return false, nil
	}

	// 路由白名单匹配，未命中则放行给其他 handler。
	route, err := s.matchRoute(ctx, method, r.URL.Path)
	if err != nil {
		return false, err
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

	// 记录请求开始时间，用于计算网关处理耗时。
	startAt := time.Now()
	var (
		// userID 用于审计落库。
		userID uint64
		// statusCode 默认给 502，若代理正常返回会被 recorder 覆盖。
		statusCode = http.StatusBadGateway
		// errCode 记录代理错误码，用于运维检索。
		errCode string
	)

	// 根据路由配置决定是否需要鉴权。
	if route.AuthRequired == 1 {
		verified, verifyErr := service.Auth().VerifyAccessToken(ctx, service.Auth().ExtractAccessToken(r))
		if verifyErr != nil {
			return true, verifyErr
		}
		userID = verified.UserID
		// 若要求注入用户上下文，则向下游传递身份相关 Header。
		if route.InjectUserContext == 1 {
			r.Header.Set("X-User-Id", fmt.Sprintf("%d", verified.UserID))
			r.Header.Set("X-Account-Status-Code", verified.AccountStatusCode)
		}
	}

	// 根据 route 的 upstream_service 解析真实上游地址。
	targetBase, err := s.resolveUpstream(ctx, route.UpstreamService)
	if err != nil {
		return true, err
	}
	// 规范化为可用于 ReverseProxy 的 URL。
	targetURL, err := parseUpstreamURL(targetBase)
	if err != nil {
		return true, err
	}

	// 支持配置化改写目标路径；未配置则保持原路径。
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
		// 按路由策略覆盖路径与 Host。
		req.URL.Path = upstreamPath
		req.Host = targetURL.Host
		// 强制透传 request_id，便于下游与网关日志对齐。
		req.Header.Set("X-Request-Id", requestID)
	}
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, proxyErr error) {
		// 代理失败时输出 502 并记录错误码。
		statusCode = http.StatusBadGateway
		errCode = "UPSTREAM_PROXY_ERROR"
		http.Error(rw, proxyErr.Error(), http.StatusBadGateway)
	}

	// 使用状态记录器捕获最终 HTTP 状态码。
	recorder := &statusRecorder{ResponseWriter: r.Response.Writer, statusCode: http.StatusOK}
	// 大 Body 场景必须严格流式透传：
	// 这里不读取 body，不做 io.ReadAll，避免把大文件吞进网关内存。
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

// matchRoute 按 method + path 从白名单路由表匹配可用路由。
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

// resolveUpstream 将逻辑服务名映射到具体配置键。
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
// 审计策略：只记录元数据，不保存请求 body，避免敏感信息泄露。
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

// isMutatingMethod 判断是否属于写操作 HTTP 方法。
func isMutatingMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

// parseUpstreamURL 解析上游地址，缺失协议时默认补 http://。
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
// 1) 完全匹配；
// 2) 路由参数匹配（如 /v1/shops/{shop_no}）；
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

// WriteHeader 在写响应头时记录状态码。
func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// StatusCode 返回最终状态码，默认回退 200。
func (r *statusRecorder) StatusCode() int {
	if r.statusCode <= 0 {
		return http.StatusOK
	}
	return r.statusCode
}
