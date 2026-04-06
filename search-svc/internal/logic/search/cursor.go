package search

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	v1 "github.com/TsingpekTao/shopa/search-svc/api/v1"
)

func encodeCursor(cursor searchCursor) string {
	// 先把双模游标结构编码成 JSON，保证 ES/MySQL 两套翻页信息一起透传给前端。
	payload, err := json.Marshal(cursor)
	// 如果编码失败，则返回空游标让调用方退回首页，避免把半成品游标继续发出去。
	if err != nil {
		return ""
	}
	// 再把 JSON 压成 URL 安全的 Base64，避免前端传参时被转义字符污染。
	return base64.RawURLEncoding.EncodeToString(payload)
}

func decodeCursor(raw string, req *v1.SearchProductsReq) (searchCursor, error) {
	// 先按当前请求构造默认游标，确保坏游标场景也能安全回退到第一页。
	cursor := searchCursor{SortCode: normalizeSortCode(req.GetSortCode()), QueryDigest: buildQueryDigest(req)}
	// 统一裁掉游标前后空白，避免前端拼接参数时带入多余空格。
	raw = strings.TrimSpace(raw)
	// 空游标直接返回默认值，表示从首页开始查。
	if raw == "" {
		return cursor, nil
	}
	// 兼容旧版纯 offset 游标，避免老客户端升级前直接报错。
	if offset, err := strconv.Atoi(raw); err == nil && offset >= 0 {
		// 旧版游标只会影响 MySQL 回退路径，所以只填充 mysql_offset。
		cursor.MySQLOffset = offset
		// 返回兼容后的游标，让服务端继续往下执行。
		return cursor, nil
	}

	// 先按新协议解 Base64，恢复出真实的双模游标结构。
	payload, err := base64.RawURLEncoding.DecodeString(raw)
	// 解码失败说明游标已损坏，此时交给上层按“回首页”策略兜底。
	if err != nil {
		return cursor, err
	}
	// 再把 JSON 反序列化成结构体，恢复 ES 与 MySQL 的翻页位点。
	if err = json.Unmarshal(payload, &cursor); err != nil {
		return cursor, err
	}
	// 当游标里没有排序码时，补回本次请求的排序码，保证后续排序一致。
	if cursor.SortCode == "" {
		cursor.SortCode = normalizeSortCode(req.GetSortCode())
	}
	// 当游标里没有查询摘要时，补回当前请求摘要，避免老数据缺字段时崩掉。
	if cursor.QueryDigest == "" {
		cursor.QueryDigest = buildQueryDigest(req)
	}
	// 当游标摘要与当前请求不一致时，说明用户换了搜索词或筛选条件，必须强制回首页。
	if cursor.QueryDigest != buildQueryDigest(req) {
		return searchCursor{SortCode: normalizeSortCode(req.GetSortCode()), QueryDigest: buildQueryDigest(req)}, nil
	}
	// 返回校验通过的游标，让主流程继续执行主搜索或回退搜索。
	return cursor, nil
}

func buildQueryDigest(req *v1.SearchProductsReq) string {
	// 空请求没有业务查询语义，直接返回空摘要即可。
	if req == nil {
		return ""
	}
	// 把关键词、店铺、类目和排序拼成稳定字符串，确保翻页只绑定同一组查询条件。
	payload := fmt.Sprintf(
		"q=%s|shop=%s|cat=%s|store_cat=%d|store_level=%d|sort=%s",
		strings.TrimSpace(req.GetQuery()),
		strings.TrimSpace(req.GetShopNo()),
		strings.TrimSpace(req.GetCategoryNo()),
		req.GetStoreCategoryId(),
		req.GetStoreCategoryLevel(),
		normalizeSortCode(req.GetSortCode()),
	)
	// 用 sha1 生成固定长度摘要，避免把原始查询条件全部塞进游标里导致过长。
	sum := sha1.Sum([]byte(payload))
	// 把摘要转成十六进制字符串，便于后续比较和排障。
	return fmt.Sprintf("%x", sum[:])
}

func normalizeSortCode(sortCode string) string {
	// 先统一大小写与空白，保证前端大小写不规范时排序语义仍然稳定。
	switch strings.ToLower(strings.TrimSpace(sortCode)) {
	// 明确保留销量降序排序码，便于 ES/MySQL 两套路径共用。
	case "sales_desc":
		return "sales_desc"
	// 明确保留价格升序排序码，便于构造稳定游标。
	case "price_asc":
		return "price_asc"
	// 明确保留价格降序排序码，便于构造稳定游标。
	case "price_desc":
		return "price_desc"
	// 其他值全部回退为默认排序，避免非法排序码把查询搞乱。
	default:
		return "default"
	}
}
