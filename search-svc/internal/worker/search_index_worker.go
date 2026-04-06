package worker

import (
	"context"
	"sync"
	"time"

	searchlogic "github.com/TsingpekTao/shopa/search-svc/internal/logic/search"
	"github.com/gogf/gf/v2/frame/g"
)

var startOnce sync.Once

func Start(ctx context.Context) {
	startOnce.Do(func() {
		go startOutboxWorker(ctx)
		go startPhysicalDeleteWorker(ctx)
	})
}

func startOutboxWorker(ctx context.Context) {
	enabled := g.Cfg().MustGet(ctx, "worker.searchOutbox.enabled", true).Bool()
	if !enabled {
		g.Log().Info(ctx, "[search-svc] search outbox worker disabled")
		return
	}

	intervalSeconds := g.Cfg().MustGet(ctx, "worker.searchOutbox.intervalSeconds", 10).Int()
	if intervalSeconds <= 0 {
		intervalSeconds = 10
	}
	interval := time.Duration(intervalSeconds) * time.Second
	batchSize := g.Cfg().MustGet(ctx, "worker.searchOutbox.batchSize", 50).Int()
	if batchSize <= 0 {
		batchSize = 50
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err := searchlogic.ProcessSearchOutboxBatch(context.Background(), batchSize); err != nil {
			g.Log().Warningf(ctx, "[search-svc] process search outbox batch failed: %+v", err)
		}
		<-ticker.C
	}
}

func startPhysicalDeleteWorker(ctx context.Context) {
	enabled := g.Cfg().MustGet(ctx, "worker.physicalDelete.enabled", true).Bool()
	if !enabled {
		g.Log().Info(ctx, "[search-svc] physical delete worker disabled")
		return
	}

	intervalSeconds := g.Cfg().MustGet(ctx, "worker.physicalDelete.intervalSeconds", 300).Int()
	if intervalSeconds <= 0 {
		intervalSeconds = 300
	}
	interval := time.Duration(intervalSeconds) * time.Second
	batchSize := g.Cfg().MustGet(ctx, "worker.physicalDelete.batchSize", 50).Int()
	if batchSize <= 0 {
		batchSize = 50
	}
	retentionHours := g.Cfg().MustGet(ctx, "worker.physicalDelete.retentionHours", 168).Int()
	if retentionHours <= 0 {
		retentionHours = 168
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err := searchlogic.ProcessSearchPhysicalDeleteBatch(context.Background(), batchSize, retentionHours); err != nil {
			g.Log().Warningf(ctx, "[search-svc] process physical delete batch failed: %+v", err)
		}
		<-ticker.C
	}
}
