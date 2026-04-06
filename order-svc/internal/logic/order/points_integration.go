package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	pointsCompensationActionCancel = "CANCEL_POINTS_RESERVATION"
)

type pointsIntegrationConf struct {
	Enabled     bool
	BaseURL     string
	Timeout     time.Duration
	PreviewPath string
	LockPath    string
	CancelPath  string
	ConfirmPath string
	GrantPath   string
}

type pointsOrderDraft struct {
	OrderNo        string            `json:"order_no"`
	UserID         uint64            `json:"user_id"`
	GoodsAmount    uint64            `json:"goods_amount"`
	FreightAmount  uint64            `json:"freight_amount"`
	DiscountAmount uint64            `json:"discount_amount"`
	PayableAmount  uint64            `json:"payable_amount"`
	SubOrders      []pointsSubDraft  `json:"sub_orders"`
	Lines          []pointsItemDraft `json:"lines"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type pointsSubDraft struct {
	ShopNo        string            `json:"shop_no"`
	GoodsAmount   uint64            `json:"goods_amount"`
	FreightAmount uint64            `json:"freight_amount"`
	PayableAmount uint64            `json:"payable_amount"`
	Items         []pointsItemDraft `json:"items"`
}

type pointsItemDraft struct {
	ShopNo      string `json:"shop_no"`
	SkuNo       string `json:"sku_no"`
	SpuNo       string `json:"spu_no"`
	Qty         uint32 `json:"qty"`
	SalePrice   uint64 `json:"sale_price"`
	LineAmount  uint64 `json:"line_amount"`
	MarketPrice uint64 `json:"market_price"`
}

type pointsSubAllocation struct {
	ShopNo               string `json:"shop_no"`
	PointsUsed           uint64 `json:"points_used"`
	PointsDiscountAmount uint64 `json:"points_discount_amount"`
}

type pointsPreviewRequest struct {
	OrderDraft    *pointsOrderDraft `json:"order_draft"`
	UsePoints     bool              `json:"use_points"`
	IntentPoints  uint64            `json:"intent_points"`
	RequestSource string            `json:"request_source"`
}

type pointsPreviewResponse struct {
	UsePoints              bool                  `json:"use_points"`
	PointsUsed             uint64                `json:"points_used"`
	PointsDiscountAmount   uint64                `json:"points_discount_amount"`
	PointsRuleSnapshotJSON string                `json:"points_rule_snapshot_json"`
	PointsRuleDigest       string                `json:"points_rule_snapshot_digest"`
	SubAllocations         []pointsSubAllocation `json:"sub_allocations"`
}

type pointsLockRequest struct {
	OrderDraft                 *pointsOrderDraft `json:"order_draft"`
	UsePoints                  bool              `json:"use_points"`
	IntentPoints               uint64            `json:"intent_points"`
	IdempotencyKey             string            `json:"idempotency_key"`
	ExpectedPointsCashAmount   uint64            `json:"expected_points_cash_amount"`
	ExpectedRuleSnapshotDigest string            `json:"expected_rule_snapshot_digest"`
	RequestSource              string            `json:"request_source"`
}

type pointsLockResponse struct {
	ReservationNo            string                `json:"reservation_no"`
	PointsUsed               uint64                `json:"points_used"`
	PointsDiscountAmount     uint64                `json:"points_discount_amount"`
	PointsRuleSnapshotJSON   string                `json:"points_rule_snapshot_json"`
	PointsRuleSnapshotDigest string                `json:"points_rule_snapshot_digest"`
	SubAllocations           []pointsSubAllocation `json:"sub_allocations"`
}

type pointsCancelRequest struct {
	ReservationNo  string `json:"reservation_no"`
	OrderNo        string `json:"order_no"`
	UserID         uint64 `json:"user_id"`
	ReasonCode     string `json:"reason_code"`
	IdempotencyKey string `json:"idempotency_key"`
	RequestSource  string `json:"request_source"`
}

type pointsConfirmRequest struct {
	ReservationNo  string `json:"reservation_no"`
	OrderNo        string `json:"order_no"`
	UserID         uint64 `json:"user_id"`
	IdempotencyKey string `json:"idempotency_key"`
	RequestSource  string `json:"request_source"`
}

type pointsGrantRequest struct {
	OrderNo        string `json:"order_no"`
	UserID         uint64 `json:"user_id"`
	PaidAmount     uint64 `json:"paid_amount"`
	IdempotencyKey string `json:"idempotency_key"`
	RequestSource  string `json:"request_source"`
}

type pointsSimpleAck struct {
	Success bool `json:"success"`
}

type gfHTTPEnvelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// loadPointsIntegrationConf 读取 order 对 points 的集成配置。
func loadPointsIntegrationConf(ctx context.Context) pointsIntegrationConf {
	var (
		baseURL   = strings.TrimRight(strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.points.baseUrl", "").String()), "/")
		timeoutMs = g.Cfg().MustGet(ctx, "upstream.points.timeoutMs", 3000).Int()
	)
	if timeoutMs <= 0 {
		timeoutMs = 3000
	}
	// 若未配置 BaseURL，则默认视为未启用，避免误调用产生长时间等待。
	return pointsIntegrationConf{
		Enabled:     g.Cfg().MustGet(ctx, "upstream.points.enabled", baseURL != "").Bool(),
		BaseURL:     baseURL,
		Timeout:     time.Duration(timeoutMs) * time.Millisecond,
		PreviewPath: strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.points.previewPath", "/v1/points/internal/order/preview").String()),
		LockPath:    strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.points.lockPath", "/v1/points/internal/order/lock").String()),
		CancelPath:  strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.points.cancelPath", "/v1/points/internal/order/cancel").String()),
		ConfirmPath: strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.points.confirmPath", "/v1/points/internal/order/confirm").String()),
		GrantPath:   strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.points.grantPath", "/v1/points/internal/order/grant").String()),
	}
}

// buildPointsOrderDraft 将订单草稿转换为积分服务需要的请求结构。
func buildPointsOrderDraft(input *createOrderInput) *pointsOrderDraft {
	if input == nil {
		return nil
	}
	var (
		subMap = make(map[string]*pointsSubDraft)
		lines  = make([]pointsItemDraft, 0, len(input.Lines))
	)
	for _, line := range input.Lines {
		lineAmount := line.SalePrice * uint64(line.Qty)
		// 逐条累加行金额，同时构建提交给 points 的行列表。
		lines = append(lines, pointsItemDraft{
			ShopNo:      line.ShopNo,
			SkuNo:       line.SkuNo,
			SpuNo:       line.SpuNo,
			Qty:         line.Qty,
			SalePrice:   line.SalePrice,
			LineAmount:  lineAmount,
			MarketPrice: line.MarketPrice,
		})
		subDraft, ok := subMap[line.ShopNo]
		if !ok {
			subDraft = &pointsSubDraft{ShopNo: line.ShopNo}
			subMap[line.ShopNo] = subDraft
		}
		// 子单聚合：同店铺的行聚合为一个 sub order，便于下游拆分账务。
		subDraft.GoodsAmount += lineAmount
		subDraft.PayableAmount += lineAmount
		subDraft.Items = append(subDraft.Items, pointsItemDraft{
			ShopNo:      line.ShopNo,
			SkuNo:       line.SkuNo,
			SpuNo:       line.SpuNo,
			Qty:         line.Qty,
			SalePrice:   line.SalePrice,
			LineAmount:  lineAmount,
			MarketPrice: line.MarketPrice,
		})
	}
	subOrders := make([]pointsSubDraft, 0, len(subMap))
	for _, subDraft := range subMap {
		subOrders = append(subOrders, *subDraft)
	}
	return &pointsOrderDraft{
		OrderNo:        input.OrderNo,
		UserID:         input.UserID,
		GoodsAmount:    input.GoodsAmount,
		FreightAmount:  input.FreightAmount,
		DiscountAmount: input.DiscountAmount,
		PayableAmount:  input.PayableAmount,
		SubOrders:      subOrders,
		Lines:          lines,
		Metadata: map[string]string{
			"buyer_remark": input.BuyerRemark,
		},
	}
}

// validatePointsExpectation 校验前端预期的积分折现结果与规则摘要。
func validatePointsExpectation(expectedCash uint64, expectedDigest string, resp *pointsPreviewResponse) error {
	if resp == nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "points preview result is nil")
	}
	if expectedCash > 0 && resp.PointsDiscountAmount != expectedCash {
		return gerror.NewCode(gcode.CodeInvalidParameter, "expected points cash amount mismatch")
	}
	if strings.TrimSpace(expectedDigest) != "" && strings.TrimSpace(resp.PointsRuleDigest) != strings.TrimSpace(expectedDigest) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "points rule snapshot digest mismatch")
	}
	return nil
}

// previewOrderPoints 调用积分服务预览本次订单的积分抵扣结果。
func (s *sOrder) previewOrderPoints(ctx context.Context, input *createOrderInput, usePoints bool, intentPoints, expectedCash uint64, expectedDigest, requestSource string) (*pointsPreviewResponse, error) {
	if !usePoints {
		// 不使用积分时直接返回空结果，避免下游无意义调用。
		return &pointsPreviewResponse{}, nil
	}
	conf := loadPointsIntegrationConf(ctx)
	if !conf.Enabled || conf.BaseURL == "" {
		return nil, gerror.NewCode(gcode.CodeNotImplemented, "points preview endpoint not configured")
	}
	resp := new(pointsPreviewResponse)
	if err := doPointsJSONRequest(ctx, conf, conf.PreviewPath, &pointsPreviewRequest{
		OrderDraft:    buildPointsOrderDraft(input),
		UsePoints:     usePoints,
		IntentPoints:  intentPoints,
		RequestSource: requestSource,
	}, resp); err != nil {
		return nil, err
	}
	if err := validatePointsExpectation(expectedCash, expectedDigest, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// lockPointsForOrder 调用积分服务锁定本次订单要占用的积分。
func (s *sOrder) lockPointsForOrder(ctx context.Context, input *createOrderInput, usePoints bool, intentPoints, expectedCash uint64, expectedDigest, idemKey, requestSource string) (*pointsLockResponse, error) {
	if !usePoints {
		// 未使用积分时，不需要向积分侧占用额度。
		return &pointsLockResponse{}, nil
	}
	conf := loadPointsIntegrationConf(ctx)
	if !conf.Enabled || conf.BaseURL == "" {
		return nil, gerror.NewCode(gcode.CodeNotImplemented, "points lock endpoint not configured")
	}
	resp := new(pointsLockResponse)
	if err := doPointsJSONRequest(ctx, conf, conf.LockPath, &pointsLockRequest{
		OrderDraft:                 buildPointsOrderDraft(input),
		UsePoints:                  usePoints,
		IntentPoints:               intentPoints,
		IdempotencyKey:             idemKey,
		ExpectedPointsCashAmount:   expectedCash,
		ExpectedRuleSnapshotDigest: expectedDigest,
		RequestSource:              requestSource,
	}, resp); err != nil {
		return nil, err
	}
	if expectedCash > 0 && resp.PointsDiscountAmount != expectedCash {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "locked points cash amount mismatch")
	}
	if strings.TrimSpace(expectedDigest) != "" && strings.TrimSpace(resp.PointsRuleSnapshotDigest) != strings.TrimSpace(expectedDigest) {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "locked points rule snapshot digest mismatch")
	}
	return resp, nil
}

// cancelLockedPoints 调用积分服务取消已锁定的积分。
func (s *sOrder) cancelLockedPoints(ctx context.Context, userID uint64, orderNo, reservationNo, reasonCode, idemKey string) error {
	if strings.TrimSpace(reservationNo) == "" {
		return nil
	}
	conf := loadPointsIntegrationConf(ctx)
	if !conf.Enabled || conf.BaseURL == "" {
		return gerror.NewCode(gcode.CodeNotImplemented, "points cancel endpoint not configured")
	}
	// 取消请求必须携带原因与幂等键，防止重复扣返。
	return doPointsJSONRequest(ctx, conf, conf.CancelPath, &pointsCancelRequest{
		ReservationNo:  reservationNo,
		OrderNo:        orderNo,
		UserID:         userID,
		ReasonCode:     reasonCode,
		IdempotencyKey: idemKey,
		RequestSource:  "order-svc",
	}, new(pointsSimpleAck))
}

// confirmLockedPoints 调用积分服务确认已锁定的积分。
func (s *sOrder) confirmLockedPoints(ctx context.Context, userID uint64, orderNo, reservationNo, idemKey string) error {
	if strings.TrimSpace(reservationNo) == "" {
		return nil
	}
	conf := loadPointsIntegrationConf(ctx)
	if !conf.Enabled || conf.BaseURL == "" {
		return gerror.NewCode(gcode.CodeNotImplemented, "points confirm endpoint not configured")
	}
	// 支付成功后确认冻结为正式扣减。
	return doPointsJSONRequest(ctx, conf, conf.ConfirmPath, &pointsConfirmRequest{
		ReservationNo:  reservationNo,
		OrderNo:        orderNo,
		UserID:         userID,
		IdempotencyKey: idemKey,
		RequestSource:  "order-svc",
	}, new(pointsSimpleAck))
}

// grantPointsByOrderCompleted 在订单完成后调用积分服务发放奖励积分。
func (s *sOrder) grantPointsByOrderCompleted(ctx context.Context, userID uint64, orderNo string, paidAmount uint64, idemKey string) error {
	conf := loadPointsIntegrationConf(ctx)
	if !conf.Enabled || conf.BaseURL == "" {
		return gerror.NewCode(gcode.CodeNotImplemented, "points grant endpoint not configured")
	}
	// 完成订单后的赠分行为依赖支付金额和幂等键，防止重复赠送。
	return doPointsJSONRequest(ctx, conf, conf.GrantPath, &pointsGrantRequest{
		OrderNo:        orderNo,
		UserID:         userID,
		PaidAmount:     paidAmount,
		IdempotencyKey: idemKey,
		RequestSource:  "order-svc",
	}, new(pointsSimpleAck))
}

// doPointsJSONRequest 向 points 服务发送统一的 JSON 请求。
func doPointsJSONRequest(ctx context.Context, conf pointsIntegrationConf, path string, reqBody any, respBody any) error {
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return gerror.Wrap(err, "marshal points request failed")
	}
	// 设置请求超时时间，防止主链路因下游故障被长期阻塞。
	requestCtx, cancel := context.WithTimeout(ctx, conf.Timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(requestCtx, http.MethodPost, conf.BaseURL+path, bytes.NewReader(payload))
	if err != nil {
		return gerror.Wrap(err, "build points request failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Name", "order-svc")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return gerror.Wrap(err, "points request failed")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return gerror.Wrap(err, "read points response failed")
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return gerror.NewCodef(gcode.CodeInternalError, "points endpoint %s returned status=%d body=%s", path, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	if respBody == nil || len(body) == 0 {
		return nil
	}
	// 优先尝试直接解包到业务响应体；失败后再按 GoFrame envelope 解包，兼容两种返回格式。
	if err = json.Unmarshal(body, respBody); err == nil {
		return nil
	}
	var envelope gfHTTPEnvelope
	if err = json.Unmarshal(body, &envelope); err != nil {
		return gerror.Wrapf(err, "unmarshal points response failed, path=%s", path)
	}
	if envelope.Code != 0 {
		return gerror.NewCodef(gcode.CodeInternalError, "points endpoint %s returned code=%d message=%s", path, envelope.Code, envelope.Message)
	}
	if respBody == nil || len(envelope.Data) == 0 || string(envelope.Data) == "null" {
		return nil
	}
	if err = json.Unmarshal(envelope.Data, respBody); err != nil {
		return gerror.Wrapf(err, "unmarshal points response data failed, path=%s", path)
	}
	return nil
}

// pointsAllocationsToMap 将子单级积分分摊结果转换为店铺索引映射。
func pointsAllocationsToMap(in []pointsSubAllocation) map[string]pointsSubAllocation {
	out := make(map[string]pointsSubAllocation, len(in))
	for _, item := range in {
		shopNo := strings.TrimSpace(item.ShopNo)
		if shopNo == "" {
			continue
		}
		out[shopNo] = item
	}
	return out
}

// buildPointsCancelIdempotencyKey 构造取消积分锁定时使用的幂等键。
func buildPointsCancelIdempotencyKey(orderNo, action string) string {
	return fmt.Sprintf("%s:%s", action, orderNo)
}
