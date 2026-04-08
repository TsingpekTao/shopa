package worker

import (
	"context"
	"sync"
	"time"

	orderlogic "github.com/TsingpekTao/shopa/order-svc/internal/logic/order"
	"github.com/gogf/gf/v2/frame/g"
)

var autoReceiveWorkerOnce sync.Once

func StartAutoReceiveWorker(ctx context.Context) {
	autoReceiveWorkerOnce.Do(func() {
		go func() {
			conf := loadAutoReceiveWorkerConf(ctx)
			if !conf.Enabled {
				g.Log().Info(ctx, "[order-svc] auto receive worker disabled")
				return
			}
			ticker := time.NewTicker(time.Duration(conf.IntervalSeconds) * time.Second)
			defer ticker.Stop()
			for {
				if err := orderlogic.ProcessAutoReceiveOrders(context.Background(), conf.BatchSize); err != nil {
					g.Log().Errorf(ctx, "[order-svc] process auto receive batch failed: %+v", err)
				}
				<-ticker.C
			}
		}()
	})
}

type autoReceiveWorkerConf struct {
	Enabled         bool
	IntervalSeconds int
	BatchSize       int
}

func loadAutoReceiveWorkerConf(ctx context.Context) autoReceiveWorkerConf {
	interval := g.Cfg().MustGet(ctx, "worker.autoReceive.intervalSeconds", 60).Int()
	batch := g.Cfg().MustGet(ctx, "worker.autoReceive.batchSize", 20).Int()
	if interval <= 0 {
		interval = 60
	}
	if batch <= 0 {
		batch = 20
	}
	return autoReceiveWorkerConf{
		Enabled:         g.Cfg().MustGet(ctx, "worker.autoReceive.enabled", true).Bool(),
		IntervalSeconds: interval,
		BatchSize:       batch,
	}
}
