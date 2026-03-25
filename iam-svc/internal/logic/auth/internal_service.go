package auth

import (
	"context"
	"strings"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/TsingpekTao/shopa/iam-svc/internal/consts"
	"github.com/TsingpekTao/shopa/iam-svc/internal/dao"
	"github.com/TsingpekTao/shopa/iam-svc/internal/errs"
	"github.com/TsingpekTao/shopa/iam-svc/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetAuthUserById 按 user_id 查询认证摘要，供网关和内部服务调用。
func (s *Service) GetAuthUserById(ctx context.Context, req *v1.GetAuthUserByIdReq) (*v1.GetAuthUserByIdRes, error) {
	if req.GetUserId() == 0 {
		return nil, errs.New(errs.CodeInvalidParam)
	}

	auth, err := s.findAuthByUserID(ctx, req.GetUserId())
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	if auth == nil {
		return &v1.GetAuthUserByIdRes{}, nil
	}

	user, err := s.buildAuthUserSummary(ctx, auth)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	return &v1.GetAuthUserByIdRes{User: user}, nil
}

// BatchGetAuthUsers 批量返回认证摘要，按请求顺序构造输出。
func (s *Service) BatchGetAuthUsers(ctx context.Context, req *v1.BatchGetAuthUsersReq) (*v1.BatchGetAuthUsersRes, error) {
	userIDs := req.GetUserIds()
	if len(userIDs) == 0 {
		return &v1.BatchGetAuthUsersRes{Users: []*v1.AuthUserSummary{}}, nil
	}

	var (
		users      []entity.IamUserAuth
		roles      []entity.IamUserRole
		membership []entity.IamMembership
	)

	authCols := dao.IamUserAuth.Columns()
	if err := dao.IamUserAuth.Ctx(ctx).WhereIn(authCols.UserId, userIDs).Scan(&users); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	if len(users) == 0 {
		return &v1.BatchGetAuthUsersRes{Users: []*v1.AuthUserSummary{}}, nil
	}

	roleCols := dao.IamUserRole.Columns()
	if err := dao.IamUserRole.Ctx(ctx).
		WhereIn(roleCols.UserId, userIDs).
		Where(roleCols.Status, consts.StatusActive).
		Scan(&roles); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	memberCols := dao.IamMembership.Columns()
	if err := dao.IamMembership.Ctx(ctx).WhereIn(memberCols.UserId, userIDs).Scan(&membership); err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	roleMap := make(map[uint64][]*v1.RoleItem, len(users))
	for _, role := range roles {
		roleMap[role.UserId] = append(roleMap[role.UserId], s.roleToProto(role))
	}

	memberMap := make(map[uint64]*v1.MembershipSummary, len(membership))
	for _, m := range membership {
		mm := m
		memberMap[m.UserId] = s.membershipToProto(&mm)
	}

	userMap := make(map[uint64]entity.IamUserAuth, len(users))
	for _, user := range users {
		userMap[user.UserId] = user
	}

	out := make([]*v1.AuthUserSummary, 0, len(userIDs))
	for _, uid := range userIDs {
		u, ok := userMap[uid]
		if !ok {
			continue
		}
		out = append(out, &v1.AuthUserSummary{
			UserId:        u.UserId,
			Phone:         u.Phone,
			Email:         u.Email,
			AccountStatus: v1.AccountStatus(u.AccountStatus),
			Roles:         roleMap[u.UserId],
			Membership:    memberMap[u.UserId],
			TokenVersion:  uint32(u.TokenVersion),
			UpdatedAt:     toProtoTimestamp(u.UpdatedAt),
		})
	}

	return &v1.BatchGetAuthUsersRes{Users: out}, nil
}

// VerifyAccessToken 供内部服务校验 access token，并返回角色/会员等鉴权上下文。
func (s *Service) VerifyAccessToken(ctx context.Context, req *v1.VerifyAccessTokenReq) (*v1.VerifyAccessTokenRes, error) {
	token := strings.TrimSpace(req.GetAccessToken())
	if token == "" {
		return &v1.VerifyAccessTokenRes{Valid: false}, nil
	}

	claims, err := s.parseAccessToken(token)
	if err != nil {
		g.Log().Debugf(ctx, "[iam-svc] verify access token failed: %+v", err)
		return &v1.VerifyAccessTokenRes{Valid: false}, nil
	}

	auth, err := s.findAuthByUserID(ctx, claims.UserID)
	if err != nil || auth == nil {
		return &v1.VerifyAccessTokenRes{Valid: false}, nil
	}
	if uint32(auth.TokenVersion) != claims.TokenVersion {
		return &v1.VerifyAccessTokenRes{Valid: false}, nil
	}
	if auth.AccountStatus == consts.AccountStatusDisabled {
		return &v1.VerifyAccessTokenRes{Valid: false}, nil
	}

	roles, err := s.findRolesByUserID(ctx, auth.UserId)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	membership, err := s.findMembershipByUserID(ctx, auth.UserId)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	roleItems := make([]*v1.RoleItem, 0, len(roles))
	for _, role := range roles {
		roleItems = append(roleItems, s.roleToProto(role))
	}

	var expiresAt *timestamppb.Timestamp
	if claims.ExpiresAt != nil {
		expiresAt = timestamppb.New(claims.ExpiresAt.Time.UTC())
	}

	return &v1.VerifyAccessTokenRes{
		Valid:         true,
		UserId:        auth.UserId,
		Sid:           claims.SID,
		TokenVersion:  uint32(auth.TokenVersion),
		AccountStatus: v1.AccountStatus(auth.AccountStatus),
		Roles:         roleItems,
		Membership:    s.membershipToProto(membership),
		ExpiresAt:     expiresAt,
	}, nil
}
