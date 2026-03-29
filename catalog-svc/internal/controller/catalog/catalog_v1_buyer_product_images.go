package catalog

import (
	"context"
	"strings"

	v1 "github.com/TsingpekTao/shopa/catalog-svc/api/catalog/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type spuImageRow struct {
	Id       uint64 `orm:"id"`
	SpuNo    string `orm:"spu_no"`
	ImageUrl string `orm:"image_url"`
}

func (c *ControllerV1) ListBuyerProductImages(ctx context.Context, req *v1.ListBuyerProductImagesReq) (*v1.ListBuyerProductImagesRes, error) {
	spuNos := parseSpuNoCSV(req.SpuNos)
	if len(spuNos) == 0 {
		return &v1.ListBuyerProductImagesRes{Items: []v1.BuyerProductImageItem{}}, nil
	}

	var rows []spuImageRow
	err := g.DB().Model("catalog_spu_image").
		Fields("id", "spu_no", "image_url").
		WhereIn("spu_no", spuNos).
		Where("status", 1).
		OrderDesc("is_primary").
		OrderAsc("sort_order").
		OrderAsc("id").
		Scan(&rows)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "doesn't exist") {
			return &v1.ListBuyerProductImagesRes{Items: []v1.BuyerProductImageItem{}}, nil
		}
		return nil, err
	}

	picked := make(map[string]v1.BuyerProductImageItem, len(rows))
	for _, row := range rows {
		if strings.TrimSpace(row.SpuNo) == "" || strings.TrimSpace(row.ImageUrl) == "" {
			continue
		}
		if _, ok := picked[row.SpuNo]; ok {
			continue
		}
		picked[row.SpuNo] = v1.BuyerProductImageItem{
			SpuNo:    row.SpuNo,
			ImageUrl: row.ImageUrl,
		}
	}

	items := make([]v1.BuyerProductImageItem, 0, len(picked))
	for _, spuNo := range spuNos {
		if item, ok := picked[spuNo]; ok {
			items = append(items, item)
		}
	}
	return &v1.ListBuyerProductImagesRes{Items: items}, nil
}

func parseSpuNoCSV(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r' || r == '\t' || r == ' '
	})
	out := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		spuNo := strings.TrimSpace(part)
		if spuNo == "" {
			continue
		}
		if _, ok := seen[spuNo]; ok {
			continue
		}
		seen[spuNo] = struct{}{}
		out = append(out, spuNo)
		if len(out) >= 200 {
			break
		}
	}
	return out
}
