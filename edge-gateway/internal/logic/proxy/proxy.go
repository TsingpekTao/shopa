package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/TsingpekTao/shopa/edge-gateway/internal/consts"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/dao"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/model/do"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/model/entity"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/guid"
	"golang.org/x/sync/singleflight"
)

type sProxy struct {
	buyerProductsCache *gcache.Cache

	buyerProductsGroup singleflight.Group
}

func New() *sProxy {

	ctx := gctx.New()

	lruCap := g.Cfg().MustGet(ctx, "gateway.proxy.buyerProductsCacheLruCap", 512).Int()
	if lruCap <= 0 {

		lruCap = 512
	}
	return &sProxy{

		buyerProductsCache: gcache.New(lruCap),
	}
}

func init() {

	service.RegisterProxy(New())
}

func (s *sProxy) HandleProxyRequest(ctx context.Context, r *ghttp.Request) (bool, error) {

	if r == nil {
		return false, nil
	}

	method := strings.ToUpper(strings.TrimSpace(r.Method))
	if !isProxyMethod(method) {

		return false, nil
	}

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

	requestID := strings.TrimSpace(r.Header.Get("X-Request-Id"))
	if requestID == "" {

		requestID = strings.ReplaceAll(guid.S(), "-", "")
		r.Header.Set("X-Request-Id", requestID)
	}

	startAt := time.Now()

	params := extractPathParams(route.PathPattern, r.URL.Path)

	action := strings.TrimSpace(route.Action)
	if action == "" {
		action = strings.TrimSpace(route.RouteCode)
	}

	resourceID := pickResourceID(route.ResourceIdPathKey, params)

	var (
		userID uint64

		statusCode = http.StatusBadGateway

		errCode string

		permissionKey = strings.TrimSpace(route.RequiredPermissionKey)
	)

	if route.AuthRequired == 1 {

		verified, verifyErr := service.Auth().VerifyAccessToken(ctx, service.Auth().ExtractAccessToken(r))
		if verifyErr != nil {
			return true, verifyErr
		}

		userID = verified.UserID

		if permissionKey != "" && !hasPermission(verified.Permissions, permissionKey) {

			return true, gerror.NewCodef(consts.CodeForbidden, "permission denied: %s", permissionKey)
		}

		if route.InjectUserContext == 1 {

			r.Header.Set("X-User-Id", fmt.Sprintf("%d", verified.UserID))
			r.Header.Set("X-Account-Status-Code", verified.AccountStatusCode)
		}
	}

	targetBase, err := s.resolveUpstream(ctx, route.UpstreamService)
	if err != nil {
		return true, err
	}

	targetURL, err := parseUpstreamURL(targetBase)
	if err != nil {
		return true, err
	}

	upstreamPath := strings.TrimSpace(route.UpstreamPathTemplate)
	if upstreamPath == "" {
		upstreamPath = r.URL.Path
	}

	upstreamPath = applyPathParams(upstreamPath, params)

	cacheEnabled := s.isBuyerProductsCacheEnabled(ctx) && isBuyerProductsListPath(method, r.URL.Path)
	if cacheEnabled {

		cacheKey := buildBuyerProductsCacheKey(method, r.URL.Path, r.URL.RawQuery)
		if payload, ok := s.getBuyerProductsCache(ctx, cacheKey); ok {

			s.writeBufferedResponse(r, payload, "HIT")
			statusCode = payload.StatusCode
		} else {

			value, doErr, _ := s.buyerProductsGroup.Do(cacheKey, func() (interface{}, error) {

				reqClone := r.Request.Clone(ctx)
				return s.proxyToBuffer(targetURL, upstreamPath, requestID, reqClone), nil
			})
			if doErr != nil {
				return true, doErr
			}
			payload, _ := value.(*cachedProxyResponse)
			if payload == nil {

				payload = &cachedProxyResponse{
					StatusCode:  http.StatusBadGateway,
					ContentType: "text/plain; charset=utf-8",
					Body:        []byte("proxy response is empty"),
				}
				errCode = "UPSTREAM_PROXY_ERROR"
			}

			s.writeBufferedResponse(r, payload, "MISS")
			statusCode = payload.StatusCode

			if payload.StatusCode == http.StatusOK {

				ttlSeconds := s.getBuyerProductsCacheTTLSeconds(ctx)
				_ = s.buyerProductsCache.Set(ctx, cacheKey, payload, time.Duration(ttlSeconds)*time.Second)
			}
			if payload.StatusCode == http.StatusBadGateway {

				errCode = "UPSTREAM_PROXY_ERROR"
			}
		}
	} else {

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		defaultDirector := proxy.Director
		proxy.Director = func(req *http.Request) {

			defaultDirector(req)

			req.URL.Path = upstreamPath
			req.Host = targetURL.Host

			req.Header.Set("X-Request-Id", requestID)
		}
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, proxyErr error) {

			statusCode = http.StatusBadGateway
			errCode = "UPSTREAM_PROXY_ERROR"
			http.Error(rw, proxyErr.Error(), http.StatusBadGateway)
		}

		recorder := &statusRecorder{ResponseWriter: r.Response.Writer, statusCode: http.StatusOK}
		proxy.ServeHTTP(recorder, r.Request)

		statusCode = recorder.StatusCode()
	}

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
		PermissionKey:   permissionKey,
		Action:          action,
		ResourceId:      resourceID,
	})
	return true, nil
}

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

func (s *sProxy) resolveUpstream(ctx context.Context, upstreamService string) (string, error) {
	key := ""
	switch strings.ToLower(strings.TrimSpace(upstreamService)) {
	case "catalog":
		key = "upstream.catalogHttp"
	case "inventory":
		key = "upstream.inventoryHttp"
	case "order":
		key = "upstream.orderHttp"
	case "cart":
		key = "upstream.cartHttp"
	case "payment":
		key = "upstream.paymentHttp"
	case "risk":
		key = "upstream.riskHttp"
	case "notification":
		key = "upstream.notificationHttp"
	case "promotion":
		key = "upstream.promotionHttp"
	case "search":
		key = "upstream.searchHttp"
	case "aftersale", "after_sale":
		key = "upstream.aftersaleHttp"
	case "review":
		key = "upstream.reviewHttp"
	case "fulfillment":
		key = "upstream.fulfillmentHttp"
	case "chat":
		key = "upstream.chatHttp"
	case "agent":
		key = "upstream.agentHttp"
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

func (s *sProxy) writeAudit(ctx context.Context, data *do.EdgeRequestAudit) {
	if data == nil {

		return
	}
	if fields, ok := data.DegradedFieldsJson.([]string); ok {

		b, _ := json.Marshal(fields)
		data.DegradedFieldsJson = string(b)
	}

	_, _ = dao.EdgeRequestAudit.Ctx(ctx).Data(data).Insert()
}

func isProxyMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

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

func builtinProxyRoutes() []entity.EdgeProxyRoute {
	return []entity.EdgeProxyRoute{

		{RouteCode: "BUILTIN_IAM_SMS_SEND", Method: http.MethodPost, PathPattern: "/v1/auth/sms/send", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/sms/send", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_IAM_REGISTER_PASSWORD", Method: http.MethodPost, PathPattern: "/v1/auth/register/password", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/register/password", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_IAM_LOGIN_PASSWORD", Method: http.MethodPost, PathPattern: "/v1/auth/login/password", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/login/password", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_IAM_LOGIN_SMS", Method: http.MethodPost, PathPattern: "/v1/auth/login/sms", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/login/sms", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_IAM_RESET_PASSWORD_SMS", Method: http.MethodPost, PathPattern: "/v1/auth/password/reset/sms", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/password/reset/sms", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_IAM_REFRESH_TOKEN", Method: http.MethodPost, PathPattern: "/v1/auth/token/refresh", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/token/refresh", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_IAM_LOGOUT", Method: http.MethodPost, PathPattern: "/v1/auth/logout", UpstreamService: "iam", UpstreamPathTemplate: "/v1/auth/logout", AuthRequired: 1, InjectUserContext: 1},

		{RouteCode: "BUILTIN_PROFILE_GET", Method: http.MethodGet, PathPattern: "/v1/me/profile", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/profile", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_PROFILE_PATCH", Method: http.MethodPatch, PathPattern: "/v1/me/profile", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/profile", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ADDRESS_LIST", Method: http.MethodGet, PathPattern: "/v1/me/addresses", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/addresses", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ADDRESS_CREATE", Method: http.MethodPost, PathPattern: "/v1/me/addresses", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/addresses", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ADDRESS_PATCH", Method: http.MethodPatch, PathPattern: "/v1/me/addresses/{addressId}", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/addresses/{addressId}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ADDRESS_DELETE", Method: http.MethodDelete, PathPattern: "/v1/me/addresses/{addressId}", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/addresses/{addressId}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ADDRESS_SET_DEFAULT", Method: http.MethodPost, PathPattern: "/v1/me/addresses/default", UpstreamService: "user_profile", UpstreamPathTemplate: "/v1/me/addresses/default", AuthRequired: 1, InjectUserContext: 1},

		{RouteCode: "BUILTIN_SELLER_APPLICATION_LIST", Method: http.MethodGet, PathPattern: "/v1/seller/applications", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/applications", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_SELLER_APPLICATION_GET", Method: http.MethodGet, PathPattern: "/v1/seller/applications/{applicationNo}", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/applications/{applicationNo}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_SELLER_APPLICATION_CREATE_DRAFT", Method: http.MethodPost, PathPattern: "/v1/seller/applications/draft", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/applications/draft", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_SELLER_APPLICATION_UPDATE_DRAFT", Method: http.MethodPatch, PathPattern: "/v1/seller/applications/{applicationNo}/draft", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/applications/{applicationNo}/draft", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_SELLER_APPLICATION_SUBMIT", Method: http.MethodPost, PathPattern: "/v1/seller/applications/{applicationNo}/submit", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/applications/{applicationNo}/submit", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_SELLER_APPLICATION_RESUBMIT", Method: http.MethodPost, PathPattern: "/v1/seller/applications/{rejectedApplicationNo}/resubmit", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/applications/{rejectedApplicationNo}/resubmit", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_SELLER_STORE_CATEGORY_LIST", Method: http.MethodGet, PathPattern: "/v1/seller/shops/{shopNo}/store-categories", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/shops/{shopNo}/store-categories", AuthRequired: 1, InjectUserContext: 1, Action: "seller.shop.store_category.list", ResourceIdPathKey: "shopNo"},
		{RouteCode: "BUILTIN_SELLER_STORE_CATEGORY_CREATE", Method: http.MethodPost, PathPattern: "/v1/seller/shops/{shopNo}/store-categories", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/shops/{shopNo}/store-categories", AuthRequired: 1, InjectUserContext: 1, Action: "seller.shop.store_category.create", ResourceIdPathKey: "shopNo"},
		{RouteCode: "BUILTIN_SELLER_STORE_CATEGORY_SORT", Method: http.MethodPost, PathPattern: "/v1/seller/shops/{shopNo}/store-categories:sort", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/shops/{shopNo}/store-categories:sort", AuthRequired: 1, InjectUserContext: 1, Action: "seller.shop.store_category.sort", ResourceIdPathKey: "shopNo"},
		{RouteCode: "BUILTIN_SELLER_STORE_CATEGORY_UPDATE", Method: http.MethodPatch, PathPattern: "/v1/seller/shops/{shopNo}/store-categories/{categoryId}", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/shops/{shopNo}/store-categories/{categoryId}", AuthRequired: 1, InjectUserContext: 1, Action: "seller.shop.store_category.update", ResourceIdPathKey: "categoryId"},
		{RouteCode: "BUILTIN_SELLER_STORE_CATEGORY_DELETE", Method: http.MethodDelete, PathPattern: "/v1/seller/shops/{shopNo}/store-categories/{categoryId}", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/shops/{shopNo}/store-categories/{categoryId}", AuthRequired: 1, InjectUserContext: 1, Action: "seller.shop.store_category.delete", ResourceIdPathKey: "categoryId"},
		{RouteCode: "BUILTIN_SELLER_PRODUCT_STORE_CATEGORY_GET", Method: http.MethodGet, PathPattern: "/v1/seller/shops/{shopNo}/products/{spuNo}/store-category", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/shops/{shopNo}/products/{spuNo}/store-category", AuthRequired: 1, InjectUserContext: 1, Action: "seller.shop.product_store_category.get", ResourceIdPathKey: "spuNo"},
		{RouteCode: "BUILTIN_SELLER_PRODUCT_STORE_CATEGORY_UPDATE", Method: http.MethodPut, PathPattern: "/v1/seller/shops/{shopNo}/products/{spuNo}/store-category", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/shops/{shopNo}/products/{spuNo}/store-category", AuthRequired: 1, InjectUserContext: 1, Action: "seller.shop.product_store_category.update", ResourceIdPathKey: "spuNo"},
		{RouteCode: "BUILTIN_SELLER_PRODUCT_STORE_CATEGORY_BATCH_GET", Method: http.MethodPost, PathPattern: "/v1/seller/shops/{shopNo}/products/store-categories:batch-get", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/seller/shops/{shopNo}/products/store-categories:batch-get", AuthRequired: 1, InjectUserContext: 1, Action: "seller.shop.product_store_category.batch_get", ResourceIdPathKey: "shopNo"},

		{RouteCode: "BUILTIN_CATALOG_LIST_MY_PRODUCTS", Method: http.MethodGet, PathPattern: "/v1/catalog/seller/products", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_GET_MY_PRODUCT", Method: http.MethodGet, PathPattern: "/v1/catalog/seller/products/{spu_no}", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products/{spu_no}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_CREATE_DRAFT", Method: http.MethodPost, PathPattern: "/v1/catalog/seller/products/draft", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products/draft", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_UPDATE_DRAFT", Method: http.MethodPut, PathPattern: "/v1/catalog/seller/products/draft", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products/draft", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_DELETE_DRAFT", Method: http.MethodPost, PathPattern: "/v1/catalog/seller/products/draft:delete", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products/draft:delete", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_UPSERT_SKU", Method: http.MethodPost, PathPattern: "/v1/catalog/seller/products/skus:upsert", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products/skus:upsert", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_SUBMIT_REVIEW", Method: http.MethodPost, PathPattern: "/v1/catalog/seller/products/review:submit", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products/review:submit", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_RESUBMIT_REVIEW", Method: http.MethodPost, PathPattern: "/v1/catalog/seller/products/review:resubmit", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/seller/products/review:resubmit", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CATALOG_LIST_BUYER_PRODUCTS", Method: http.MethodGet, PathPattern: "/v1/catalog/buyer/products", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/buyer/products", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_CATALOG_LIST_BUYER_PRODUCT_IMAGES", Method: http.MethodGet, PathPattern: "/v1/catalog/buyer/product-images", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/buyer/product-images", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_CATALOG_SEARCH_BUYER_PRODUCTS", Method: http.MethodGet, PathPattern: "/v1/catalog/buyer/products/search", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/buyer/products/search", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_CATALOG_GET_BUYER_PRODUCT", Method: http.MethodGet, PathPattern: "/v1/catalog/buyer/products/{spu_no}", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/buyer/products/{spu_no}", AuthRequired: 0, InjectUserContext: 0},

		{RouteCode: "BUILTIN_INVENTORY_ADJUST", Method: http.MethodPost, PathPattern: "/v1/inventory/seller/stock:adjust", UpstreamService: "inventory", UpstreamPathTemplate: "/v1/inventory/seller/stock:adjust", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_INVENTORY_GET_SKU", Method: http.MethodGet, PathPattern: "/v1/inventory/sku/{sku_no}", UpstreamService: "inventory", UpstreamPathTemplate: "/v1/inventory/sku/{sku_no}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_INVENTORY_BATCH_GET_SKU", Method: http.MethodPost, PathPattern: "/v1/inventory/sku:batch-get", UpstreamService: "inventory", UpstreamPathTemplate: "/v1/inventory/sku:batch-get", AuthRequired: 1, InjectUserContext: 1},

		{RouteCode: "BUILTIN_CART_ADD_ITEM", Method: http.MethodPost, PathPattern: "/v1/cart/items:add", UpstreamService: "cart", UpstreamPathTemplate: "/v1/cart/items:add", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CART_UPDATE_QTY", Method: http.MethodPost, PathPattern: "/v1/cart/items:qty", UpstreamService: "cart", UpstreamPathTemplate: "/v1/cart/items:qty", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CART_TOGGLE_CHECKED", Method: http.MethodPost, PathPattern: "/v1/cart/items:check", UpstreamService: "cart", UpstreamPathTemplate: "/v1/cart/items:check", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CART_BATCH_CHECKED", Method: http.MethodPost, PathPattern: "/v1/cart/items:batch-check", UpstreamService: "cart", UpstreamPathTemplate: "/v1/cart/items:batch-check", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CART_REMOVE_ITEMS", Method: http.MethodPost, PathPattern: "/v1/cart/items:remove", UpstreamService: "cart", UpstreamPathTemplate: "/v1/cart/items:remove", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CART_CLEAR_INVALID", Method: http.MethodPost, PathPattern: "/v1/cart/items:clear-invalid", UpstreamService: "cart", UpstreamPathTemplate: "/v1/cart/items:clear-invalid", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CART_GET_MY_CART", Method: http.MethodGet, PathPattern: "/v1/cart/me", UpstreamService: "cart", UpstreamPathTemplate: "/v1/cart/me", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CART_PREPARE_CHECKOUT", Method: http.MethodPost, PathPattern: "/v1/cart/checkout:prepare", UpstreamService: "cart", UpstreamPathTemplate: "/v1/cart/checkout:prepare", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CART_INTERNAL_CONSUME_CHECKOUT", Method: http.MethodPost, PathPattern: "/v1/cart/internal/checkout:consume", UpstreamService: "cart", UpstreamPathTemplate: "/v1/cart/internal/checkout:consume", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CART_INTERNAL_MARK_ORDERED", Method: http.MethodPost, PathPattern: "/v1/cart/internal/items:ordered", UpstreamService: "cart", UpstreamPathTemplate: "/v1/cart/internal/items:ordered", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CART_INTERNAL_UPSERT_SKU_PROJECTION", Method: http.MethodPost, PathPattern: "/v1/cart/internal/sku-projection:batch-upsert", UpstreamService: "cart", UpstreamPathTemplate: "/v1/cart/internal/sku-projection:batch-upsert", AuthRequired: 1, InjectUserContext: 1},

		{RouteCode: "BUILTIN_ORDER_CREATE_FROM_CART", Method: http.MethodPost, PathPattern: "/v1/order/buyer/orders:create-from-cart", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/buyer/orders:create-from-cart", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_CREATE_BUY_NOW", Method: http.MethodPost, PathPattern: "/v1/order/buyer/orders:create-buy-now", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/buyer/orders:create-buy-now", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_REQUEST_PAY", Method: http.MethodPost, PathPattern: "/v1/order/buyer/orders:request-pay", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/buyer/orders:request-pay", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_CANCEL_MY_ORDER", Method: http.MethodPost, PathPattern: "/v1/order/buyer/orders:cancel", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/buyer/orders:cancel", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_CONFIRM_MY_RECEIPT", Method: http.MethodPost, PathPattern: "/v1/order/buyer/orders:confirm-received", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/buyer/orders:confirm-received", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_GET_MY_ORDER_DETAIL", Method: http.MethodGet, PathPattern: "/v1/order/buyer/orders/{order_no}", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/buyer/orders/{order_no}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_UPDATE_MY_ORDER_ADDRESS", Method: http.MethodPatch, PathPattern: "/v1/order/buyer/orders/{order_no}/address", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/buyer/orders/{order_no}/address", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_LIST_MY_ORDERS", Method: http.MethodGet, PathPattern: "/v1/order/buyer/orders", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/buyer/orders", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_LIST_SHOP_ORDERS", Method: http.MethodPost, PathPattern: "/v1/order/seller/list", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/seller/list", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_GET_SHOP_ORDER_DETAIL", Method: http.MethodGet, PathPattern: "/v1/order/seller/sub/{sub_order_no}", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/seller/sub/{sub_order_no}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_MARK_SUB_ORDER_SHIPPED", Method: http.MethodPost, PathPattern: "/v1/order/seller/sub/ship", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/seller/sub/ship", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_INTERNAL_CLOSE_UNPAID", Method: http.MethodPost, PathPattern: "/v1/order/internal/orders:close-if-unpaid", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/internal/orders:close-if-unpaid", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_INTERNAL_PAY_CALLBACK", Method: http.MethodPost, PathPattern: "/v1/order/internal/payments:callback", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/internal/payments:callback", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_ORDER_INTERNAL_SNAPSHOT", Method: http.MethodGet, PathPattern: "/v1/order/internal/orders/{order_no}/snapshot", UpstreamService: "order", UpstreamPathTemplate: "/v1/order/internal/orders/{order_no}/snapshot", AuthRequired: 1, InjectUserContext: 1},

		{RouteCode: "BUILTIN_AFTERSALE_BUYER_CREATE", Method: http.MethodPost, PathPattern: "/v1/aftersale/buyer/cases:create", UpstreamService: "aftersale", UpstreamPathTemplate: "/v1/aftersale/buyer/cases:create", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_AFTERSALE_BUYER_CANCEL", Method: http.MethodPost, PathPattern: "/v1/aftersale/buyer/cases:cancel", UpstreamService: "aftersale", UpstreamPathTemplate: "/v1/aftersale/buyer/cases:cancel", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_AFTERSALE_BUYER_DETAIL", Method: http.MethodGet, PathPattern: "/v1/aftersale/buyer/cases/{after_sale_no}", UpstreamService: "aftersale", UpstreamPathTemplate: "/v1/aftersale/buyer/cases/{after_sale_no}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_AFTERSALE_BUYER_LIST", Method: http.MethodGet, PathPattern: "/v1/aftersale/buyer/cases", UpstreamService: "aftersale", UpstreamPathTemplate: "/v1/aftersale/buyer/cases", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_AFTERSALE_SELLER_LIST", Method: http.MethodGet, PathPattern: "/v1/aftersale/seller/shops/{shop_no}/cases", UpstreamService: "aftersale", UpstreamPathTemplate: "/v1/aftersale/seller/shops/{shop_no}/cases", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_AFTERSALE_SELLER_DETAIL", Method: http.MethodGet, PathPattern: "/v1/aftersale/seller/cases/{after_sale_no}", UpstreamService: "aftersale", UpstreamPathTemplate: "/v1/aftersale/seller/cases/{after_sale_no}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_AFTERSALE_SELLER_APPROVE", Method: http.MethodPost, PathPattern: "/v1/aftersale/seller/cases:approve", UpstreamService: "aftersale", UpstreamPathTemplate: "/v1/aftersale/seller/cases:approve", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_AFTERSALE_SELLER_REJECT", Method: http.MethodPost, PathPattern: "/v1/aftersale/seller/cases:reject", UpstreamService: "aftersale", UpstreamPathTemplate: "/v1/aftersale/seller/cases:reject", AuthRequired: 1, InjectUserContext: 1},

		{RouteCode: "BUILTIN_REVIEW_BUYER_CREATE", Method: http.MethodPost, PathPattern: "/v1/review/buyer/reviews:create", UpstreamService: "review", UpstreamPathTemplate: "/v1/review/buyer/reviews:create", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_REVIEW_BUYER_APPEND", Method: http.MethodPost, PathPattern: "/v1/review/buyer/reviews:append", UpstreamService: "review", UpstreamPathTemplate: "/v1/review/buyer/reviews:append", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_REVIEW_BUYER_LIST", Method: http.MethodGet, PathPattern: "/v1/review/buyer/reviews", UpstreamService: "review", UpstreamPathTemplate: "/v1/review/buyer/reviews", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_REVIEW_SELLER_REPLY", Method: http.MethodPost, PathPattern: "/v1/review/seller/reviews:reply", UpstreamService: "review", UpstreamPathTemplate: "/v1/review/seller/reviews:reply", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_REVIEW_PUBLIC_LIST", Method: http.MethodGet, PathPattern: "/v1/review/public/spu/{spu_no}/reviews", UpstreamService: "review", UpstreamPathTemplate: "/v1/review/public/spu/{spu_no}/reviews", AuthRequired: 0, InjectUserContext: 0},
		{RouteCode: "BUILTIN_REVIEW_PUBLIC_SUMMARY", Method: http.MethodGet, PathPattern: "/v1/review/public/spu/{spu_no}/summary", UpstreamService: "review", UpstreamPathTemplate: "/v1/review/public/spu/{spu_no}/summary", AuthRequired: 0, InjectUserContext: 0},

		{RouteCode: "BUILTIN_FULFILLMENT_SELLER_CREATE", Method: http.MethodPost, PathPattern: "/v1/fulfillment/seller/shipments:create", UpstreamService: "fulfillment", UpstreamPathTemplate: "/v1/fulfillment/seller/shipments:create", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_FULFILLMENT_SELLER_SHIP", Method: http.MethodPost, PathPattern: "/v1/fulfillment/seller/shipments:ship", UpstreamService: "fulfillment", UpstreamPathTemplate: "/v1/fulfillment/seller/shipments:ship", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_FULFILLMENT_SELLER_LIST", Method: http.MethodGet, PathPattern: "/v1/fulfillment/seller/shops/{shop_no}/shipments", UpstreamService: "fulfillment", UpstreamPathTemplate: "/v1/fulfillment/seller/shops/{shop_no}/shipments", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_FULFILLMENT_SELLER_DETAIL", Method: http.MethodGet, PathPattern: "/v1/fulfillment/seller/shipments/{shipment_no}", UpstreamService: "fulfillment", UpstreamPathTemplate: "/v1/fulfillment/seller/shipments/{shipment_no}", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_FULFILLMENT_BUYER_LOGISTICS", Method: http.MethodGet, PathPattern: "/v1/fulfillment/buyer/orders/{order_no}/logistics", UpstreamService: "fulfillment", UpstreamPathTemplate: "/v1/fulfillment/buyer/orders/{order_no}/logistics", AuthRequired: 1, InjectUserContext: 1},

		{RouteCode: "BUILTIN_CHAT_BUYER_CREATE_CONVERSATION", Method: http.MethodPost, PathPattern: "/v1/chat/buyer/conversations:get-or-create", UpstreamService: "chat", UpstreamPathTemplate: "/v1/chat/buyer/conversations:get-or-create", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CHAT_BUYER_SEND", Method: http.MethodPost, PathPattern: "/v1/chat/buyer/messages:send", UpstreamService: "chat", UpstreamPathTemplate: "/v1/chat/buyer/messages:send", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CHAT_BUYER_CONVERSATIONS", Method: http.MethodGet, PathPattern: "/v1/chat/buyer/conversations", UpstreamService: "chat", UpstreamPathTemplate: "/v1/chat/buyer/conversations", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CHAT_BUYER_MESSAGES", Method: http.MethodGet, PathPattern: "/v1/chat/buyer/conversations/{conversation_no}/messages", UpstreamService: "chat", UpstreamPathTemplate: "/v1/chat/buyer/conversations/{conversation_no}/messages", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CHAT_BUYER_MARK_READ", Method: http.MethodPost, PathPattern: "/v1/chat/buyer/conversations:mark-read", UpstreamService: "chat", UpstreamPathTemplate: "/v1/chat/buyer/conversations:mark-read", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CHAT_BUYER_UNREAD_SUMMARY", Method: http.MethodGet, PathPattern: "/v1/chat/buyer/unread-summary", UpstreamService: "chat", UpstreamPathTemplate: "/v1/chat/buyer/unread-summary", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CHAT_SELLER_CONVERSATIONS", Method: http.MethodGet, PathPattern: "/v1/chat/seller/shops/{shop_no}/conversations", UpstreamService: "chat", UpstreamPathTemplate: "/v1/chat/seller/shops/{shop_no}/conversations", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CHAT_SELLER_MESSAGES", Method: http.MethodGet, PathPattern: "/v1/chat/seller/shops/{shop_no}/conversations/{conversation_no}/messages", UpstreamService: "chat", UpstreamPathTemplate: "/v1/chat/seller/shops/{shop_no}/conversations/{conversation_no}/messages", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CHAT_SELLER_SEND", Method: http.MethodPost, PathPattern: "/v1/chat/seller/messages:send", UpstreamService: "chat", UpstreamPathTemplate: "/v1/chat/seller/messages:send", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CHAT_SELLER_MARK_READ", Method: http.MethodPost, PathPattern: "/v1/chat/seller/conversations:mark-read", UpstreamService: "chat", UpstreamPathTemplate: "/v1/chat/seller/conversations:mark-read", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CHAT_INTERNAL_SYSTEM_NOTICE", Method: http.MethodPost, PathPattern: "/v1/chat/internal/system-notices:publish", UpstreamService: "chat", UpstreamPathTemplate: "/v1/chat/internal/system-notices:publish", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_CHAT_INTERNAL_SNAPSHOT", Method: http.MethodGet, PathPattern: "/v1/chat/internal/conversations/{conversation_no}/snapshot", UpstreamService: "chat", UpstreamPathTemplate: "/v1/chat/internal/conversations/{conversation_no}/snapshot", AuthRequired: 1, InjectUserContext: 1},

		{RouteCode: "BUILTIN_AGENT_BUYER_CREATE_CONVERSATION", Method: http.MethodPost, PathPattern: "/v1/agent/buyer/conversations:get-or-create", UpstreamService: "agent", UpstreamPathTemplate: "/v1/agent/buyer/conversations:get-or-create", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_AGENT_BUYER_SEND", Method: http.MethodPost, PathPattern: "/v1/agent/buyer/messages:send", UpstreamService: "agent", UpstreamPathTemplate: "/v1/agent/buyer/messages:send", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_AGENT_BUYER_RUN_STATUS", Method: http.MethodGet, PathPattern: "/v1/agent/buyer/conversations/{conversation_no}/run-status", UpstreamService: "agent", UpstreamPathTemplate: "/v1/agent/buyer/conversations/{conversation_no}/run-status", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_AGENT_BUYER_MESSAGES", Method: http.MethodGet, PathPattern: "/v1/agent/buyer/conversations/{conversation_no}/messages", UpstreamService: "agent", UpstreamPathTemplate: "/v1/agent/buyer/conversations/{conversation_no}/messages", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_AGENT_BUYER_ESCALATE", Method: http.MethodPost, PathPattern: "/v1/agent/buyer/conversations/{conversation_no}:escalate", UpstreamService: "agent", UpstreamPathTemplate: "/v1/agent/buyer/conversations/{conversation_no}:escalate", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_AGENT_BUYER_FEEDBACK", Method: http.MethodPost, PathPattern: "/v1/agent/buyer/feedback:submit", UpstreamService: "agent", UpstreamPathTemplate: "/v1/agent/buyer/feedback:submit", AuthRequired: 1, InjectUserContext: 1},

		{RouteCode: "BUILTIN_MEDIA_UPLOAD_INIT", Method: http.MethodPost, PathPattern: "/v1/media/upload/init", UpstreamService: "media", UpstreamPathTemplate: "/v1/media/upload/init", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_MEDIA_UPLOAD_COMPLETE", Method: http.MethodPost, PathPattern: "/v1/media/upload/complete", UpstreamService: "media", UpstreamPathTemplate: "/v1/media/upload/complete", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_MEDIA_ASSET_PROCESS_STATUS", Method: http.MethodGet, PathPattern: "/v1/media/assets/{assetId}/process-status", UpstreamService: "media", UpstreamPathTemplate: "/v1/media/assets/{assetId}/process-status", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_MEDIA_ASSET_READ_URL", Method: http.MethodGet, PathPattern: "/v1/media/assets/{assetId}/read-url", UpstreamService: "media", UpstreamPathTemplate: "/v1/media/assets/{assetId}/read-url", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_MEDIA_BIZ_ASSETS", Method: http.MethodGet, PathPattern: "/v1/media/biz-assets", UpstreamService: "media", UpstreamPathTemplate: "/v1/media/biz-assets", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_MEDIA_BINDING_REPLACE", Method: http.MethodPost, PathPattern: "/v1/media/bindings/replace", UpstreamService: "media", UpstreamPathTemplate: "/v1/media/bindings/replace", AuthRequired: 1, InjectUserContext: 1},
		{RouteCode: "BUILTIN_MEDIA_BINDING_BATCH_UNBIND", Method: http.MethodPost, PathPattern: "/v1/media/bindings/batch-unbind", UpstreamService: "media", UpstreamPathTemplate: "/v1/media/bindings/batch-unbind", AuthRequired: 1, InjectUserContext: 1},

		{RouteCode: "BUILTIN_ADMIN_SELLER_APPLICATIONS", Method: http.MethodGet, PathPattern: "/v1/admin/seller/applications", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/admin/seller/applications", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "merchant:review:view", Action: "merchant.review.list"},
		{RouteCode: "BUILTIN_ADMIN_SELLER_APPLICATION_DETAIL", Method: http.MethodGet, PathPattern: "/v1/admin/seller/applications/{applicationNo}", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/admin/seller/applications/{applicationNo}", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "merchant:review:view", Action: "merchant.review.detail", ResourceIdPathKey: "applicationNo"},
		{RouteCode: "BUILTIN_ADMIN_SELLER_APPLICATION_APPROVE", Method: http.MethodPost, PathPattern: "/v1/admin/seller/applications/{applicationNo}/approve", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/admin/seller/applications/{applicationNo}/approve", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "merchant:review:approve", Action: "merchant.review.approve", ResourceIdPathKey: "applicationNo"},
		{RouteCode: "BUILTIN_ADMIN_SELLER_APPLICATION_REJECT", Method: http.MethodPost, PathPattern: "/v1/admin/seller/applications/{applicationNo}/reject", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/admin/seller/applications/{applicationNo}/reject", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "merchant:review:reject", Action: "merchant.review.reject", ResourceIdPathKey: "applicationNo"},
		{RouteCode: "BUILTIN_ADMIN_SHOP_FREEZE", Method: http.MethodPost, PathPattern: "/v1/admin/seller/shops/{shopNo}/freeze", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/admin/seller/shops/{shopNo}/freeze", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "shop:manage:freeze", Action: "shop.manage.freeze", ResourceIdPathKey: "shopNo"},
		{RouteCode: "BUILTIN_ADMIN_SHOP_CLOSE", Method: http.MethodPost, PathPattern: "/v1/admin/seller/shops/{shopNo}/close", UpstreamService: "seller_shop", UpstreamPathTemplate: "/v1/admin/seller/shops/{shopNo}/close", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "shop:manage:close", Action: "shop.manage.close", ResourceIdPathKey: "shopNo"},

		{RouteCode: "BUILTIN_ADMIN_PRODUCT_REVIEW_TASKS", Method: http.MethodGet, PathPattern: "/v1/catalog/admin/review/tasks", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/admin/review/tasks", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "product:review:view", Action: "product.review.list"},
		{RouteCode: "BUILTIN_ADMIN_PRODUCT_REVIEW_DETAIL", Method: http.MethodGet, PathPattern: "/v1/catalog/admin/review/{spu_no}", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/admin/review/{spu_no}", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "product:review:view", Action: "product.review.detail", ResourceIdPathKey: "spu_no"},
		{RouteCode: "BUILTIN_ADMIN_PRODUCT_APPROVE", Method: http.MethodPost, PathPattern: "/v1/catalog/admin/review/approve", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/admin/review/approve", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "product:review:approve", Action: "product.review.approve"},
		{RouteCode: "BUILTIN_ADMIN_PRODUCT_REJECT", Method: http.MethodPost, PathPattern: "/v1/catalog/admin/review/reject", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/admin/review/reject", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "product:review:reject", Action: "product.review.reject"},
		{RouteCode: "BUILTIN_ADMIN_PRODUCT_FREEZE", Method: http.MethodPost, PathPattern: "/v1/catalog/admin/review/freeze", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/admin/review/freeze", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "product:review:freeze", Action: "product.review.freeze"},
		{RouteCode: "BUILTIN_ADMIN_PRODUCT_UNFREEZE", Method: http.MethodPost, PathPattern: "/v1/catalog/admin/review/unfreeze", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/admin/review/unfreeze", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "product:review:unfreeze", Action: "product.review.unfreeze"},
		{RouteCode: "BUILTIN_ADMIN_PRODUCT_FORCE_OFF_SHELF", Method: http.MethodPost, PathPattern: "/v1/catalog/admin/review/force-off-shelf", UpstreamService: "catalog", UpstreamPathTemplate: "/v1/catalog/admin/review/force-off-shelf", AuthRequired: 1, InjectUserContext: 1, RequiredPermissionKey: "product:review:force_off_shelf", Action: "product.review.force_off_shelf"},
	}
}

func isBuyerProductsListPath(method, path string) bool {
	return strings.EqualFold(strings.TrimSpace(method), http.MethodGet) &&
		strings.EqualFold(strings.TrimSpace(path), "/v1/catalog/buyer/products")
}

func (s *sProxy) isBuyerProductsCacheEnabled(ctx context.Context) bool {
	return g.Cfg().MustGet(ctx, "gateway.proxy.buyerProductsCacheEnabled", true).Bool()
}

func (s *sProxy) getBuyerProductsCacheTTLSeconds(ctx context.Context) int {
	ttl := g.Cfg().MustGet(ctx, "gateway.proxy.buyerProductsCacheTtlSeconds", 30).Int()
	if ttl <= 0 {
		return 30
	}
	return ttl
}

func buildBuyerProductsCacheKey(method, path, rawQuery string) string {
	return strings.Join([]string{
		strings.ToUpper(strings.TrimSpace(method)),
		strings.TrimSpace(path),
		canonicalQueryString(rawQuery),
	}, "|")
}

func canonicalQueryString(rawQuery string) string {
	values, err := url.ParseQuery(rawQuery)
	if err != nil || len(values) == 0 {
		return ""
	}

	keys := make([]string, 0, len(values))
	normalized := make(url.Values, len(values))
	for key, items := range values {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		copied := append([]string(nil), items...)
		sort.Strings(copied)
		normalized[key] = copied
		keys = append(keys, key)
	}
	sort.Strings(keys)

	canonical := make(url.Values, len(keys))
	for _, key := range keys {
		canonical[key] = normalized[key]
	}
	return canonical.Encode()
}

func (s *sProxy) getBuyerProductsCache(ctx context.Context, key string) (*cachedProxyResponse, bool) {
	if s.buyerProductsCache == nil || strings.TrimSpace(key) == "" {
		return nil, false
	}
	value, err := s.buyerProductsCache.Get(ctx, key)
	if err != nil || value == nil {
		return nil, false
	}
	payload, ok := value.Val().(*cachedProxyResponse)
	if !ok || payload == nil {
		return nil, false
	}
	return payload, true
}

func (s *sProxy) writeBufferedResponse(r *ghttp.Request, payload *cachedProxyResponse, cacheStatus string) {
	if r == nil || payload == nil {
		return
	}
	if payload.ContentType != "" {
		r.Response.Header().Set("Content-Type", payload.ContentType)
	}
	if strings.TrimSpace(cacheStatus) != "" {
		r.Response.Header().Set("X-Shopa-Edge-Cache", cacheStatus)
	}
	statusCode := payload.StatusCode
	if statusCode <= 0 {
		statusCode = http.StatusOK
	}
	r.Response.Writer.WriteHeader(statusCode)
	if len(payload.Body) > 0 {
		_, _ = r.Response.Writer.Write(payload.Body)
	}
}

func (s *sProxy) proxyToBuffer(
	targetURL *url.URL,
	upstreamPath string,
	requestID string,
	req *http.Request,
) *cachedProxyResponse {
	if targetURL == nil || req == nil {

		return &cachedProxyResponse{
			StatusCode:  http.StatusBadGateway,
			ContentType: "text/plain; charset=utf-8",
			Body:        []byte("invalid proxy request"),
		}
	}

	recorder := newBufferedProxyResponseRecorder()
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	defaultDirector := proxy.Director
	proxy.Director = func(proxyReq *http.Request) {

		defaultDirector(proxyReq)
		proxyReq.URL.Path = upstreamPath
		proxyReq.Host = targetURL.Host
		proxyReq.Header.Set("X-Request-Id", requestID)
	}
	proxy.ErrorHandler = func(rw http.ResponseWriter, proxyReq *http.Request, proxyErr error) {

		http.Error(rw, proxyErr.Error(), http.StatusBadGateway)
	}

	proxy.ServeHTTP(recorder, req)

	return recorder.ToCachedResponse()
}

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

func extractPathParams(pattern, actual string) map[string]string {
	params := map[string]string{}
	if !strings.Contains(pattern, "{") || !strings.Contains(pattern, "}") {
		return params
	}
	pSeg := splitPath(pattern)
	aSeg := splitPath(actual)
	if len(pSeg) != len(aSeg) {
		return params
	}
	for i := 0; i < len(pSeg); i++ {
		segment := pSeg[i]
		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			key := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(segment, "{"), "}"))
			if key != "" {
				params[key] = aSeg[i]
			}
		}
	}
	return params
}

func splitPath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return []string{}
	}
	return strings.Split(trimmed, "/")
}

func applyPathParams(template string, params map[string]string) string {
	if len(params) == 0 {
		return template
	}
	out := template
	for key, value := range params {
		out = strings.ReplaceAll(out, "{"+key+"}", value)
	}
	return out
}

func pickResourceID(resourceIDPathKey string, params map[string]string) string {
	resourceIDPathKey = strings.TrimSpace(resourceIDPathKey)
	if resourceIDPathKey != "" {
		if v := strings.TrimSpace(params[resourceIDPathKey]); v != "" {
			return v
		}
	}
	if len(params) == 0 {
		return ""
	}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if v := strings.TrimSpace(params[key]); v != "" {
			return v
		}
	}
	return ""
}

func hasPermission(granted []string, required string) bool {
	required = strings.TrimSpace(required)
	if required == "" {
		return true
	}
	for _, item := range granted {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if item == required || item == "*:*:*" {
			return true
		}
		if wildcardPermissionMatch(item, required) {
			return true
		}
	}
	return false
}

func wildcardPermissionMatch(pattern, target string) bool {
	if !strings.Contains(pattern, "*") {
		return false
	}
	patternParts := strings.Split(pattern, ":")
	targetParts := strings.Split(target, ":")
	if len(patternParts) != len(targetParts) {
		return false
	}
	for i := range patternParts {
		if patternParts[i] == "*" {
			continue
		}
		if patternParts[i] != targetParts[i] {
			return false
		}
	}
	return true
}

type cachedProxyResponse struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

type bufferedProxyResponseRecorder struct {
	header     http.Header
	body       bytes.Buffer
	statusCode int
}

func newBufferedProxyResponseRecorder() *bufferedProxyResponseRecorder {
	return &bufferedProxyResponseRecorder{
		header:     make(http.Header),
		statusCode: http.StatusOK,
	}
}

func (r *bufferedProxyResponseRecorder) Header() http.Header {

	return r.header
}

func (r *bufferedProxyResponseRecorder) WriteHeader(code int) {

	r.statusCode = code
}

func (r *bufferedProxyResponseRecorder) Write(data []byte) (int, error) {
	if r.statusCode <= 0 {

		r.statusCode = http.StatusOK
	}

	return r.body.Write(data)
}

func (r *bufferedProxyResponseRecorder) ToCachedResponse() *cachedProxyResponse {

	bodyBytes := append([]byte(nil), r.body.Bytes()...)
	return &cachedProxyResponse{
		StatusCode:  r.statusCode,
		ContentType: strings.TrimSpace(r.header.Get("Content-Type")),
		Body:        bodyBytes,
	}
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {

	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) StatusCode() int {
	if r.statusCode <= 0 {

		return http.StatusOK
	}
	return r.statusCode
}
