package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	sellershopv1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/dao"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/seller-shop-svc/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	storeCategoryAssignmentEventType = "ShopStoreCategoryAssignmentChanged"
	storeCategoryOutboxStatusNew     = 1
	storeCategoryOutboxStatusSent    = 2
	storeCategoryOutboxStatusDead    = 3
	storeCategoryLastErrorMaxLen     = 500
)

var workerOnce sync.Once

type sellerOutboxEnvelope struct {
	EventID    string                      `json:"event_id"`
	RequestID  string                      `json:"request_id"`
	OccurredAt string                      `json:"occurred_at"`
	Data       storeCategoryAssignmentBody `json:"data"`
}

type storeCategoryAssignmentBody struct {
	ShopNo            string   `json:"shop_no"`
	SpuNo             string   `json:"spu_no"`
	StoreCategoryId   uint64   `json:"store_category_id"`
	StoreCategoryL1   uint64   `json:"store_category_l1"`
	StoreCategoryL2   uint64   `json:"store_category_l2"`
	StoreCategoryPath []uint64 `json:"store_category_path"`
	UpdatedAt         string   `json:"updated_at"`
}

type searchPatchResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func Start(ctx context.Context) {
	workerOnce.Do(func() {
		go startStoreCategoryAssignmentWorker(ctx)
		go startStoreCategoryReconcileWorker(ctx)
	})
}

func startStoreCategoryAssignmentWorker(ctx context.Context) {
	if !g.Cfg().MustGet(ctx, "worker.storeCategoryAssignment.enabled", true).Bool() {
		g.Log().Info(ctx, "[seller-shop-svc] store category assignment worker disabled")
		return
	}

	intervalSeconds := g.Cfg().MustGet(ctx, "worker.storeCategoryAssignment.intervalSeconds", 10).Int()
	if intervalSeconds <= 0 {
		intervalSeconds = 10
	}
	batchSize := g.Cfg().MustGet(ctx, "worker.storeCategoryAssignment.batchSize", 50).Int()
	if batchSize <= 0 {
		batchSize = 50
	}

	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		if err := processStoreCategoryAssignmentBatch(context.Background(), batchSize); err != nil {
			g.Log().Warningf(ctx, "[seller-shop-svc] process store category assignment batch failed: %+v", err)
		}
		<-ticker.C
	}
}

func startStoreCategoryReconcileWorker(ctx context.Context) {
	if !g.Cfg().MustGet(ctx, "worker.storeCategoryReconcile.enabled", true).Bool() {
		g.Log().Info(ctx, "[seller-shop-svc] store category reconcile worker disabled")
		return
	}

	runHour := g.Cfg().MustGet(ctx, "worker.storeCategoryReconcile.runHour", 3).Int()
	if runHour < 0 || runHour > 23 {
		runHour = 3
	}

	for {
		wait := durationUntilNextHour(runHour)
		timer := time.NewTimer(wait)
		<-timer.C
		timer.Stop()

		result, err := service.SellerShop().ReconcileStoreCategoryCounts(context.Background(), &sellershopv1.ReconcileStoreCategoryCountsReq{})
		if err != nil {
			g.Log().Warningf(ctx, "[seller-shop-svc] reconcile store category counts failed: %+v", err)
			continue
		}
		g.Log().Infof(ctx, "[seller-shop-svc] reconciled store category counts, updated_categories=%d", result.UpdatedCategories)
	}
}

func durationUntilNextHour(targetHour int) time.Duration {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), targetHour, 0, 0, 0, now.Location())
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return next.Sub(now)
}

func processStoreCategoryAssignmentBatch(ctx context.Context, batchSize int) error {
	cols := dao.SellerOutboxEvent.Columns()
	now := gtime.Now()

	var rows []*entity.SellerOutboxEvent
	if err := dao.SellerOutboxEvent.Ctx(ctx).
		Where(cols.EventType, storeCategoryAssignmentEventType).
		Where(cols.Status, storeCategoryOutboxStatusNew).
		WhereLTE(cols.AvailableAt, now).
		OrderAsc(cols.Id).
		Limit(batchSize).
		Scan(&rows); err != nil {
		return err
	}

	for _, row := range rows {
		if row == nil {
			continue
		}
		if err := dispatchStoreCategoryAssignment(ctx, row); err != nil {
			return err
		}
	}
	return nil
}

func dispatchStoreCategoryAssignment(ctx context.Context, row *entity.SellerOutboxEvent) error {
	var envelope sellerOutboxEnvelope
	if err := json.Unmarshal([]byte(row.PayloadJson), &envelope); err != nil {
		return markStoreCategoryOutboxDead(ctx, row, fmt.Errorf("invalid payload json: %w", err))
	}
	if strings.TrimSpace(envelope.Data.SpuNo) == "" {
		return markStoreCategoryOutboxDead(ctx, row, fmt.Errorf("missing spu_no in store category assignment payload"))
	}

	if err := pushStoreCategoryPatch(ctx, envelope.Data); err != nil {
		return markStoreCategoryOutboxRetry(ctx, row, err)
	}
	return markStoreCategoryOutboxSent(ctx, row)
}

func pushStoreCategoryPatch(ctx context.Context, payload storeCategoryAssignmentBody) error {
	baseURL := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.search", "http://127.0.0.1:8018").String())
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8018"
	}
	timeoutSeconds := g.Cfg().MustGet(ctx, "worker.storeCategoryAssignment.requestTimeoutSeconds", 5).Int()
	if timeoutSeconds <= 0 {
		timeoutSeconds = 5
	}

	body := map[string]any{
		"spuNo":             payload.SpuNo,
		"storeCategoryId":   payload.StoreCategoryId,
		"storeCategoryL1":   payload.StoreCategoryL1,
		"storeCategoryL2":   payload.StoreCategoryL2,
		"storeCategoryPath": payload.StoreCategoryPath,
	}
	if updatedAt := parseRFC3339Timestamp(payload.UpdatedAt); updatedAt != nil {
		body["updatedAt"] = updatedAt
	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal store category patch request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(baseURL, "/")+"/v1/internal/search/docs/store-category", bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("build store category patch request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("call search-svc store category patch: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("search-svc store category patch http=%d body=%s", resp.StatusCode, truncateWorkerBody(respBody))
	}
	if len(respBody) == 0 {
		return nil
	}

	var wrapper searchPatchResponse
	if err := json.Unmarshal(respBody, &wrapper); err != nil {
		return nil
	}
	if wrapper.Code != 0 {
		return fmt.Errorf("search-svc store category patch rejected code=%d message=%s", wrapper.Code, wrapper.Message)
	}
	return nil
}

func parseRFC3339Timestamp(raw string) *timestamppb.Timestamp {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, raw)
	if err != nil {
		return nil
	}
	return timestamppb.New(parsed)
}

func markStoreCategoryOutboxSent(ctx context.Context, row *entity.SellerOutboxEvent) error {
	cols := dao.SellerOutboxEvent.Columns()
	now := gtime.Now()
	_, err := dao.SellerOutboxEvent.Ctx(ctx).
		Where(cols.Id, row.Id).
		Data(do.SellerOutboxEvent{
			Status:      storeCategoryOutboxStatusSent,
			SentAt:      now,
			LastError:   "",
			UpdatedAt:   now,
			AvailableAt: nil,
		}).
		Update()
	return err
}

func markStoreCategoryOutboxRetry(ctx context.Context, row *entity.SellerOutboxEvent, lastErr error) error {
	failCount := int(row.FailCount) + 1
	status := storeCategoryOutboxStatusNew
	if failCount >= 20 {
		status = storeCategoryOutboxStatusDead
	}
	nextRetry := gtime.NewFromTime(time.Now().Add(storeCategoryBackoffDuration(failCount)))
	cols := dao.SellerOutboxEvent.Columns()
	_, err := dao.SellerOutboxEvent.Ctx(ctx).
		Where(cols.Id, row.Id).
		Data(do.SellerOutboxEvent{
			Status:      status,
			FailCount:   failCount,
			LastError:   truncateWorkerString(lastErr.Error(), storeCategoryLastErrorMaxLen),
			AvailableAt: nextRetry,
			UpdatedAt:   gtime.Now(),
		}).
		Update()
	if err != nil {
		return err
	}
	return lastErr
}

func markStoreCategoryOutboxDead(ctx context.Context, row *entity.SellerOutboxEvent, lastErr error) error {
	cols := dao.SellerOutboxEvent.Columns()
	_, err := dao.SellerOutboxEvent.Ctx(ctx).
		Where(cols.Id, row.Id).
		Data(do.SellerOutboxEvent{
			Status:      storeCategoryOutboxStatusDead,
			FailCount:   row.FailCount,
			LastError:   truncateWorkerString(lastErr.Error(), storeCategoryLastErrorMaxLen),
			AvailableAt: nil,
			UpdatedAt:   gtime.Now(),
		}).
		Update()
	if err != nil {
		return err
	}
	return lastErr
}

func storeCategoryBackoffDuration(retryCount int) time.Duration {
	switch {
	case retryCount <= 1:
		return 5 * time.Second
	case retryCount <= 3:
		return 15 * time.Second
	case retryCount <= 6:
		return 30 * time.Second
	case retryCount <= 10:
		return 2 * time.Minute
	default:
		return 5 * time.Minute
	}
}

func truncateWorkerString(value string, limit int) string {
	if limit <= 0 || len(value) <= limit {
		return value
	}
	return value[:limit]
}

func truncateWorkerBody(body []byte) string {
	return truncateWorkerString(strings.TrimSpace(string(body)), 1000)
}
