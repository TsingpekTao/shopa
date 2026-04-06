package v1

import (
	pb "github.com/TsingpekTao/shopa/points-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type InitPointsAccountIfAbsentReq struct {
	g.Meta `path:"/v1/internal/points/accounts/init" method:"post" tags:"PointsInternal" summary:"Initialize points account if absent"`
	pb.InitPointsAccountIfAbsentReq
}

type InitPointsAccountIfAbsentRes = pb.InitPointsAccountIfAbsentRes

type GetPointsByUserIdReq struct {
	g.Meta `path:"/v1/internal/points/users/{userId}" method:"get" tags:"PointsInternal" summary:"Get points balance by user id"`
	UserId uint64 `json:"userId" in:"path" v:"required#userId is required"`
}

type GetPointsByUserIdRes = pb.GetPointsByUserIdRes

type PointsOrderDraftItem struct {
	ShopNo      string `json:"shop_no"`
	SkuNo       string `json:"sku_no"`
	SpuNo       string `json:"spu_no"`
	Qty         uint32 `json:"qty"`
	SalePrice   uint64 `json:"sale_price"`
	LineAmount  uint64 `json:"line_amount"`
	MarketPrice uint64 `json:"market_price"`
}

type PointsOrderDraftSub struct {
	ShopNo        string                 `json:"shop_no"`
	GoodsAmount   uint64                 `json:"goods_amount"`
	FreightAmount uint64                 `json:"freight_amount"`
	PayableAmount uint64                 `json:"payable_amount"`
	Items         []PointsOrderDraftItem `json:"items"`
}

type PointsOrderDraft struct {
	OrderNo        string                 `json:"order_no"`
	UserID         uint64                 `json:"user_id"`
	GoodsAmount    uint64                 `json:"goods_amount"`
	FreightAmount  uint64                 `json:"freight_amount"`
	DiscountAmount uint64                 `json:"discount_amount"`
	PayableAmount  uint64                 `json:"payable_amount"`
	SubOrders      []PointsOrderDraftSub  `json:"sub_orders"`
	Lines          []PointsOrderDraftItem `json:"lines"`
	Metadata       map[string]string      `json:"metadata,omitempty"`
}

type PointsSubAllocation struct {
	ShopNo               string `json:"shop_no"`
	PointsUsed           uint64 `json:"points_used"`
	PointsDiscountAmount uint64 `json:"points_discount_amount"`
}

type PreviewOrderReq struct {
	g.Meta        `path:"/v1/points/internal/order/preview" method:"post" tags:"PointsInternal" summary:"Preview points deduction for order draft"`
	OrderDraft    *PointsOrderDraft `json:"order_draft"`
	UsePoints     bool              `json:"use_points"`
	IntentPoints  uint64            `json:"intent_points"`
	RequestSource string            `json:"request_source"`
}

type PreviewOrderRes struct {
	UsePoints              bool                  `json:"use_points"`
	PointsUsed             uint64                `json:"points_used"`
	PointsDiscountAmount   uint64                `json:"points_discount_amount"`
	PointsRuleSnapshotJSON string                `json:"points_rule_snapshot_json"`
	PointsRuleDigest       string                `json:"points_rule_snapshot_digest"`
	SubAllocations         []PointsSubAllocation `json:"sub_allocations"`
}

type LockOrderReq struct {
	g.Meta                     `path:"/v1/points/internal/order/lock" method:"post" tags:"PointsInternal" summary:"Lock points for order draft"`
	OrderDraft                 *PointsOrderDraft `json:"order_draft"`
	UsePoints                  bool              `json:"use_points"`
	IntentPoints               uint64            `json:"intent_points"`
	IdempotencyKey             string            `json:"idempotency_key"`
	ExpectedPointsCashAmount   uint64            `json:"expected_points_cash_amount"`
	ExpectedRuleSnapshotDigest string            `json:"expected_rule_snapshot_digest"`
	RequestSource              string            `json:"request_source"`
}

type LockOrderRes struct {
	ReservationNo            string                `json:"reservation_no"`
	PointsUsed               uint64                `json:"points_used"`
	PointsDiscountAmount     uint64                `json:"points_discount_amount"`
	PointsRuleSnapshotJSON   string                `json:"points_rule_snapshot_json"`
	PointsRuleSnapshotDigest string                `json:"points_rule_snapshot_digest"`
	SubAllocations           []PointsSubAllocation `json:"sub_allocations"`
}

type CancelOrderReq struct {
	g.Meta         `path:"/v1/points/internal/order/cancel" method:"post" tags:"PointsInternal" summary:"Cancel locked points reservation"`
	ReservationNo  string `json:"reservation_no"`
	OrderNo        string `json:"order_no"`
	UserID         uint64 `json:"user_id"`
	ReasonCode     string `json:"reason_code"`
	IdempotencyKey string `json:"idempotency_key"`
	RequestSource  string `json:"request_source"`
}

type ConfirmOrderReq struct {
	g.Meta         `path:"/v1/points/internal/order/confirm" method:"post" tags:"PointsInternal" summary:"Confirm locked points reservation"`
	ReservationNo  string `json:"reservation_no"`
	OrderNo        string `json:"order_no"`
	UserID         uint64 `json:"user_id"`
	IdempotencyKey string `json:"idempotency_key"`
	RequestSource  string `json:"request_source"`
}

type GrantOrderReq struct {
	g.Meta         `path:"/v1/points/internal/order/grant" method:"post" tags:"PointsInternal" summary:"Grant points after order completed"`
	OrderNo        string `json:"order_no"`
	UserID         uint64 `json:"user_id"`
	PaidAmount     uint64 `json:"paid_amount"`
	IdempotencyKey string `json:"idempotency_key"`
	RequestSource  string `json:"request_source"`
}

type SimpleAckRes struct {
	Success bool `json:"success"`
}

type ReturnRefundReq struct {
	g.Meta         `path:"/v1/points/internal/refund/return" method:"post" tags:"PointsInternal" summary:"Return used points by refund"`
	RefundNo       string `json:"refund_no"`
	OrderNo        string `json:"order_no"`
	SubOrderNo     string `json:"sub_order_no"`
	ShopNo         string `json:"shop_no"`
	UserID         uint64 `json:"user_id"`
	PointsToReturn uint64 `json:"points_to_return"`
	CashAmountCent uint64 `json:"cash_amount_cent"`
	IdempotencyKey string `json:"idempotency_key"`
	RequestSource  string `json:"request_source"`
}

type ReturnRefundRes struct {
	Success            bool   `json:"success"`
	PointsReturnAmount uint64 `json:"points_return_amount"`
}

type ReverseRefundReq struct {
	g.Meta               `path:"/v1/points/internal/refund/reverse" method:"post" tags:"PointsInternal" summary:"Reverse granted points by refund"`
	RefundNo             string `json:"refund_no"`
	OrderNo              string `json:"order_no"`
	SubOrderNo           string `json:"sub_order_no"`
	ShopNo               string `json:"shop_no"`
	UserID               uint64 `json:"user_id"`
	PointsToReverse      uint64 `json:"points_to_reverse"`
	ApprovedRefundAmount uint64 `json:"approved_refund_amount"`
	IdempotencyKey       string `json:"idempotency_key"`
	RequestSource        string `json:"request_source"`
}

type ReverseRefundRes struct {
	Success                bool   `json:"success"`
	PointsReverseAmount    uint64 `json:"points_reverse_amount"`
	PointsCashOffsetAmount uint64 `json:"points_cash_offset_amount"`
	AccountDebtAfter       uint64 `json:"account_debt_after"`
}
