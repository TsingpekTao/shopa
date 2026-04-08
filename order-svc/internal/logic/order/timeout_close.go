package order

import (
	"context"
	"strings"
	"sync"

	v1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	"github.com/TsingpekTao/shopa/order-svc/internal/dao"
	"github.com/TsingpekTao/shopa/order-svc/internal/model/entity"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type unpaidCloseWorkerConf struct {
	Enabled         bool
	IntervalSeconds int
	BatchSize       int
}

type expiredUnpaidOrderCloser func(ctx context.Context, orderNo string) (bool, error)

var unpaidCloseWorkerOnce sync.Once

func ProcessExpiredUnpaidOrders(ctx context.Context, limit int) error {
	return New().processExpiredUnpaidOrders(ctx, limit)
}

func (s *sOrder) processExpiredUnpaidOrders(ctx context.Context, limit int) error {
	conf := loadUnpaidCloseWorkerConf(ctx)
	limit = normalizeUnpaidCloseBatchSize(limit, conf)

	orderNos, err := s.listExpiredUnpaidOrderNos(ctx, gtime.Now(), limit)
	if err != nil {
		return err
	}
	processed, closed := runExpiredUnpaidOrderBatch(ctx, orderNos, func(ctx context.Context, orderNo string) (bool, error) {
		res, err := s.CloseOrderIfUnpaid(ctx, &v1.CloseOrderIfUnpaidReq{
			OrderNo:    orderNo,
			ReasonCode: v1.CancelReasonCode_CANCEL_REASON_CODE_TIMEOUT_CLOSE,
		})
		if err != nil {
			return false, err
		}
		return res.GetClosed(), nil
	})
	if processed > 0 {
		g.Log().Infof(ctx, "[order-svc] processed expired unpaid orders, processed=%d closed=%d", processed, closed)
	}
	return nil
}

func (s *sOrder) listExpiredUnpaidOrderNos(ctx context.Context, now *gtime.Time, limit int) ([]string, error) {
	if now == nil {
		now = gtime.Now()
	}
	var rows []*entity.OrderMain
	cols := dao.OrderMain.Columns()
	if err := dao.OrderMain.Ctx(ctx).
		Fields(cols.OrderNo).
		Where(cols.OrderStatus, uint(v1.OrderStatus_ORDER_STATUS_PENDING_PAY)).
		WhereLTE(cols.PayDeadlineAt, now).
		WhereNull(cols.DeletedAt).
		OrderAsc(cols.PayDeadlineAt).
		OrderAsc(cols.Id).
		Limit(limit).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query expired unpaid orders failed")
	}
	orderNos := make([]string, 0, len(rows))
	for _, row := range rows {
		if row == nil || strings.TrimSpace(row.OrderNo) == "" {
			continue
		}
		orderNos = append(orderNos, strings.TrimSpace(row.OrderNo))
	}
	return orderNos, nil
}

func runExpiredUnpaidOrderBatch(ctx context.Context, orderNos []string, closeFn expiredUnpaidOrderCloser) (processed int, closed int) {
	if closeFn == nil {
		return 0, 0
	}
	for _, rawOrderNo := range orderNos {
		orderNo := strings.TrimSpace(rawOrderNo)
		if orderNo == "" {
			continue
		}
		processed++
		wasClosed, err := closeFn(ctx, orderNo)
		if err != nil {
			g.Log().Errorf(ctx, "[order-svc] close expired unpaid order failed, order_no=%s err=%+v", orderNo, err)
			continue
		}
		if wasClosed {
			closed++
		}
	}
	return processed, closed
}

func normalizeUnpaidCloseBatchSize(limit int, conf unpaidCloseWorkerConf) int {
	if limit > 0 {
		return limit
	}
	if conf.BatchSize > 0 {
		return conf.BatchSize
	}
	return 20
}

func loadUnpaidCloseWorkerConf(ctx context.Context) unpaidCloseWorkerConf {
	var (
		interval = g.Cfg().MustGet(ctx, "worker.unpaidClose.intervalSeconds", 10).Int()
		batch    = g.Cfg().MustGet(ctx, "worker.unpaidClose.batchSize", 20).Int()
	)
	if interval <= 0 {
		interval = 10
	}
	if batch <= 0 {
		batch = 20
	}
	return unpaidCloseWorkerConf{
		Enabled:         g.Cfg().MustGet(ctx, "worker.unpaidClose.enabled", true).Bool(),
		IntervalSeconds: interval,
		BatchSize:       batch,
	}
}
