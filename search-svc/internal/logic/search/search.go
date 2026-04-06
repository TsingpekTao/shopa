package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/search-svc/api/v1"
	"github.com/TsingpekTao/shopa/search-svc/internal/dao"
	"github.com/TsingpekTao/shopa/search-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/search-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/search-svc/internal/service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sSearch struct{}

func New() *sSearch {
	// 返回搜索逻辑实现实例，供 service 层在启动时注册使用。
	return &sSearch{}
}

func init() {
	// 在包初始化阶段注册搜索服务实现，确保 controller 能找到对应逻辑对象。
	service.RegisterSearch(New())
}

func (s *sSearch) SearchProducts(ctx context.Context, req *v1.SearchProductsReq) (*v1.SearchProductsRes, error) {
	// 空请求直接返回空列表，避免后续读取请求字段时出现空指针。
	if req == nil {
		return &v1.SearchProductsRes{List: []*v1.SpuCard{}}, nil
	}

	var (
		// 记录请求起始时间，用于后续按引擎维度打搜索耗时指标。
		startedAt = time.Now()
		// 归一化分页大小，统一收敛非法入参和超大分页风险。
		pageSize = normalizePageSize(req.GetPageSize())
		// 预留游标变量，用于承接 ES/MySQL 双模翻页状态。
		cursor searchCursor
		// 复用 err 变量承接游标解析错误，避免重复定义。
		err error
	)
	// 先解析前端传入的游标，恢复本次查询应该使用的翻页位点。
	if cursor, err = decodeCursor(req.GetNextCursor(), req); err != nil {
		// 如果游标损坏，则安全回退到首页，但仍保留当前查询的排序与摘要约束。
		cursor = searchCursor{SortCode: normalizeSortCode(req.GetSortCode()), QueryDigest: buildQueryDigest(req)}
	}

	// 优先获取 ES 引擎，因为 ES 是正式主检索路径，相关性和深分页能力更强。
	engine, esErr := getElasticsearchEngine(ctx)
	// 只有 ES 初始化成功时，才尝试走主检索路径。
	if esErr == nil {
		// 先走 ES 搜索，优先保证线上正式搜索体验。
		esResult, searchErr := engine.search(ctx, req, cursor, pageSize)
		// ES 查询成功时，直接返回主路径结果，不再触发 MySQL 回退。
		if searchErr == nil {
			// 把 ES 主路径耗时写入指标，方便区分主路径性能。
			recordQueryLatency(ctx, "es", startedAt)
			// 把真实搜索词写入统计表，支撑后续热词和联想词能力。
			touchKeywordStat(ctx, req.GetQuery())
			// 返回 ES 结果，并显式标记本次请求不是降级结果。
			return &v1.SearchProductsRes{
				List:       esResult.List,
				NextCursor: esResult.NextCursor,
				HasMore:    esResult.HasMore,
				Partial:    false,
			}, nil
		}
		// 记录 ES 查询错误，供后续降级原因和日志统一使用。
		esErr = searchErr
		// 对 ES 查询错误分类打点，便于运维快速识别 timeout/5xx/index_missing。
		recordESQueryError(ctx, classifyESError(searchErr))
	}

	// 一旦走到这里，就说明 ES 不可用或查询失败，需要记录一次降级事件。
	recordFallback(ctx, classifyESError(esErr))
	// 启动 MySQL 快照表回退查询，优先保可用性而不是保完全一致的相关性。
	mysqlResult, mysqlErr := searchProductsByMySQL(ctx, req, cursor, pageSize)
	// 如果 MySQL 回退也失败，则本次搜索彻底失败。
	if mysqlErr != nil {
		// 如果 ES 也失败过，则把两段错误一起拼起来，方便排障时看到完整链路。
		if esErr != nil {
			return nil, gerror.Wrapf(mysqlErr, "mysql fallback failed after elasticsearch error: %v", esErr)
		}
		// 如果只有 MySQL 自身失败，则直接返回 MySQL 错误。
		return nil, mysqlErr
	}

	// 记录 MySQL 回退路径的查询耗时，便于观察降级状态下的数据库压力。
	recordQueryLatency(ctx, "mysql_fallback", startedAt)
	// 即使走了回退路径，也要把本次搜索词计入关键词统计，保证运营数据不断档。
	touchKeywordStat(ctx, req.GetQuery())
	// 返回回退结果，并明确告诉前端当前处于降级查询状态。
	return &v1.SearchProductsRes{
		List:           mysqlResult.List,
		NextCursor:     mysqlResult.NextCursor,
		HasMore:        mysqlResult.HasMore,
		Partial:        true,
		DegradedFields: []string{"search_engine"},
	}, nil
}

func (s *sSearch) SuggestKeywords(ctx context.Context, req *v1.SuggestKeywordsReq) (*v1.SuggestKeywordsRes, error) {
	// 空请求直接返回空建议词列表，避免后续取前缀时出现空指针。
	if req == nil {
		return &v1.SuggestKeywordsRes{Keywords: []string{}}, nil
	}

	// 先读取前端传入的建议词上限。
	limit := int(req.GetLimit())
	// 如果上限非法，则回退到默认值，避免返回 0 条结果。
	if limit <= 0 {
		limit = 10
	}
	// 如果上限过大，则强行收敛到 20，避免低价值的大量建议拖慢响应。
	if limit > 20 {
		limit = 20
	}

	// 先从 MySQL 热词表和快照表获取建议词，因为这条路径复杂度最低且稳定。
	keywords, err := suggestKeywordsByMySQL(ctx, strings.TrimSpace(req.GetPrefix()), limit)
	// MySQL 建议词失败时直接返回错误，因为基础兜底路径已经不可用。
	if err != nil {
		return nil, err
	}

	// 当 MySQL 建议词数量不足时，再尝试用 ES 前缀查询补齐召回。
	if len(keywords) < limit {
		// 先获取 ES 引擎，避免 ES 不可用时直接把建议词接口拖死。
		engine, esErr := getElasticsearchEngine(ctx)
		// 只有 ES 可用时，才继续补全剩余候选词。
		if esErr == nil {
			// 只补足差额，避免 ES 结果把 MySQL 热词结果全部冲掉。
			esKeywords, suggestErr := engine.suggest(ctx, req.GetPrefix(), limit-len(keywords))
			// ES 建议成功时，把两路结果去重合并，保证前缀词展示更丰富。
			if suggestErr == nil {
				keywords = mergeUniqueStrings(keywords, esKeywords, limit)
			}
		}
	}

	// 返回合并后的建议词列表。
	return &v1.SuggestKeywordsRes{Keywords: keywords}, nil
}

func (s *sSearch) BatchGetSpuCards(ctx context.Context, req *v1.BatchGetSpuCardsReq) (*v1.BatchGetSpuCardsRes, error) {
	// 没有传任何 spu_no 时直接返回空列表，避免无意义查询。
	if req == nil || len(req.GetSpuNos()) == 0 {
		return &v1.BatchGetSpuCardsRes{List: []*v1.SpuCard{}}, nil
	}

	// 先做去空、去重与顺序规整，保证批量查询结果稳定可控。
	spuNos := normalizeSpuNos(req.GetSpuNos())
	// 如果规整后没有有效主键，则直接返回空列表。
	if len(spuNos) == 0 {
		return &v1.BatchGetSpuCardsRes{List: []*v1.SpuCard{}}, nil
	}

	// 优先走 ES mget，因为正式读路径应尽量复用搜索引擎的数据视图。
	engine, esErr := getElasticsearchEngine(ctx)
	// 只有 ES 引擎可用时，才尝试主路径读取。
	if esErr == nil {
		// 直接按传入主键顺序批量读取卡片，避免额外查询复杂度。
		list, err := engine.mget(ctx, spuNos)
		// ES 命中成功时，直接返回结果。
		if err == nil {
			return &v1.BatchGetSpuCardsRes{List: list}, nil
		}
		// 记录 ES 查询错误，便于观测批量卡片接口是否受 ES 故障影响。
		recordESQueryError(ctx, classifyESError(err))
		// 记录一次回退事件，提醒运维当前已经降级到 MySQL。
		recordFallback(ctx, classifyESError(err))
	}

	// ES 不可用时回退到 MySQL 快照表，优先保证详情页/列表页基础卡片可读。
	list, err := batchGetSpuCardsByMySQL(ctx, spuNos)
	// MySQL 回退失败时，才真正把错误抛给上层。
	if err != nil {
		// 如果 ES 之前也失败过，则把 ES 错误一起挂到返回信息里方便排障。
		if esErr != nil {
			return nil, gerror.Wrapf(err, "mysql batch fallback failed after elasticsearch error: %v", esErr)
		}
		// 如果只有 MySQL 自身失败，则直接返回 MySQL 错误。
		return nil, err
	}
	// 返回 MySQL 回退成功的卡片列表。
	return &v1.BatchGetSpuCardsRes{List: list}, nil
}

func (s *sSearch) UpsertSpuDoc(ctx context.Context, req *v1.UpsertSpuDocReq) (*v1.UpsertSpuDocRes, error) {
	// upsert 至少需要 spu_no 作为文档主键，否则无法写快照也无法写 ES。
	if req == nil || strings.TrimSpace(req.GetSpuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}

	// 先把请求转成统一搜索文档结构，保证 MySQL 与 ES 两条写路径共享同一份数据模型。
	current, err := loadSnapshotBySpuNo(ctx, req.GetSpuNo())
	if err != nil {
		return nil, err
	}
	doc := buildSearchDocFromUpsertReq(req, current)
	// 先写 MySQL 快照表，把数据库作为治理与回退真相源。
	applied, err := upsertSnapshot(ctx, doc)
	// 快照写失败时直接返回，因为真相源都没落下就不能继续写 ES。
	if err != nil {
		return nil, err
	}
	// 如果由于外部版本号过旧而没有真正应用快照，则说明这是过期事件，直接视为成功即可。
	if !applied {
		return &v1.UpsertSpuDocRes{Ok: true}, nil
	}

	// 获取 ES 引擎，尝试把最新快照同步到搜索索引。
	engine, esErr := getElasticsearchEngine(ctx)
	// 只有 ES 可用时，才尝试直接写索引。
	if esErr == nil {
		// 先尝试直写 ES，优先缩短搜索可见延迟。
		if err = engine.indexDoc(ctx, doc); err == nil {
			return &v1.UpsertSpuDocRes{Ok: true}, nil
		}
		// 如果是外部版本冲突，说明旧事件被拒绝，这属于正常防乱序行为，应直接视为成功。
		if isESConflict(err) {
			recordESWriteConflict(ctx, outboxEventTypeUpsert)
			return &v1.UpsertSpuDocRes{Ok: true}, nil
		}
		// ES 直写失败时，把事件写进 outbox，交给异步 worker 兜底重试。
		if queueErr := enqueueOutboxEvent(ctx, outboxEventTypeUpsert, doc); queueErr != nil {
			return nil, gerror.Wrapf(queueErr, "queue outbox after elasticsearch upsert failure failed: %v", err)
		}
		// 打一条降级日志，提示当前写链路已经退化成 MySQL + Outbox。
		g.Log().Warningf(ctx, "[search-svc] elasticsearch upsert degraded to outbox, spu_no=%s err=%+v", doc.SpuNo, err)
		return &v1.UpsertSpuDocRes{Ok: true}, nil
	}

	// 如果 ES 引擎压根不可用，也要保证 outbox 事件落盘，等待后续恢复后重放。
	if err = enqueueOutboxEvent(ctx, outboxEventTypeUpsert, doc); err != nil {
		return nil, err
	}
	// 只要快照和 outbox 写成功，就可以把本次请求视为成功接收。
	return &v1.UpsertSpuDocRes{Ok: true}, nil
}

func (s *sSearch) PatchSpuStoreCategory(ctx context.Context, req *v1.PatchSpuStoreCategoryReq) (*v1.PatchSpuStoreCategoryRes, error) {
	if req == nil || strings.TrimSpace(req.SpuNo) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}

	current, err := loadSnapshotBySpuNo(ctx, req.SpuNo)
	if err != nil {
		return nil, err
	}
	// 店内分类补丁只允许覆盖已有的搜索快照，避免反向创建缺少 shop_no/title 的空壳文档。
	if current == nil {
		g.Log().Warningf(ctx, "[search-svc] skip store-category patch for missing snapshot, spu_no=%s", strings.TrimSpace(req.SpuNo))
		return &v1.PatchSpuStoreCategoryRes{Ok: true}, nil
	}

	doc := searchDocFromEntity(current)
	doc.StoreCategoryId = req.StoreCategoryId
	doc.StoreCategoryL1 = req.StoreCategoryL1
	doc.StoreCategoryL2 = req.StoreCategoryL2
	doc.StoreCategoryPath = append([]uint64(nil), req.StoreCategoryPath...)

	updatedAt := time.Now()
	if req.UpdatedAt != nil && !req.UpdatedAt.AsTime().IsZero() {
		updatedAt = req.UpdatedAt.AsTime()
	}
	doc.SourceUpdatedAt = &updatedAt
	doc.SourceVersion = uint64(updatedAt.UnixMilli())
	if current.SourceVersion >= doc.SourceVersion {
		doc.SourceVersion = current.SourceVersion + 1
	}

	applied, err := upsertSnapshot(ctx, doc)
	if err != nil {
		return nil, err
	}
	if !applied {
		return &v1.PatchSpuStoreCategoryRes{Ok: true}, nil
	}

	engine, esErr := getElasticsearchEngine(ctx)
	if esErr == nil {
		if err = engine.indexDoc(ctx, doc); err == nil {
			return &v1.PatchSpuStoreCategoryRes{Ok: true}, nil
		}
		if isESConflict(err) {
			recordESWriteConflict(ctx, outboxEventTypeUpsert)
			return &v1.PatchSpuStoreCategoryRes{Ok: true}, nil
		}
		if queueErr := enqueueOutboxEvent(ctx, outboxEventTypeUpsert, doc); queueErr != nil {
			return nil, gerror.Wrapf(queueErr, "queue outbox after store-category patch failure failed: %v", err)
		}
		g.Log().Warningf(ctx, "[search-svc] elasticsearch store-category patch degraded to outbox, spu_no=%s err=%+v", doc.SpuNo, err)
		return &v1.PatchSpuStoreCategoryRes{Ok: true}, nil
	}

	if err = enqueueOutboxEvent(ctx, outboxEventTypeUpsert, doc); err != nil {
		return nil, err
	}
	return &v1.PatchSpuStoreCategoryRes{Ok: true}, nil
}

func (s *sSearch) DeleteSpuDoc(ctx context.Context, req *v1.DeleteSpuDocReq) (*v1.DeleteSpuDocRes, error) {
	// 删除操作同样必须指定 spu_no，否则无法定位要软删的文档。
	if req == nil || strings.TrimSpace(req.GetSpuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}

	// 先构造软删除文档快照，把当前状态、版本号和删除标记统一收口。
	doc, err := buildDeleteDoc(ctx, req.GetSpuNo())
	// 如果连删除快照都构造失败，则直接返回错误。
	if err != nil {
		return nil, err
	}

	// 先把 MySQL 快照表标记为软删除，保证治理层和回退层先看到删除事实。
	applied, err := softDeleteSnapshot(ctx, doc)
	// 快照软删失败时直接返回，因为真相源还没落下。
	if err != nil {
		return nil, err
	}
	// 如果因为版本过旧没有真正应用，则说明旧删除事件被拒绝，直接视为成功即可。
	if !applied {
		return &v1.DeleteSpuDocRes{Ok: true}, nil
	}

	// 获取 ES 引擎，尝试把软删除状态同步到搜索索引。
	engine, esErr := getElasticsearchEngine(ctx)
	// 只有 ES 可用时，才尝试同步删除状态。
	if esErr == nil {
		// 先做 ES 软删，避免文档在搜索结果里继续可见。
		if err = engine.softDeleteDoc(ctx, doc); err == nil {
			return &v1.DeleteSpuDocRes{Ok: true}, nil
		}
		// 如果是版本冲突，则说明这是过期删除事件，应直接视为成功。
		if isESConflict(err) {
			recordESWriteConflict(ctx, outboxEventTypeDelete)
			return &v1.DeleteSpuDocRes{Ok: true}, nil
		}
		// ES 软删失败时，把删除事件写进 outbox，交给异步 worker 重试。
		if queueErr := enqueueOutboxEvent(ctx, outboxEventTypeDelete, doc); queueErr != nil {
			return nil, gerror.Wrapf(queueErr, "queue outbox after elasticsearch delete failure failed: %v", err)
		}
		// 记录删除链路已经降级到 outbox 的告警日志。
		g.Log().Warningf(ctx, "[search-svc] elasticsearch delete degraded to outbox, spu_no=%s err=%+v", doc.SpuNo, err)
		return &v1.DeleteSpuDocRes{Ok: true}, nil
	}

	// ES 引擎不可用时，也要先把 outbox 事件落库，确保删除最终一致。
	if err = enqueueOutboxEvent(ctx, outboxEventTypeDelete, doc); err != nil {
		return nil, err
	}
	// 只要快照和 outbox 已经落下，就认为本次删除请求接收成功。
	return &v1.DeleteSpuDocRes{Ok: true}, nil
}

func (s *sSearch) RebuildIndex(ctx context.Context, req *v1.RebuildIndexReq) (*v1.RebuildIndexRes, error) {
	// 为重建任务生成唯一 job_no，方便后续后台追踪和页面展示。
	jobNo := fmt.Sprintf("RB%s", time.Now().Format("20060102150405"))
	// 统一记录当前时间，避免一条作业内出现多个不一致时间戳。
	now := gtime.Now()
	// 先把重建任务写入任务表，保证后台异步执行前就有审计记录。
	if _, err := dao.SearchRebuildJob.Ctx(ctx).Data(do.SearchRebuildJob{
		JobNo:      jobNo,
		ReasonCode: strings.TrimSpace(req.GetReasonCode()),
		Status:     "ACCEPTED",
		CreatedAt:  now,
		UpdatedAt:  now,
	}).Insert(); err != nil {
		return nil, gerror.Wrap(err, "insert search_rebuild_job failed")
	}

	// 使用 goroutine 异步执行重建，避免接口同步阻塞太久。
	go runRebuildJob(jobNo)
	// 立即返回 ACCEPTED，让调用方后续按 job_no 查询任务状态。
	return &v1.RebuildIndexRes{JobNo: jobNo, StatusCode: "ACCEPTED"}, nil
}

func runRebuildJob(jobNo string) {
	// 后台任务使用独立上下文，避免受原始 HTTP/RPC 请求取消影响。
	ctx := context.Background()
	// 复用任务表列定义，减少硬编码字段名出错风险。
	cols := dao.SearchRebuildJob.Columns()
	// 记录作业真正开始执行的时间点。
	now := gtime.Now()
	// 先把任务状态改成 RUNNING，避免任务一直停留在 ACCEPTED。
	_, _ = dao.SearchRebuildJob.Ctx(ctx).Where(cols.JobNo, jobNo).Data(do.SearchRebuildJob{
		Status:    "RUNNING",
		StartedAt: now,
		UpdatedAt: now,
	}).Update()

	// 重建依赖 ES，新建索引和别名切换都必须先拿到 ES 引擎。
	engine, err := getElasticsearchEngine(ctx)
	// ES 初始化失败时，把作业标记成 FAILED 并立刻结束。
	if err != nil {
		_, _ = dao.SearchRebuildJob.Ctx(ctx).Where(cols.JobNo, jobNo).Data(do.SearchRebuildJob{
			Status:     "FAILED",
			FinishedAt: gtime.Now(),
			UpdatedAt:  gtime.Now(),
		}).Update()
		g.Log().Errorf(ctx, "[search-svc] rebuild job %s init elasticsearch failed: %+v", jobNo, err)
		return
	}

	// 为本次重建生成新的版本化索引名，避免直接操作线上索引。
	indexName := engine.nextVersionedIndexName()
	// 先创建新索引及映射，确保后续导入有独立目标。
	if err = engine.createRebuildIndex(ctx, indexName); err != nil {
		_, _ = dao.SearchRebuildJob.Ctx(ctx).Where(cols.JobNo, jobNo).Data(do.SearchRebuildJob{
			Status:     "FAILED",
			FinishedAt: gtime.Now(),
			UpdatedAt:  gtime.Now(),
		}).Update()
		g.Log().Errorf(ctx, "[search-svc] rebuild job %s create index failed: %+v", jobNo, err)
		return
	}

	// 统计当前待重建的有效文档数，用于最后回填任务统计。
	totalCount, countErr := dao.SearchSpuDoc.Ctx(ctx).Where(dao.SearchSpuDoc.Columns().Deleted, 0).Count()
	// 统计失败不阻断重建，只记录 warning 方便排障。
	if countErr != nil {
		g.Log().Warningf(ctx, "[search-svc] rebuild job %s count docs failed: %+v", jobNo, countErr)
	}

	var (
		// lastID 记录上一批扫描到的主键，支撑按主键顺序滚动扫描。
		lastID uint64
		// successCount 统计成功导入到新索引的文档数。
		successCount uint64
		// failCount 统计导入失败次数，用于决定最终任务状态。
		failCount uint64
	)
	// 按批次循环导入快照表数据，直到没有更多记录为止。
	for {
		// 每批按主键增量拉取一段数据，避免 offset 深分页扫库。
		rows, listErr := listRebuildDocs(ctx, lastID, 200)
		// 如果拉取快照失败，则记录失败并中断整个重建流程。
		if listErr != nil {
			failCount++
			g.Log().Errorf(ctx, "[search-svc] rebuild job %s list docs failed: %+v", jobNo, listErr)
			break
		}
		// 当前批次没有任何数据时，说明导入阶段已经结束。
		if len(rows) == 0 {
			break
		}
		// 逐条把快照数据写入目标索引，保证单条失败不影响整批继续跑。
		for _, row := range rows {
			// 推进 lastID，确保下一批从当前记录之后继续扫描。
			lastID = row.Id
			// 把数据库实体转换成统一搜索文档结构，复用正式写链路的数据模型。
			doc := searchDocFromEntity(row)
			// 直接写入新索引实体名，而不是线上别名，避免重建污染线上流量。
			if indexErr := engine.indexDocToTarget(ctx, indexName, false, doc); indexErr != nil {
				// 单条写入失败时累计失败数，方便最后统一给出任务结果。
				failCount++
				// 版本冲突也需要打指标，用于观察源数据是否存在乱序或重复消费。
				if isESConflict(indexErr) {
					recordESWriteConflict(ctx, outboxEventTypeUpsert)
				}
				g.Log().Errorf(ctx, "[search-svc] rebuild job %s index doc failed, spu_no=%s err=%+v", jobNo, doc.SpuNo, indexErr)
				continue
			}
			// 单条写入成功后累计成功数，供任务结果统计使用。
			successCount++
		}
	}

	// 默认认为本次重建成功，只有出现失败时再下调状态。
	statusCode := "SUCCEEDED"
	// 只要有任意文档导入失败，就把任务整体标记为 FAILED。
	if failCount > 0 {
		statusCode = "FAILED"
	}
	// 只有导入全部成功时，才允许把线上别名切到新索引。
	if statusCode == "SUCCEEDED" {
		// 原子切换别名，保证搜索读写视图整体切换到新索引。
		if err = engine.swapRebuildAlias(ctx, indexName); err != nil {
			// 别名切换失败要把任务置为 FAILED，提醒人工介入。
			statusCode = "FAILED"
			g.Log().Errorf(ctx, "[search-svc] rebuild job %s swap alias failed: %+v", jobNo, err)
		}
	}

	// 把最终状态和统计结果回写任务表，供管理端查看。
	_, _ = dao.SearchRebuildJob.Ctx(ctx).Where(cols.JobNo, jobNo).Data(do.SearchRebuildJob{
		Status:       statusCode,
		TotalCount:   totalCount,
		SuccessCount: successCount,
		FailCount:    failCount,
		FinishedAt:   gtime.Now(),
		UpdatedAt:    gtime.Now(),
	}).Update()
}

func ProcessSearchOutboxBatch(ctx context.Context, batchSize int) error {
	// 未指定批量大小时，使用默认值控制单批消费量。
	if batchSize <= 0 {
		batchSize = 50
	}

	// 复用 outbox 表列定义，避免查询条件和更新条件写错字段。
	cols := dao.SearchOutboxEvent.Columns()
	// 预留切片承接本批待消费的 outbox 事件。
	var rows []*entity.SearchOutboxEvent
	// 只拉取 pending 且到达 next_retry_at 的事件，避免提前重试。
	err := dao.SearchOutboxEvent.Ctx(ctx).
		Where(cols.Status, outboxStatusPending).
		Wheref("(%s IS NULL OR %s <= ?)", cols.NextRetryAt, cols.NextRetryAt, gtime.Now()).
		OrderAsc(cols.Id).
		Limit(batchSize).
		Scan(&rows)
	if err != nil {
		return gerror.Wrap(err, "query search outbox events failed")
	}
	// 当前没有待消费事件时直接返回，避免 worker 空转。
	if len(rows) == 0 {
		return nil
	}

	// 消费前先确认 ES 引擎可用，否则本批无需继续执行。
	engine, err := getElasticsearchEngine(ctx)
	if err != nil {
		return err
	}

	// 逐条消费事件，保证单条失败不会阻断整批事件处理。
	for _, row := range rows {
		// 防御性跳过空指针事件，避免脏数据导致 worker 崩溃。
		if row == nil {
			continue
		}
		// 处理单条事件，并在失败时打 warning 日志供后续排查。
		if err = processSingleOutboxEvent(ctx, engine, row); err != nil {
			g.Log().Warningf(ctx, "[search-svc] process outbox event failed, event_id=%s err=%+v", row.EventId, err)
		}
	}
	// 本批事件处理完成后返回 nil，让调度器继续下一轮轮询。
	return nil
}

func ProcessSearchPhysicalDeleteBatch(ctx context.Context, batchSize int, retentionHours int) error {
	// 未传批量大小时，使用默认批次控制单轮物理删除吞吐。
	if batchSize <= 0 {
		batchSize = 50
	}
	// 未传保留时长时，默认保留 7 天后再物理删除，给乱序事件留缓冲窗口。
	if retentionHours <= 0 {
		retentionHours = 168
	}

	// 执行物理删除前先确认 ES 引擎可用，否则无法真正删掉索引文档。
	engine, err := getElasticsearchEngine(ctx)
	if err != nil {
		return err
	}

	// 计算软删除文档的截止时间，只有超过窗口的文档才允许物理删除。
	cutoff := gtime.NewFromTime(time.Now().Add(-time.Duration(retentionHours) * time.Hour))
	// 复用快照表列定义，避免条件构造时写错字段名。
	cols := dao.SearchSpuDoc.Columns()
	// 预留切片承接本批待物理删除的候选文档。
	var rows []*entity.SearchSpuDoc
	// 只扫描已软删、带 deleted_at 且超过保留窗口的文档。
	err = dao.SearchSpuDoc.Ctx(ctx).
		Where(cols.Deleted, 1).
		WhereNotNull(cols.DeletedAt).
		WhereLTE(cols.DeletedAt, cutoff).
		OrderAsc(cols.Id).
		Limit(batchSize).
		Scan(&rows)
	if err != nil {
		return gerror.Wrap(err, "query physical delete candidates failed")
	}

	// 统计本批真正完成物理删除的数量，供指标埋点使用。
	deletedCount := 0
	// 逐条执行 ES 物理删除，确保单条失败能立即返回并阻断后续危险操作。
	for _, row := range rows {
		// 防御性跳过空指针行，避免脏数据导致 worker 崩溃。
		if row == nil {
			continue
		}
		// 先从 ES 物理删除文档，并带外部版本控制防止旧事件把文档删错。
		if err = engine.physicalDeleteDoc(ctx, row.SpuNo, row.SourceVersion); err != nil {
			// 如果删除失败且属于版本冲突，也要打指标，说明删除顺序与最新版本有竞争。
			if isESConflict(err) {
				recordESWriteConflict(ctx, outboxEventTypeDelete)
			}
			return err
		}
		// ES 物理删除成功后，把 deleted_at 清空，表示该软删记录已完成回收。
		if _, updateErr := dao.SearchSpuDoc.DB().Ctx(ctx).Exec(ctx,
			fmt.Sprintf("UPDATE %s SET %s = NULL, %s = ? WHERE %s = ?", dao.SearchSpuDoc.Table(), cols.DeletedAt, cols.UpdatedAt, cols.Id),
			gtime.Now(),
			row.Id,
		); updateErr != nil {
			return gerror.Wrap(updateErr, "mark physical delete candidate completed failed")
		}
		// 每成功回收一条文档，就累加删除计数。
		deletedCount++
	}

	// 把本批物理删除成功数量写入指标，便于观测回收任务进展。
	recordPhysicalDelete(ctx, deletedCount)
	return nil
}

func processSingleOutboxEvent(ctx context.Context, engine *esEngine, row *entity.SearchOutboxEvent) error {
	// 先把 outbox payload 解析出来，恢复本次需要同步的 ES 事件内容。
	var payload outboxPayload
	// payload 解码失败说明事件体已经损坏，需要直接标死避免无限重试。
	if err := json.Unmarshal([]byte(row.PayloadJson), &payload); err != nil {
		return markOutboxDead(ctx, row, gerror.Wrap(err, "decode outbox payload failed"))
	}

	// 预留执行错误，便于根据事件类型统一走成功、重试或死亡分支。
	var execErr error
	// 按事件类型选择真正的 ES 操作，保持 outbox 消费逻辑和业务事件一一对应。
	switch payload.EventType {
	// upsert 事件走 ES 文档写入流程。
	case outboxEventTypeUpsert:
		execErr = engine.indexDoc(ctx, payload.Doc)
	// delete 事件走 ES 软删除流程。
	case outboxEventTypeDelete:
		execErr = engine.softDeleteDoc(ctx, payload.Doc)
	// 未知事件类型直接视为非法数据，避免被重复消费。
	default:
		execErr = gerror.Newf("unsupported outbox event type %s", payload.EventType)
	}

	// 执行成功或仅发生版本冲突时，都视为无需再重试的终态。
	if execErr == nil || isESConflict(execErr) {
		// 版本冲突属于正常防乱序结果，也要打指标供观测。
		if isESConflict(execErr) {
			recordESWriteConflict(ctx, payload.EventType)
		}
		// 把事件标记成 done，避免后续重复消费。
		return markOutboxDone(ctx, row)
	}
	// 其他错误按退避策略重试，直到达到 dead 条件。
	return markOutboxRetry(ctx, row, execErr)
}

func markOutboxDone(ctx context.Context, row *entity.SearchOutboxEvent) error {
	// 复用列定义，避免状态回写时写错字段名。
	cols := dao.SearchOutboxEvent.Columns()
	// 把当前事件标记为 done，并清空 next_retry_at，表示不再需要 worker 处理。
	_, err := dao.SearchOutboxEvent.Ctx(ctx).Where(cols.Id, row.Id).Data(do.SearchOutboxEvent{
		Status:      outboxStatusDone,
		RetryCount:  row.RetryCount,
		NextRetryAt: nil,
		UpdatedAt:   gtime.Now(),
	}).Update()
	// 返回数据库更新结果，便于上层感知状态回写是否成功。
	return err
}

func markOutboxRetry(ctx context.Context, row *entity.SearchOutboxEvent, lastErr error) error {
	// 每次失败都把重试次数加一，便于后续退避和死信判断。
	retryCount := int(row.RetryCount) + 1
	// 默认把事件保持为 pending，等待下一次 worker 轮询。
	status := outboxStatusPending
	// 当重试次数过多时，直接把事件打成 dead，避免无限重放拖垮系统。
	if retryCount >= 20 {
		status = outboxStatusDead
	}
	// 根据当前重试次数计算下一次重试时间，给下游故障留恢复窗口。
	nextRetry := gtime.NewFromTime(time.Now().Add(outboxBackoffDuration(retryCount)))
	// 复用列定义，保证状态更新字段统一可靠。
	cols := dao.SearchOutboxEvent.Columns()
	// 回写新的状态、重试次数和下一次重试时间。
	_, err := dao.SearchOutboxEvent.Ctx(ctx).Where(cols.Id, row.Id).Data(do.SearchOutboxEvent{
		Status:      status,
		RetryCount:  retryCount,
		NextRetryAt: nextRetry,
		UpdatedAt:   gtime.Now(),
	}).Update()
	if err != nil {
		return gerror.Wrapf(err, "mark outbox retry failed after error: %v", lastErr)
	}
	// 数据库状态更新成功后，继续把原始业务错误返回给上层日志链路。
	return lastErr
}

func markOutboxDead(ctx context.Context, row *entity.SearchOutboxEvent, lastErr error) error {
	// 复用列定义，保证 dead 状态回写字段正确。
	cols := dao.SearchOutboxEvent.Columns()
	// 直接把事件标记为 dead，并清空 next_retry_at，防止继续被 worker 消费。
	_, err := dao.SearchOutboxEvent.Ctx(ctx).Where(cols.Id, row.Id).Data(do.SearchOutboxEvent{
		Status:      outboxStatusDead,
		RetryCount:  row.RetryCount,
		NextRetryAt: nil,
		UpdatedAt:   gtime.Now(),
	}).Update()
	if err != nil {
		return gerror.Wrapf(err, "mark outbox dead failed after error: %v", lastErr)
	}
	// 保留原始错误返回，方便调用方和日志系统看到最初失败原因。
	return lastErr
}

func enqueueOutboxEvent(ctx context.Context, eventType string, doc *searchDoc) error {
	// 先把事件类型和文档快照编码成 payload，供 worker 后续重放。
	payload, err := json.Marshal(outboxPayload{EventType: eventType, Doc: doc})
	// payload 编码失败时直接返回，因为连重放数据都不完整。
	if err != nil {
		return gerror.Wrap(err, "marshal search outbox payload failed")
	}
	// 统一记录当前时间，既用作 next_retry_at，也用作创建更新时间。
	now := gtime.Now()
	// 生成唯一 event_id，便于后续日志、审计和人工排查。
	eventID := fmt.Sprintf("SE%s", time.Now().Format("20060102150405.000000000"))
	// 把事件落入 outbox 表，作为 ES 异步同步与重试的真相源。
	_, err = dao.SearchOutboxEvent.Ctx(ctx).Data(do.SearchOutboxEvent{
		EventId:       eventID,
		AggregateType: "SEARCH_SPU_DOC",
		AggregateId:   doc.SpuNo,
		EventType:     eventType,
		PayloadJson:   string(payload),
		Status:        outboxStatusPending,
		RetryCount:    0,
		NextRetryAt:   now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}).Insert()
	if err != nil {
		return gerror.Wrap(err, "insert search outbox event failed")
	}
	// 事件写入成功后返回 nil，让上层继续走“已接收”语义。
	return nil
}

func searchProductsByMySQL(ctx context.Context, req *v1.SearchProductsReq, cursor searchCursor, pageSize int) (*mysqlSearchResult, error) {
	cols := dao.SearchSpuDoc.Columns()
	model := dao.SearchSpuDoc.Ctx(ctx).
		Where(cols.Deleted, 0).
		WhereIn(cols.OnShelfStatusCode, []string{"ON", "ON_SHELF", "4"})

	if query := strings.TrimSpace(req.GetQuery()); query != "" {
		like := "%" + query + "%"
		model = model.Wheref("(%s LIKE ? OR %s LIKE ? OR %s LIKE ? OR %s LIKE ?)", cols.Title, cols.SpuNo, cols.ShopName, cols.AttrsJson, like, like, like, like)
	}
	if shopNo := strings.TrimSpace(req.GetShopNo()); shopNo != "" {
		model = model.Where(cols.ShopNo, shopNo)
	}
	if categoryNo := strings.TrimSpace(req.GetCategoryNo()); categoryNo != "" {
		model = model.Where(cols.CategoryNo, categoryNo)
	}
	if req.GetStoreCategoryId() > 0 {
		switch req.GetStoreCategoryLevel() {
		case 1:
			model = model.Where(cols.StoreCategoryL1, req.GetStoreCategoryId())
		case 2:
			model = model.Where(cols.StoreCategoryId, req.GetStoreCategoryId())
		default:
			model = model.Where(cols.StoreCategoryId, req.GetStoreCategoryId())
		}
	}

	model = applyMySQLSort(model, req)

	var rows []*entity.SearchSpuDoc
	if err := model.Offset(cursor.MySQLOffset).Limit(pageSize + 1).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query search_spu_doc fallback failed")
	}

	hasMore := len(rows) > pageSize
	if hasMore {
		rows = rows[:pageSize]
	}

	list := make([]*v1.SpuCard, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, toSpuCard(row))
	}

	nextCursor := ""
	if hasMore {
		nextCursor = encodeCursor(searchCursor{
			EngineHint:  "mysql",
			MySQLOffset: cursor.MySQLOffset + len(list),
			SortCode:    cursor.SortCode,
			QueryDigest: cursor.QueryDigest,
		})
	}
	return &mysqlSearchResult{List: list, NextCursor: nextCursor, HasMore: hasMore}, nil
}

func batchGetSpuCardsByMySQL(ctx context.Context, spuNos []string) ([]*v1.SpuCard, error) {
	// 复用快照表列定义，避免字段名硬编码分散在查询里。
	cols := dao.SearchSpuDoc.Columns()
	// 预留结果切片承接批量扫描出的快照记录。
	var rows []*entity.SearchSpuDoc
	// 先按 spu_no 批量查快照表，并过滤掉已删除文档。
	err := dao.SearchSpuDoc.Ctx(ctx).
		WhereIn(cols.SpuNo, spuNos).
		Where(cols.Deleted, 0).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "batch query search_spu_doc fallback failed")
	}

	// 先把扫描结果按 spu_no 建索引，便于后面按入参顺序重排。
	bySpu := make(map[string]*entity.SearchSpuDoc, len(rows))
	for _, row := range rows {
		// 防御性跳过空行，避免脏数据导致 map 赋值异常。
		if row == nil {
			continue
		}
		// 以 spu_no 为键缓存快照，供后续顺序回填。
		bySpu[row.SpuNo] = row
	}

	// 按原始请求顺序构造返回结果，避免数据库 IN 查询打乱顺序影响前端。
	list := make([]*v1.SpuCard, 0, len(spuNos))
	for _, spuNo := range spuNos {
		// 只回填真实查到的文档，缺失主键直接跳过。
		if row, ok := bySpu[spuNo]; ok {
			list = append(list, toSpuCard(row))
		}
	}
	// 返回按入参顺序重排后的卡片列表。
	return list, nil
}

func suggestKeywordsByMySQL(ctx context.Context, prefix string, limit int) ([]string, error) {
	// 先清洗前缀，避免空白字符影响后续 LIKE 匹配。
	prefix = strings.TrimSpace(prefix)
	// 预留建议词切片，优先承接热词表结果。
	keywords := make([]string, 0, limit)

	// 预留热词表扫描结果，便于先走低成本热词召回。
	var statRows []*entity.SearchKeywordStat
	// 先构造热词表查询模型，默认从 search_keyword_stat 读取热搜词。
	statModel := dao.SearchKeywordStat.Ctx(ctx)
	// 有前缀时只取匹配前缀的热词，提升联想词相关性。
	if prefix != "" {
		statModel = statModel.WhereLike(dao.SearchKeywordStat.Columns().Keyword, prefix+"%")
	}
	// 按搜索次数和最近搜索时间排序，优先返回最热门且最近仍有效的关键词。
	if err := statModel.OrderDesc(dao.SearchKeywordStat.Columns().SearchCount).OrderDesc(dao.SearchKeywordStat.Columns().LastSearchedAt).Limit(limit * 2).Scan(&statRows); err != nil {
		return nil, gerror.Wrap(err, "query search_keyword_stat failed")
	}

	// 用 seen 做去重，避免热词表和标题补全阶段出现重复词。
	seen := make(map[string]struct{}, limit)
	for _, row := range statRows {
		// 防御性跳过空行，避免脏数据造成空指针问题。
		if row == nil {
			continue
		}
		// 统一清洗关键词，避免空白字符被当成有效建议词返回。
		keyword := strings.TrimSpace(row.Keyword)
		// 空关键词没有展示价值，直接跳过。
		if keyword == "" {
			continue
		}
		// 已返回过的词不再重复加入，保证前端列表整洁。
		if _, ok := seen[keyword]; ok {
			continue
		}
		// 标记当前关键词已经被选中，防止后续重复。
		seen[keyword] = struct{}{}
		// 把热词加入结果，优先保证热门词召回。
		keywords = append(keywords, keyword)
		// 达到目标数量后立即返回，避免继续查标题补全浪费资源。
		if len(keywords) >= limit {
			return keywords, nil
		}
	}

	// 热词不足时，再从商品快照标题里补充联想词。
	var docRows []*entity.SearchSpuDoc
	docModel := dao.SearchSpuDoc.Ctx(ctx).
		Fields(dao.SearchSpuDoc.Columns().Title).
		Where(dao.SearchSpuDoc.Columns().Deleted, 0).
		WhereIn(dao.SearchSpuDoc.Columns().OnShelfStatusCode, []string{"ON", "ON_SHELF", "4"})
	if prefix != "" {
		docModel = docModel.WhereLike(dao.SearchSpuDoc.Columns().Title, "%"+prefix+"%")
	}
	// 按销量和最近更新时间排序，优先用热门且较新的商品标题补全建议词。
	if err := docModel.OrderDesc(dao.SearchSpuDoc.Columns().SalesCount).OrderDesc(dao.SearchSpuDoc.Columns().SourceUpdatedAt).Limit(limit * 3).Scan(&docRows); err != nil {
		// 标题补全失败时直接返回已有热词结果，避免联想词接口因为补充路径失败而整体不可用。
		return keywords, nil
	}

	// 逐条扫描标题候选，并继续沿用 seen 做全局去重。
	for _, row := range docRows {
		// 防御性跳过空行，避免脏数据影响补全过程。
		if row == nil {
			continue
		}
		// 清洗标题字符串，避免空白标题进入联想词结果。
		title := strings.TrimSpace(row.Title)
		// 空标题没有补全价值，直接跳过。
		if title == "" {
			continue
		}
		// 已经出现过的标题不再重复返回。
		if _, ok := seen[title]; ok {
			continue
		}
		// 标记该标题已选中，避免重复加入结果集。
		seen[title] = struct{}{}
		// 把标题加入建议词列表，补足热词不足的部分。
		keywords = append(keywords, title)
		// 达到目标数量后立即停止，控制查询结果规模。
		if len(keywords) >= limit {
			break
		}
	}
	// 返回 MySQL 热词+标题补全合并后的建议词结果。
	return keywords, nil
}

func applyMySQLSort(model *gdb.Model, req *v1.SearchProductsReq) *gdb.Model {
	// 复用快照表列定义，保证排序字段和表结构同步。
	cols := dao.SearchSpuDoc.Columns()
	// 按统一排序码切换 MySQL 回退排序，尽量贴近 ES 主路径体验。
	switch normalizeSortCode(req.GetSortCode()) {
	// 销量排序优先 sales_count，其次用更新时间和 spu_no 保证稳定顺序。
	case "sales_desc":
		return model.OrderDesc(cols.SalesCount).OrderDesc(cols.SourceUpdatedAt).OrderAsc(cols.SpuNo)
	// 价格升序时优先 min_price，再用更新时间和主键兜底。
	case "price_asc":
		return model.OrderAsc(cols.MinPrice).OrderDesc(cols.SourceUpdatedAt).OrderAsc(cols.SpuNo)
	// 价格降序时优先 min_price 降序，再用更新时间和主键兜底。
	case "price_desc":
		return model.OrderDesc(cols.MinPrice).OrderDesc(cols.SourceUpdatedAt).OrderAsc(cols.SpuNo)
	// 默认排序按销量、更新时间和主键排序，尽量贴近无关键词列表页体验。
	default:
		return model.OrderDesc(cols.SalesCount).OrderDesc(cols.SourceUpdatedAt).OrderAsc(cols.SpuNo)
	}
}

func buildSearchDocFromUpsertReq(req *v1.UpsertSpuDocReq, current *entity.SearchSpuDoc) *searchDoc {
	sourceUpdatedAt := time.Now()
	if req.GetUpdatedAt() != nil {
		sourceUpdatedAt = req.GetUpdatedAt().AsTime()
	}

	doc := &searchDoc{
		SpuNo:             strings.TrimSpace(req.GetSpuNo()),
		Title:             strings.TrimSpace(req.GetTitle()),
		ShopNo:            strings.TrimSpace(req.GetShopNo()),
		ShopName:          strings.TrimSpace(req.GetShopName()),
		CategoryNo:        strings.TrimSpace(req.GetCategoryNo()),
		CoverAssetId:      req.GetCoverAssetId(),
		CoverUrl:          strings.TrimSpace(req.GetCoverUrl()),
		MinPrice:          req.GetMinPrice(),
		MaxPrice:          req.GetMaxPrice(),
		StockTotal:        req.GetStockTotal(),
		SalesCount:        req.GetSalesCount(),
		AvgScoreX100:      req.GetAvgScoreX100(),
		ReviewTotal:       req.GetReviewTotal(),
		ShopStatusCode:    strings.TrimSpace(req.GetShopStatusCode()),
		OnShelfStatusCode: strings.TrimSpace(req.GetOnShelfStatusCode()),
		AttrsJson:         strings.TrimSpace(req.GetAttrsJson()),
		SourceUpdatedAt:   &sourceUpdatedAt,
		SourceVersion:     deriveSourceVersion(req.GetUpdatedAt()),
		Deleted:           false,
		DeletedAt:         nil,
	}

	if hasStoreCategoryPayload(req) {
		doc.StoreCategoryId = req.GetStoreCategoryId()
		doc.StoreCategoryL1 = req.GetStoreCategoryL1()
		doc.StoreCategoryL2 = req.GetStoreCategoryL2()
		doc.StoreCategoryPath = append([]uint64(nil), req.GetStoreCategoryPath()...)
	} else if current != nil {
		doc.StoreCategoryId = current.StoreCategoryId
		doc.StoreCategoryL1 = current.StoreCategoryL1
		doc.StoreCategoryL2 = current.StoreCategoryL2
		doc.StoreCategoryPath = parseStoreCategoryPathJSON(current.StoreCategoryPathJson)
	}

	return doc
}

func buildDeleteDoc(ctx context.Context, spuNo string) (*searchDoc, error) {
	current, err := loadSnapshotBySpuNo(ctx, spuNo)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	doc := &searchDoc{
		SpuNo:             strings.TrimSpace(spuNo),
		OnShelfStatusCode: "OFF",
		SourceUpdatedAt:   &now,
		SourceVersion:     uint64(now.UnixMilli()),
		Deleted:           true,
		DeletedAt:         &now,
		StoreCategoryPath: []uint64{},
	}
	if current != nil {
		doc.Title = current.Title
		doc.ShopNo = current.ShopNo
		doc.ShopName = current.ShopName
		doc.CategoryNo = current.CategoryNo
		doc.StoreCategoryId = current.StoreCategoryId
		doc.StoreCategoryL1 = current.StoreCategoryL1
		doc.StoreCategoryL2 = current.StoreCategoryL2
		doc.StoreCategoryPath = parseStoreCategoryPathJSON(current.StoreCategoryPathJson)
		doc.CoverAssetId = current.CoverAssetId
		doc.CoverUrl = current.CoverUrl
		doc.MinPrice = current.MinPrice
		doc.MaxPrice = current.MaxPrice
		doc.StockTotal = current.StockTotal
		doc.SalesCount = current.SalesCount
		doc.AvgScoreX100 = current.AvgScoreX100
		doc.ReviewTotal = current.ReviewTotal
		doc.ShopStatusCode = current.ShopStatusCode
		doc.AttrsJson = current.AttrsJson
		if current.SourceUpdatedAt != nil {
			copied := current.SourceUpdatedAt.Time
			doc.SourceUpdatedAt = &copied
		}
		if current.SourceVersion >= doc.SourceVersion {
			doc.SourceVersion = current.SourceVersion + 1
		}
	}
	return doc, nil
}

func deriveSourceVersion(updatedAt *timestamppb.Timestamp) uint64 {
	// 上游传了更新时间时，优先把它转成毫秒级外部版本号。
	if updatedAt != nil {
		// 只有非零时间才有资格作为真实 source_version，避免零值污染版本序列。
		if ts := updatedAt.AsTime(); !ts.IsZero() {
			return uint64(ts.UnixMilli())
		}
	}
	// 上游没传时间时，用当前时间毫秒值兜底，保证版本号单调向前。
	return uint64(time.Now().UnixMilli())
}

func loadSnapshotBySpuNo(ctx context.Context, spuNo string) (*entity.SearchSpuDoc, error) {
	// 预留实体指针承接数据库扫描结果。
	var row *entity.SearchSpuDoc
	// 按 spu_no 精确读取一条快照，用于写前比对和构造删除文档。
	err := dao.SearchSpuDoc.Ctx(ctx).Where(dao.SearchSpuDoc.Columns().SpuNo, strings.TrimSpace(spuNo)).Limit(1).Scan(&row)
	if err != nil {
		return nil, gerror.Wrap(err, "load search_spu_doc by spu_no failed")
	}
	// 返回命中的快照实体；如果没查到则返回 nil,nil 语义。
	return row, nil
}

func upsertSnapshot(ctx context.Context, doc *searchDoc) (bool, error) {
	// 空文档没有任何写入意义，直接返回 false,nil 表示未应用。
	if doc == nil {
		return false, nil
	}
	// 用 applied 标记本次快照是否真的落库，便于上层判断是否需要继续写 ES。
	var applied bool
	// 用事务把“读当前版本 + 判版本 + 写快照”收成原子操作，防止并发乱序覆盖。
	err := dao.SearchSpuDoc.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 先在事务内读取当前 spu 的快照，确保版本判断基于一致视图。
		current, err := loadSnapshotBySpuNoTx(ctx, tx, doc.SpuNo)
		if err != nil {
			return err
		}
		// 当前快照版本更大或相等时，说明这是过期/重复事件，直接忽略即可。
		if current != nil && current.SourceVersion >= doc.SourceVersion {
			applied = false
			return nil
		}

		// 把业务文档转换成 DO 结构，统一供 insert/update 复用。
		data := buildSnapshotDO(doc)
		// 当前不存在快照时，走 insert 创建新记录。
		if current == nil {
			// 创建新记录时补上 created_at，保证审计时间完整。
			data.CreatedAt = gtime.Now()
			_, err = tx.Model(dao.SearchSpuDoc.Table()).Ctx(ctx).Data(data).Insert()
		} else {
			// 当前已有快照时，按主键更新成新版本内容。
			_, err = tx.Model(dao.SearchSpuDoc.Table()).Ctx(ctx).Where(dao.SearchSpuDoc.Columns().Id, current.Id).Data(data).Update()
		}
		// 快照落库失败时直接终止事务，避免上层误判为已应用。
		if err != nil {
			return gerror.Wrap(err, "upsert search_spu_doc snapshot failed")
		}
		// 快照真正落库成功后，把 applied 标记成 true。
		applied = true
		return nil
	})
	// 返回是否应用成功以及事务结果，供上层决定是否继续写 ES。
	return applied, err
}

func softDeleteSnapshot(ctx context.Context, doc *searchDoc) (bool, error) {
	// 空文档没有任何软删意义，直接返回 false,nil。
	if doc == nil {
		return false, nil
	}
	// 用 applied 标记本次软删是否真正写入快照表。
	var applied bool
	// 用事务收敛软删的版本判断与写库操作，防止并发乱序覆盖新数据。
	err := dao.SearchSpuDoc.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 先在事务里读取当前快照，确保版本比较基于一致数据。
		current, err := loadSnapshotBySpuNoTx(ctx, tx, doc.SpuNo)
		if err != nil {
			return err
		}
		// 当前快照版本更高时，说明当前删除事件已过期，应直接忽略。
		if current != nil && current.SourceVersion >= doc.SourceVersion {
			applied = false
			return nil
		}

		// 把删除文档转换成 DO 结构，供 insert/update 共用。
		data := buildSnapshotDO(doc)
		// 没有现有快照时，直接插入一条删除状态快照。
		if current == nil {
			// 新插入的删除快照也需要补 created_at，保证审计时间完整。
			data.CreatedAt = gtime.Now()
			_, err = tx.Model(dao.SearchSpuDoc.Table()).Ctx(ctx).Data(data).Insert()
		} else {
			// 已存在快照时，按主键更新为软删除状态。
			_, err = tx.Model(dao.SearchSpuDoc.Table()).Ctx(ctx).Where(dao.SearchSpuDoc.Columns().Id, current.Id).Data(data).Update()
		}
		// 数据库软删回写失败时终止事务，避免上层误判。
		if err != nil {
			return gerror.Wrap(err, "soft delete search_spu_doc snapshot failed")
		}
		// 快照软删成功后把 applied 标记为 true。
		applied = true
		return nil
	})
	// 返回是否真正应用删除以及事务执行结果。
	return applied, err
}

func buildSnapshotDO(doc *searchDoc) do.SearchSpuDoc {
	var (
		sourceUpdatedAt *gtime.Time
		deletedAt       *gtime.Time
		attrsJSON       any
	)
	if doc.SourceUpdatedAt != nil {
		sourceUpdatedAt = gtime.NewFromTime(*doc.SourceUpdatedAt)
	}
	if doc.DeletedAt != nil {
		deletedAt = gtime.NewFromTime(*doc.DeletedAt)
	}
	attrsJSON = normalizeOptionalJSONValue(doc.AttrsJson)
	return do.SearchSpuDoc{
		SpuNo:                 doc.SpuNo,
		Title:                 doc.Title,
		ShopNo:                doc.ShopNo,
		ShopName:              doc.ShopName,
		CategoryNo:            doc.CategoryNo,
		StoreCategoryId:       doc.StoreCategoryId,
		StoreCategoryL1:       doc.StoreCategoryL1,
		StoreCategoryL2:       doc.StoreCategoryL2,
		StoreCategoryPathJson: marshalStoreCategoryPathJSON(doc.StoreCategoryPath),
		CoverAssetId:          doc.CoverAssetId,
		CoverUrl:              doc.CoverUrl,
		MinPrice:              doc.MinPrice,
		MaxPrice:              doc.MaxPrice,
		StockTotal:            doc.StockTotal,
		SalesCount:            doc.SalesCount,
		AvgScoreX100:          doc.AvgScoreX100,
		ReviewTotal:           doc.ReviewTotal,
		ShopStatusCode:        doc.ShopStatusCode,
		OnShelfStatusCode:     doc.OnShelfStatusCode,
		AttrsJson:             attrsJSON,
		SourceVersion:         doc.SourceVersion,
		SourceUpdatedAt:       sourceUpdatedAt,
		Deleted:               boolToUint(doc.Deleted),
		DeletedAt:             deletedAt,
		UpdatedAt:             gtime.Now(),
	}
}

func normalizeOptionalJSONValue(raw string) any {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return gdb.Raw("NULL")
	}
	return raw
}

func loadSnapshotBySpuNoTx(ctx context.Context, tx gdb.TX, spuNo string) (*entity.SearchSpuDoc, error) {
	// 预留实体指针承接事务内查询结果。
	var row *entity.SearchSpuDoc
	// 在事务上下文中按 spu_no 精确读取当前快照，供版本判断复用。
	err := tx.Model(dao.SearchSpuDoc.Table()).Ctx(ctx).Where(dao.SearchSpuDoc.Columns().SpuNo, strings.TrimSpace(spuNo)).Limit(1).Scan(&row)
	if err != nil {
		return nil, gerror.Wrap(err, "load search_spu_doc snapshot in tx failed")
	}
	// 返回事务内命中的快照实体；如果不存在则返回 nil,nil。
	return row, nil
}

func listRebuildDocs(ctx context.Context, lastID uint64, limit int) ([]*entity.SearchSpuDoc, error) {
	// 复用快照表列定义，保证扫描条件与表结构一致。
	cols := dao.SearchSpuDoc.Columns()
	// 先构造只读有效快照的扫描模型，并按主键升序分页。
	model := dao.SearchSpuDoc.Ctx(ctx).
		Where(cols.Deleted, 0).
		OrderAsc(cols.Id).
		Limit(limit)
	// 有上次扫描主键时，只拉取其后的数据，避免重复扫描。
	if lastID > 0 {
		model = model.WhereGT(cols.Id, lastID)
	}
	// 预留结果切片承接当前批次快照数据。
	var rows []*entity.SearchSpuDoc
	// 执行批量扫描，供重建任务逐批导入 ES。
	if err := model.Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list rebuild docs failed")
	}
	// 返回当前批次快照列表。
	return rows, nil
}

func searchDocFromEntity(row *entity.SearchSpuDoc) *searchDoc {
	if row == nil {
		return nil
	}
	storeCategoryPath := parseStoreCategoryPathJSON(row.StoreCategoryPathJson)
	doc := &searchDoc{
		SpuNo:             row.SpuNo,
		Title:             row.Title,
		ShopNo:            row.ShopNo,
		ShopName:          row.ShopName,
		CategoryNo:        row.CategoryNo,
		StoreCategoryId:   row.StoreCategoryId,
		StoreCategoryL1:   row.StoreCategoryL1,
		StoreCategoryL2:   row.StoreCategoryL2,
		StoreCategoryPath: storeCategoryPath,
		CoverAssetId:      row.CoverAssetId,
		CoverUrl:          row.CoverUrl,
		MinPrice:          row.MinPrice,
		MaxPrice:          row.MaxPrice,
		StockTotal:        row.StockTotal,
		SalesCount:        row.SalesCount,
		AvgScoreX100:      row.AvgScoreX100,
		ReviewTotal:       row.ReviewTotal,
		ShopStatusCode:    row.ShopStatusCode,
		OnShelfStatusCode: row.OnShelfStatusCode,
		AttrsJson:         row.AttrsJson,
		SourceVersion:     row.SourceVersion,
		Deleted:           row.Deleted > 0,
	}
	if row.SourceUpdatedAt != nil {
		ts := row.SourceUpdatedAt.Time
		doc.SourceUpdatedAt = &ts
	}
	if row.DeletedAt != nil {
		ts := row.DeletedAt.Time
		doc.DeletedAt = &ts
	}
	return doc
}

func hasStoreCategoryPayload(req *v1.UpsertSpuDocReq) bool {
	if req == nil {
		return false
	}
	return req.GetStoreCategoryId() > 0 || req.GetStoreCategoryL1() > 0 || req.GetStoreCategoryL2() > 0 || len(req.GetStoreCategoryPath()) > 0
}

func marshalStoreCategoryPathJSON(path []uint64) string {
	if len(path) == 0 {
		return "[]"
	}
	payload, err := json.Marshal(path)
	if err != nil {
		return "[]"
	}
	return string(payload)
}

func parseStoreCategoryPathJSON(raw string) []uint64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []uint64{}
	}
	var path []uint64
	if err := json.Unmarshal([]byte(raw), &path); err != nil {
		return []uint64{}
	}
	return path
}

func toSpuCard(row *entity.SearchSpuDoc) *v1.SpuCard {
	// 空实体没有可返回内容，直接返回 nil。
	if row == nil {
		return nil
	}
	// 预留更新时间字段，优先使用 source_updated_at，保证前端看到的是源数据时间。
	var updatedAt *timestamppb.Timestamp
	// 有源更新时间时优先使用它，保持响应与 ES 文档一致。
	if row.SourceUpdatedAt != nil {
		updatedAt = timestamppb.New(row.SourceUpdatedAt.Time)
		// 没有源更新时间时，再回退到数据库更新时间兜底。
	} else if row.UpdatedAt != nil {
		updatedAt = timestamppb.New(row.UpdatedAt.Time)
	}
	// 把快照实体完整映射成对外商品卡片结构。
	return &v1.SpuCard{
		SpuNo:             row.SpuNo,
		Title:             row.Title,
		CoverAssetId:      row.CoverAssetId,
		CoverUrl:          row.CoverUrl,
		MinPrice:          row.MinPrice,
		MaxPrice:          row.MaxPrice,
		ShopNo:            row.ShopNo,
		ShopName:          row.ShopName,
		ShopStatusCode:    row.ShopStatusCode,
		OnShelfStatusCode: row.OnShelfStatusCode,
		SalesCount:        row.SalesCount,
		StockTotal:        row.StockTotal,
		AvgScoreX100:      row.AvgScoreX100,
		ReviewTotal:       row.ReviewTotal,
		UpdatedAt:         updatedAt,
	}
}

func toSpuCardFromDoc(doc *searchDoc) *v1.SpuCard {
	// 空文档无法生成卡片，直接返回 nil。
	if doc == nil {
		return nil
	}
	// 预留更新时间字段，尽量复用文档里的 source_updated_at。
	var updatedAt *timestamppb.Timestamp
	// 文档带源更新时间时，直接转成 protobuf 时间戳回给前端。
	if doc.SourceUpdatedAt != nil {
		updatedAt = timestamppb.New(*doc.SourceUpdatedAt)
	}
	// 把统一搜索文档映射成卡片结构，供 ES 主路径直接返回。
	return &v1.SpuCard{
		SpuNo:             doc.SpuNo,
		Title:             doc.Title,
		CoverAssetId:      doc.CoverAssetId,
		CoverUrl:          doc.CoverUrl,
		MinPrice:          doc.MinPrice,
		MaxPrice:          doc.MaxPrice,
		ShopNo:            doc.ShopNo,
		ShopName:          doc.ShopName,
		ShopStatusCode:    doc.ShopStatusCode,
		OnShelfStatusCode: doc.OnShelfStatusCode,
		SalesCount:        doc.SalesCount,
		StockTotal:        doc.StockTotal,
		AvgScoreX100:      doc.AvgScoreX100,
		ReviewTotal:       doc.ReviewTotal,
		UpdatedAt:         updatedAt,
	}
}

func touchKeywordStat(ctx context.Context, query string) {
	// 先清洗搜索词，避免空白查询污染热词统计。
	query = strings.TrimSpace(query)
	// 空查询没有统计价值，直接返回。
	if query == "" {
		return
	}
	// 复用热词表列定义，避免更新字段时写错名称。
	cols := dao.SearchKeywordStat.Columns()
	// 统一记录当前时间，用于 last_searched_at 和 updated_at。
	now := gtime.Now()
	// 先尝试按 keyword 原子累加搜索次数，命中已有热词时走更新路径。
	result, err := dao.SearchKeywordStat.Ctx(ctx).Where(cols.Keyword, query).Data(do.SearchKeywordStat{
		SearchCount:    gdb.Raw(fmt.Sprintf("%s + 1", cols.SearchCount)),
		LastSearchedAt: now,
		UpdatedAt:      now,
	}).Update()
	if err != nil {
		g.Log().Warningf(ctx, "[search-svc] update search_keyword_stat failed, keyword=%s err=%+v", query, err)
		return
	}
	// 读取更新影响行数，用于判断当前关键词是否已存在。
	affected, _ := result.RowsAffected()
	// 已有记录被成功更新时，直接返回，不再走插入逻辑。
	if affected > 0 {
		return
	}
	// 没有命中现有关键词时，插入一条新的热词记录作为初始统计。
	_, err = dao.SearchKeywordStat.Ctx(ctx).Data(do.SearchKeywordStat{
		Keyword:        query,
		SearchCount:    1,
		LastSearchedAt: now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}).Insert()
	if err != nil {
		g.Log().Warningf(ctx, "[search-svc] insert search_keyword_stat failed, keyword=%s err=%+v", query, err)
	}
}

func outboxBackoffDuration(retryCount int) time.Duration {
	// 按重试次数分段退避，避免下游故障时 outbox 高频重试压垮 ES。
	switch {
	// 前两次重试使用最短退避，尽量快速消化偶发抖动。
	case retryCount <= 1:
		return 5 * time.Second
	// 轻度重试阶段拉长到 15 秒，给短暂故障留恢复时间。
	case retryCount <= 3:
		return 15 * time.Second
	// 中度重试阶段拉长到 30 秒，避免失败事件挤占正常流量。
	case retryCount <= 6:
		return 30 * time.Second
	// 高重试阶段拉长到 2 分钟，进一步降低故障放大效应。
	case retryCount <= 10:
		return 2 * time.Minute
	// 极端场景统一退避 5 分钟，防止死循环式重试。
	default:
		return 5 * time.Minute
	}
}

func mergeUniqueStrings(existing []string, incoming []string, limit int) []string {
	// 用 seen 去重，避免热词与 ES 建议词拼接后出现重复项。
	seen := make(map[string]struct{}, len(existing)+len(incoming))
	// 预留合并结果切片，并按 limit 控制最终返回规模。
	merged := make([]string, 0, limit)
	// 先保留 existing 的顺序，优先维持主路径结果的排序稳定性。
	for _, item := range existing {
		// 清洗字符串，避免空白项进入结果。
		normalized := strings.TrimSpace(item)
		if normalized == "" {
			continue
		}
		// 已出现过的词不再重复加入。
		if _, ok := seen[normalized]; ok {
			continue
		}
		// 标记当前词已被选择。
		seen[normalized] = struct{}{}
		// 把词加入结果列表。
		merged = append(merged, normalized)
		// 达到上限后立即返回，避免继续处理无意义数据。
		if len(merged) >= limit {
			return merged
		}
	}
	// 再按顺序补充 incoming 的新词，尽量在不破坏主顺序的前提下补全结果。
	for _, item := range incoming {
		// 清洗字符串，避免空白项进入结果。
		normalized := strings.TrimSpace(item)
		if normalized == "" {
			continue
		}
		// 已出现过的词不再重复加入。
		if _, ok := seen[normalized]; ok {
			continue
		}
		// 标记当前词已被选择。
		seen[normalized] = struct{}{}
		// 把新词补进结果列表。
		merged = append(merged, normalized)
		// 达到上限后立即返回，控制结果大小。
		if len(merged) >= limit {
			return merged
		}
	}
	// 返回合并去重后的建议词列表。
	return merged
}

func boolToUint(value bool) uint {
	// 把布尔值转换成数据库层常用的 0/1 表示，便于快照表写库。
	if value {
		return 1
	}
	// false 一律映射成 0，保持 deleted 等布尔字段存储一致。
	return 0
}

func normalizePageSize(pageSize int32) int {
	// 先把 protobuf 的 int32 转成本地 int，便于后续分页计算。
	size := int(pageSize)
	// 非法页大小回退到默认值 20，避免查 0 条或负数页。
	if size <= 0 {
		return 20
	}
	// 过大的页大小强制裁剪到 100，避免单次查询过重。
	if size > 100 {
		return 100
	}
	// 返回规整后的页大小。
	return size
}

func normalizeSpuNos(spuNos []string) []string {
	// 预留输出切片，保持入参顺序并过滤无效主键。
	out := make([]string, 0, len(spuNos))
	// 用 seen 去重，避免批量查询同一 spu 多次命中。
	seen := make(map[string]struct{}, len(spuNos))
	// 按入参顺序逐个清洗和去重。
	for _, spuNo := range spuNos {
		// 清洗空白字符，避免“空格主键”进入查询。
		normalized := strings.TrimSpace(spuNo)
		// 空主键没有任何业务意义，直接跳过。
		if normalized == "" {
			continue
		}
		// 重复主键不再重复加入，避免查询和返回结果冗余。
		if _, ok := seen[normalized]; ok {
			continue
		}
		// 标记当前 spu_no 已经加入结果。
		seen[normalized] = struct{}{}
		// 把规整后的主键加入输出列表。
		out = append(out, normalized)
	}
	// 返回去空去重后的 spu_no 列表。
	return out
}
