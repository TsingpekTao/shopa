package worker

import (
	"context"
	"sync"
	"time"

	orderlogic "github.com/TsingpekTao/shopa/order-svc/internal/logic/order"
	"github.com/gogf/gf/v2/frame/g"
)

var unpaidCloseWorkerOnce sync.Once

// StartUnpaidCloseWorker 启动超时未支付订单关闭 worker。
func StartUnpaidCloseWorker(ctx context.Context) {
	unpaidCloseWorkerOnce.Do(func() {
		go func() {
			conf := loadUnpaidCloseWorkerConf(ctx)
			if !conf.Enabled {
				g.Log().Info(ctx, "[order-svc] unpaid close worker disabled")
				return
			}
			ticker := time.NewTicker(time.Duration(conf.IntervalSeconds) * time.Second)
			defer ticker.Stop()
			for {
				if err := orderlogic.ProcessExpiredUnpaidOrders(context.Background(), conf.BatchSize); err != nil {
					g.Log().Errorf(ctx, "[order-svc] process expired unpaid order batch failed: %+v", err)
				}
				<-ticker.C
			}
		}()
	})
}

type unpaidCloseWorkerConf struct {
	Enabled         bool
	IntervalSeconds int
	BatchSize       int
}

func loadUnpaidCloseWorkerConf(ctx context.Context) unpaidCloseWorkerConf {
	interval := g.Cfg().MustGet(ctx, "worker.unpaidClose.intervalSeconds", 10).Int()
	batch := g.Cfg().MustGet(ctx, "worker.unpaidClose.batchSize", 20).Int()
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
