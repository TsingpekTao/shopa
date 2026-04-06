package inventory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	catalogv1 "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	inventoryv1 "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
	"github.com/TsingpekTao/shopa/inventory-svc/internal/dao"
	"github.com/TsingpekTao/shopa/inventory-svc/internal/model/entity"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const defaultCatalogHTTPBaseURL = "http://127.0.0.1:8004"

type catalogProjectionHTTPReq struct {
	Items []*catalogv1.UpsertSkuStockProjectionItem `json:"items"`
}

type goFrameResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// syncCatalogProjectionByAdjustResults 在库存调整成功后，将最新库存状态同步到 catalog-svc。
func (s *sInventory) syncCatalogProjectionByAdjustResults(ctx context.Context, results []*inventoryv1.AdjustResultItem) error {
	return s.syncCatalogProjectionBySkuNos(ctx, adjustResultSkuNos(results))
}

// syncCatalogProjectionByReservation 在预占/确认/取消成功后，将相关 SKU 的库存状态同步到 catalog-svc。
func (s *sInventory) syncCatalogProjectionByReservation(ctx context.Context, reservation *inventoryv1.ReservationRecord) error {
	if reservation == nil {
		return nil
	}
	return s.syncCatalogProjectionBySkuNos(ctx, reservedItemSkuNos(reservation.GetItems()))
}

// syncCatalogProjectionBySkuNos 统一读取库存真相源并把投影结果推送给 catalog-svc。
func (s *sInventory) syncCatalogProjectionBySkuNos(ctx context.Context, skuNos []string) error {
	normalized := uniqueNonEmptyStrings(skuNos)
	if len(normalized) == 0 {
		return nil
	}

	var rows []*entity.InventoryStock
	err := dao.InventoryStock.Ctx(ctx).
		WhereIn(dao.InventoryStock.Columns().SkuNo, normalized).
		Scan(&rows)
	if err != nil {
		return gerror.Wrap(err, "query inventory stock for catalog projection failed")
	}
	if len(rows) == 0 {
		return nil
	}

	items := make([]*catalogv1.UpsertSkuStockProjectionItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || strings.TrimSpace(row.SkuNo) == "" || strings.TrimSpace(row.SpuNo) == "" {
			continue
		}
		items = append(items, &catalogv1.UpsertSkuStockProjectionItem{
			SkuNo:         row.SkuNo,
			SpuNo:         row.SpuNo,
			StockStatus:   toCatalogStockStatus(row.StockStatus),
			StockVersion:  row.StockVersion,
			SourceEventId: projectionSourceEventID(row),
			OccurredAt:    projectionOccurredAt(row),
		})
	}
	if len(items) == 0 {
		return nil
	}

	cfgValue, _ := g.Cfg().Get(ctx, "upstream.catalogHttp")
	baseURL := strings.TrimSpace(cfgValue.String())
	if baseURL == "" {
		baseURL = defaultCatalogHTTPBaseURL
	}

	payload, err := json.Marshal(catalogProjectionHTTPReq{Items: items})
	if err != nil {
		return gerror.Wrap(err, "marshal catalog stock projection request failed")
	}

	callCtx, callCancel := context.WithTimeout(ctx, 3*time.Second)
	defer callCancel()

	request, err := http.NewRequestWithContext(
		callCtx,
		http.MethodPost,
		strings.TrimRight(baseURL, "/")+"/v1/catalog/internal/stock-projection:upsert",
		bytes.NewReader(payload),
	)
	if err != nil {
		return gerror.Wrap(err, "create catalog stock projection request failed")
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return gerror.Wrap(err, "push inventory stock projection to catalog failed")
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return gerror.Newf("catalog projection http status %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}

	var envelope goFrameResponse[catalogv1.UpsertSkuStockProjectionRes]
	if err = json.Unmarshal(body, &envelope); err != nil {
		return gerror.Wrap(err, "decode catalog stock projection response failed")
	}
	if envelope.Code != 0 {
		return gerror.Newf("catalog projection failed: %s", envelope.Message)
	}
	return nil
}

// logCatalogProjectionSyncFailure 统一记录投影同步失败日志，避免主业务结果被投影问题反向打断。
func (s *sInventory) logCatalogProjectionSyncFailure(ctx context.Context, scene string, err error) {
	if err == nil {
		return
	}
	g.Log().Warningf(ctx, "sync catalog stock projection failed after %s: %+v", scene, err)
}

func adjustResultSkuNos(results []*inventoryv1.AdjustResultItem) []string {
	out := make([]string, 0, len(results))
	for _, item := range results {
		if item == nil {
			continue
		}
		out = append(out, item.GetSkuNo())
	}
	return out
}

func reservedItemSkuNos(items []*inventoryv1.ReservedItem) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		out = append(out, item.GetSkuNo())
	}
	return out
}

func uniqueNonEmptyStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func toCatalogStockStatus(status uint) catalogv1.StockStatus {
	switch inventoryv1.StockStatus(status) {
	case inventoryv1.StockStatus_STOCK_STATUS_IN_STOCK:
		return catalogv1.StockStatus_STOCK_STATUS_IN_STOCK
	case inventoryv1.StockStatus_STOCK_STATUS_OUT_OF_STOCK:
		return catalogv1.StockStatus_STOCK_STATUS_OUT_OF_STOCK
	default:
		return catalogv1.StockStatus_STOCK_STATUS_UNSPECIFIED
	}
}

func projectionSourceEventID(row *entity.InventoryStock) string {
	if row == nil {
		return ""
	}
	if strings.TrimSpace(row.LastEventId) != "" {
		return row.LastEventId
	}
	return fmt.Sprintf("ip:%s:%d", row.SkuNo, row.StockVersion)
}

func projectionOccurredAt(row *entity.InventoryStock) *timestamppb.Timestamp {
	if row != nil && row.UpdatedAt != nil {
		return timestamppb.New(row.UpdatedAt.Time)
	}
	return timestamppb.Now()
}
