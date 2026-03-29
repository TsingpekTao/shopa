package v1

import "github.com/gogf/gf/v2/frame/g"

type GetOverviewReq struct {
	g.Meta `path:"/v1/admin/me/overview" method:"get" tags:"Admin" summary:"Get admin overview from gateway"`
}

type GetOverviewRes struct {
	UserID            uint64   `json:"userId"`
	AccountStatusCode string   `json:"accountStatusCode"`
	Roles             []string `json:"roles"`
	Permissions       []string `json:"permissions"`
}

type GetDashboardOverviewReq struct {
	g.Meta `path:"/v1/admin/dashboard/overview" method:"get" tags:"Admin" summary:"Get admin dashboard overview"`
}

type DashboardMetric struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value int64  `json:"value"`
}

type GetDashboardOverviewRes struct {
	Metrics        []DashboardMetric `json:"metrics"`
	Partial        bool              `json:"partial"`
	DegradedFields []string          `json:"degradedFields"`
}

type GetShopInsightsReq struct {
	g.Meta `path:"/v1/admin/shops/{shopNo}/insights" method:"get" tags:"Admin" summary:"Get admin shop insights"`
	ShopNo string `json:"shopNo" in:"path" v:"required"`
}

type GetShopInsightsRes struct {
	ShopNo         string   `json:"shopNo"`
	ShopName       string   `json:"shopName"`
	ShopStatusCode string   `json:"shopStatusCode"`
	Partial        bool     `json:"partial"`
	DegradedFields []string `json:"degradedFields"`
}

type ListConversationsReq struct {
	g.Meta `path:"/v1/admin/cs/conversations" method:"get" tags:"Admin" summary:"List customer service conversations placeholder"`
}

type ConversationItem struct {
	ConversationNo string `json:"conversationNo"`
	Title          string `json:"title"`
	LastMessageAt  string `json:"lastMessageAt"`
	StatusCode     string `json:"statusCode"`
}

type ListConversationsRes struct {
	Items          []ConversationItem `json:"items"`
	Partial        bool               `json:"partial"`
	DegradedFields []string           `json:"degradedFields"`
}
