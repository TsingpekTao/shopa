package auth

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/TsingpekTao/shopa/iam-svc/internal/consts"
	"github.com/TsingpekTao/shopa/iam-svc/internal/dao"
	"github.com/TsingpekTao/shopa/iam-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/iam-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// nextUserID 基于 snowflake 生成分布式唯一 user_id。
func (s *Service) nextUserID() uint64 {
	return uint64(s.node.Generate().Int64())
}

// findAuthByPhone 按手机号查询认证主表。
func (s *Service) findAuthByPhone(ctx context.Context, phone string) (*entity.IamUserAuth, error) {
	var (
		cols      = dao.IamUserAuth.Columns()
		record, e = dao.IamUserAuth.Ctx(ctx).Where(cols.Phone, phone).One()
	)
	if e != nil {
		return nil, e
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var auth entity.IamUserAuth
	if e = record.Struct(&auth); e != nil {
		return nil, e
	}
	return &auth, nil
}

// findAuthByIdentifier 按登录标识查询：包含 '@' 走邮箱，否则走手机号。
func (s *Service) findAuthByIdentifier(ctx context.Context, identifier string) (*entity.IamUserAuth, error) {
	var (
		cols   = dao.IamUserAuth.Columns()
		record gdb.Record
		err    error
	)

	if strings.Contains(identifier, "@") {
		record, err = dao.IamUserAuth.Ctx(ctx).Where(cols.Email, identifier).One()
	} else {
		record, err = dao.IamUserAuth.Ctx(ctx).Where(cols.Phone, identifier).One()
	}
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var auth entity.IamUserAuth
	if err = record.Struct(&auth); err != nil {
		return nil, err
	}
	return &auth, nil
}

// findAuthByUserID 按 user_id 查询认证主表。
func (s *Service) findAuthByUserID(ctx context.Context, userID uint64) (*entity.IamUserAuth, error) {
	var (
		cols      = dao.IamUserAuth.Columns()
		record, e = dao.IamUserAuth.Ctx(ctx).Where(cols.UserId, userID).One()
	)
	if e != nil {
		return nil, e
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var auth entity.IamUserAuth
	if e = record.Struct(&auth); e != nil {
		return nil, e
	}
	return &auth, nil
}

// findRolesByUserID 查询用户当前有效角色列表。
func (s *Service) findRolesByUserID(ctx context.Context, userID uint64) ([]entity.IamUserRole, error) {
	var (
		cols  = dao.IamUserRole.Columns()
		roles []entity.IamUserRole
		err   = dao.IamUserRole.Ctx(ctx).
			Where(cols.UserId, userID).
			Where(cols.Status, consts.StatusActive).
			OrderAsc(cols.Id).
			Scan(&roles)
	)
	return roles, err
}

// findMembershipByUserID 查询会员等级与积分摘要。
func (s *Service) findMembershipByUserID(ctx context.Context, userID uint64) (*entity.IamMembership, error) {
	var (
		cols      = dao.IamMembership.Columns()
		record, e = dao.IamMembership.Ctx(ctx).Where(cols.UserId, userID).One()
	)
	if e != nil {
		return nil, e
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var membership entity.IamMembership
	if e = record.Struct(&membership); e != nil {
		return nil, e
	}
	return &membership, nil
}

// findRefreshSessionBySID 按会话 sid 查询 refresh session。
func (s *Service) findRefreshSessionBySID(ctx context.Context, sid string) (*entity.IamRefreshSession, error) {
	var (
		cols      = dao.IamRefreshSession.Columns()
		record, e = dao.IamRefreshSession.Ctx(ctx).Where(cols.Sid, sid).One()
	)
	if e != nil {
		return nil, e
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var session entity.IamRefreshSession
	if e = record.Struct(&session); e != nil {
		return nil, e
	}
	return &session, nil
}

// roleToProto 将角色实体映射为对外 proto。
func (s *Service) roleToProto(role entity.IamUserRole) *v1.RoleItem {
	return &v1.RoleItem{
		RoleCode:  v1.RoleCode(role.RoleCode),
		ScopeType: v1.ScopeType(role.ScopeType),
		ScopeId:   role.ScopeId,
	}
}

// membershipToProto 将会员实体映射为对外 proto。
func (s *Service) membershipToProto(m *entity.IamMembership) *v1.MembershipSummary {
	if m == nil {
		return &v1.MembershipSummary{}
	}
	return &v1.MembershipSummary{
		LevelCode: m.LevelCode,
		Points:    m.Points,
		ExpireAt:  toProtoTimestamp(m.ExpireAt),
	}
}

// buildSessionSummary 聚合会话返回所需的角色、会员与最近登录信息。
func (s *Service) buildSessionSummary(ctx context.Context, auth *entity.IamUserAuth) (*v1.SessionSummary, error) {
	roles, err := s.findRolesByUserID(ctx, auth.UserId)
	if err != nil {
		return nil, err
	}
	membership, err := s.findMembershipByUserID(ctx, auth.UserId)
	if err != nil {
		return nil, err
	}

	roleItems := make([]*v1.RoleItem, 0, len(roles))
	for _, role := range roles {
		roleItems = append(roleItems, s.roleToProto(role))
	}

	return &v1.SessionSummary{
		UserId:        auth.UserId,
		AccountStatus: v1.AccountStatus(auth.AccountStatus),
		Roles:         roleItems,
		Membership:    s.membershipToProto(membership),
		LastLoginAt:   toProtoTimestamp(auth.LastLoginAt),
		LastLoginIp:   auth.LastLoginIp,
	}, nil
}

// buildAuthUserSummary 聚合内部查询接口所需的认证摘要。
func (s *Service) buildAuthUserSummary(ctx context.Context, auth *entity.IamUserAuth) (*v1.AuthUserSummary, error) {
	roles, err := s.findRolesByUserID(ctx, auth.UserId)
	if err != nil {
		return nil, err
	}
	membership, err := s.findMembershipByUserID(ctx, auth.UserId)
	if err != nil {
		return nil, err
	}

	roleItems := make([]*v1.RoleItem, 0, len(roles))
	for _, role := range roles {
		roleItems = append(roleItems, s.roleToProto(role))
	}

	return &v1.AuthUserSummary{
		UserId:        auth.UserId,
		Phone:         auth.Phone,
		Email:         auth.Email,
		AccountStatus: v1.AccountStatus(auth.AccountStatus),
		Roles:         roleItems,
		Membership:    s.membershipToProto(membership),
		TokenVersion:  uint32(auth.TokenVersion),
		UpdatedAt:     toProtoTimestamp(auth.UpdatedAt),
	}, nil
}

// buildTokenAuthResult 把 token pair 包装到 AuthResult.oneof 中。
func (s *Service) buildTokenAuthResult(pair *v1.TokenPair) *v1.AuthResult {
	return &v1.AuthResult{
		Result: &v1.AuthResult_TokenPair{TokenPair: pair},
	}
}

// buildMFAAuthResult 把 MFA challenge 包装到 AuthResult.oneof 中。
func (s *Service) buildMFAAuthResult(challengeID string, expireAt time.Time) *v1.AuthResult {
	return &v1.AuthResult{
		Result: &v1.AuthResult_MfaChallenge{
			MfaChallenge: &v1.MfaChallenge{
				ChallengeId: challengeID,
				ExpireAt:    timestamppb.New(expireAt.UTC()),
			},
		},
	}
}

// insertLoginLog 写入登录审计日志，不影响主流程返回。
func (s *Service) insertLoginLog(ctx context.Context, userID uint64, identifier string, channel v1.LoginChannel, success bool, failReason string, meta riskMeta) {
	_, _ = dao.IamLoginLog.Ctx(ctx).Data(do.IamLoginLog{
		UserId:      userID,
		Identifier:  identifier,
		Channel:     uint(channel),
		Success:     boolToInt(success),
		FailReason:  failReason,
		Ip:          meta.ClientIP,
		Geo:         "",
		Ua:          meta.UserAgent,
		Fingerprint: meta.Fingerprint,
	}).Insert()
}

// insertSmsLog 写入短信发送审计日志。
func (s *Service) insertSmsLog(ctx context.Context, scene v1.SmsScene, target string, success bool, code int, meta riskMeta) {
	_, _ = dao.IamSmsLog.Ctx(ctx).Data(do.IamSmsLog{
		Scene:       uint(scene),
		Target:      target,
		Provider:    consts.SmsProviderMock,
		BizId:       "",
		Ip:          meta.ClientIP,
		Ua:          meta.UserAgent,
		Fingerprint: meta.Fingerprint,
		Success:     boolToInt(success),
		ErrorCode:   code,
	}).Insert()
}

// insertOutboxUserRegistered 将注册事件写入 outbox，交给后台 worker 异步投递。
func (s *Service) insertOutboxUserRegistered(ctx context.Context, tx gdb.TX, userID uint64, initName string) error {
	eventID := newEventID()
	occurredAt := time.Now().UTC()

	payload, err := json.Marshal(registerEventPayload{
		EventID:         eventID,
		EventVersion:    "v1",
		UserID:          userID,
		InitDisplayName: initName,
		RegisterChannel: "password_sms",
		OccurredAt:      occurredAt.Format(time.RFC3339Nano),
	})
	if err != nil {
		return err
	}

	_, err = tx.Model(dao.IamOutboxEvent.Table()).Data(do.IamOutboxEvent{
		EventId:     eventID,
		EventType:   "UserRegisteredV1",
		PayloadJson: string(payload),
		Status:      consts.OutboxStatusNew,
		AvailableAt: gtime.NewFromTime(occurredAt),
		FailCount:   0,
		LastError:   "",
	}).Insert()
	return err
}

// boolToInt 把布尔值转成数据库常用 0/1 表达。
func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

// hashSmsCode 对短信验证码做带场景和手机号的摘要，避免明文入库。
func (s *Service) hashSmsCode(scene v1.SmsScene, phone, code string) string {
	plain := strings.Join([]string{strconv.FormatInt(int64(scene), 10), phone, code, s.jwt.Secret}, "|")
	return sha256Hex(plain)
}

// hashRefreshToken 对 refresh token 做不可逆摘要存储。
func (s *Service) hashRefreshToken(raw string) string {
	return sha256Hex(raw)
}
