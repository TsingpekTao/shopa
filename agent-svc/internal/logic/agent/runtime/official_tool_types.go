package runtime

import "context"

type PolicyKnowledgeHit struct {
	SourceTypeCode string `json:"source_type_code,omitempty"`
	SourceID       string `json:"source_id,omitempty"`
	SourceVersion  uint64 `json:"source_version,omitempty"`
	Title          string `json:"title,omitempty"`
	Snippet        string `json:"snippet,omitempty"`
}

type PolicyKnowledgeRepository interface {
	SearchPolicyKnowledge(ctx context.Context, query string, limit uint32) ([]PolicyKnowledgeHit, error)
}

type ProductSearchFilter struct {
	Query      string `json:"query,omitempty"`
	ShopNo     string `json:"shop_no,omitempty"`
	CategoryNo string `json:"category_no,omitempty"`
	SortCode   string `json:"sort_code,omitempty"`
	Limit      uint32 `json:"limit,omitempty"`
}

type ProductSearchItem struct {
	SpuNo      string `json:"spu_no,omitempty"`
	Title      string `json:"title,omitempty"`
	CoverURL   string `json:"cover_url,omitempty"`
	MinPrice   uint64 `json:"min_price,omitempty"`
	MaxPrice   uint64 `json:"max_price,omitempty"`
	ShopName   string `json:"shop_name,omitempty"`
	ReasonText string `json:"reason_text,omitempty"`
}

type ProductSearchRepository interface {
	SearchProducts(ctx context.Context, filter ProductSearchFilter) ([]ProductSearchItem, error)
}
