package payment

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	orderv1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	v1 "github.com/TsingpekTao/shopa/payment-svc/api/v1"
	"github.com/TsingpekTao/shopa/payment-svc/internal/dao"
	"github.com/TsingpekTao/shopa/payment-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/payment-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/payment-svc/internal/service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// sPayment 定义支付领域逻辑实现。
type sPayment struct{}

type paymentLaunch struct {
	PayURL         string
	PayPayloadJSON string
}

type internalOrderCallbackClient interface {
	HandlePayCallback(ctx context.Context, req *orderv1.HandlePayCallbackReq, opts ...grpc.CallOption) (*orderv1.HandlePayCallbackRes, error)
}

type alipayConfig struct {
	GatewayURL    string
	AppID         string
	PrivateKey    string
	PublicKey     string
	ReturnURL     string
	NotifyURL     string
	SubjectPrefix string
}

const (
	// actionCreateIntent 标记“创建支付意图”的幂等动作。
	actionCreateIntent = "CREATE_PAYMENT_INTENT"
	// actionCloseIntent 标记“关闭支付意图”的幂等动作。
	actionCloseIntent = "CLOSE_PAYMENT_INTENT"
	// actionCreateRefund 标记“创建退款任务”的幂等动作。
	actionCreateRefund = "CREATE_REFUND_TASK"
	// actionExecuteRefund 标记“执行退款任务”的幂等动作。
	actionExecuteRefund = "EXECUTE_REFUND_TASK"
	// actionRunRecon 标记“发起日对账”的幂等动作。
	actionRunRecon = "RUN_DAILY_RECONCILIATION"
)

var (
	newInternalOrderClient = defaultInternalOrderClient

	internalOrderClientMu   sync.Mutex
	internalOrderClientAddr string
	internalOrderClientConn *grpc.ClientConn
	internalOrderClientInst internalOrderCallbackClient
)

// New 创建支付逻辑实例。
func New() *sPayment {
	// 返回无状态逻辑对象供 service 层复用。
	return &sPayment{}
}

func init() {
	// 在初始化阶段注册 payment 领域服务实现。
	service.RegisterPayment(New())
}

// CreatePaymentIntent 创建支付意图并返回收银参数。
func (s *sPayment) CreatePaymentIntent(ctx context.Context, req *v1.CreatePaymentIntentReq) (*v1.CreatePaymentIntentRes, error) {
	// 校验关键入参，避免生成非法支付单。
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || req.GetPayableAmount() == 0 {
		// 参数不完整时返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/payable_amount are required")
	}
	// 先尝试读取幂等快照，命中则直接回放历史结果。
	if hit, payload, err := s.replayIdempotency(ctx, req.GetIdempotencyKey(), actionCreateIntent); err != nil {
		// 幂等查询失败时中断流程并返回错误。
		return nil, err
	} else if hit {
		// 从幂等快照读取 payment_no，用于构造回放响应。
		paymentNo := gconv.String(payload["payment_no"])
		// 只有在快照中存在 payment_no 时才回查详情。
		if strings.TrimSpace(paymentNo) != "" {
			// 按 payment_no 查询支付意图快照。
			intent, getErr := s.getIntentByPaymentNo(ctx, paymentNo)
			// 查询失败时返回具体错误。
			if getErr != nil {
				return nil, getErr
			}
			// 返回幂等回放结果，避免重复创建支付单。
			launch, launchErr := s.buildPaymentLaunch(ctx, req, intent)
			if launchErr != nil {
				return nil, launchErr
			}
			return &v1.CreatePaymentIntentRes{Intent: intent, PayUrl: launch.PayURL, PayPayloadJson: launch.PayPayloadJSON, IdempotentReplay: true}, nil
		}
	}

	// 生成支付单号作为主业务键。
	paymentNo := newBizNo("PAY")
	// 解析订单过期时间用于对齐网关有效期。
	orderExpireAt := pbTsToGTime(req.GetOrderExpireAt())
	// 默认把网关过期时间设置为订单过期时间。
	gatewayExpireAt := orderExpireAt
	// 当上游未传订单过期时间时使用“当前+15分钟”兜底。
	if orderExpireAt == nil {
		// 生成默认过期时间，保证支付窗口可控。
		gatewayExpireAt = gtime.New(time.Now().Add(15 * time.Minute))
		// 同步写回订单过期时间，保持两者一致。
		orderExpireAt = gatewayExpireAt
	}
	if existing, err := s.getIntentEntityByOrderNo(ctx, req.GetOrderNo()); err != nil {
		if !isIdempotencyMissError(err) {
			return nil, err
		}
	} else {
		if canReuseExistingIntentForCreate(existing) {
			if err := s.refreshExistingIntentForCreate(ctx, existing, req, orderExpireAt, gatewayExpireAt); err != nil {
				return nil, err
			}
			intent, err := s.getIntentByOrderNo(ctx, req.GetOrderNo())
			if err != nil {
				return nil, err
			}
			launch, err := s.buildPaymentLaunch(ctx, req, intent)
			if err != nil {
				return nil, err
			}
			if err = s.saveIdempotency(ctx, req.GetIdempotencyKey(), actionCreateIntent, intent.GetPaymentNo(), map[string]any{
				"payment_no":       intent.GetPaymentNo(),
				"pay_url":          launch.PayURL,
				"pay_payload_json": launch.PayPayloadJSON,
			}); err != nil {
				g.Log().Warningf(ctx, "save idempotency failed: %v", err)
			}
			return &v1.CreatePaymentIntentRes{Intent: intent, PayUrl: launch.PayURL, PayPayloadJson: launch.PayPayloadJSON, IdempotentReplay: false}, nil
		}
		if v1.PaymentIntentStatus(existing.Status) == v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAID {
			return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "payment already completed")
		}
		return nil, gerror.NewCodef(gcode.CodeBusinessValidationFailed, "payment intent not reusable, status=%s", v1.PaymentIntentStatus(existing.Status).String())
	}
	intent := &v1.PaymentIntent{
		PaymentNo:       paymentNo,
		OrderNo:         req.GetOrderNo(),
		UserId:          req.GetUserId(),
		PayChannel:      req.GetPayChannel(),
		Status:          v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAYING,
		PayableAmount:   req.GetPayableAmount(),
		RefundedAmount:  0,
		CurrencyCode:    safeCurrency(req.GetCurrencyCode()),
		OrderExpireAt:   gtimeToPB(orderExpireAt),
		GatewayExpireAt: gtimeToPB(gatewayExpireAt),
		Version:         1,
	}
	launch, err := s.buildPaymentLaunch(ctx, req, intent)
	if err != nil {
		return nil, err
	}

	// 组装 payment_intent 入库数据。
	data := do.PaymentIntent{
		// 写入支付单号作为唯一标识。
		PaymentNo: paymentNo,
		// 写入来源订单号用于跨服务关联。
		OrderNo: req.GetOrderNo(),
		// 写入用户 ID 用于鉴权追踪。
		UserId: req.GetUserId(),
		// 写入支付渠道枚举值。
		PayChannel: uint(req.GetPayChannel()),
		// 初始化状态为 CREATED。
		Status: uint(v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAYING),
		// 写入应付金额（分）。
		PayableAmount: req.GetPayableAmount(),
		// 初始化已退款金额为 0。
		RefundedAmount: 0,
		// 规范化币种编码，默认 CNY。
		CurrencyCode: safeCurrency(req.GetCurrencyCode()),
		// 持久化订单超时关单时间。
		OrderExpireAt: orderExpireAt,
		// 持久化网关侧失效时间。
		GatewayExpireAt: gatewayExpireAt,
		// 初始化版本号用于后续 CAS 更新。
		Version: 1,
	}
	// 执行插入创建支付意图记录。
	if _, err := dao.PaymentIntent.Ctx(ctx).Data(data).Insert(); err != nil {
		// 插入失败时返回数据库错误。
		return nil, gerror.Wrap(err, "insert payment_intent failed")
	}

	// 查询刚创建的支付意图用于回包。
	intent, err = s.getIntentByPaymentNo(ctx, paymentNo)
	// 查询失败时返回错误。
	if err != nil {
		return nil, err
	}
	// 记录幂等快照用于客户端重试回放。
	if err = s.saveIdempotency(ctx, req.GetIdempotencyKey(), actionCreateIntent, paymentNo, map[string]any{"payment_no": paymentNo, "pay_url": launch.PayURL, "pay_payload_json": launch.PayPayloadJSON}); err != nil {
		// 幂等快照失败不影响主流程，只记录告警。
		g.Log().Warningf(ctx, "save idempotency failed: %v", err)
	}
	// 返回创建结果与 mock 支付链接。
	return &v1.CreatePaymentIntentRes{Intent: intent, PayUrl: launch.PayURL, PayPayloadJson: launch.PayPayloadJSON, IdempotentReplay: false}, nil
}

// QueryPaymentIntent 查询支付意图快照。
func (s *sPayment) QueryPaymentIntent(ctx context.Context, req *v1.QueryPaymentIntentReq) (*v1.QueryPaymentIntentRes, error) {
	// 校验至少提供 payment_no 或 order_no 之一。
	if req == nil || (strings.TrimSpace(req.GetPaymentNo()) == "" && strings.TrimSpace(req.GetOrderNo()) == "") {
		// 缺少查询主键时返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "payment_no or order_no is required")
	}
	// 优先按 payment_no 查询，避免 order_no 语义歧义。
	if strings.TrimSpace(req.GetPaymentNo()) != "" {
		// 按支付单号查询并返回结果。
		intent, err := s.getIntentByPaymentNo(ctx, req.GetPaymentNo())
		// 直接返回查询结果。
		return &v1.QueryPaymentIntentRes{Intent: intent}, err
	}
	// 退化为按订单号查询支付单。
	intent, err := s.getIntentByOrderNo(ctx, req.GetOrderNo())
	// 返回订单维度查询结果。
	return &v1.QueryPaymentIntentRes{Intent: intent}, err
}

// HandleGatewayCallback 处理支付网关回调。
func (s *sPayment) HandleGatewayCallback(ctx context.Context, req *v1.HandleGatewayCallbackReq) (*v1.HandleGatewayCallbackRes, error) {
	// 校验回调事件 ID，确保回调天然幂等。
	if req == nil || strings.TrimSpace(req.GetCallbackEventId()) == "" {
		// 缺少回调主键时返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "callback_event_id is required")
	}
	// 先写回调日志并利用唯一键实现去重。
	r, err := dao.PaymentCallbackLog.Ctx(ctx).Data(do.PaymentCallbackLog{
		// 记录网关事件 ID。
		CallbackEventId: req.GetCallbackEventId(),
		// 记录支付单号。
		PaymentNo: req.GetPaymentNo(),
		// 记录订单号。
		OrderNo: req.GetOrderNo(),
		// 记录支付渠道，便于后续按渠道排障。
		PayChannel: uint(req.GetPayChannel()),
		// 记录网关状态码，保留审计证据。
		GatewayStatusCode: req.GetGatewayStatusCode(),
		// 记录本次支付金额。
		PaidAmount: req.GetPaidAmount(),
		// 记录原始回调载荷。
		RawPayload: req.GetRawPayload(),
		// 记录回调签名原文。
		Signature: req.GetSignature(),
	}).InsertIgnore()
	// 日志落库失败时直接返回错误。
	if err != nil {
		return nil, gerror.Wrap(err, "insert callback log failed")
	}
	// 默认标记为“非重复回调”。
	idempotentHit := false
	// 当唯一键冲突导致 0 行受影响时表示重复回调。
	if r != nil {
		// 读取受影响行数判断是否命中幂等。
		if rows, _ := r.RowsAffected(); rows == 0 {
			// 标记回调幂等命中。
			idempotentHit = true
		}
	}

	// 查询支付意图快照用于状态流转。
	intent, err := s.getIntentByReq(ctx, req.GetPaymentNo(), req.GetOrderNo())
	// 查询失败时返回错误。
	if err != nil {
		return nil, err
	}
	// 已支付状态下直接返回，避免重复更新。
	if intent.GetStatus() == v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAID {
		// 返回当前快照与幂等标记。
		return &v1.HandleGatewayCallbackRes{Intent: intent, CallbackIdempotentHit: idempotentHit}, nil
	}

	// 组装支付成功更新字段。
	update := do.PaymentIntent{
		// 更新状态为 PAID。
		Status: uint(v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAID),
		// 写入外部渠道交易号。
		ExternalTradeNo: req.GetExternalTradeNo(),
		// 写入支付成功时间。
		PaidAt: pbTsToGTime(req.GetPaidAt()),
		// 递增版本号用于并发控制。
		Version: gdb.Raw("version+1"),
	}
	// 使用 payment_no + 允许前置状态做 CAS 更新。
	result, err := dao.PaymentIntent.Ctx(ctx).
		Where(dao.PaymentIntent.Columns().PaymentNo, intent.GetPaymentNo()).
		WhereIn(dao.PaymentIntent.Columns().Status, []uint{
			uint(v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_CREATED),
			uint(v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAYING),
		}).
		Data(update).
		Update()
	// 更新失败时返回数据库错误。
	if err != nil {
		return nil, gerror.Wrap(err, "update payment status failed")
	}
	// 读取受影响行数判断 CAS 是否成功。
	rows, _ := result.RowsAffected()
	// CAS 成功时回查最新快照。
	if rows > 0 {
		// 查询更新后的支付快照。
		intent, err = s.getIntentByReq(ctx, req.GetPaymentNo(), req.GetOrderNo())
		// 查询失败时返回错误。
		if err != nil {
			return nil, err
		}
	}
	if err := s.notifyOrderPayCallback(ctx, req, intent); err != nil {
		return nil, err
	}
	// 返回回调处理结果。
	return &v1.HandleGatewayCallbackRes{Intent: intent, CallbackIdempotentHit: idempotentHit}, nil
}

// ClosePaymentIntent 关闭支付意图。
func (s *sPayment) ClosePaymentIntent(ctx context.Context, req *v1.ClosePaymentIntentReq) (*v1.ClosePaymentIntentRes, error) {
	// 校验至少提供一个业务主键。
	if req == nil || (strings.TrimSpace(req.GetPaymentNo()) == "" && strings.TrimSpace(req.GetOrderNo()) == "") {
		// 参数不完整时返回错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "payment_no or order_no is required")
	}
	// 先尝试回放幂等结果，避免重复关单。
	if hit, payload, err := s.replayIdempotency(ctx, req.GetIdempotencyKey(), actionCloseIntent); err != nil {
		// 幂等查询失败时返回错误。
		return nil, err
	} else if hit {
		// 从幂等快照读取 payment_no 进行回查。
		paymentNo := gconv.String(payload["payment_no"])
		// 当快照存在 payment_no 时直接回放当前状态。
		if strings.TrimSpace(paymentNo) != "" {
			// 查询并回放当前支付快照。
			intent, getErr := s.getIntentByPaymentNo(ctx, paymentNo)
			// 查询失败时返回错误。
			if getErr != nil {
				return nil, getErr
			}
			// 依据当前状态返回对应关单结果。
			if intent.GetStatus() == v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAID {
				// 已支付场景返回 ALREADY_PAID。
				return &v1.ClosePaymentIntentRes{Intent: intent, CloseStatus: v1.GatewayCloseStatus_GATEWAY_CLOSE_STATUS_ALREADY_PAID}, nil
			}
			// 其它情况按 ALREADY_CLOSED 回放。
			return &v1.ClosePaymentIntentRes{Intent: intent, CloseStatus: v1.GatewayCloseStatus_GATEWAY_CLOSE_STATUS_ALREADY_CLOSED}, nil
		}
	}

	// 查询支付意图用于状态判断。
	intent, err := s.getIntentByReq(ctx, req.GetPaymentNo(), req.GetOrderNo())
	// 查询失败时返回错误。
	if err != nil {
		return nil, err
	}
	// 已支付订单不允许关单，直接返回状态。
	if intent.GetStatus() == v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAID {
		// 返回已支付状态。
		return &v1.ClosePaymentIntentRes{Intent: intent, CloseStatus: v1.GatewayCloseStatus_GATEWAY_CLOSE_STATUS_ALREADY_PAID}, nil
	}
	// 已关闭订单直接幂等返回。
	if intent.GetStatus() == v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_CLOSED {
		// 返回已关闭状态。
		return &v1.ClosePaymentIntentRes{Intent: intent, CloseStatus: v1.GatewayCloseStatus_GATEWAY_CLOSE_STATUS_ALREADY_CLOSED}, nil
	}

	// 组装关单更新字段。
	update := do.PaymentIntent{
		// 把状态改为 CLOSED。
		Status: uint(v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_CLOSED),
		// 写入关单时间戳。
		ClosedAt: gtime.Now(),
		// 递增版本号以支持 CAS。
		Version: gdb.Raw("version+1"),
	}
	// 只允许 CREATED/PAYING 状态迁移到 CLOSED。
	result, err := dao.PaymentIntent.Ctx(ctx).
		Where(dao.PaymentIntent.Columns().PaymentNo, intent.GetPaymentNo()).
		WhereIn(dao.PaymentIntent.Columns().Status, []uint{
			uint(v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_CREATED),
			uint(v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAYING),
		}).
		Data(update).
		Update()
	// 更新失败时返回错误。
	if err != nil {
		return nil, gerror.Wrap(err, "close payment intent failed")
	}
	// 读取受影响行数判断 CAS 命中情况。
	rows, _ := result.RowsAffected()
	// 无行更新时视为并发状态已变化，回查最新状态。
	if rows == 0 {
		// 回查最新支付快照。
		intent, err = s.getIntentByPaymentNo(ctx, intent.GetPaymentNo())
		// 回查失败时返回错误。
		if err != nil {
			return nil, err
		}
		// 按当前状态返回最合理的关闭结果。
		if intent.GetStatus() == v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAID {
			// 并发已支付场景返回 ALREADY_PAID。
			return &v1.ClosePaymentIntentRes{Intent: intent, CloseStatus: v1.GatewayCloseStatus_GATEWAY_CLOSE_STATUS_ALREADY_PAID}, nil
		}
		// 并发已关闭场景返回 ALREADY_CLOSED。
		return &v1.ClosePaymentIntentRes{Intent: intent, CloseStatus: v1.GatewayCloseStatus_GATEWAY_CLOSE_STATUS_ALREADY_CLOSED}, nil
	}
	// 回查关单后的最新快照。
	intent, err = s.getIntentByPaymentNo(ctx, intent.GetPaymentNo())
	// 回查失败时返回错误。
	if err != nil {
		return nil, err
	}
	// 保存关单幂等快照用于重试回放。
	if err = s.saveIdempotency(ctx, req.GetIdempotencyKey(), actionCloseIntent, intent.GetPaymentNo(), map[string]any{"payment_no": intent.GetPaymentNo()}); err != nil {
		// 幂等快照失败不影响主流程，仅记录告警。
		g.Log().Warningf(ctx, "save idempotency failed: %v", err)
	}
	// 返回关单成功结果。
	return &v1.ClosePaymentIntentRes{Intent: intent, CloseStatus: v1.GatewayCloseStatus_GATEWAY_CLOSE_STATUS_CLOSED}, nil
}

// CreateRefundTask 创建退款任务并执行防超退校验。
func (s *sPayment) CreateRefundTask(ctx context.Context, req *v1.CreateRefundTaskReq) (*v1.CreateRefundTaskRes, error) {
	// 校验核心参数，保证退款请求完整。
	if req == nil || strings.TrimSpace(req.GetAfterSaleNo()) == "" || strings.TrimSpace(req.GetPaymentNo()) == "" || req.GetRefundAmount() == 0 {
		// 参数非法时返回错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "after_sale_no/payment_no/refund_amount are required")
	}
	// 先尝试回放幂等结果，避免重复创建退款任务。
	if hit, payload, err := s.replayIdempotency(ctx, req.GetIdempotencyKey(), actionCreateRefund); err != nil {
		// 幂等查询失败时返回错误。
		return nil, err
	} else if hit {
		// 从幂等快照读取历史退款任务号。
		refundTaskNo := gconv.String(payload["refund_task_no"])
		// 历史任务号存在时直接回放成功状态。
		if strings.TrimSpace(refundTaskNo) != "" {
			// 返回幂等回放结果。
			return &v1.CreateRefundTaskRes{RefundTaskNo: refundTaskNo, Status: v1.RefundTaskStatus_REFUND_TASK_STATUS_PENDING}, nil
		}
	}

	// 查询支付意图，校验退款上下文。
	intent, err := s.getIntentByPaymentNo(ctx, req.GetPaymentNo())
	// 查询失败时返回错误。
	if err != nil {
		return nil, err
	}
	// 若本次退款会超过应付金额则直接阻断。
	if intent.GetRefundedAmount()+req.GetRefundAmount() > intent.GetPayableAmount() {
		// 返回业务校验失败避免资损。
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "refund exceeds payable amount")
	}

	// 生成退款任务号。
	refundTaskNo := newBizNo("RFD")
	// 插入退款任务并利用 after_sale_no 唯一键防重。
	r, err := dao.RefundTask.Ctx(ctx).Data(do.RefundTask{
		// 写入退款任务号。
		RefundTaskNo: refundTaskNo,
		// 写入售后单号作为唯一防重键。
		AfterSaleNo: req.GetAfterSaleNo(),
		// 写入订单号用于追踪。
		OrderNo: req.GetOrderNo(),
		// 写入支付单号用于资金关联。
		PaymentNo: req.GetPaymentNo(),
		// 写入本次退款金额。
		RefundAmount: req.GetRefundAmount(),
		// 初始化状态为 PENDING。
		Status: uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_PENDING),
	}).InsertIgnore()
	// 任务写入失败时返回错误。
	if err != nil {
		return nil, gerror.Wrap(err, "insert refund task failed")
	}
	// 处理 after_sale_no 重复写入场景。
	if r != nil {
		// 判断是否发生唯一键冲突。
		if rows, _ := r.RowsAffected(); rows == 0 {
			// 查询已存在的退款任务并回放。
			var old entity.RefundTask
			// 按售后单号回查旧任务。
			if err = dao.RefundTask.Ctx(ctx).Where(dao.RefundTask.Columns().AfterSaleNo, req.GetAfterSaleNo()).Scan(&old); err != nil {
				// 回查失败时返回错误。
				return nil, gerror.Wrap(err, "query refund task failed")
			}
			// 返回已存在任务实现自然幂等。
			return &v1.CreateRefundTaskRes{RefundTaskNo: old.RefundTaskNo, Status: v1.RefundTaskStatus(old.Status)}, nil
		}
	}

	// 计算退款后的目标支付状态。
	nextStatus := uint(v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_REFUNDING)
	// 当累计退款等于应付金额时标记为全额已退。
	if intent.GetRefundedAmount()+req.GetRefundAmount() == intent.GetPayableAmount() {
		// 切换到 REFUNDED 状态。
		nextStatus = uint(v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_REFUNDED)
	}
	// 使用 CAS 更新已退款金额并保证“已退<=应付”。
	cols := dao.PaymentIntent.Columns()
	result, err := dao.PaymentIntent.Ctx(ctx).
		Where(cols.PaymentNo, req.GetPaymentNo()).
		Where(fmt.Sprintf("%s + ? <= %s", cols.RefundedAmount, cols.PayableAmount), req.GetRefundAmount()).
		Data(do.PaymentIntent{
			// 原子累加已退款金额。
			RefundedAmount: gdb.Raw(fmt.Sprintf("%s + %d", cols.RefundedAmount, req.GetRefundAmount())),
			// 同步推进支付状态到退款中/已退款。
			Status: nextStatus,
			// 递增版本号支持并发控制。
			Version: gdb.Raw("version+1"),
		}).
		Update()
	// 更新失败时返回错误。
	if err != nil {
		return nil, gerror.Wrap(err, "cas update refunded_amount failed")
	}
	// 当 CAS 未命中时返回业务校验失败。
	if rows, _ := result.RowsAffected(); rows == 0 {
		// 返回防超退阻断错误。
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "refund cas check failed")
	}
	// 保存创建退款任务幂等快照。
	if err = s.saveIdempotency(ctx, req.GetIdempotencyKey(), actionCreateRefund, refundTaskNo, map[string]any{"refund_task_no": refundTaskNo}); err != nil {
		// 幂等快照失败仅告警不阻塞主流程。
		g.Log().Warningf(ctx, "save idempotency failed: %v", err)
	}
	// 返回退款任务创建成功。
	return &v1.CreateRefundTaskRes{RefundTaskNo: refundTaskNo, Status: v1.RefundTaskStatus_REFUND_TASK_STATUS_PENDING}, nil
}

// ExecuteRefundTask 执行退款任务（当前为 mock 成功）。
func (s *sPayment) ExecuteRefundTask(ctx context.Context, req *v1.ExecuteRefundTaskReq) (*v1.ExecuteRefundTaskRes, error) {
	// 校验退款任务号。
	if req == nil || strings.TrimSpace(req.GetRefundTaskNo()) == "" {
		// 参数缺失时返回错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund_task_no is required")
	}
	// 先尝试回放幂等结果，避免重复执行任务。
	if hit, payload, err := s.replayIdempotency(ctx, req.GetIdempotencyKey(), actionExecuteRefund); err != nil {
		// 幂等查询失败时返回错误。
		return nil, err
	} else if hit {
		// 回放历史退款任务号。
		refundTaskNo := gconv.String(payload["refund_task_no"])
		// 回放历史已退款金额。
		latestRefundedAmount := gconv.Uint64(payload["latest_refunded_amount"])
		// 快照存在时直接返回历史成功结果。
		if strings.TrimSpace(refundTaskNo) != "" {
			// 返回幂等回放结果。
			return &v1.ExecuteRefundTaskRes{RefundTaskNo: refundTaskNo, Status: v1.RefundTaskStatus_REFUND_TASK_STATUS_SUCCEEDED, LatestRefundedAmount: latestRefundedAmount}, nil
		}
	}

	// 查询退款任务实体。
	var task entity.RefundTask
	// 按退款任务号查库。
	if err := dao.RefundTask.Ctx(ctx).Where(dao.RefundTask.Columns().RefundTaskNo, req.GetRefundTaskNo()).Scan(&task); err != nil {
		// 查询失败时返回错误。
		return nil, gerror.Wrap(err, "query refund task failed")
	}
	// 任务不存在时返回 not found。
	if task.RefundTaskNo == "" {
		// 返回任务不存在错误。
		return nil, gerror.NewCode(gcode.CodeNotFound, "refund task not found")
	}
	// 已成功任务直接幂等返回。
	if task.Status == uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_SUCCEEDED) {
		// 查询支付快照获取最新已退款金额。
		intent, err := s.getIntentByPaymentNo(ctx, task.PaymentNo)
		// 查询失败时返回错误。
		if err != nil {
			return nil, err
		}
		// 返回历史成功执行结果。
		return &v1.ExecuteRefundTaskRes{RefundTaskNo: req.GetRefundTaskNo(), Status: v1.RefundTaskStatus_REFUND_TASK_STATUS_SUCCEEDED, LatestRefundedAmount: intent.GetRefundedAmount()}, nil
	}

	// 把任务状态从 PENDING/FAILED 原子切换到 PROCESSING。
	_, err := dao.RefundTask.Ctx(ctx).
		Where(dao.RefundTask.Columns().RefundTaskNo, req.GetRefundTaskNo()).
		WhereIn(dao.RefundTask.Columns().Status, []uint{
			uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_PENDING),
			uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_FAILED),
			uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_PROCESSING),
		}).
		Data(do.RefundTask{Status: uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_PROCESSING)}).
		Update()
	// 状态推进失败时返回错误。
	if err != nil {
		return nil, gerror.Wrap(err, "mark refund task processing failed")
	}
	// mock 支付网关退款成功后写回任务状态。
	if _, err = dao.RefundTask.Ctx(ctx).Where(dao.RefundTask.Columns().RefundTaskNo, req.GetRefundTaskNo()).Data(do.RefundTask{Status: uint(v1.RefundTaskStatus_REFUND_TASK_STATUS_SUCCEEDED)}).Update(); err != nil {
		// 更新失败时返回错误。
		return nil, gerror.Wrap(err, "update refund task failed")
	}

	// 查询支付快照以返回最新退款金额。
	intent, err := s.getIntentByPaymentNo(ctx, task.PaymentNo)
	// 查询失败时返回错误。
	if err != nil {
		return nil, err
	}
	// 根据累计退款金额刷新支付状态。
	payStatus := uint(v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_REFUNDING)
	// 全额退款时切换为 REFUNDED。
	if intent.GetRefundedAmount() >= intent.GetPayableAmount() {
		// 更新目标状态为已退款。
		payStatus = uint(v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_REFUNDED)
	}
	// 写回支付状态并递增版本号。
	if _, err = dao.PaymentIntent.Ctx(ctx).Where(dao.PaymentIntent.Columns().PaymentNo, task.PaymentNo).Data(do.PaymentIntent{Status: payStatus, Version: gdb.Raw("version+1")}).Update(); err != nil {
		// 写回失败时返回错误。
		return nil, gerror.Wrap(err, "update payment status for refund failed")
	}
	// 保存执行退款任务幂等快照。
	if err = s.saveIdempotency(ctx, req.GetIdempotencyKey(), actionExecuteRefund, req.GetRefundTaskNo(), map[string]any{"refund_task_no": req.GetRefundTaskNo(), "latest_refunded_amount": intent.GetRefundedAmount()}); err != nil {
		// 幂等快照失败仅告警不影响主链路。
		g.Log().Warningf(ctx, "save idempotency failed: %v", err)
	}
	// 返回任务执行结果。
	return &v1.ExecuteRefundTaskRes{RefundTaskNo: req.GetRefundTaskNo(), Status: v1.RefundTaskStatus_REFUND_TASK_STATUS_SUCCEEDED, LatestRefundedAmount: intent.GetRefundedAmount()}, nil
}

// RunDailyReconciliation 创建每日对账任务。
func (s *sPayment) RunDailyReconciliation(ctx context.Context, req *v1.RunDailyReconciliationReq) (*v1.RunDailyReconciliationRes, error) {
	// 先尝试回放幂等结果，避免重复创建同次请求任务。
	if hit, payload, err := s.replayIdempotency(ctx, req.GetIdempotencyKey(), actionRunRecon); err != nil {
		// 幂等查询失败时返回错误。
		return nil, err
	} else if hit {
		// 从快照读取历史对账任务号。
		reconTaskNo := gconv.String(payload["recon_task_no"])
		// 当历史任务存在时回放成功结果。
		if strings.TrimSpace(reconTaskNo) != "" {
			// 返回幂等回放结果。
			return &v1.RunDailyReconciliationRes{ReconTaskNo: reconTaskNo, TotalRecords: gconv.Uint64(payload["total_records"]), DiffRecords: gconv.Uint64(payload["diff_records"])}, nil
		}
	}

	// 生成对账任务号。
	reconTaskNo := newBizNo("RCN")
	// 解析入参日期字符串。
	reconDate, _ := gtime.StrToTime(req.GetReconDate())
	// 未传日期时默认使用当天日期。
	if reconDate == nil {
		// 使用当前时间作为对账日期兜底。
		reconDate = gtime.Now()
	}
	// 写入对账任务主记录。
	if _, err := dao.PaymentReconciliationTask.Ctx(ctx).Data(do.PaymentReconciliationTask{ReconTaskNo: reconTaskNo, ReconDate: reconDate, PayChannel: uint(req.GetPayChannel()), Status: "PENDING", TotalRecords: 0, DiffRecords: 0}).Insert(); err != nil {
		// 插入失败时返回错误。
		return nil, gerror.Wrap(err, "insert reconciliation task failed")
	}
	// 保存幂等快照供客户端重试回放。
	if err := s.saveIdempotency(ctx, req.GetIdempotencyKey(), actionRunRecon, reconTaskNo, map[string]any{"recon_task_no": reconTaskNo, "total_records": 0, "diff_records": 0}); err != nil {
		// 幂等快照失败仅告警不影响主流程。
		g.Log().Warningf(ctx, "save idempotency failed: %v", err)
	}
	// 返回对账任务创建成功。
	return &v1.RunDailyReconciliationRes{ReconTaskNo: reconTaskNo, TotalRecords: 0, DiffRecords: 0}, nil
}

// ListReconciliationDiffs 查询对账差异列表。
func (s *sPayment) ListReconciliationDiffs(ctx context.Context, req *v1.ListReconciliationDiffsReq) (*v1.ListReconciliationDiffsRes, error) {
	// 计算安全分页大小，限制单页最大 100。
	pageSize := req.GetPageSize()
	// 当 page_size 非法时回退默认值 20。
	if pageSize <= 0 || pageSize > 100 {
		// 应用默认分页大小。
		pageSize = 20
	}
	// 初始化查询模型并按自增 ID 倒序。
	m := dao.PaymentReconciliationRecord.Ctx(ctx).Where(dao.PaymentReconciliationRecord.Columns().ReconTaskNo, req.GetReconTaskNo()).OrderDesc(dao.PaymentReconciliationRecord.Columns().Id).Limit(int(pageSize) + 1)
	// 当携带游标时应用游标过滤条件。
	if strings.TrimSpace(req.GetNextCursor()) != "" {
		// 只读取比游标更早的数据记录。
		m = m.WhereLT(dao.PaymentReconciliationRecord.Columns().Id, gconv.Int64(req.GetNextCursor()))
	}
	// 执行差异记录查询。
	var rows []*entity.PaymentReconciliationRecord
	// 扫描结果到实体切片。
	if err := m.Scan(&rows); err != nil {
		// 查询失败时返回错误。
		return nil, gerror.Wrap(err, "query reconciliation records failed")
	}
	// 默认标记为无下一页。
	hasMore := false
	// 超过 pageSize 说明存在下一页数据。
	if len(rows) > int(pageSize) {
		// 设置 hasMore 标记。
		hasMore = true
		// 裁剪为当前页大小。
		rows = rows[:pageSize]
	}
	// 预分配响应列表容量减少扩容。
	list := make([]*v1.ReconciliationDiff, 0, len(rows))
	// 逐条转换数据库实体到 protobuf 结构。
	for _, row := range rows {
		// 追加单条差异记录到响应列表。
		list = append(list, &v1.ReconciliationDiff{DiffNo: row.DiffNo, ReconTaskNo: row.ReconTaskNo, DiffType: v1.ReconciliationDiffType(row.DiffType), PaymentNo: row.PaymentNo, OrderNo: row.OrderNo, LocalAmount: row.LocalAmount, GatewayAmount: row.GatewayAmount, Status: row.Status, DetailJson: row.DetailJson, CreatedAt: gtimeToPB(row.CreatedAt), ResolvedAt: gtimeToPB(row.ResolvedAt)})
	}
	// 默认下一页游标为空。
	nextCursor := ""
	// 仅在存在下一页且当前页非空时生成游标。
	if hasMore && len(rows) > 0 {
		// 用最后一条记录 ID 作为 next_cursor。
		nextCursor = gconv.String(rows[len(rows)-1].Id)
	}
	// 返回分页查询结果。
	return &v1.ListReconciliationDiffsRes{List: list, NextCursor: nextCursor, HasMore: hasMore}, nil
}

// ResolveReconciliationDiff 把差异单标记为已处理。
func (s *sPayment) ResolveReconciliationDiff(ctx context.Context, req *v1.ResolveReconciliationDiffReq) (*v1.ResolveReconciliationDiffRes, error) {
	// 校验差异单号参数。
	if req == nil || strings.TrimSpace(req.GetDiffNo()) == "" {
		// 参数缺失时返回错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "diff_no is required")
	}
	// 更新差异记录状态与处理时间。
	if _, err := dao.PaymentReconciliationRecord.Ctx(ctx).Where(dao.PaymentReconciliationRecord.Columns().DiffNo, req.GetDiffNo()).Data(do.PaymentReconciliationRecord{Status: "RESOLVED", ResolvedAt: gtime.Now()}).Update(); err != nil {
		// 更新失败时返回错误。
		return nil, gerror.Wrap(err, "resolve reconciliation diff failed")
	}
	// 返回处理完成结果。
	return &v1.ResolveReconciliationDiffRes{DiffNo: req.GetDiffNo(), Status: "RESOLVED"}, nil
}

// GetPaymentSnapshot 查询支付快照供内部服务读取。
func (s *sPayment) GetPaymentSnapshot(ctx context.Context, req *v1.GetPaymentSnapshotReq) (*v1.GetPaymentSnapshotRes, error) {
	// 使用统一查询逻辑按 payment_no/order_no 获取快照。
	intent, err := s.getIntentByReq(ctx, req.GetPaymentNo(), req.GetOrderNo())
	// 查询失败时返回错误。
	if err != nil {
		return nil, err
	}
	// 返回支付快照结果。
	return &v1.GetPaymentSnapshotRes{Intent: intent}, nil
}

// replayIdempotency 查询幂等快照并反序列化为 map。
func (s *sPayment) replayIdempotency(ctx context.Context, idemKey string, actionCode string) (bool, map[string]any, error) {
	// 当幂等键为空时直接跳过幂等处理。
	if strings.TrimSpace(idemKey) == "" {
		// 返回未命中幂等。
		return false, nil, nil
	}
	// 查询幂等记录。
	var row entity.PaymentIdempotency
	// 按幂等键+动作码读取唯一快照。
	if err := dao.PaymentIdempotency.Ctx(ctx).Where(do.PaymentIdempotency{IdemKey: idemKey, BizCode: actionCode}).Scan(&row); err != nil {
		if isIdempotencyMissError(err) {
			return false, nil, nil
		}
		// 查询失败时返回错误。
		return false, nil, gerror.Wrap(err, "query idempotency failed")
	}
	// 未命中记录时返回 false。
	if row.Id == 0 {
		// 返回未命中结果。
		return false, nil, nil
	}
	// 初始化空 payload 容器。
	payload := map[string]any{}
	// 当 response_json 非空时尝试反序列化。
	if strings.TrimSpace(row.ResponseJson) != "" {
		// 反序列化失败时返回错误以避免脏回放。
		if err := json.Unmarshal([]byte(row.ResponseJson), &payload); err != nil {
			// 返回解析异常。
			return false, nil, gerror.Wrap(err, "unmarshal idempotency payload failed")
		}
	}
	// 返回幂等命中与快照内容。
	return true, payload, nil
}

func isIdempotencyMissError(err error) bool {
	if err == nil {
		return false
	}
	if gerror.HasCode(err, gcode.CodeNotFound) {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), sql.ErrNoRows.Error())
}

// saveIdempotency 写入幂等快照。
func (s *sPayment) saveIdempotency(ctx context.Context, idemKey string, actionCode string, bizNo string, payload map[string]any) error {
	// 幂等键为空时不落幂等表。
	if strings.TrimSpace(idemKey) == "" {
		// 直接返回成功。
		return nil
	}
	// 序列化快照用于回放。
	body, err := json.Marshal(payload)
	// 序列化失败时返回错误。
	if err != nil {
		return gerror.Wrap(err, "marshal idempotency payload failed")
	}
	// 插入幂等记录并忽略重复。
	_, err = dao.PaymentIdempotency.Ctx(ctx).Data(do.PaymentIdempotency{IdemKey: idemKey, BizCode: actionCode, BizNo: bizNo, ResponseJson: string(body)}).InsertIgnore()
	// 返回写入结果。
	return err
}

// getIntentByReq 按 payment_no 或 order_no 查询支付快照。
func (s *sPayment) getIntentByReq(ctx context.Context, paymentNo string, orderNo string) (*v1.PaymentIntent, error) {
	// 优先按 payment_no 查询，确保命中唯一记录。
	if strings.TrimSpace(paymentNo) != "" {
		// 按 payment_no 查询快照。
		return s.getIntentByPaymentNo(ctx, paymentNo)
	}
	// 退化为按 order_no 查询快照。
	return s.getIntentByOrderNo(ctx, orderNo)
}

// getIntentByPaymentNo 按支付单号查询支付意图。
func (s *sPayment) getIntentEntityByOrderNo(ctx context.Context, orderNo string) (*entity.PaymentIntent, error) {
	var row entity.PaymentIntent
	if err := dao.PaymentIntent.Ctx(ctx).Where(dao.PaymentIntent.Columns().OrderNo, orderNo).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query payment_intent entity by order_no failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "payment intent not found")
	}
	return &row, nil
}

func canReuseExistingIntentForCreate(row *entity.PaymentIntent) bool {
	if row == nil {
		return false
	}
	switch v1.PaymentIntentStatus(row.Status) {
	case v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_CREATED, v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAYING:
		return true
	default:
		return false
	}
}

func (s *sPayment) refreshExistingIntentForCreate(ctx context.Context, row *entity.PaymentIntent, req *v1.CreatePaymentIntentReq, orderExpireAt, gatewayExpireAt *gtime.Time) error {
	if row == nil || row.Id == 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "payment intent row is required")
	}
	cols := dao.PaymentIntent.Columns()
	_, err := dao.PaymentIntent.Ctx(ctx).
		Where(cols.Id, row.Id).
		Where(cols.OrderNo, row.OrderNo).
		Data(do.PaymentIntent{
			PayChannel:      uint(req.GetPayChannel()),
			Status:          uint(v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAYING),
			PayableAmount:   req.GetPayableAmount(),
			CurrencyCode:    safeCurrency(req.GetCurrencyCode()),
			OrderExpireAt:   orderExpireAt,
			GatewayExpireAt: gatewayExpireAt,
			Version:         gdb.Raw(cols.Version + " + 1"),
		}).
		Update()
	return gerror.Wrap(err, "refresh existing payment_intent failed")
}

func (s *sPayment) getIntentByPaymentNo(ctx context.Context, paymentNo string) (*v1.PaymentIntent, error) {
	// 查询数据库实体。
	var row entity.PaymentIntent
	// 按 payment_no 扫描记录。
	if err := dao.PaymentIntent.Ctx(ctx).Where(dao.PaymentIntent.Columns().PaymentNo, paymentNo).Scan(&row); err != nil {
		// 查询失败时返回错误。
		return nil, gerror.Wrap(err, "query payment_intent by payment_no failed")
	}
	// 记录不存在时返回 not found。
	if row.Id == 0 {
		// 返回资源不存在错误。
		return nil, gerror.NewCode(gcode.CodeNotFound, "payment intent not found")
	}
	// 转换为 protobuf 结构返回。
	return toPBIntent(&row), nil
}

// getIntentByOrderNo 按订单号查询支付意图。
func (s *sPayment) getIntentByOrderNo(ctx context.Context, orderNo string) (*v1.PaymentIntent, error) {
	// 查询数据库实体。
	var row entity.PaymentIntent
	// 按 order_no 扫描记录。
	if err := dao.PaymentIntent.Ctx(ctx).Where(dao.PaymentIntent.Columns().OrderNo, orderNo).Scan(&row); err != nil {
		// 查询失败时返回错误。
		return nil, gerror.Wrap(err, "query payment_intent by order_no failed")
	}
	// 记录不存在时返回 not found。
	if row.Id == 0 {
		// 返回资源不存在错误。
		return nil, gerror.NewCode(gcode.CodeNotFound, "payment intent not found")
	}
	// 转换为 protobuf 结构返回。
	return toPBIntent(&row), nil
}

// toPBIntent 把数据库实体转换为 protobuf 响应。
func buildOrderPayCallbackRequest(req *v1.HandleGatewayCallbackReq, intent *v1.PaymentIntent) *orderv1.HandlePayCallbackReq {
	if req == nil || intent == nil {
		return nil
	}
	return &orderv1.HandlePayCallbackReq{
		PaymentEventId: firstNonEmptyString(req.GetCallbackEventId(), fmt.Sprintf("payment-callback-%s", intent.GetPaymentNo())),
		PayNo:          firstNonEmptyString(req.GetPaymentNo(), intent.GetPaymentNo()),
		OrderNo:        firstNonEmptyString(req.GetOrderNo(), intent.GetOrderNo()),
		PayChannel:     orderv1.PayChannel(req.GetPayChannel()),
		PayStatusCode:  strings.ToUpper(strings.TrimSpace(req.GetGatewayStatusCode())),
		ChannelTradeNo: strings.TrimSpace(req.GetExternalTradeNo()),
		PaidAmount:     req.GetPaidAmount(),
		PaidAt:         req.GetPaidAt(),
		RawPayload:     req.GetRawPayload(),
		IdempotencyKey: fmt.Sprintf("payment-callback:%s", firstNonEmptyString(req.GetCallbackEventId(), intent.GetPaymentNo())),
	}
}

func (s *sPayment) notifyOrderPayCallback(ctx context.Context, req *v1.HandleGatewayCallbackReq, intent *v1.PaymentIntent) error {
	callbackReq := buildOrderPayCallbackRequest(req, intent)
	if callbackReq == nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "order callback request is required")
	}
	client, err := newInternalOrderClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.HandlePayCallback(ctx, callbackReq)
	return gerror.Wrap(err, "notify order pay callback failed")
}

func defaultInternalOrderClient(ctx context.Context) (internalOrderCallbackClient, error) {
	addr := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.orderGrpc", "127.0.0.1:9008").String())
	if addr == "" {
		addr = "127.0.0.1:9008"
	}

	internalOrderClientMu.Lock()
	defer internalOrderClientMu.Unlock()

	if internalOrderClientInst != nil && internalOrderClientConn != nil && internalOrderClientAddr == addr {
		return internalOrderClientInst, nil
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(timeoutCtx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, gerror.Wrapf(err, "dial order-svc failed, addr=%s", addr)
	}

	if internalOrderClientConn != nil {
		_ = internalOrderClientConn.Close()
	}
	internalOrderClientConn = conn
	internalOrderClientAddr = addr
	internalOrderClientInst = orderv1.NewInternalOrderServiceClient(conn)
	return internalOrderClientInst, nil
}

func toPBIntent(row *entity.PaymentIntent) *v1.PaymentIntent {
	// 保护性判空避免空指针。
	if row == nil {
		// 空实体直接返回 nil。
		return nil
	}
	// 执行字段映射并返回响应对象。
	return &v1.PaymentIntent{
		PaymentNo:       row.PaymentNo,
		OrderNo:         row.OrderNo,
		UserId:          row.UserId,
		PayChannel:      v1.PayChannel(row.PayChannel),
		Status:          v1.PaymentIntentStatus(row.Status),
		PayableAmount:   row.PayableAmount,
		RefundedAmount:  row.RefundedAmount,
		CurrencyCode:    row.CurrencyCode,
		OrderExpireAt:   gtimeToPB(row.OrderExpireAt),
		GatewayExpireAt: gtimeToPB(row.GatewayExpireAt),
		ExternalTradeNo: row.ExternalTradeNo,
		PaidAt:          gtimeToPB(row.PaidAt),
		ClosedAt:        gtimeToPB(row.ClosedAt),
		CreatedAt:       gtimeToPB(row.CreatedAt),
		UpdatedAt:       gtimeToPB(row.UpdatedAt),
		Version:         row.Version,
	}
}

// pbTsToGTime 把 protobuf 时间转换为 gtime。
func pbTsToGTime(ts *timestamppb.Timestamp) *gtime.Time {
	// 空时间戳直接返回 nil。
	if ts == nil {
		// 返回空指针表示未传时间。
		return nil
	}
	// 按 protobuf 时间创建 gtime 对象。
	return gtime.New(ts.AsTime())
}

// gtimeToPB 把 gtime 转换为 protobuf 时间。
func gtimeToPB(t *gtime.Time) *timestamppb.Timestamp {
	// 空时间直接返回 nil。
	if t == nil {
		// 返回空时间戳。
		return nil
	}
	// 生成 protobuf 时间对象。
	return timestamppb.New(t.Time)
}

// safeCurrency 规范化币种编码。
func safeCurrency(input string) string {
	// 去除输入两端空白字符。
	ccy := strings.TrimSpace(input)
	// 当调用方未传币种时默认 CNY。
	if ccy == "" {
		// 返回默认币种代码。
		return "CNY"
	}
	// 把币种统一转成大写。
	return strings.ToUpper(ccy)
}

func (s *sPayment) buildPaymentLaunch(ctx context.Context, req *v1.CreatePaymentIntentReq, intent *v1.PaymentIntent) (*paymentLaunch, error) {
	if intent == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "payment intent is required")
	}
	switch intent.GetPayChannel() {
	case v1.PayChannel_PAY_CHANNEL_MOCK:
		return &paymentLaunch{
			PayURL:         fmt.Sprintf("mock://pay/%s", intent.GetPaymentNo()),
			PayPayloadJSON: `{"gateway":"mock"}`,
		}, nil
	case v1.PayChannel_PAY_CHANNEL_ALIPAY:
		return buildAlipayPagePaymentLaunch(loadAlipayConfig(ctx), req, intent)
	default:
		return nil, gerror.NewCodef(gcode.CodeNotSupported, "pay channel %s not supported yet", intent.GetPayChannel().String())
	}
}

func loadAlipayConfig(ctx context.Context) alipayConfig {
	return alipayConfig{
		GatewayURL:    cfgEnvString(ctx, "payment.alipay.gatewayUrl", "https://openapi.alipay.com/gateway.do", "SHOPA_ALIPAY_GATEWAY_URL", "ALIPAY_GATEWAY_URL"),
		AppID:         cfgEnvString(ctx, "payment.alipay.appId", "", "SHOPA_ALIPAY_APP_ID", "ALIPAY_APP_ID"),
		PrivateKey:    cfgEnvString(ctx, "payment.alipay.privateKey", "", "SHOPA_ALIPAY_PRIVATE_KEY", "ALIPAY_PRIVATE_KEY"),
		PublicKey:     cfgEnvString(ctx, "payment.alipay.publicKey", "", "SHOPA_ALIPAY_PUBLIC_KEY", "ALIPAY_PUBLIC_KEY"),
		ReturnURL:     cfgEnvString(ctx, "payment.alipay.returnUrl", "", "SHOPA_ALIPAY_RETURN_URL", "ALIPAY_RETURN_URL"),
		NotifyURL:     cfgEnvString(ctx, "payment.alipay.notifyUrl", "", "SHOPA_ALIPAY_NOTIFY_URL", "ALIPAY_NOTIFY_URL"),
		SubjectPrefix: cfgEnvString(ctx, "payment.alipay.subjectPrefix", "Shopa", "SHOPA_ALIPAY_SUBJECT_PREFIX", "ALIPAY_SUBJECT_PREFIX"),
	}
}

func cfgEnvString(ctx context.Context, cfgKey, defaultValue string, envKeys ...string) string {
	for _, envKey := range envKeys {
		if value := strings.TrimSpace(os.Getenv(envKey)); value != "" {
			return value
		}
	}
	return strings.TrimSpace(g.Cfg().MustGet(ctx, cfgKey, defaultValue).String())
}

func buildAlipayPagePaymentLaunch(cfg alipayConfig, req *v1.CreatePaymentIntentReq, intent *v1.PaymentIntent) (*paymentLaunch, error) {
	if strings.TrimSpace(cfg.AppID) == "" || strings.TrimSpace(cfg.PrivateKey) == "" {
		return nil, gerror.NewCode(gcode.CodeInternalError, "alipay config missing payment.alipay.appId/payment.alipay.privateKey")
	}

	now := time.Now().In(alipayTimeLocation())
	subject := firstNonEmptyString(strings.TrimSpace(req.GetSubject()), fmt.Sprintf("%s %s", firstNonEmptyString(cfg.SubjectPrefix, "Shopa"), strings.TrimSpace(intent.GetOrderNo())))
	returnURL := firstNonEmptyString(strings.TrimSpace(req.GetReturnUrl()), cfg.ReturnURL)
	notifyURL := firstNonEmptyString(strings.TrimSpace(req.GetNotifyUrl()), cfg.NotifyURL)

	bizContentBody, err := json.Marshal(map[string]string{
		"out_trade_no": strings.TrimSpace(intent.GetPaymentNo()),
		"product_code": "FAST_INSTANT_TRADE_PAY",
		"subject":      subject,
		"total_amount": amountFenToYuanString(intent.GetPayableAmount()),
		"time_expire":  formatAlipayTimestamp(intent.GetGatewayExpireAt()),
	})
	if err != nil {
		return nil, gerror.Wrap(err, "marshal alipay biz_content failed")
	}

	params := map[string]string{
		"app_id":      cfg.AppID,
		"biz_content": string(bizContentBody),
		"charset":     "utf-8",
		"format":      "JSON",
		"method":      "alipay.trade.page.pay",
		"notify_url":  notifyURL,
		"return_url":  returnURL,
		"sign_type":   "RSA2",
		"timestamp":   now.Format("2006-01-02 15:04:05"),
		"version":     "1.0",
	}
	sign, err := signAlipayParams(params, cfg.PrivateKey)
	if err != nil {
		return nil, err
	}
	params["sign"] = sign

	values := url.Values{}
	for key, value := range params {
		if strings.TrimSpace(value) == "" {
			continue
		}
		values.Set(key, value)
	}
	payURL := strings.TrimRight(cfg.GatewayURL, "?") + "?" + values.Encode()

	payloadBody, err := json.Marshal(map[string]string{
		"gateway": "alipay",
		"method":  "GET",
		"url":     payURL,
	})
	if err != nil {
		return nil, gerror.Wrap(err, "marshal alipay payload failed")
	}
	return &paymentLaunch{
		PayURL:         payURL,
		PayPayloadJSON: string(payloadBody),
	}, nil
}

func signAlipayParams(params map[string]string, privateKeyPEM string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", gerror.NewCode(gcode.CodeInternalError, "invalid alipay private key pem")
	}
	keyAny, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", gerror.Wrap(err, "parse alipay private key failed")
	}
	privateKey, ok := keyAny.(*rsa.PrivateKey)
	if !ok {
		return "", gerror.NewCode(gcode.CodeInternalError, "alipay private key is not rsa")
	}
	signContent := buildAlipaySignContent(params)
	hash := sha256.Sum256([]byte(signContent))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", gerror.Wrap(err, "sign alipay params failed")
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func buildAlipaySignContent(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for key, value := range params {
		if key == "sign" || strings.TrimSpace(value) == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", key, params[key]))
	}
	return strings.Join(parts, "&")
}

func amountFenToYuanString(amount uint64) string {
	return fmt.Sprintf("%.2f", float64(amount)/100)
}

func formatAlipayTimestamp(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().In(alipayTimeLocation()).Format("2006-01-02 15:04:05")
}

func alipayTimeLocation() *time.Location {
	return time.FixedZone("CST", 8*3600)
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

// newBizNo 生成简单业务单号。
func newBizNo(prefix string) string {
	// 使用前缀 + 纳秒时间生成近似唯一单号。
	return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
}
