package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	agentruntime "github.com/TsingpekTao/shopa/agent-svc/internal/logic/agent/runtime"
)

type buyerOrderSnapshotHTTPRepository struct {
	baseURL string
	client  *http.Client
}

type buyerOrderDetailEnvelope struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

func newBuyerOrderSnapshotHTTPRepository(baseURL string, client *http.Client) *buyerOrderSnapshotHTTPRepository {
	// 统一清洗基础地址，避免拼接路径时出现重复斜杠。
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	// 如果调用方没有传 http client，就给一个带超时的默认 client。
	if client == nil {
		client = &http.Client{Timeout: 3 * time.Second}
	}
	// 返回新的买家订单只读仓储，供 runtime 安全查询调用。
	return &buyerOrderSnapshotHTTPRepository{
		baseURL: baseURL,
		client:  client,
	}
}

func (r *buyerOrderSnapshotHTTPRepository) QueryOrderSnapshot(ctx context.Context, filter agentruntime.OrderOwnershipFilter) (*agentruntime.OrderSnapshot, error) {
	// 如果仓储未正确初始化基础地址，就直接按下游不可用返回。
	if r == nil || strings.TrimSpace(r.baseURL) == "" {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}
	// 如果关键过滤条件缺失，就直接按上下文不足返回。
	if filter.UserID == 0 || strings.TrimSpace(filter.OrderNo) == "" {
		return nil, &agentruntime.OrderLookupError{Code: "INSUFFICIENT_CONTEXT"}
	}
	// 用买家侧订单详情接口拼接强归属查询地址，避免落到无鉴权的内部接口。
	endpoint := fmt.Sprintf("%s/v1/order/buyer/orders/%s", r.baseURL, url.PathEscape(strings.TrimSpace(filter.OrderNo)))
	// 创建带上下文的 GET 请求，保证取消和超时可以透传。
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	// 如果请求构造失败，就按下游不可用统一映射。
	if err != nil {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}
	// 把安全身份放到请求头里，让 order-svc 自己做 buyer 归属校验。
	req.Header.Set("X-User-Id", strconv.FormatUint(filter.UserID, 10))
	// 把请求编号透传给下游，便于问题追踪。
	if strings.TrimSpace(filter.RequestID) != "" {
		req.Header.Set("X-Request-Id", strings.TrimSpace(filter.RequestID))
	}
	// 再带上一个固定服务名，便于下游识别调用来源。
	req.Header.Set("X-Service-Name", "agent-svc")
	// 执行买家侧订单详情查询请求。
	resp, err := r.client.Do(req)
	// 网络错误统一折叠成下游不可用。
	if err != nil {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}
	// 确保响应体被及时关闭，避免连接泄漏。
	defer resp.Body.Close()
	// 401/403 统一映射成权限失败，不暴露更多细节。
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, &agentruntime.OrderLookupError{Code: "PERMISSION_DENIED"}
	}
	// 404 统一映射成安全 not found。
	if resp.StatusCode == http.StatusNotFound {
		return nil, &agentruntime.OrderLookupError{Code: "NOT_FOUND"}
	}
	// 其他 5xx 或异常状态都按下游不可用处理。
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}
	// 解析 GoFrame HTTP 响应包裹结构，读取其中的 data.order。
	var envelope buyerOrderDetailEnvelope
	if err = json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}
	// 如果接口显式返回失败 code，也统一映射成安全语义。
	if envelope.Code != 0 && envelope.Code != 200 {
		if strings.Contains(strings.ToLower(envelope.Message), "not found") {
			return nil, &agentruntime.OrderLookupError{Code: "NOT_FOUND"}
		}
		if strings.Contains(strings.ToLower(envelope.Message), "not authorized") {
			return nil, &agentruntime.OrderLookupError{Code: "PERMISSION_DENIED"}
		}
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}
	// 从 data 节点中取出 order 对象，作为订单事实来源。
	orderMap, ok := pickMap(envelope.Data, "order")
	if !ok {
		return nil, &agentruntime.OrderLookupError{Code: "NOT_FOUND"}
	}
	// 把下游订单详情映射成 agent runtime 需要的统一订单快照。
	return &agentruntime.OrderSnapshot{
		OrderNo:            pickString(orderMap, "order_no", "orderNo"),
		SubOrderNo:         pickString(orderMap, "sub_order_no", "subOrderNo"),
		MainStatus:         strings.ToUpper(pickString(orderMap, "order_status", "orderStatus", "main_status", "mainStatus")),
		PaymentStatus:      strings.ToUpper(pickString(orderMap, "payment_status", "paymentStatus")),
		FulfillmentStatus:  normalizeFulfillmentStatus(orderMap),
		LogisticsStatus:    strings.ToUpper(pickString(orderMap, "logistics_status", "logisticsStatus", "delivery_status", "deliveryStatus")),
		AfterSaleStatus:    strings.ToUpper(pickString(orderMap, "after_sale_status", "afterSaleStatus", "refund_status", "refundStatus")),
		LatestUpdateTime:   pickString(orderMap, "updated_at", "updatedAt"),
		OwnershipConfirmed: true,
	}, nil
}

func (r *buyerOrderSnapshotHTTPRepository) ListRecentOrders(ctx context.Context, filter agentruntime.RecentOrderListFilter) ([]agentruntime.RecentOrderCandidate, error) {
	if r == nil || strings.TrimSpace(r.baseURL) == "" {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}
	if filter.UserID == 0 {
		return nil, &agentruntime.OrderLookupError{Code: "INSUFFICIENT_CONTEXT"}
	}

	limit := filter.Limit
	if limit == 0 {
		limit = 5
	}

	endpoint := fmt.Sprintf("%s/v1/order/buyer/orders?page_size=%d", r.baseURL, limit)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}
	req.Header.Set("X-User-Id", strconv.FormatUint(filter.UserID, 10))
	if strings.TrimSpace(filter.RequestID) != "" {
		req.Header.Set("X-Request-Id", strings.TrimSpace(filter.RequestID))
	}
	req.Header.Set("X-Service-Name", "agent-svc")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, &agentruntime.OrderLookupError{Code: "PERMISSION_DENIED"}
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}

	var envelope buyerOrderDetailEnvelope
	if err = json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}
	if envelope.Code != 0 && envelope.Code != 200 {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}

	orderList, ok := pickArray(envelope.Data, "orders", "list")
	if !ok {
		return nil, nil
	}

	candidates := make([]agentruntime.RecentOrderCandidate, 0, len(orderList))
	for _, rawOrder := range orderList {
		orderMap, mapOK := rawOrder.(map[string]any)
		if !mapOK {
			continue
		}
		orderNo := pickString(orderMap, "order_no", "orderNo")
		if strings.TrimSpace(orderNo) == "" {
			continue
		}

		candidate := agentruntime.RecentOrderCandidate{
			OrderNo:          orderNo,
			DisplayTitle:     extractRecentOrderTitle(orderMap),
			MainStatus:       strings.ToUpper(pickString(orderMap, "order_status", "orderStatus", "main_status", "mainStatus")),
			PaymentStatus:    strings.ToUpper(pickString(orderMap, "payment_status", "paymentStatus")),
			LatestUpdateTime: pickString(orderMap, "updated_at", "updatedAt"),
		}

		snapshot, detailErr := r.QueryOrderSnapshot(ctx, agentruntime.OrderOwnershipFilter{
			UserID:         filter.UserID,
			ShopNo:         filter.ShopNo,
			OrderNo:        orderNo,
			RequestID:      filter.RequestID,
			ConversationNo: filter.ConversationNo,
			RunNo:          filter.RunNo,
		})
		if detailErr == nil && snapshot != nil {
			candidate.SubOrderNo = snapshot.SubOrderNo
			candidate.FulfillmentStatus = snapshot.FulfillmentStatus
			candidate.LogisticsStatus = snapshot.LogisticsStatus
			candidate.AfterSaleStatus = snapshot.AfterSaleStatus
			if strings.TrimSpace(candidate.LatestUpdateTime) == "" {
				candidate.LatestUpdateTime = snapshot.LatestUpdateTime
			}
		}

		candidates = append(candidates, candidate)
	}

	return candidates, nil
}

func pickMap(source map[string]any, keys ...string) (map[string]any, bool) {
	// 逐个尝试候选字段名，兼容 snake_case 和 camelCase。
	for _, key := range keys {
		// 如果候选字段存在并且本身是 map，就直接返回。
		if raw, ok := source[key]; ok {
			if typed, innerOK := raw.(map[string]any); innerOK {
				return typed, true
			}
		}
	}
	// 没找到可用 map 时返回失败。
	return nil, false
}

func pickArray(source map[string]any, keys ...string) ([]any, bool) {
	for _, key := range keys {
		if raw, ok := source[key]; ok {
			if typed, innerOK := raw.([]any); innerOK {
				return typed, true
			}
		}
	}
	return nil, false
}

func pickString(source map[string]any, keys ...string) string {
	// 逐个尝试候选字段名，兼容不同下游序列化风格。
	for _, key := range keys {
		// 如果字段不存在，就继续尝试下一个候选名。
		raw, ok := source[key]
		if !ok || raw == nil {
			continue
		}
		// 把命中的字段统一转成字符串，便于上层规则消费。
		switch typed := raw.(type) {
		case string:
			return strings.TrimSpace(typed)
		default:
			return strings.TrimSpace(fmt.Sprintf("%v", typed))
		}
	}
	// 没找到值时返回空串，让上层走缺省分支。
	return ""
}

func normalizeFulfillmentStatus(orderMap map[string]any) string {
	// 优先读取已经结构化好的 fulfillment_status 字段。
	if fulfillment := strings.ToUpper(pickString(orderMap, "fulfillment_status", "fulfillmentStatus")); fulfillment != "" {
		return fulfillment
	}
	// 如果只有 delivery_status，就把常见值折叠到统一履约状态。
	switch strings.ToUpper(pickString(orderMap, "delivery_status", "deliveryStatus", "logistics_status", "logisticsStatus")) {
	case "UNSHIPPED", "NOT_SHIPPED", "PENDING_SHIPMENT":
		return "UNSHIPPED"
	case "SHIPPED", "IN_TRANSIT":
		return "SHIPPED"
	case "DELIVERED", "SIGNED":
		return "DELIVERED"
	default:
		return ""
	}
}

func extractRecentOrderTitle(orderMap map[string]any) string {
	subOrders, ok := pickArray(orderMap, "sub_orders", "subOrders")
	if ok {
		for _, rawSubOrder := range subOrders {
			subOrderMap, mapOK := rawSubOrder.(map[string]any)
			if !mapOK {
				continue
			}
			items, itemsOK := pickArray(subOrderMap, "items")
			if !itemsOK {
				continue
			}
			for _, rawItem := range items {
				itemMap, itemOK := rawItem.(map[string]any)
				if !itemOK {
					continue
				}
				title := pickString(itemMap, "spu_title", "spuTitle", "sku_name", "skuName")
				if strings.TrimSpace(title) != "" {
					return title
				}
			}
		}
	}
	return pickString(orderMap, "buyer_remark", "buyerRemark")
}
