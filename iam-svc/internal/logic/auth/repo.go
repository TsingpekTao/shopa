package auth

import (
	"context"
	"encoding/json"
	"sort"
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

// nextUserID 通过 snowflake 生成分布式唯一 user_id。
// 该 ID 生成不依赖数据库自增，适合多实例并发注册。
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

// findAuthByIdentifier 按登录标识查询用户。
// 路由规则：包含 "@" 走邮箱，否则走手机号。
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

// findRolesByUserID 查询用户当前有效角色集合。
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

// findPermissionsByRoles 查询角色集合映射的权限点列表。
func (s *Service) findPermissionsByRoles(ctx context.Context, roles []entity.IamUserRole) ([]string, error) {
	if len(roles) == 0 {
		return []string{}, nil
	}

	roleCodes := make([]uint, 0, len(roles))
	seenRoleCode := make(map[uint]struct{}, len(roles))
	for _, role := range roles {
		if _, ok := seenRoleCode[role.RoleCode]; ok {
			continue
		}
		seenRoleCode[role.RoleCode] = struct{}{}
		roleCodes = append(roleCodes, role.RoleCode)
	}
	if len(roleCodes) == 0 {
		return []string{}, nil
	}

	var records []struct {
		PermissionKey string `json:"permission_key"`
	}
	err := dao.IamUserRole.DB().
		Model("iam_role_permission rp").
		Ctx(ctx).
		LeftJoin("iam_permission p", "p.permission_key = rp.permission_key").
		Fields("rp.permission_key").
		WhereIn("rp.role_code", roleCodes).
		Where("rp.status", consts.StatusActive).
		Where("p.status", consts.StatusActive).
		Scan(&records)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return []string{}, nil
	}

	permissionSet := make(map[string]struct{}, len(records))
	permissions := make([]string, 0, len(records))
	for _, row := range records {
		key := strings.TrimSpace(row.PermissionKey)
		if key == "" {
			continue
		}
		if _, ok := permissionSet[key]; ok {
			continue
		}
		permissionSet[key] = struct{}{}
		permissions = append(permissions, key)
	}
	sort.Strings(permissions)
	return permissions, nil
}

// findMembershipByUserID 查询会员等级与积分信息。
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

// findRefreshSessionBySID 按 sid 查询 refresh 会话。
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

// roleToProto 将角色实体转换为对外协议对象。
func (s *Service) roleToProto(role entity.IamUserRole) *v1.RoleItem {
	return &v1.RoleItem{
		RoleCode:  v1.RoleCode(role.RoleCode),
		ScopeType: v1.ScopeType(role.ScopeType),
		ScopeId:   role.ScopeId,
	}
}

// membershipToProto 将会员实体转换为对外协议对象。
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

// buildSessionSummary 聚合会话摘要所需数据：角色、会员、最近登录信息。
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

// buildAuthUserSummary 聚合内部查询接口需要的认证摘要。
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

// buildTokenAuthResult 将 tokenPair 包装到 AuthResult.oneof。
func (s *Service) buildTokenAuthResult(pair *v1.TokenPair) *v1.AuthResult {
	return &v1.AuthResult{
		Result: &v1.AuthResult_TokenPair{TokenPair: pair},
	}
}

// buildMFAAuthResult 将 challenge 包装到 AuthResult.oneof。
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

// insertLoginLog 写入登录审计日志。
// 该链路“失败不阻断主流程”，因此忽略写库错误，避免日志系统影响登录可用性。
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
func (s *Service) insertSmsLog(ctx context.Context, scene v1.SmsScene, target, provider, bizID string, success bool, code int, meta riskMeta) {
	if strings.TrimSpace(provider) == "" {
		provider = consts.SmsProviderMock
	}
	_, _ = dao.IamSmsLog.Ctx(ctx).Data(do.IamSmsLog{
		Scene:       uint(scene),
		Target:      target,
		Provider:    provider,
		BizId:       bizID,
		Ip:          meta.ClientIP,
		Ua:          meta.UserAgent,
		Fingerprint: meta.Fingerprint,
		Success:     boolToInt(success),
		ErrorCode:   code,
	}).Insert()
}

// insertOutboxUserRegistered 把“用户注册”事件写入 outbox。
// 关键点：和主事务共用同一个 tx，确保“主业务成功”与“事件可投递”原子一致。
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

// boolToInt 把布尔值转换为数据库常见 0/1 表达。
func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

// hashSmsCode 对短信验证码做带场景与手机号的摘要。
// 这样即便同一验证码被复用，不同场景/手机号也不会产生相同摘要。
func (s *Service) hashSmsCode(scene v1.SmsScene, phone, code string) string {
	plain := strings.Join([]string{strconv.FormatInt(int64(scene), 10), phone, code, s.jwt.Secret}, "|")
	return sha256Hex(plain)
}

// hashRefreshToken 对 refresh token 做不可逆摘要后存储。
func (s *Service) hashRefreshToken(raw string) string {
	return sha256Hex(raw)
}
