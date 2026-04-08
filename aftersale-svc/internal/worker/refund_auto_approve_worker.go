package worker

import (
	"context"
	"sync"
	"time"

	aftersalelogic "github.com/TsingpekTao/shopa/aftersale-svc/internal/logic/aftersale"
	"github.com/gogf/gf/v2/frame/g"
)

var refundAutoApproveWorkerOnce sync.Once

// StartRefundAutoApproveWorker 定时扫描超时未审核的退款批次，自动同意并继续执行退款。
func StartRefundAutoApproveWorker(ctx context.Context) {
	refundAutoApproveWorkerOnce.Do(func() {
		go func() {
			conf := loadRefundAutoApproveWorkerConf(ctx)
			if !conf.Enabled {
				g.Log().Info(ctx, "[aftersale-svc] refund auto approve worker disabled")
				return
			}
			ticker := time.NewTicker(time.Duration(conf.IntervalSeconds) * time.Second)
			defer ticker.Stop()
			for {
				if err := aftersalelogic.ProcessExpiredRefundReviewBatches(context.Background(), conf.BatchSize); err != nil {
					g.Log().Errorf(ctx, "[aftersale-svc] process expired refund batches failed: %+v", err)
				}
				<-ticker.C
			}
		}()
	})
}

type refundAutoApproveWorkerConf struct {
	Enabled         bool
	IntervalSeconds int
	BatchSize       int
}

func loadRefundAutoApproveWorkerConf(ctx context.Context) refundAutoApproveWorkerConf {
	interval := g.Cfg().MustGet(ctx, "worker.refundAutoApprove.intervalSeconds", 30).Int()
	batch := g.Cfg().MustGet(ctx, "worker.refundAutoApprove.batchSize", 20).Int()
	if interval <= 0 {
		interval = 30
	}
	if batch <= 0 {
		batch = 20
	}
	return refundAutoApproveWorkerConf{
		Enabled:         g.Cfg().MustGet(ctx, "worker.refundAutoApprove.enabled", true).Bool(),
		IntervalSeconds: interval,
		BatchSize:       batch,
	}
}
