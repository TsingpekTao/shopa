package search

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	noopmetric "go.opentelemetry.io/otel/metric/noop"
)

var (
	// 用 once 保证指标对象只初始化一次，避免并发场景重复注册同名指标。
	metricsOnce sync.Once

	// 记录 ES 故障后回退 MySQL 的次数，用于触发静默降级告警。
	fallbackCounter metric.Int64Counter
	// 记录搜索请求耗时，便于区分 ES 主路径和 MySQL 回退路径的性能差异。
	queryLatency metric.Int64Histogram
	// 记录 ES 查询错误总量，便于快速判断故障类型和爆炸半径。
	esQueryErrors metric.Int64Counter
	// 记录 ES 外部版本冲突次数，便于观测乱序事件是否被正确拦截。
	esWriteConflicts metric.Int64Counter
	// 记录物理删除条数，便于观测软删回收任务是否正常推进。
	physicalDeletes metric.Int64Counter
)

func initMetrics() {
	// 用 once 把初始化收敛成一次，避免每次请求都重新向 OTel 申请指标对象。
	metricsOnce.Do(func() {
		// 使用固定 service meter 名称，方便指标平台按服务维度聚合。
		meter := otel.Meter("shopa/search-svc")
		// 同时准备 noop meter，确保本地未接 OTel 时也不会因为 nil 指标崩溃。
		noopMeter := noopmetric.NewMeterProvider().Meter("shopa/search-svc")

		// 注册回退计数器，用于观测 ES 失效后是否发生降级。
		fallbackCounter, _ = meter.Int64Counter("search_fallback_count_total")
		// 如果真实 meter 没有返回可用指标，就切到 noop 实现保持调用安全。
		if fallbackCounter == nil {
			fallbackCounter, _ = noopMeter.Int64Counter("search_fallback_count_total")
		}
		// 注册查询耗时直方图，用于比较不同引擎的响应时间。
		queryLatency, _ = meter.Int64Histogram("search_query_latency_ms")
		// 如果真实直方图不可用，就切到 noop 直方图避免埋点报错。
		if queryLatency == nil {
			queryLatency, _ = noopMeter.Int64Histogram("search_query_latency_ms")
		}
		// 注册 ES 查询错误计数器，用于定位 timeout/5xx/index_missing 等故障。
		esQueryErrors, _ = meter.Int64Counter("search_es_query_error_total")
		// 如果真实计数器不可用，就退回 noop 计数器保持容错。
		if esQueryErrors == nil {
			esQueryErrors, _ = noopMeter.Int64Counter("search_es_query_error_total")
		}
		// 注册 ES 写冲突计数器，用于观测 external version 是否拦住旧事件覆盖。
		esWriteConflicts, _ = meter.Int64Counter("search_es_write_conflict_total")
		// 如果真实计数器不可用，就退回 noop 计数器保持服务可用。
		if esWriteConflicts == nil {
			esWriteConflicts, _ = noopMeter.Int64Counter("search_es_write_conflict_total")
		}
		// 注册物理删除计数器，用于观测软删文档是否按计划被回收。
		physicalDeletes, _ = meter.Int64Counter("search_physical_delete_total")
		// 如果真实计数器不可用，就退回 noop 计数器保持埋点逻辑无副作用。
		if physicalDeletes == nil {
			physicalDeletes, _ = noopMeter.Int64Counter("search_physical_delete_total")
		}
	})
}

func recordFallback(ctx context.Context, reason string) {
	// 先确保指标对象完成初始化，避免首次调用时出现空指针。
	initMetrics()
	// 给回退事件打 reason 标签，方便运维区分 timeout、5xx、index_missing 等根因。
	fallbackCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("reason", reason)))
}

func recordQueryLatency(ctx context.Context, engine string, startedAt time.Time) {
	// 先确保直方图已注册完成，避免首次请求埋点失败。
	initMetrics()
	// 记录从请求开始到当前时刻的毫秒耗时，并附带 engine 标签区分 ES 与 MySQL。
	queryLatency.Record(ctx, time.Since(startedAt).Milliseconds(), metric.WithAttributes(attribute.String("engine", engine)))
}

func recordESQueryError(ctx context.Context, reason string) {
	// 先确保错误计数器已准备好，避免故障期间连埋点都失效。
	initMetrics()
	// 记录一次 ES 查询错误，并带上分类标签供告警与排障使用。
	esQueryErrors.Add(ctx, 1, metric.WithAttributes(attribute.String("reason", reason)))
}

func recordESWriteConflict(ctx context.Context, eventType string) {
	// 先确保冲突计数器已初始化，避免 outbox/重建链路报错时漏埋点。
	initMetrics()
	// 记录一次 ES 外部版本冲突，并标记触发冲突的事件类型。
	esWriteConflicts.Add(ctx, 1, metric.WithAttributes(attribute.String("event_type", eventType)))
}

func recordPhysicalDelete(ctx context.Context, deletedCount int) {
	// 先确保物理删除计数器已初始化，避免清理任务埋点失败。
	initMetrics()
	// 没有实际删除任何文档时直接返回，避免写入无意义的 0 值埋点。
	if deletedCount <= 0 {
		return
	}
	// 把本批次成功物理删除的文档条数累加进指标，便于观察回收吞吐。
	physicalDeletes.Add(ctx, int64(deletedCount))
}
