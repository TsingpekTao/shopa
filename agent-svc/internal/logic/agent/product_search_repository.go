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

type buyerProductSearchHTTPRepository struct {
	baseURL string
	client  *http.Client
}

func newBuyerProductSearchHTTPRepository(baseURL string, client *http.Client) *buyerProductSearchHTTPRepository {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if client == nil {
		client = &http.Client{Timeout: 3 * time.Second}
	}
	return &buyerProductSearchHTTPRepository{
		baseURL: baseURL,
		client:  client,
	}
}

func (r *buyerProductSearchHTTPRepository) SearchProducts(ctx context.Context, filter agentruntime.ProductSearchFilter) ([]agentruntime.ProductSearchItem, error) {
	if r == nil || strings.TrimSpace(r.baseURL) == "" {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}

	limit := filter.Limit
	if limit == 0 {
		limit = 5
	}

	values := url.Values{}
	values.Set("query", strings.TrimSpace(filter.Query))
	values.Set("pageSize", strconv.FormatUint(uint64(limit), 10))
	if strings.TrimSpace(filter.ShopNo) != "" {
		values.Set("shopNo", strings.TrimSpace(filter.ShopNo))
	}
	if strings.TrimSpace(filter.CategoryNo) != "" {
		values.Set("categoryNo", strings.TrimSpace(filter.CategoryNo))
	}
	if strings.TrimSpace(filter.SortCode) != "" {
		values.Set("sortCode", strings.TrimSpace(filter.SortCode))
	}

	endpoint := fmt.Sprintf("%s/v1/search/products?%s", r.baseURL, values.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}
	req.Header.Set("X-Service-Name", "agent-svc")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, &agentruntime.OrderLookupError{Code: "DOWNSTREAM_UNAVAILABLE"}
	}
	defer resp.Body.Close()

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

	rawList, ok := pickArray(envelope.Data, "list")
	if !ok {
		return nil, nil
	}

	items := make([]agentruntime.ProductSearchItem, 0, len(rawList))
	for _, rawItem := range rawList {
		itemMap, mapOK := rawItem.(map[string]any)
		if !mapOK {
			continue
		}
		items = append(items, agentruntime.ProductSearchItem{
			SpuNo:      pickString(itemMap, "spu_no", "spuNo"),
			Title:      pickString(itemMap, "title"),
			CoverURL:   pickString(itemMap, "cover_url", "coverUrl"),
			MinPrice:   parseUint64Value(itemMap, "min_price", "minPrice"),
			MaxPrice:   parseUint64Value(itemMap, "max_price", "maxPrice"),
			ShopName:   pickString(itemMap, "shop_name", "shopName"),
			ReasonText: buildProductReasonText(strings.TrimSpace(filter.Query), itemMap),
		})
	}
	return items, nil
}

func parseUint64Value(source map[string]any, keys ...string) uint64 {
	for _, key := range keys {
		raw, ok := source[key]
		if !ok || raw == nil {
			continue
		}
		switch typed := raw.(type) {
		case float64:
			return uint64(typed)
		case int:
			return uint64(typed)
		case int64:
			return uint64(typed)
		case uint64:
			return typed
		case string:
			if value, err := strconv.ParseUint(strings.TrimSpace(typed), 10, 64); err == nil {
				return value
			}
		}
	}
	return 0
}

func buildProductReasonText(query string, item map[string]any) string {
	title := strings.TrimSpace(pickString(item, "title"))
	if query == "" {
		return ""
	}
	if title == "" {
		return "这款和你刚才提到的需求比较接近。"
	}
	return fmt.Sprintf("%s 和你刚才提到的“%s”需求比较接近。", title, query)
}
