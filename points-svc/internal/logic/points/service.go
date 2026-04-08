package points

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/TsingpekTao/shopa/points-svc/api/v1"
	"github.com/TsingpekTao/shopa/points-svc/internal/dao"
	"github.com/TsingpekTao/shopa/points-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/points-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/points-svc/internal/service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	// 账户状态码沿用库表与协议中的固定枚举值。
	accountStatusActive   = "ACTIVE"
	accountStatusFrozen   = "FROZEN"
	accountStatusDisabled = "DISABLED"

	// 积分预留状态用于区分锁定、确认和取消三种订单抵扣阶段。
	reservationStatusLocked    = "LOCKED"
	reservationStatusConfirmed = "CONFIRMED"
	reservationStatusCanceled  = "CANCELED"
	reservationStatusExpired   = "EXPIRED"

	// 过期桶状态描述桶内积分当前是否仍可继续参与扣减。
	bucketStatusActive  = "ACTIVE"
	bucketStatusExpired = "EXPIRED"
	bucketStatusClosed  = "CLOSED"

	// 积分桶来源类型用于追踪积分是由发放、调整还是退款宽限生成。
	bucketSourceGrant       = "GRANT"
	bucketSourceAdjust      = "ADJUST"
	bucketSourceRefundGrace = "REFUND_GRACE"
)

// sPoints 聚合积分账户、规则、预留、流水和退款相关的核心逻辑实现。
type sPoints struct{}

// ruleSnapshot 是写入预留单和授予明细的规则快照，确保后续核销使用同一份规则上下文。
type ruleSnapshot struct {
	RuleCode            string `json:"rule_code"`
	RuleName            string `json:"rule_name"`
	MinOrderAmountCent  int64  `json:"min_order_amount_cent"`
	MaxDeductionRateBps uint32 `json:"max_deduction_rate_bps"`
	DeductPointsPerCent uint64 `json:"deduct_points_per_cent"`
	GrantPointsPerCent  uint64 `json:"grant_points_per_cent"`
	RefundGraceDays     uint32 `json:"refund_grace_days"`
}

// ledgerInput 统一描述积分流水落库时需要记录的业务字段。
type ledgerInput struct {
	UserID          uint64
	EntryTypeCode   string
	BizType         string
	BizNo           string
	ReservationNo   string
	RelatedBucketNo string
	PointsDelta     int64
	AvailableAfter  int64
	FrozenAfter     int64
	DebtAfter       uint64
	CashAmountCent  int64
	Remark          string
	Extra           map[string]any
}

// New 创建积分逻辑实现，并交给 service 层完成注册暴露。
func New() *sPoints {
	return &sPoints{}
}

// init 在包加载时注册积分服务实现，保持与 GoFrame 默认启动方式一致。
func init() {
	service.RegisterPoints(New())
}

// InitPointsAccountIfAbsent 为指定用户补齐积分账户；若账户已存在则直接返回当前状态。
func (s *sPoints) InitPointsAccountIfAbsent(ctx context.Context, req *pb.InitPointsAccountIfAbsentReq) (*pb.InitPointsAccountIfAbsentRes, error) {
	if req.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	created, account, err := s.initAccountIfAbsent(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &pb.InitPointsAccountIfAbsentRes{Created: created, Balance: toCompatBalance(account), StatusCode: account.StatusCode}, nil
}

// GetPointsByUserId 查询用户当前可兼容旧协议的积分余额与账户状态。
func (s *sPoints) GetPointsByUserId(ctx context.Context, req *pb.GetPointsByUserIdReq) (*pb.GetPointsByUserIdRes, error) {
	if req.GetUserId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	account, err := s.ensureAccount(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &pb.GetPointsByUserIdRes{UserId: account.UserId, Balance: toCompatBalance(account), Status: toCompatStatus(account.StatusCode)}, nil
}

// InitAccountFromRegisterEvent 供注册事件消费方调用，确保新用户具备积分账户。
func (s *sPoints) InitAccountFromRegisterEvent(ctx context.Context, userID uint64) error {
	if userID == 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	_, _, err := s.initAccountIfAbsent(ctx, userID)
	return err
}

// initAccountIfAbsent 使用幂等插入初始化账户，并在首次创建时补一条初始化流水。
func (s *sPoints) initAccountIfAbsent(ctx context.Context, userID uint64) (bool, *entity.PointsAccount, error) {
	created := false
	r, err := dao.PointsAccount.Ctx(ctx).Data(do.PointsAccount{
		UserId:              userID,
		AvailableBalance:    0,
		FrozenBalance:       0,
		StatusCode:          accountStatusActive,
		TotalEarnedPoints:   0,
		TotalUsedPoints:     0,
		TotalExpiredPoints:  0,
		TotalAdjustedPoints: 0,
	}).InsertIgnore()
	if err != nil {
		return false, nil, gerror.Wrap(err, "insert points_account failed")
	}
	if r != nil {
		if rows, rowsErr := r.RowsAffected(); rowsErr == nil && rows > 0 {
			created = true
		}
	}
	account, err := s.getAccount(ctx, userID)
	if err != nil {
		return false, nil, err
	}
	if account == nil {
		return false, nil, gerror.NewCode(gcode.CodeInternalError, "points account missing after init")
	}
	if created {
		_, err = dao.PointsLedger.Ctx(ctx).Data(do.PointsLedger{
			LedgerNo:        generateBizNo("PL"),
			UserId:          userID,
			EntryTypeCode:   "INIT",
			BizType:         "REGISTER_INIT",
			BizNo:           fmt.Sprintf("register:%d", userID),
			PointsDelta:     0,
			AvailableAfter:  account.AvailableBalance,
			FrozenAfter:     account.FrozenBalance,
			DebtAfter:       debtFromAvailable(account.AvailableBalance),
			CashAmountCent:  0,
			Remark:          "register init account",
			ReservationNo:   "",
			RelatedBucketNo: "",
		}).InsertIgnore()
		if err != nil {
			return false, nil, gerror.Wrap(err, "insert init ledger failed")
		}
	}
	return created, account, nil
}

// ensureAccount 保证调用方拿到可用账户，不存在时自动初始化。
func (s *sPoints) ensureAccount(ctx context.Context, userID uint64) (*entity.PointsAccount, error) {
	_, account, err := s.initAccountIfAbsent(ctx, userID)
	return account, err
}

// getAccount 按用户编号读取账户，不带事务并允许账户不存在。
func (s *sPoints) getAccount(ctx context.Context, userID uint64) (*entity.PointsAccount, error) {
	var row entity.PointsAccount
	if err := dao.PointsAccount.Ctx(ctx).Where(dao.PointsAccount.Columns().UserId, userID).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query points_account failed")
	}
	if row.UserId == 0 {
		return nil, nil
	}
	return &row, nil
}

// getAccountTx 在事务内读取账户，用于需要原子更新余额的流程。
func (s *sPoints) getAccountTx(ctx context.Context, tx gdb.TX, userID uint64) (*entity.PointsAccount, error) {
	var row entity.PointsAccount
	if err := tx.Model(dao.PointsAccount.Table()).Where(dao.PointsAccount.Columns().UserId, userID).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query points_account failed")
	}
	if row.UserId == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "points account not found")
	}
	return &row, nil
}

// loadRule 读取当前生效的积分规则，并返回序列化后的快照和摘要供外部校验。
func (s *sPoints) loadRule(ctx context.Context) (*entity.PointsRuleConfig, *ruleSnapshot, string, string, error) {
	var row entity.PointsRuleConfig
	now := gtime.Now()
	model := dao.PointsRuleConfig.Ctx(ctx).
		Where(dao.PointsRuleConfig.Columns().StatusCode, accountStatusActive).
		WhereLTE(dao.PointsRuleConfig.Columns().EffectiveAt, now).
		Wheref("(%s IS NULL OR %s > ?)", dao.PointsRuleConfig.Columns().ExpireAt, dao.PointsRuleConfig.Columns().ExpireAt, now).
		OrderDesc(dao.PointsRuleConfig.Columns().EffectiveAt).
		OrderDesc(dao.PointsRuleConfig.Columns().UpdatedAt)
	if err := model.Scan(&row); err != nil {
		return nil, nil, "", "", gerror.Wrap(err, "query points_rule_config failed")
	}
	if row.RuleCode == "" {
		row = entity.PointsRuleConfig{
			RuleCode:            "DEFAULT_RULE",
			RuleName:            "Default points rule",
			StatusCode:          accountStatusActive,
			MinOrderAmountCent:  100,
			MaxDeductionRateBps: 500,
			DeductPointsPerCent: 1,
			GrantPointsPerCent:  1,
			RefundGraceDays:     7,
		}
	}
	snapshot := &ruleSnapshot{RuleCode: row.RuleCode, RuleName: row.RuleName, MinOrderAmountCent: row.MinOrderAmountCent, MaxDeductionRateBps: uint32(row.MaxDeductionRateBps), DeductPointsPerCent: row.DeductPointsPerCent, GrantPointsPerCent: row.GrantPointsPerCent, RefundGraceDays: uint32(row.RefundGraceDays)}
	payload, err := json.Marshal(snapshot)
	if err != nil {
		return nil, nil, "", "", gerror.Wrap(err, "marshal rule snapshot failed")
	}
	return &row, snapshot, string(payload), hashJSON(string(payload)), nil
}

// insertLedgerTx 在事务中写入积分流水，统一收口各类余额变更记录。
func (s *sPoints) insertLedgerTx(ctx context.Context, tx gdb.TX, in *ledgerInput) (string, error) {
	if in == nil {
		return "", nil
	}
	extraJSON := "null"
	if len(in.Extra) > 0 {
		if encoded, err := json.Marshal(in.Extra); err == nil {
			extraJSON = string(encoded)
		}
	}
	ledgerNo := generateBizNo("PL")
	_, err := tx.Model(dao.PointsLedger.Table()).Data(do.PointsLedger{
		LedgerNo:        ledgerNo,
		UserId:          in.UserID,
		EntryTypeCode:   in.EntryTypeCode,
		BizType:         in.BizType,
		BizNo:           in.BizNo,
		ReservationNo:   in.ReservationNo,
		RelatedBucketNo: in.RelatedBucketNo,
		PointsDelta:     in.PointsDelta,
		AvailableAfter:  in.AvailableAfter,
		FrozenAfter:     in.FrozenAfter,
		DebtAfter:       in.DebtAfter,
		CashAmountCent:  in.CashAmountCent,
		Remark:          truncate(in.Remark, 255),
		ExtraJson:       extraJSON,
	}).Insert()
	if err != nil {
		return "", gerror.Wrap(err, "insert points_ledger failed")
	}
	_ = ctx
	return ledgerNo, nil
}

// toCompatBalance 将账户可用余额收敛为旧协议需要的非负整数表示。
func toCompatBalance(account *entity.PointsAccount) uint64 {
	if account == nil || account.AvailableBalance <= 0 {
		return 0
	}
	return uint64(account.AvailableBalance)
}

// toCompatStatus 把内部账户状态映射成兼容旧协议的数字枚举。
func toCompatStatus(status string) uint32 {
	if status == accountStatusActive {
		return 1
	}
	return 2
}

// debtFromAvailable 从可用余额推导出欠账积分，便于对外展示负债信息。
func debtFromAvailable(available int64) uint64 {
	if available >= 0 {
		return 0
	}
	return uint64(-available)
}

// hashJSON 生成规则快照摘要，用于前后端或跨服务做幂等校验。
func hashJSON(value string) string {
	sum := sha1.Sum([]byte(value))
	return hex.EncodeToString(sum[:])
}

// formatTime 统一把 GoFrame 时间对象转为协议侧使用的字符串格式。
func formatTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// truncate 控制备注类文本长度，避免超出数据库字段限制。
func truncate(value string, limit int) string {
	if limit <= 0 || len(value) <= limit {
		return value
	}
	return value[:limit]
}

// generateBizNo 生成积分域内的业务单号，沿用时间戳加随机尾号的简单策略。
func generateBizNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.UnixNano()%1000000)
}
