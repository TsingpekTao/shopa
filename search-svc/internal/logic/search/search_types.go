package search

import (
	"time"

	v1 "github.com/TsingpekTao/shopa/search-svc/api/v1"
)

const (
	// 标记 ES 查询超时，便于 fallback 指标按原因维度告警。
	esReasonTimeout = "timeout"
	// 标记 ES 返回 5xx，便于区分集群内部故障。
	esReason5xx = "5xx"
	// 标记 ES 整体不可用，便于统一处理网络与初始化异常。
	esReasonUnavailable = "unavailable"
	// 标记索引或别名缺失，便于快速识别发布或重建异常。
	esReasonIndexMiss = "index_missing"
	// 标记 ES 外部版本冲突，便于识别乱序事件被拒绝的场景。
	esReasonConflict = "version_conflict"

	// outbox 的 pending 状态表示事件等待 worker 消费。
	outboxStatusPending = 0
	// outbox 的 done 状态表示事件已经成功落到 ES。
	outboxStatusDone = 1
	// outbox 的 dead 状态表示事件重试耗尽，需要人工介入。
	outboxStatusDead = 2

	// upsert 事件用于通知 worker 把快照同步到 ES。
	outboxEventTypeUpsert = "SEARCH_DOC_UPSERT"
	// delete 事件用于通知 worker 把软删状态同步到 ES。
	outboxEventTypeDelete = "SEARCH_DOC_DELETE"
)

type searchCursor struct {
	// engine_hint 告诉服务端本次游标主要来自哪个检索引擎，便于调试与优先解析。
	EngineHint string `json:"engine_hint,omitempty"`
	// es_search_after 保存 ES 深分页排序位点，避免 from/size 带来的深分页性能问题。
	ESSearchAfter []any `json:"es_search_after,omitempty"`
	// mysql_offset 保存 MySQL 回退路径的偏移量，保证 ES 宕机时还能继续翻页。
	MySQLOffset int `json:"mysql_offset,omitempty"`
	// sort_code 绑定排序语义，防止用户翻页中途切换排序导致游标失真。
	SortCode string `json:"sort_code,omitempty"`
	// query_digest 绑定查询快照，防止用户换搜索词后继续沿用旧游标。
	QueryDigest string `json:"query_digest,omitempty"`
}

type searchDoc struct {
	// spu_no 是 ES 与 MySQL 快照共同使用的主键。
	SpuNo string `json:"spu_no"`
	// title 是主召回字段，承载最高搜索权重。
	Title string `json:"title"`
	// shop_no 用于店铺过滤和展示。
	ShopNo string `json:"shop_no"`
	// shop_name 既用于展示，也参与低权重召回。
	ShopName string `json:"shop_name"`
	// category_no 用于类目过滤。
	CategoryNo string `json:"category_no"`
	// store_category_id 是绑定的叶子店内分类 ID。
	StoreCategoryId uint64 `json:"store_category_id"`
	// store_category_l1 是一级店内分类 ID。
	StoreCategoryL1 uint64 `json:"store_category_l1"`
	// store_category_l2 是二级店内分类 ID。
	StoreCategoryL2 uint64 `json:"store_category_l2"`
	// store_category_path 承载层级路径，用于一级分类直接召回。
	StoreCategoryPath []uint64 `json:"store_category_path,omitempty"`
	// cover_asset_id 用于关联媒体资源。
	CoverAssetId uint64 `json:"cover_asset_id"`
	// cover_url 用于直接回显商品封面图。
	CoverUrl string `json:"cover_url"`
	// min_price 用于排序与价格展示。
	MinPrice uint64 `json:"min_price"`
	// max_price 用于价格区间展示。
	MaxPrice uint64 `json:"max_price"`
	// stock_total 用于库存信息展示。
	StockTotal uint64 `json:"stock_total"`
	// sales_count 用于销量排序与展示。
	SalesCount uint64 `json:"sales_count"`
	// avg_score_x100 用于评分展示，避免浮点误差。
	AvgScoreX100 uint64 `json:"avg_score_x100"`
	// review_total 用于评论总数展示。
	ReviewTotal uint64 `json:"review_total"`
	// shop_status_code 表示店铺经营状态，影响过滤逻辑。
	ShopStatusCode string `json:"shop_status_code"`
	// on_shelf_status_code 表示商品上架状态，影响是否可被搜索命中。
	OnShelfStatusCode string `json:"on_shelf_status_code"`
	// attrs_json 保存可检索属性文本，补充标题召回。
	AttrsJson string `json:"attrs_json"`
	// source_updated_at 记录源数据更新时间，参与排序与版本治理。
	SourceUpdatedAt *time.Time `json:"source_updated_at,omitempty"`
	// source_version 是写入 ES 的外部版本号，用于抵御乱序覆盖。
	SourceVersion uint64 `json:"source_version"`
	// deleted 表示文档是否处于软删除状态。
	Deleted bool `json:"deleted"`
	// deleted_at 记录软删除发生时间，用于后续物理删除窗口判断。
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type esSearchResult struct {
	// list 是 ES 主搜索路径返回的商品卡片。
	List []*v1.SpuCard
	// next_cursor 是下一页游标，内部同时携带 ES 与 MySQL 位点。
	NextCursor string
	// has_more 表示是否还有下一页。
	HasMore bool
}

type mysqlSearchResult struct {
	// list 是 MySQL 回退路径返回的商品卡片。
	List []*v1.SpuCard
	// next_cursor 是下一页游标，主要承载 MySQL offset，同时保留双模结构。
	NextCursor string
	// has_more 表示是否还有下一页。
	HasMore bool
}

type outboxPayload struct {
	// event_type 标识本次同步是 upsert 还是 delete。
	EventType string `json:"event_type"`
	// doc 承载需要写入或删除的搜索文档快照。
	Doc *searchDoc `json:"doc,omitempty"`
}
