package search

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	v1 "github.com/TsingpekTao/shopa/search-svc/api/v1"
	"github.com/TsingpekTao/shopa/search-svc/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// sSearch 定义搜索逻辑实现结构。
type sSearch struct{}

// docRecord 定义搜索文档内存对象。
type docRecord struct {
	SpuNo             string
	Title             string
	ShopNo            string
	CategoryNo        string
	MinPrice          uint64
	MaxPrice          uint64
	StockTotal        uint64
	AvgScoreX100      uint64
	ReviewTotal       uint64
	ShopStatusCode    string
	OnShelfStatusCode string
	ShopName          string
	CoverAssetID      uint64
	CoverURL          string
	SalesCount        uint64
	UpdatedAt         time.Time
}

var (
	// searchMu 保护文档索引并发安全。
	searchMu sync.Mutex
	// docs 保存 spu_no -> docRecord 映射。
	docs = make(map[string]*docRecord)
)

// New 创建搜索逻辑对象。
func New() *sSearch {
	// 返回新的逻辑实现实例。
	return &sSearch{}
}

func init() {
	// 在初始化阶段注册搜索服务实现。
	service.RegisterSearch(New())
}

// SearchProducts 根据关键字与过滤条件返回商品卡片列表。
func (s *sSearch) SearchProducts(ctx context.Context, req *v1.SearchProductsReq) (*v1.SearchProductsRes, error) {
	// 防止未使用上下文参数导致编译告警。
	_ = ctx
	// 请求为空时返回空结果。
	if req == nil {
		// 返回空列表避免空指针风险。
		return &v1.SearchProductsRes{List: []*v1.SpuCard{}}, nil
	}
	// 进入临界区读取索引文档。
	searchMu.Lock()
	// 函数返回前释放锁。
	defer searchMu.Unlock()
	// 初始化匹配结果切片。
	matched := make([]*docRecord, 0, len(docs))
	// 归一化查询关键字。
	query := strings.ToLower(strings.TrimSpace(req.GetQuery()))
	// 遍历内存文档并执行轻量过滤。
	for _, d := range docs {
		// 当指定店铺筛选且不匹配时跳过。
		if req.GetShopNo() != "" && d.ShopNo != req.GetShopNo() {
			// 跳过不匹配店铺文档。
			continue
		}
		// 当指定类目筛选且不匹配时跳过。
		if req.GetCategoryNo() != "" && d.CategoryNo != req.GetCategoryNo() {
			// 跳过不匹配类目文档。
			continue
		}
		// 当关键字为空时直接视为匹配。
		if query == "" {
			// 追加到匹配列表。
			matched = append(matched, d)
			// 继续下一个文档。
			continue
		}
		// 标题或 spu_no 命中关键字时视为匹配。
		if strings.Contains(strings.ToLower(d.Title), query) || strings.Contains(strings.ToLower(d.SpuNo), query) {
			// 追加到匹配列表。
			matched = append(matched, d)
		}
	}
	// 按更新时间降序排序，保证结果稳定。
	sort.SliceStable(matched, func(i, j int) bool { return matched[i].UpdatedAt.After(matched[j].UpdatedAt) })
	// 解析游标作为起始偏移。
	start := parseCursor(0, req.GetNextCursor())
	// 非法游标回退到 0。
	if start < 0 || start > len(matched) {
		// 重置为首页起始位置。
		start = 0
	}
	// 设置默认分页大小。
	pageSize := int(req.GetPageSize())
	// 当分页参数非法时使用默认值 20。
	if pageSize <= 0 {
		// 将分页大小回退到默认值。
		pageSize = 20
	}
	// 当分页大小过大时做上限保护。
	if pageSize > 100 {
		// 将分页大小裁剪到 100。
		pageSize = 100
	}
	// 计算当前页结束位置。
	end := start + pageSize
	// 当结束位置超出数组长度时裁剪到上界。
	if end > len(matched) {
		// 防止切片越界。
		end = len(matched)
	}
	// 初始化响应卡片列表。
	list := make([]*v1.SpuCard, 0, end-start)
	// 遍历窗口内文档并转换为响应对象。
	for i := start; i < end; i++ {
		// 读取当前文档对象。
		d := matched[i]
		// 追加转换后的卡片对象。
		list = append(list, toSpuCard(d))
	}
	// 判断是否有下一页。
	hasMore := end < len(matched)
	// 初始化下一页游标。
	nextCursor := ""
	// 当有下一页时设置 next_cursor。
	if hasMore {
		// 将下一个偏移量编码为游标字符串。
		nextCursor = fmt.Sprintf("%d", end)
	}
	// 返回检索结果。
	return &v1.SearchProductsRes{List: list, NextCursor: nextCursor, HasMore: hasMore, Partial: false}, nil
}

// SuggestKeywords 根据前缀返回建议词。
func (s *sSearch) SuggestKeywords(ctx context.Context, req *v1.SuggestKeywordsReq) (*v1.SuggestKeywordsRes, error) {
	// 防止未使用上下文参数导致编译告警。
	_ = ctx
	// 请求为空时返回空列表。
	if req == nil {
		// 返回空建议词集合。
		return &v1.SuggestKeywordsRes{Keywords: []string{}}, nil
	}
	// 归一化前缀字符串。
	prefix := strings.ToLower(strings.TrimSpace(req.GetPrefix()))
	// 计算返回上限，默认 10 条。
	limit := int(req.GetLimit())
	// 当上限参数非法时使用默认值。
	if limit <= 0 {
		// 将限制数量设置为默认值。
		limit = 10
	}
	// 进入临界区读取文档集合。
	searchMu.Lock()
	// 函数返回前释放锁。
	defer searchMu.Unlock()
	// 初始化建议词切片。
	out := make([]string, 0, limit)
	// 初始化去重集合。
	seen := make(map[string]struct{})
	// 遍历文档标题生成建议词。
	for _, d := range docs {
		// 标题为空时跳过。
		if strings.TrimSpace(d.Title) == "" {
			// 跳过空标题文档。
			continue
		}
		// 前缀不匹配时跳过。
		if prefix != "" && !strings.Contains(strings.ToLower(d.Title), prefix) {
			// 跳过不匹配前缀文档。
			continue
		}
		// 已出现的标题不重复加入。
		if _, ok := seen[d.Title]; ok {
			// 跳过重复建议词。
			continue
		}
		// 记录建议词去重标记。
		seen[d.Title] = struct{}{}
		// 追加建议词到结果列表。
		out = append(out, d.Title)
		// 达到上限后提前结束循环。
		if len(out) >= limit {
			// 提前终止迭代以控制开销。
			break
		}
	}
	// 返回建议词结果。
	return &v1.SuggestKeywordsRes{Keywords: out}, nil
}

// BatchGetSpuCards 批量读取 spu 卡片信息。
func (s *sSearch) BatchGetSpuCards(ctx context.Context, req *v1.BatchGetSpuCardsReq) (*v1.BatchGetSpuCardsRes, error) {
	// 防止未使用上下文参数导致编译告警。
	_ = ctx
	// 请求为空时返回空列表。
	if req == nil {
		// 返回空卡片集合。
		return &v1.BatchGetSpuCardsRes{List: []*v1.SpuCard{}}, nil
	}
	// 进入临界区读取文档映射。
	searchMu.Lock()
	// 函数返回前释放锁。
	defer searchMu.Unlock()
	// 初始化响应切片。
	list := make([]*v1.SpuCard, 0, len(req.GetSpuNos()))
	// 遍历请求中的 spu_no 列表。
	for _, spuNo := range req.GetSpuNos() {
		// 按 spu_no 查询文档。
		if d, ok := docs[spuNo]; ok {
			// 追加命中的卡片对象。
			list = append(list, toSpuCard(d))
		}
	}
	// 返回批量卡片结果。
	return &v1.BatchGetSpuCardsRes{List: list}, nil
}

// UpsertSpuDoc 写入或更新搜索文档。
func (s *sSearch) UpsertSpuDoc(ctx context.Context, req *v1.UpsertSpuDocReq) (*v1.UpsertSpuDocRes, error) {
	// 防止未使用上下文参数导致编译告警。
	_ = ctx
	// 校验关键参数必须存在。
	if req == nil || strings.TrimSpace(req.GetSpuNo()) == "" {
		// 参数非法时返回错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}
	// 进入临界区更新文档映射。
	searchMu.Lock()
	// 函数返回前释放锁。
	defer searchMu.Unlock()
	// 计算更新时间，默认使用当前时间。
	updatedAt := time.Now()
	// 当请求带更新时间时使用请求值。
	if req.GetUpdatedAt() != nil {
		// 转换 protobuf 时间为 Go 时间。
		updatedAt = req.GetUpdatedAt().AsTime()
	}
	// 写入或覆盖文档记录。
	docs[req.GetSpuNo()] = &docRecord{
		SpuNo:             req.GetSpuNo(),
		Title:             req.GetTitle(),
		ShopNo:            req.GetShopNo(),
		CategoryNo:        req.GetCategoryNo(),
		MinPrice:          req.GetMinPrice(),
		MaxPrice:          req.GetMaxPrice(),
		StockTotal:        req.GetStockTotal(),
		AvgScoreX100:      req.GetAvgScoreX100(),
		ReviewTotal:       req.GetReviewTotal(),
		ShopStatusCode:    req.GetShopStatusCode(),
		OnShelfStatusCode: req.GetOnShelfStatusCode(),
		ShopName:          req.GetShopName(),
		CoverAssetID:      req.GetCoverAssetId(),
		CoverURL:          req.GetCoverUrl(),
		SalesCount:        req.GetSalesCount(),
		UpdatedAt:         updatedAt,
	}
	// 返回写入成功结果。
	return &v1.UpsertSpuDocRes{Ok: true}, nil
}

// DeleteSpuDoc 删除指定 spu 的搜索文档。
func (s *sSearch) DeleteSpuDoc(ctx context.Context, req *v1.DeleteSpuDocReq) (*v1.DeleteSpuDocRes, error) {
	// 防止未使用上下文参数导致编译告警。
	_ = ctx
	// 校验 spu_no 必填。
	if req == nil || strings.TrimSpace(req.GetSpuNo()) == "" {
		// 参数非法时返回错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "spu_no is required")
	}
	// 进入临界区删除文档对象。
	searchMu.Lock()
	// 函数返回前释放锁。
	defer searchMu.Unlock()
	// 从映射中移除对应 spu_no。
	delete(docs, req.GetSpuNo())
	// 返回删除成功结果。
	return &v1.DeleteSpuDocRes{Ok: true}, nil
}

// RebuildIndex 触发索引重建任务。
func (s *sSearch) RebuildIndex(ctx context.Context, req *v1.RebuildIndexReq) (*v1.RebuildIndexRes, error) {
	// 防止未使用上下文参数导致编译告警。
	_ = ctx
	// 生成重建任务号。
	jobNo := fmt.Sprintf("RB%s", time.Now().Format("20060102150405.000000"))
	// 基础版本直接返回任务已受理状态。
	return &v1.RebuildIndexRes{JobNo: jobNo, StatusCode: "ACCEPTED"}, nil
}

// toSpuCard 将文档对象转换为响应卡片对象。
func toSpuCard(d *docRecord) *v1.SpuCard {
	// 空对象直接返回 nil。
	if d == nil {
		// 防御式返回避免空指针。
		return nil
	}
	// 返回组装后的卡片对象。
	return &v1.SpuCard{
		SpuNo:             d.SpuNo,
		Title:             d.Title,
		CoverAssetId:      d.CoverAssetID,
		CoverUrl:          d.CoverURL,
		MinPrice:          d.MinPrice,
		MaxPrice:          d.MaxPrice,
		ShopNo:            d.ShopNo,
		ShopName:          d.ShopName,
		ShopStatusCode:    d.ShopStatusCode,
		OnShelfStatusCode: d.OnShelfStatusCode,
		SalesCount:        d.SalesCount,
		StockTotal:        d.StockTotal,
		AvgScoreX100:      d.AvgScoreX100,
		ReviewTotal:       d.ReviewTotal,
		UpdatedAt:         timestamppb.New(d.UpdatedAt),
	}
}

// parseCursor 将游标字符串转换为偏移量。
func parseCursor(def int, raw string) int {
	// 清理游标字符串首尾空白。
	raw = strings.TrimSpace(raw)
	// 空游标直接返回默认值。
	if raw == "" {
		// 返回默认偏移量。
		return def
	}
	// 声明偏移量变量。
	var v int
	// 执行字符串解析。
	_, err := fmt.Sscanf(raw, "%d", &v)
	// 解析失败时返回默认值。
	if err != nil {
		// 返回默认偏移量。
		return def
	}
	// 返回解析后的偏移值。
	return v
}

