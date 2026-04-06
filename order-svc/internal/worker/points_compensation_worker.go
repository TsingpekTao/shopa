package worker

import (
	"context"
	"sync"
	"time"

	orderlogic "github.com/TsingpekTao/shopa/order-svc/internal/logic/order"
	"github.com/gogf/gf/v2/frame/g"
)

var pointsCompensationWorkerOnce sync.Once

// StartPointsCompensationWorker 启动本地积分补偿任务 worker。
func StartPointsCompensationWorker(ctx context.Context) {
	pointsCompensationWorkerOnce.Do(func() {
		go func() {
			conf := loadPointsCompWorkerConf(ctx)
			if !conf.Enabled {
				g.Log().Info(ctx, "[order-svc] points compensation worker disabled")
				return
			}
			ticker := time.NewTicker(time.Duration(conf.IntervalSeconds) * time.Second)
			defer ticker.Stop()
			for {
				if err := orderlogic.ProcessPendingPointsCompensationTasks(context.Background(), conf.BatchSize); err != nil {
					g.Log().Errorf(ctx, "[order-svc] process points compensation batch failed: %+v", err)
				}
				<-ticker.C
			}
		}()
	})
}

type pointsCompWorkerConf struct {
	Enabled         bool
	IntervalSeconds int
	BatchSize       int
}

func loadPointsCompWorkerConf(ctx context.Context) pointsCompWorkerConf {
	return pointsCompWorkerConf{
		Enabled:         g.Cfg().MustGet(ctx, "worker.pointsCompensation.enabled", true).Bool(),
		IntervalSeconds: g.Cfg().MustGet(ctx, "worker.pointsCompensation.intervalSeconds", 10).Int(),
		BatchSize:       g.Cfg().MustGet(ctx, "worker.pointsCompensation.batchSize", 20).Int(),
	}
}

