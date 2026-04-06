package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	v1 "github.com/TsingpekTao/shopa/search-svc/api/v1"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var (
	// esEngineOnce 保证 ES 引擎只做一次惰性初始化，避免并发场景重复建连。
	esEngineOnce sync.Once
	// sharedESEngine 缓存已经构造好的 ES 适配器，供整个 logic 包复用。
	sharedESEngine *esEngine
	// sharedESEngineErr 缓存初始化阶段的错误，便于后续请求直接感知失败原因。
	sharedESEngineErr error
)

// esEngine 封装 search.go 依赖的 Elasticsearch 查询、写入与重建能力。
type esEngine struct {
	// client 是 go-elasticsearch 原生客户端，负责底层 HTTP 通讯。
	client *elasticsearch.Client
	// readAlias 是搜索主读别名，线上查询始终优先走它。
	readAlias string
	// writeAlias 是增量写入别名，正常 upsert/delete 都写到它。
	writeAlias string
	// indexPrefix 用于生成版本化索引名和别名前缀。
	indexPrefix string
	// requestTimeout 统一约束单次 ES 请求的最大耗时。
	requestTimeout time.Duration
}

// esHTTPError 把 ES 的 HTTP 状态码和 error.type 抽出来，供分类与冲突识别使用。
type esHTTPError struct {
	// statusCode 是 ES 返回的 HTTP 状态码。
	statusCode int
	// errorType 是 ES error.type，例如 index_not_found_exception。
	errorType string
	// reason 是 ES error.reason，便于日志排障。
	reason string
	// body 兜底保留原始响应体，避免信息丢失。
	body string
}

// Error 实现 error 接口，方便该错误在上层直接打印。
func (e *esHTTPError) Error() string {
	// 先处理空指针场景，避免日志链路再触发 panic。
	if e == nil {
		return "elasticsearch error"
	}
	// 优先输出结构化的状态码、错误类型和原因，方便快速排障。
	if e.errorType != "" || e.reason != "" {
		return fmt.Sprintf("elasticsearch status=%d type=%s reason=%s", e.statusCode, e.errorType, e.reason)
	}
	// 没有解析出结构化错误时，再回退到原始 body。
	if strings.TrimSpace(e.body) != "" {
		return fmt.Sprintf("elasticsearch status=%d body=%s", e.statusCode, e.body)
	}
	// 最后保底只输出状态码，至少不丢最基础的错误信息。
	return fmt.Sprintf("elasticsearch status=%d", e.statusCode)
}

// getElasticsearchEngine 按需初始化并返回 ES 适配器实例。
func getElasticsearchEngine(ctx context.Context) (*esEngine, error) {
	// 只在首次访问时真正构造 ES 客户端，后续请求复用缓存对象。
	esEngineOnce.Do(func() {
		// 把初始化逻辑单独收敛到 newElasticsearchEngine，保持入口清晰。
		sharedESEngine, sharedESEngineErr = newElasticsearchEngine(ctx)
	})
	// 初始化失败时直接返回缓存的错误，让上层走 MySQL 回退。
	if sharedESEngineErr != nil {
		return nil, sharedESEngineErr
	}
	// 初始化成功则返回共享引擎实例。
	return sharedESEngine, nil
}

// newElasticsearchEngine 读取配置并真正构造 ES 客户端。
func newElasticsearchEngine(ctx context.Context) (*esEngine, error) {
	// 从配置里读取地址列表，并兜底到本地单节点地址。
	rawAddresses := g.Cfg().MustGet(ctx, "elasticsearch.addresses", []string{"http://127.0.0.1:9200"}).Strings()
	// 预留清洗后的地址列表，避免配置里混入空字符串。
	addresses := make([]string, 0, len(rawAddresses))
	// 逐个清理地址，保证最终传给客户端的是有效配置。
	for _, item := range rawAddresses {
		// 去掉配置值前后的空白字符。
		normalized := strings.TrimSpace(item)
		// 空地址没有意义，直接跳过。
		if normalized == "" {
			continue
		}
		// 把清洗后的地址加入最终配置切片。
		addresses = append(addresses, normalized)
	}
	// 如果没有任何可用地址，就直接返回初始化错误。
	if len(addresses) == 0 {
		return nil, gerror.New("elasticsearch.addresses is empty")
	}

	// 读取索引前缀，默认使用 search-svc 本地配置值。
	indexPrefix := strings.TrimSpace(g.Cfg().MustGet(ctx, "elasticsearch.indexPrefix", "shopa_search").String())
	// 如果配置里给了空前缀，则回退到默认前缀，避免生成非法索引名。
	if indexPrefix == "" {
		indexPrefix = "shopa_search"
	}

	// 读取请求超时字符串，统一用于 search/write/rebuild 全部请求。
	timeoutRaw := strings.TrimSpace(g.Cfg().MustGet(ctx, "elasticsearch.requestTimeout", "5s").String())
	// 先按 Go duration 语义解析配置值。
	requestTimeout, err := time.ParseDuration(timeoutRaw)
	// 非法超时配置统一回退到 5 秒，避免启动时被配置打死。
	if err != nil || requestTimeout <= 0 {
		requestTimeout = 5 * time.Second
	}

	// 构造底层 HTTP Transport，控制连接复用和响应头超时。
	transport := &http.Transport{
		// 透传环境变量代理配置，方便本地或测试环境调试。
		Proxy: http.ProxyFromEnvironment,
		// 适度提升空闲连接池，减少高频查询的重复建连开销。
		MaxIdleConns: 32,
		// 每个主机保留适量空闲连接，兼顾吞吐与资源占用。
		MaxIdleConnsPerHost: 16,
		// 响应头超时沿用统一请求超时，避免连接长期悬挂。
		ResponseHeaderTimeout: requestTimeout,
	}

	// 用配置构造 go-elasticsearch 客户端。
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		// 配置好的地址列表作为客户端节点集合。
		Addresses: addresses,
		// 用户名来自配置，留空时表示不启用鉴权。
		Username: strings.TrimSpace(g.Cfg().MustGet(ctx, "elasticsearch.username", "").String()),
		// 密码来自配置，留空时表示不启用鉴权。
		Password: strings.TrimSpace(g.Cfg().MustGet(ctx, "elasticsearch.password", "").String()),
		// 自定义 Transport 用于复用连接与控制超时。
		Transport: transport,
	})
	// 客户端构造失败时直接返回错误，让上层走回退。
	if err != nil {
		return nil, gerror.Wrap(err, "init elasticsearch client failed")
	}

	// 返回封装好的 ES 引擎实例，供 search.go 复用。
	return &esEngine{
		// 挂载底层客户端。
		client: client,
		// 读别名固定为 <prefix>_spu_read。
		readAlias: fmt.Sprintf("%s_spu_read", strings.ToLower(indexPrefix)),
		// 写别名固定为 <prefix>_spu_write。
		writeAlias: fmt.Sprintf("%s_spu_write", strings.ToLower(indexPrefix)),
		// 保留索引前缀用于生成版本化索引名。
		indexPrefix: strings.ToLower(indexPrefix),
		// 保留统一超时配置。
		requestTimeout: requestTimeout,
	}, nil
}

// search 执行 search.go 主链路依赖的 ES 商品检索。
func (e *esEngine) search(ctx context.Context, req *v1.SearchProductsReq, cursor searchCursor, pageSize int) (*esSearchResult, error) {
	// 空请求直接返回空结果，避免继续访问 nil 请求字段。
	if req == nil {
		return &esSearchResult{List: []*v1.SpuCard{}}, nil
	}
	// 按请求条件构造 ES 查询体。
	body, err := json.Marshal(e.buildSearchBody(req, cursor, pageSize))
	// 查询体编码失败说明本地逻辑有问题，直接返回错误。
	if err != nil {
		return nil, gerror.Wrap(err, "marshal elasticsearch search body failed")
	}
	// 为本次查询挂上统一的超时控制，避免请求无限阻塞。
	reqCtx, cancel := context.WithTimeout(ctx, e.requestTimeout)
	// 请求结束后及时释放计时器资源。
	defer cancel()

	// 根据当前是否有关键词决定是否追踪分数。
	trackScores := strings.TrimSpace(req.GetQuery()) != ""
	// 组装 Search API 请求。
	esReq := esapi.SearchRequest{
		// 查询固定走读别名，保证线上检索与重建隔离。
		Index: []string{e.readAlias},
		// 请求体承载完整 bool/filter/sort/search_after 结构。
		Body: bytes.NewReader(body),
		// 禁止忽略缺索引，确保上层能感知 index_missing 并打降级指标。
		AllowNoIndices: boolPtr(false),
		// 禁止悄悄忽略不可用索引，避免静默故障。
		IgnoreUnavailable: boolPtr(false),
		// 有关键词时跟踪 score，便于带 _score 的排序稳定返回。
		TrackScores: boolPtr(trackScores),
		// 查询超时与统一引擎超时保持一致。
		Timeout: e.requestTimeout,
	}
	// 发送请求到 ES。
	res, err := esReq.Do(reqCtx, e.client)
	// 网络或上下文错误直接返回，由上层统一分类。
	if err != nil {
		return nil, err
	}

	// 预留响应结构，用于读取 hits 与 sort 值。
	var payload struct {
		// Hits 是 ES 查询结果顶层命中集合。
		Hits struct {
			// Hits 里每一项是一条商品文档命中。
			Hits []struct {
				// Source 承载真实商品文档。
				Source searchDoc `json:"_source"`
				// Sort 承载 search_after 所需的排序游标值。
				Sort []any `json:"sort"`
			} `json:"hits"`
		} `json:"hits"`
	}
	// 解码响应并顺便做统一错误处理。
	if err = decodeESResponse(res, &payload); err != nil {
		return nil, err
	}

	// 先取出全部命中结果，后续再判断是否有下一页。
	hits := payload.Hits.Hits
	// 多取一条即可判断是否还有下一页。
	hasMore := len(hits) > pageSize
	// 有下一页时裁掉最后一条探测记录，避免回给前端。
	if hasMore {
		hits = hits[:pageSize]
	}

	// 预留卡片切片承接本页命中结果。
	list := make([]*v1.SpuCard, 0, len(hits))
	// 逐条把 ES 文档转换成对外 SpuCard。
	for _, hit := range hits {
		// 复制一份文档，避免直接引用循环变量地址。
		doc := hit.Source
		// 追加转换后的卡片结果。
		list = append(list, toSpuCardFromDoc(&doc))
	}

	// 默认没有下一页时返回空游标。
	nextCursor := ""
	// 有下一页且当前页至少有一条有效命中时，生成双模游标。
	if hasMore && len(hits) > 0 {
		// 取当前页最后一条命中的 sort 值作为下一页 search_after。
		lastSort := hits[len(hits)-1].Sort
		// 同时把 ES 和 MySQL 的分页位点一起编码进 opaque cursor。
		nextCursor = encodeCursor(searchCursor{
			// 标记当前游标主要来源于 ES 主检索路径。
			EngineHint: "es",
			// 写入 search_after 值以支持 ES 深分页。
			ESSearchAfter: lastSort,
			// 同时推进 MySQL offset，支持 ES 宕机后的无缝回退。
			MySQLOffset: cursor.MySQLOffset + len(hits),
			// 绑定排序语义，防止用户中途切换 sort 继续沿用旧游标。
			SortCode: normalizeSortCode(req.GetSortCode()),
			// 绑定查询摘要，防止用户换词后继续使用旧游标。
			QueryDigest: buildQueryDigest(req),
		})
	}

	// 返回 ES 主检索结果给 search.go 主流程。
	return &esSearchResult{
		// 返回当前页卡片列表。
		List: list,
		// 返回下一页游标。
		NextCursor: nextCursor,
		// 返回是否仍有更多结果。
		HasMore: hasMore,
	}, nil
}

// suggest 用 ES 前缀匹配补充建议词，作为 MySQL 热词不足时的补召回。
func (e *esEngine) suggest(ctx context.Context, prefix string, limit int) ([]string, error) {
	// 清理前缀空白，避免无意义查询。
	prefix = strings.TrimSpace(prefix)
	// 空前缀没有建议词补全价值，直接返回空结果。
	if prefix == "" || limit <= 0 {
		return []string{}, nil
	}

	// 构造 prefix 补全查询体，优先用标题召回。
	body, err := json.Marshal(map[string]any{
		// 适度多取一点候选，再在本地去重裁剪。
		"size": limit * 2,
		// 查询主体使用 bool，统一挂 deleted/on_shelf 过滤。
		"query": map[string]any{
			// bool 既负责过滤，也负责 should 召回。
			"bool": map[string]any{
				// filter 固定约束只查可售且未删除商品。
				"filter": []any{
					map[string]any{"term": map[string]any{"deleted": false}},
					map[string]any{"terms": map[string]any{"on_shelf_status_code": []string{"ON", "ON_SHELF", "4"}}},
				},
				// should 用前缀语义召回标题、店铺名和属性文本。
				"should": []any{
					map[string]any{"match_phrase_prefix": map[string]any{"title": map[string]any{"query": prefix, "boost": 5}}},
					map[string]any{"match_phrase_prefix": map[string]any{"shop_name": map[string]any{"query": prefix, "boost": 2}}},
					map[string]any{"match_phrase_prefix": map[string]any{"attrs_json": map[string]any{"query": prefix, "boost": 1}}},
				},
				// 至少命中一个 should 才算有效建议词候选。
				"minimum_should_match": 1,
			},
		},
		// 候选结果按销量、更新时间和主键稳定排序。
		"sort": []any{
			map[string]any{"sales_count": "desc"},
			map[string]any{"source_updated_at": "desc"},
			map[string]any{"spu_no": "asc"},
		},
	})
	// 编码失败则直接返回错误。
	if err != nil {
		return nil, gerror.Wrap(err, "marshal elasticsearch suggest body failed")
	}

	// 单独为建议词查询挂统一超时。
	reqCtx, cancel := context.WithTimeout(ctx, e.requestTimeout)
	// 请求结束后及时释放资源。
	defer cancel()

	// 组装 Search API 请求执行建议词补召回。
	res, err := esapi.SearchRequest{
		// 依旧走读别名，保证建议词与主搜索视图一致。
		Index: []string{e.readAlias},
		// 请求体承载 prefix 查询。
		Body: bytes.NewReader(body),
		// 不允许忽略缺索引，缺索引要让上层感知并回退。
		AllowNoIndices: boolPtr(false),
		// 不允许忽略不可用索引。
		IgnoreUnavailable: boolPtr(false),
		// 请求超时沿用统一配置。
		Timeout: e.requestTimeout,
	}.Do(reqCtx, e.client)
	// 网络错误直接返回。
	if err != nil {
		return nil, err
	}

	// 预留响应结构，只关心 _source 标题和店铺名。
	var payload struct {
		// Hits 顶层命中集合。
		Hits struct {
			// Hits 里的每个命中只需要关心 _source。
			Hits []struct {
				// Source 承载命中文档正文。
				Source searchDoc `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	// 统一解码响应并处理 HTTP 错误。
	if err = decodeESResponse(res, &payload); err != nil {
		return nil, err
	}

	// 预留建议词列表。
	keywords := make([]string, 0, limit)
	// 逐条把标题和店铺名压入候选，再由 mergeUniqueStrings 去重。
	for _, hit := range payload.Hits.Hits {
		// 先追加标题，保证商品标题优先级更高。
		keywords = mergeUniqueStrings(keywords, []string{hit.Source.Title}, limit)
		// 已经达到上限时就不用继续扫描了。
		if len(keywords) >= limit {
			break
		}
		// 再把店铺名作为较低优先级补充建议词。
		keywords = mergeUniqueStrings(keywords, []string{hit.Source.ShopName}, limit)
		// 达到上限后立即停止。
		if len(keywords) >= limit {
			break
		}
	}

	// 返回 ES 建议词补召回结果。
	return keywords, nil
}

// mget 按 spu_no 批量读取 ES 文档，并保持入参顺序不变。
func (e *esEngine) mget(ctx context.Context, spuNos []string) ([]*v1.SpuCard, error) {
	// 没有任何主键时直接返回空列表。
	if len(spuNos) == 0 {
		return []*v1.SpuCard{}, nil
	}

	// 组装 mget 请求体，只需要 ids 即可保持顺序。
	body, err := json.Marshal(map[string]any{
		// ids 顺序与入参一致，ES 响应也会按这个顺序返回 docs。
		"ids": spuNos,
	})
	// 编码失败则直接返回。
	if err != nil {
		return nil, gerror.Wrap(err, "marshal elasticsearch mget body failed")
	}

	// 为 mget 挂统一超时控制。
	reqCtx, cancel := context.WithTimeout(ctx, e.requestTimeout)
	// 请求结束后释放资源。
	defer cancel()

	// 发起 ES mget 请求。
	res, err := esapi.MgetRequest{
		// mget 依旧走读别名，保证读取视图一致。
		Index: e.readAlias,
		// 请求体承载 ids 列表。
		Body: bytes.NewReader(body),
	}.Do(reqCtx, e.client)
	// 网络错误直接返回。
	if err != nil {
		return nil, err
	}

	// 预留 mget 响应结构。
	var payload struct {
		// Docs 是按 ids 顺序返回的文档结果数组。
		Docs []struct {
			// Found 标识该主键对应文档是否存在。
			Found bool `json:"found"`
			// Source 承载文档正文。
			Source searchDoc `json:"_source"`
		} `json:"docs"`
	}
	// 解码响应并统一处理错误。
	if err = decodeESResponse(res, &payload); err != nil {
		return nil, err
	}

	// 预留返回卡片切片，容量与 docs 数量保持一致。
	list := make([]*v1.SpuCard, 0, len(payload.Docs))
	// 逐条按响应顺序恢复卡片。
	for _, item := range payload.Docs {
		// 未找到的文档直接跳过，保持与 MySQL 回退语义一致。
		if !item.Found {
			continue
		}
		// 已软删文档不再返回给上层业务。
		if item.Source.Deleted {
			continue
		}
		// 复制文档内容后转换成卡片结构。
		doc := item.Source
		list = append(list, toSpuCardFromDoc(&doc))
	}

	// 返回按入参顺序恢复后的卡片列表。
	return list, nil
}

// indexDoc 用写别名执行正常文档 upsert。
func (e *esEngine) indexDoc(ctx context.Context, doc *searchDoc) error {
	// 正常写链路要求强制写到 alias，避免误写到旧索引。
	return e.indexDocToTarget(ctx, e.writeAlias, true, doc)
}

// softDeleteDoc 通过写别名把软删除状态写入 ES。
func (e *esEngine) softDeleteDoc(ctx context.Context, doc *searchDoc) error {
	// 软删除本质上仍然是一次带 deleted=true 的 upsert。
	return e.indexDocToTarget(ctx, e.writeAlias, true, doc)
}

// indexDocToTarget 把文档写入指定 alias 或版本化索引。
func (e *esEngine) indexDocToTarget(ctx context.Context, target string, requireAlias bool, doc *searchDoc) error {
	// 空文档没有写入意义，直接报参数错误。
	if doc == nil {
		return gerror.New("elasticsearch doc is nil")
	}
	// 主键为空时无法建立文档 ID，直接拒绝写入。
	if strings.TrimSpace(doc.SpuNo) == "" {
		return gerror.New("elasticsearch doc spu_no is empty")
	}

	// 先把统一搜索文档编码成 JSON。
	body, err := json.Marshal(doc)
	// 编码失败说明本地结构有问题，直接返回。
	if err != nil {
		return gerror.Wrap(err, "marshal elasticsearch doc failed")
	}
	// 把 uint64 版本号安全转成 ES SDK 需要的 int 指针。
	version, err := toESVersion(doc.SourceVersion)
	// 版本号非法时直接返回，避免 ES 收到错误外部版本。
	if err != nil {
		return err
	}

	// 为本次写入挂统一超时控制。
	reqCtx, cancel := context.WithTimeout(ctx, e.requestTimeout)
	// 请求结束后及时释放资源。
	defer cancel()

	// 构造 Index API 请求。
	res, err := esapi.IndexRequest{
		// target 可以是写别名，也可以是 rebuild 阶段的新索引实体名。
		Index: target,
		// 文档 ID 固定用 spu_no，保证同一商品天然幂等覆盖。
		DocumentID: doc.SpuNo,
		// 请求体承载完整文档内容。
		Body: bytes.NewReader(body),
		// 正常写链路强制要求 alias，rebuild 明确传 false。
		RequireAlias: boolPtr(requireAlias),
		// 使用 external versioning 防止旧事件覆盖新文档。
		Version: &version,
		// 版本类型固定 external，与方案设计保持一致。
		VersionType: "external",
		// 不要求同步 refresh，避免写放大。
		Refresh: "false",
	}.Do(reqCtx, e.client)
	// 网络错误直接上抛。
	if err != nil {
		return err
	}

	// 统一解析响应错误。
	return decodeESResponse(res, nil)
}

// physicalDeleteDoc 对已经软删超窗口的文档做物理删除。
func (e *esEngine) physicalDeleteDoc(ctx context.Context, spuNo string, sourceVersion uint64) error {
	// 空主键没有删除意义，直接返回参数错误。
	if strings.TrimSpace(spuNo) == "" {
		return gerror.New("elasticsearch physical delete spu_no is empty")
	}
	// 把外部版本号转成 ES SDK 所需的 int 指针值。
	version, err := toESVersion(sourceVersion)
	// 版本号非法时直接返回。
	if err != nil {
		return err
	}

	// 为物理删除挂统一超时控制。
	reqCtx, cancel := context.WithTimeout(ctx, e.requestTimeout)
	// 请求结束后释放资源。
	defer cancel()

	// 发起 Delete API 请求。
	res, err := esapi.DeleteRequest{
		// 物理删除固定作用于写别名指向的当前索引视图。
		Index: e.writeAlias,
		// 删除主键固定为 spu_no。
		DocumentID: strings.TrimSpace(spuNo),
		// 删除也要启用 external versioning，防止旧事件误删新文档。
		Version: &version,
		// 版本类型固定 external。
		VersionType: "external",
		// 正常删除不要求同步 refresh。
		Refresh: "false",
	}.Do(reqCtx, e.client)
	// 网络错误直接返回。
	if err != nil {
		return err
	}

	// 统一解析 ES 响应。
	err = decodeESResponse(res, nil)
	// 没有错误说明物理删除成功。
	if err == nil {
		return nil
	}

	// 尝试把错误断言为结构化 ES HTTP 错误。
	var httpErr *esHTTPError
	// 如果能断言成功且是“文档不存在”，则按幂等成功处理。
	if errors.As(err, &httpErr) && httpErr.statusCode == http.StatusNotFound && strings.Contains(httpErr.errorType, "document_missing_exception") {
		return nil
	}
	// 其它错误按原样返回给上层。
	return err
}

// nextVersionedIndexName 生成 rebuild 用的新索引实体名。
func (e *esEngine) nextVersionedIndexName() string {
	// 用毫秒时间戳确保索引名单调递增且足够直观。
	return fmt.Sprintf("%s_spu_v%d", e.indexPrefix, time.Now().UnixMilli())
}

// createRebuildIndex 创建新的版本化索引并写入映射。
func (e *esEngine) createRebuildIndex(ctx context.Context, indexName string) error {
	// 创建请求需要明确的索引名。
	if strings.TrimSpace(indexName) == "" {
		return gerror.New("rebuild index name is empty")
	}

	// 构造索引 settings 和 mappings，保持与设计方案一致。
	body, err := json.Marshal(map[string]any{
		// settings 负责底层分片、副本与分析器配置。
		"settings": map[string]any{
			// 本地开发先使用单分片，简化资源占用。
			"number_of_shards": 1,
			// 本地默认不配副本，避免单节点集群黄灯噪声。
			"number_of_replicas": 0,
		},
		// mappings 固定约束搜索字段结构。
		"mappings": map[string]any{
			// 一期先关闭动态字段写入，避免脏字段污染索引结构。
			"dynamic": false,
			// properties 显式声明全部可检索字段。
			"properties": map[string]any{
				"spu_no":               map[string]any{"type": "keyword"},
				"title":                map[string]any{"type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart", "fields": map[string]any{"keyword": map[string]any{"type": "keyword", "ignore_above": 256}}},
				"shop_no":              map[string]any{"type": "keyword"},
				"shop_name":            map[string]any{"type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart", "fields": map[string]any{"keyword": map[string]any{"type": "keyword", "ignore_above": 256}}},
				"category_no":          map[string]any{"type": "keyword"},
				"store_category_id":    map[string]any{"type": "long"},
				"store_category_l1":    map[string]any{"type": "long"},
				"store_category_l2":    map[string]any{"type": "long"},
				"store_category_path":  map[string]any{"type": "long"},
				"cover_asset_id":       map[string]any{"type": "long"},
				"cover_url":            map[string]any{"type": "keyword", "ignore_above": 2048},
				"min_price":            map[string]any{"type": "long"},
				"max_price":            map[string]any{"type": "long"},
				"stock_total":          map[string]any{"type": "long"},
				"sales_count":          map[string]any{"type": "long"},
				"avg_score_x100":       map[string]any{"type": "long"},
				"review_total":         map[string]any{"type": "long"},
				"shop_status_code":     map[string]any{"type": "keyword"},
				"on_shelf_status_code": map[string]any{"type": "keyword"},
				"attrs_json":           map[string]any{"type": "text", "analyzer": "ik_max_word", "search_analyzer": "ik_smart"},
				"source_updated_at":    map[string]any{"type": "date"},
				"source_version":       map[string]any{"type": "long"},
				"deleted":              map[string]any{"type": "boolean"},
				"deleted_at":           map[string]any{"type": "date"},
			},
		},
	})
	// 索引定义编码失败则直接返回。
	if err != nil {
		return gerror.Wrap(err, "marshal elasticsearch create index body failed")
	}

	// 为创建索引请求挂统一超时。
	reqCtx, cancel := context.WithTimeout(ctx, e.requestTimeout)
	// 请求结束后释放资源。
	defer cancel()

	// 发起 Create Index 请求。
	res, err := esapi.IndicesCreateRequest{
		// 指定新索引实体名。
		Index: indexName,
		// 请求体承载 settings 和 mappings。
		Body: bytes.NewReader(body),
		// 创建索引也沿用统一超时。
		Timeout: e.requestTimeout,
	}.Do(reqCtx, e.client)
	// 网络错误直接返回。
	if err != nil {
		return err
	}

	// 统一解析响应错误。
	return decodeESResponse(res, nil)
}

// swapRebuildAlias 把线上读写别名原子切到新的版本化索引。
func (e *esEngine) swapRebuildAlias(ctx context.Context, indexName string) error {
	// 新索引名为空时无法切别名，直接返回参数错误。
	if strings.TrimSpace(indexName) == "" {
		return gerror.New("swap rebuild alias index name is empty")
	}

	// 先读取当前读别名指向的索引集合。
	readIndices, err := e.listAliasIndices(ctx, e.readAlias)
	// 读别名查询失败时直接返回，避免误切流量。
	if err != nil {
		return err
	}
	// 再读取当前写别名指向的索引集合。
	writeIndices, err := e.listAliasIndices(ctx, e.writeAlias)
	// 写别名查询失败时直接返回。
	if err != nil {
		return err
	}

	// 预留别名动作切片，后续统一走 _aliases 原子切换。
	actions := make([]map[string]any, 0, len(readIndices)+len(writeIndices)+2)
	// 逐条移除旧读别名绑定。
	for _, item := range readIndices {
		// 先判断该旧索引是否就是新索引本身。
		if item == indexName {
			continue
		}
		// 追加 remove read alias 动作。
		actions = append(actions, map[string]any{"remove": map[string]any{"index": item, "alias": e.readAlias}})
	}
	// 逐条移除旧写别名绑定。
	for _, item := range writeIndices {
		// 如果旧写索引已经是目标索引，就不用重复移除。
		if item == indexName {
			continue
		}
		// 追加 remove write alias 动作。
		actions = append(actions, map[string]any{"remove": map[string]any{"index": item, "alias": e.writeAlias}})
	}
	// 把读别名加到新索引上。
	actions = append(actions, map[string]any{"add": map[string]any{"index": indexName, "alias": e.readAlias}})
	// 把写别名加到新索引上，并显式声明为 write index。
	actions = append(actions, map[string]any{"add": map[string]any{"index": indexName, "alias": e.writeAlias, "is_write_index": true}})

	// 把别名动作编码成请求体。
	body, err := json.Marshal(map[string]any{"actions": actions})
	// 编码失败则直接返回。
	if err != nil {
		return gerror.Wrap(err, "marshal elasticsearch alias actions failed")
	}

	// 为别名切换挂统一超时，避免控制面请求长期阻塞。
	reqCtx, cancel := context.WithTimeout(ctx, e.requestTimeout)
	// 请求结束后释放资源。
	defer cancel()

	// 发起 Update Aliases 请求实现原子切换。
	res, err := esapi.IndicesUpdateAliasesRequest{
		// 请求体承载 remove/add 动作集合。
		Body: bytes.NewReader(body),
		// 别名切换超时沿用统一配置。
		Timeout: e.requestTimeout,
	}.Do(reqCtx, e.client)
	// 网络错误直接返回。
	if err != nil {
		return err
	}

	// 统一解析响应错误。
	return decodeESResponse(res, nil)
}

// listAliasIndices 查询某个 alias 当前绑定到哪些索引上。
func (e *esEngine) listAliasIndices(ctx context.Context, alias string) ([]string, error) {
	// 别名为空时没有查询意义，直接返回空结果。
	if strings.TrimSpace(alias) == "" {
		return []string{}, nil
	}

	// 为 alias 查询挂统一超时。
	reqCtx, cancel := context.WithTimeout(ctx, e.requestTimeout)
	// 请求结束后释放资源。
	defer cancel()

	// 发起 Get Alias 请求读取当前绑定关系。
	res, err := esapi.IndicesGetAliasRequest{
		// 只查询当前 alias 名称。
		Name: []string{alias},
		// 缺 alias 时不要静默吞掉，交给我们自己区分 404 语义。
		AllowNoIndices: boolPtr(false),
		// 同样不忽略不可用索引。
		IgnoreUnavailable: boolPtr(false),
	}.Do(reqCtx, e.client)
	// 网络错误直接返回。
	if err != nil {
		return nil, err
	}

	// 如果 alias 不存在，ES 会返回 404，这里按“当前无绑定”处理。
	if res.StatusCode == http.StatusNotFound {
		// 及时关闭响应体，避免连接泄漏。
		_ = res.Body.Close()
		// 没有绑定索引时直接返回空切片。
		return []string{}, nil
	}

	// 预留 map 承接 alias 响应，键就是索引名。
	payload := make(map[string]json.RawMessage)
	// 解码响应并统一处理其它 HTTP 错误。
	if err = decodeESResponse(res, &payload); err != nil {
		return nil, err
	}

	// 预留索引名切片承接 map 键。
	indices := make([]string, 0, len(payload))
	// 遍历响应 map，把索引名收集出来。
	for indexName := range payload {
		// 忽略空键，避免脏数据污染结果。
		if strings.TrimSpace(indexName) == "" {
			continue
		}
		// 把当前索引名加入结果集合。
		indices = append(indices, indexName)
	}
	// 为了让 alias 切换日志和行为稳定，统一按字典序排序。
	sort.Strings(indices)
	// 返回 alias 绑定的索引集合。
	return indices, nil
}

// buildSearchBody 按当前请求构造 ES 检索 DSL。
func (e *esEngine) buildSearchBody(req *v1.SearchProductsReq, cursor searchCursor, pageSize int) map[string]any {
	// 先构造基础 filter，强制过滤掉已删和不可售文档。
	filters := []any{
		map[string]any{"term": map[string]any{"deleted": false}},
		map[string]any{"terms": map[string]any{"on_shelf_status_code": []string{"ON", "ON_SHELF", "4"}}},
	}
	// 有店铺过滤条件时追加精确过滤。
	if shopNo := strings.TrimSpace(req.GetShopNo()); shopNo != "" {
		filters = append(filters, map[string]any{"term": map[string]any{"shop_no": shopNo}})
	}
	// 有类目过滤条件时追加精确过滤。
	if categoryNo := strings.TrimSpace(req.GetCategoryNo()); categoryNo != "" {
		filters = append(filters, map[string]any{"term": map[string]any{"category_no": categoryNo}})
	}
	if req.GetStoreCategoryId() > 0 {
		switch req.GetStoreCategoryLevel() {
		case 1:
			filters = append(filters, map[string]any{"term": map[string]any{"store_category_path": req.GetStoreCategoryId()}})
		case 2:
			filters = append(filters, map[string]any{"term": map[string]any{"store_category_id": req.GetStoreCategoryId()}})
		default:
			filters = append(filters, map[string]any{"term": map[string]any{"store_category_id": req.GetStoreCategoryId()}})
		}
	}

	// 清洗关键词，决定是否启用全文检索路径。
	queryText := strings.TrimSpace(req.GetQuery())
	// 预留 query 对象，后续根据是否有关键词分支构造。
	var query map[string]any
	// 有关键词时启用 multi_match + should 精确命中。
	if queryText != "" {
		// 构造带全文召回与精确编号命中的 bool 查询。
		query = map[string]any{
			"bool": map[string]any{
				"filter": filters,
				"must": []any{
					map[string]any{"multi_match": map[string]any{
						"query":       queryText,
						"type":        "best_fields",
						"fields":      []string{"title^5", "shop_name^2", "attrs_json^1"},
						"tie_breaker": 0.3,
					}},
				},
				"should": []any{
					map[string]any{"term": map[string]any{"spu_no": map[string]any{"value": queryText, "boost": 10}}},
				},
				"minimum_should_match": 0,
			},
		}
	} else {
		// 无关键词时只保留 filter，按列表页模式检索即可。
		query = map[string]any{
			"bool": map[string]any{
				"filter": filters,
			},
		}
	}

	// 构造完整的 Search API 请求体。
	body := map[string]any{
		// 多取一条用于判断是否还有下一页。
		"size": pageSize + 1,
		// 真实查询体。
		"query": query,
		// 排序规则固定追加 spu_no 作为 tie-breaker。
		"sort": e.buildSort(normalizeSortCode(req.GetSortCode()), queryText != ""),
	}
	// 有 search_after 时追加深分页位点。
	if len(cursor.ESSearchAfter) > 0 {
		body["search_after"] = cursor.ESSearchAfter
	}
	// 返回最终请求体。
	return body
}

// buildSort 统一生成 search_after 兼容的稳定排序规则。
func (e *esEngine) buildSort(sortCode string, hasQuery bool) []any {
	// 根据 sortCode 切换不同排序分支。
	switch sortCode {
	case "sales_desc":
		return []any{
			map[string]any{"sales_count": "desc"},
			map[string]any{"_score": "desc"},
			map[string]any{"source_updated_at": "desc"},
			map[string]any{"spu_no": "asc"},
		}
	case "price_asc":
		return []any{
			map[string]any{"min_price": "asc"},
			map[string]any{"_score": "desc"},
			map[string]any{"source_updated_at": "desc"},
			map[string]any{"spu_no": "asc"},
		}
	case "price_desc":
		return []any{
			map[string]any{"min_price": "desc"},
			map[string]any{"_score": "desc"},
			map[string]any{"source_updated_at": "desc"},
			map[string]any{"spu_no": "asc"},
		}
	default:
		if hasQuery {
			return []any{
				map[string]any{"_score": "desc"},
				map[string]any{"sales_count": "desc"},
				map[string]any{"source_updated_at": "desc"},
				map[string]any{"spu_no": "asc"},
			}
		}
		return []any{
			map[string]any{"sales_count": "desc"},
			map[string]any{"source_updated_at": "desc"},
			map[string]any{"spu_no": "asc"},
		}
	}
}

// classifyESError 把 ES 错误归类成 metrics 约定的 reason code。
func classifyESError(err error) string {
	// 空错误理论上不该出现，但这里兜底到 unavailable。
	if err == nil {
		return esReasonUnavailable
	}
	// 版本冲突优先单独识别，便于写冲突指标统计。
	if isESConflict(err) {
		return esReasonConflict
	}
	// deadline exceeded 归类为 timeout。
	if errors.Is(err, context.DeadlineExceeded) {
		return esReasonTimeout
	}
	// 可识别的 net.Error 且带 Timeout 语义时，也归到 timeout。
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return esReasonTimeout
	}
	// 尝试把错误还原为结构化 ES HTTP 错误。
	var httpErr *esHTTPError
	if errors.As(err, &httpErr) {
		if httpErr.statusCode == http.StatusNotFound && strings.Contains(httpErr.errorType, "index_not_found_exception") {
			return esReasonIndexMiss
		}
		if httpErr.statusCode == http.StatusRequestTimeout || httpErr.statusCode == http.StatusGatewayTimeout {
			return esReasonTimeout
		}
		if httpErr.statusCode >= http.StatusInternalServerError {
			return esReason5xx
		}
		return esReasonUnavailable
	}
	// 文本里能看出缺索引时，同样归类为 index_missing。
	if strings.Contains(err.Error(), "index_not_found_exception") {
		return esReasonIndexMiss
	}
	// 其它情况统一归到 unavailable。
	return esReasonUnavailable
}

// isESConflict 判断当前错误是否属于 external version 冲突。
func isESConflict(err error) bool {
	// 空错误肯定不是冲突。
	if err == nil {
		return false
	}
	// 先尝试断言结构化 HTTP 错误。
	var httpErr *esHTTPError
	if errors.As(err, &httpErr) {
		if httpErr.statusCode == http.StatusConflict {
			return true
		}
		if strings.Contains(httpErr.errorType, "version_conflict_engine_exception") {
			return true
		}
		if strings.Contains(httpErr.reason, "version_conflict_engine_exception") {
			return true
		}
	}
	// 再用字符串兜底识别，兼容被上层 wrap 的场景。
	return strings.Contains(err.Error(), "version_conflict_engine_exception")
}

// decodeESResponse 统一处理 ES HTTP 错误并在成功时解码响应体。
func decodeESResponse(res *esapi.Response, target any) error {
	// 空响应说明底层通信链路异常。
	if res == nil {
		return gerror.New("elasticsearch response is nil")
	}
	// 无论成功还是失败，最终都要关闭响应体释放连接。
	defer res.Body.Close()
	// 先把响应体完整读出来，便于失败时构造结构化错误。
	body, err := io.ReadAll(res.Body)
	// 读取响应体失败时直接返回。
	if err != nil {
		return gerror.Wrap(err, "read elasticsearch response body failed")
	}
	// ES 返回错误状态码时，尝试抽取 error.type 和 reason。
	if res.IsError() {
		// 预留错误响应结构，提取 error.type 和 error.reason。
		var payload struct {
			// Error 承载 ES 标准错误体。
			Error struct {
				// Type 是错误类型字符串。
				Type string `json:"type"`
				// Reason 是错误原因描述。
				Reason string `json:"reason"`
			} `json:"error"`
		}
		// 尽力解析错误响应，不让 JSON 解析失败阻塞真实错误返回。
		_ = json.Unmarshal(body, &payload)
		// 返回结构化 ES HTTP 错误，供上层分类与冲突识别使用。
		return &esHTTPError{
			statusCode: res.StatusCode,
			errorType:  strings.TrimSpace(payload.Error.Type),
			reason:     strings.TrimSpace(payload.Error.Reason),
			body:       strings.TrimSpace(string(body)),
		}
	}
	// 成功响应且调用方不关心 body 时，直接返回 nil。
	if target == nil {
		return nil
	}
	// 空成功体无需反序列化，直接返回。
	if len(body) == 0 {
		return nil
	}
	// 把成功响应体解码到目标结构。
	if err = json.Unmarshal(body, target); err != nil {
		return gerror.Wrap(err, "decode elasticsearch response body failed")
	}
	// 解码成功返回 nil。
	return nil
}

// toESVersion 把 uint64 版本号安全转换成 ES SDK 需要的 int。
func toESVersion(value uint64) (int, error) {
	// 计算当前平台 int 可表示的最大值。
	maxInt := uint64(^uint(0) >> 1)
	// 超出 int 范围时直接报错，避免发生截断。
	if value > maxInt {
		return 0, gerror.Newf("elasticsearch version overflow: %d", value)
	}
	// 在安全范围内时直接转换即可。
	return int(value), nil
}

// boolPtr 统一构造 *bool，简化 ES 请求参数赋值。
func boolPtr(value bool) *bool {
	// 返回局部变量地址即可，Go 会自动处理逃逸。
	return &value
}
