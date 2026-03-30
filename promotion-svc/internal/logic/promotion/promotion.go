package promotion

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	v1 "github.com/TsingpekTao/shopa/promotion-svc/api/v1"
	"github.com/TsingpekTao/shopa/promotion-svc/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	// lockSeq 用于生成进程内自增序列，避免并发锁号冲突。
	lockSeq uint64
)

// sPromotion 提供促销服务接口的最小可运行实现。
type sPromotion struct{}

func init() {
	// 在包初始化阶段注册实现，确保 service.Promotion() 可用。
	service.RegisterPromotion(New())
}

// New 返回促销逻辑实例，供 service 层统一调用。
func New() *sPromotion {
	// 返回空结构体实例承载接口方法。
	return &sPromotion{}
}

func (s *sPromotion) ListAvailableCoupons(_ context.Context, req *v1.ListAvailableCouponsReq) (*v1.ListAvailableCouponsRes, error) {
	// 校验请求对象，防止空指针访问。
	if req == nil {
		// 返回标准参数错误，提示调用方修正请求。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter)
	}
	// 当前最小实现返回空列表，保证链路可用。
	return &v1.ListAvailableCouponsRes{List: make([]*v1.Coupon, 0)}, nil
}

func (s *sPromotion) CalculateOrderDiscount(_ context.Context, req *v1.CalculateOrderDiscountReq) (*v1.CalculateOrderDiscountRes, error) {
	// 校验请求对象，避免空指针异常。
	if req == nil {
		// 返回参数错误码，便于上游统一处理。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter)
	}
	// 当前版本未启用真实优惠规则，折扣金额固定为 0。
	totalDiscountAmount := uint64(0)
	// 可支付金额初始值取商品总额。
	payableAmount := req.GetGoodsAmount()
	// 防御式校验，保证不会出现无符号下溢。
	if payableAmount >= totalDiscountAmount {
		// 正常场景下按总额减折扣计算实付。
		payableAmount -= totalDiscountAmount
	} else {
		// 异常场景兜底归零，避免出现非法金额。
		payableAmount = 0
	}
	// 返回结构化金额结果，满足结算链路入参约定。
	return &v1.CalculateOrderDiscountRes{
		GoodsAmount:         req.GetGoodsAmount(),
		TotalDiscountAmount: totalDiscountAmount,
		PayableAmount:       payableAmount,
		DiscountLines:       make([]*v1.DiscountLine, 0),
	}, nil
}

func (s *sPromotion) LockCouponForOrder(_ context.Context, req *v1.LockCouponForOrderReq) (*v1.LockCouponForOrderRes, error) {
	// 校验请求对象，防止后续字段访问异常。
	if req == nil {
		// 请求为空时直接返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter)
	}
	// 校验订单号必填，保证锁记录可追踪。
	if strings.TrimSpace(req.GetOrderNo()) == "" {
		// 缺失订单号时拒绝处理，避免产生孤儿锁。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter)
	}
	// 通过原子自增生成并发安全的序列号。
	seq := atomic.AddUint64(&lockSeq, 1)
	// 组合时间戳与序列号生成可读锁号。
	lockNo := fmt.Sprintf("LOCK-%d-%d", time.Now().UnixMilli(), seq)
	// 追加订单号后缀，方便联调排查。
	lockNo = fmt.Sprintf("%s-%s", lockNo, strings.TrimSpace(req.GetOrderNo()))
	// 返回锁定成功状态，供下游继续确认流程。
	return &v1.LockCouponForOrderRes{
		LockNo: lockNo,
		Status: v1.CouponLockStatus_COUPON_LOCK_STATUS_LOCKED,
	}, nil
}

func (s *sPromotion) ConfirmCouponUsage(_ context.Context, req *v1.ConfirmCouponUsageReq) (*v1.ConfirmCouponUsageRes, error) {
	// 校验请求对象，避免空请求进入状态变更链路。
	if req == nil {
		// 参数不合法时返回统一错误码。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter)
	}
	// 校验锁号必填，确保确认动作具备目标对象。
	if strings.TrimSpace(req.GetLockNo()) == "" {
		// 锁号为空时拒绝确认，防止错误状态流转。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter)
	}
	// 返回已确认状态，完成最小可运行闭环。
	return &v1.ConfirmCouponUsageRes{
		LockNo: strings.TrimSpace(req.GetLockNo()),
		Status: v1.CouponLockStatus_COUPON_LOCK_STATUS_CONFIRMED,
	}, nil
}

func (s *sPromotion) ReleaseCouponLock(_ context.Context, req *v1.ReleaseCouponLockReq) (*v1.ReleaseCouponLockRes, error) {
	// 校验请求对象，避免无效释放请求。
	if req == nil {
		// 请求体为空时返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter)
	}
	// 校验锁号必填，确保释放动作准确定位。
	if strings.TrimSpace(req.GetLockNo()) == "" {
		// 锁号为空时拒绝释放，防止误操作。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter)
	}
	// 返回已释放状态，支持下游取消和补偿流程。
	return &v1.ReleaseCouponLockRes{
		LockNo: strings.TrimSpace(req.GetLockNo()),
		Status: v1.CouponLockStatus_COUPON_LOCK_STATUS_RELEASED,
	}, nil
}

func (s *sPromotion) CreateCampaign(_ context.Context, req *v1.CreateCampaignReq) (*v1.CreateCampaignRes, error) {
	// 校验请求对象，避免空请求落库。
	if req == nil {
		// 参数不合法时直接返回错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter)
	}
	// 读取并清洗活动编号，去除首尾空白。
	campaignNo := strings.TrimSpace(req.GetCampaignNo())
	// 如果调用方未提供活动编号，则生成兜底编号。
	if campaignNo == "" {
		// 通过毫秒时间戳生成简易活动号。
		campaignNo = fmt.Sprintf("CAMPAIGN-%d", time.Now().UnixMilli())
	}
	// 返回活动编号，供上游后续持久化和查询。
	return &v1.CreateCampaignRes{CampaignNo: campaignNo}, nil
}
