package order

import (
	"context"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	"github.com/TsingpekTao/shopa/order-svc/internal/dao"
	"github.com/TsingpekTao/shopa/order-svc/internal/model/entity"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type autoReceiveWorkerConf struct {
	Enabled          bool
	IntervalSeconds  int
	BatchSize        int
	ReceiveAfterDays int
}

type autoReceiveOrderCompleter func(ctx context.Context, orderNo string) (bool, error)

func ProcessAutoReceiveOrders(ctx context.Context, limit int) error {
	return New().processAutoReceiveOrders(ctx, limit)
}

func (s *sOrder) processAutoReceiveOrders(ctx context.Context, limit int) error {
	conf := loadAutoReceiveWorkerConf(ctx)
	limit = normalizeAutoReceiveBatchSize(limit, conf)
	window := normalizeAutoReceiveWindow(conf.ReceiveAfterDays)

	orderNos, err := s.listAutoReceivableOrderNos(ctx, gtime.Now(), limit, window)
	if err != nil {
		return err
	}
	processed, completed := runAutoReceiveBatch(ctx, orderNos, func(ctx context.Context, orderNo string) (bool, error) {
		res, err := s.CompleteOrder(ctx, &v1.CompleteOrderReq{
			OrderNo:              orderNo,
			CompletionNote:       "auto receive after shipped window elapsed",
			CompletionSourceCode: "AUTO_RECEIVE",
			IdempotencyKey:       buildAutoReceiveIdempotencyKey(orderNo),
		})
		if err != nil {
			return false, err
		}
		return res.GetOrderStatus() == v1.OrderStatus_ORDER_STATUS_COMPLETED, nil
	})
	if processed > 0 {
		g.Log().Infof(ctx, "[order-svc] processed auto receive batch, processed=%d completed=%d", processed, completed)
	}
	return nil
}

func (s *sOrder) listAutoReceivableOrderNos(ctx context.Context, now *gtime.Time, limit int, window time.Duration) ([]string, error) {
	if now == nil {
		now = gtime.Now()
	}
	cutoff := gtime.NewFromTime(now.Time.Add(-window))
	var rows []*entity.OrderSub
	cols := dao.OrderSub.Columns()
	if err := dao.OrderSub.Ctx(ctx).
		Fields(cols.OrderNo, cols.UpdatedAt).
		Where(cols.SubStatus, uint(v1.SubOrderStatus_SUB_ORDER_STATUS_SHIPPED)).
		WhereLTE(cols.UpdatedAt, cutoff).
		WhereNull(cols.DeletedAt).
		OrderAsc(cols.UpdatedAt).
		OrderAsc(cols.Id).
		Limit(limit * 3).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query auto receivable orders failed")
	}

	seen := make(map[string]struct{}, len(rows))
	orderNos := make([]string, 0, min(limit, len(rows)))
	for _, row := range rows {
		if row == nil || !shouldAutoReceiveSubOrder(row.UpdatedAt, now, window) {
			continue
		}
		orderNo := strings.TrimSpace(row.OrderNo)
		if orderNo == "" {
			continue
		}
		if _, ok := seen[orderNo]; ok {
			continue
		}
		seen[orderNo] = struct{}{}
		orderNos = append(orderNos, orderNo)
		if len(orderNos) >= limit {
			break
		}
	}
	return orderNos, nil
}

func shouldAutoReceiveSubOrder(shippedAt, now *gtime.Time, window time.Duration) bool {
	if shippedAt == nil {
		return false
	}
	if now == nil {
		now = gtime.Now()
	}
	if window <= 0 {
		window = 7 * 24 * time.Hour
	}
	return !shippedAt.Time.After(now.Time.Add(-window))
}

func runAutoReceiveBatch(ctx context.Context, orderNos []string, completeFn autoReceiveOrderCompleter) (processed int, completed int) {
	if completeFn == nil {
		return 0, 0
	}
	for _, rawOrderNo := range orderNos {
		orderNo := strings.TrimSpace(rawOrderNo)
		if orderNo == "" {
			continue
		}
		processed++
		done, err := completeFn(ctx, orderNo)
		if err != nil {
			g.Log().Errorf(ctx, "[order-svc] auto receive order failed, order_no=%s err=%+v", orderNo, err)
			continue
		}
		if done {
			completed++
		}
	}
	return processed, completed
}

func normalizeAutoReceiveBatchSize(limit int, conf autoReceiveWorkerConf) int {
	if limit > 0 {
		return limit
	}
	if conf.BatchSize > 0 {
		return conf.BatchSize
	}
	return 20
}

func normalizeAutoReceiveWindow(days int) time.Duration {
	if days <= 0 {
		days = 7
	}
	return time.Duration(days) * 24 * time.Hour
}

func loadAutoReceiveWorkerConf(ctx context.Context) autoReceiveWorkerConf {
	var (
		interval = g.Cfg().MustGet(ctx, "worker.autoReceive.intervalSeconds", 60).Int()
		batch    = g.Cfg().MustGet(ctx, "worker.autoReceive.batchSize", 20).Int()
		days     = g.Cfg().MustGet(ctx, "worker.autoReceive.receiveAfterDays", 7).Int()
	)
	if interval <= 0 {
		interval = 60
	}
	if batch <= 0 {
		batch = 20
	}
	if days <= 0 {
		days = 7
	}
	return autoReceiveWorkerConf{
		Enabled:          g.Cfg().MustGet(ctx, "worker.autoReceive.enabled", true).Bool(),
		IntervalSeconds:  interval,
		BatchSize:        batch,
		ReceiveAfterDays: days,
	}
}

func buildAutoReceiveIdempotencyKey(orderNo string) string {
	return "auto_receive_" + strings.TrimSpace(orderNo)
}

func min(left, right int) int {
	if left < right {
		return left
	}
	return right
}
