package auth

import (
	"context"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/TsingpekTao/shopa/iam-svc/internal/consts"
	"github.com/TsingpekTao/shopa/iam-svc/internal/dao"
	"github.com/TsingpekTao/shopa/iam-svc/internal/errs"
	"github.com/TsingpekTao/shopa/iam-svc/internal/model/do"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/protobuf/types/known/emptypb"
)

// RefreshToken 执行 refresh token 轮换并签发新 token 对。
// 安全关键点：
// 1) 先做多重合法性校验（签名、过期、撤销状态、哈希一致、tokenVersion）。
// 2) 在事务内“撤销旧 session + 插入新 session”，确保会话链路原子一致。
func (s *Service) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (*v1.RefreshTokenRes, error) {
	refreshToken := strings.TrimSpace(req.GetRefreshToken())
	if refreshToken == "" {
		return nil, errs.New(errs.CodeInvalidParam)
	}

	// 第一步：解析 refresh token（含签名与基础时效校验）。
	claims, err := s.parseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 第二步：读取 refresh_session 并校验会话状态。
	session, err := s.findRefreshSessionBySID(ctx, claims.SID)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	if session == nil || session.UserId != claims.UserID {
		return nil, errs.New(errs.CodeRefreshTokenInvalid)
	}
	if session.RevokedAt != nil && !session.RevokedAt.IsZero() {
		return nil, errs.New(errs.CodeRefreshTokenInvalid)
	}
	if session.ExpiresAt == nil || session.ExpiresAt.Timestamp() <= time.Now().UTC().Unix() {
		return nil, errs.New(errs.CodeRefreshTokenInvalid)
	}
	if session.RefreshTokenHash != s.hashRefreshToken(refreshToken) {
		return nil, errs.New(errs.CodeRefreshTokenInvalid)
	}

	// 第三步：读取账号态并做 tokenVersion 对齐校验。
	auth, err := s.findAuthByUserID(ctx, session.UserId)
	if err != nil || auth == nil {
		return nil, errs.New(errs.CodeRefreshTokenInvalid)
	}
	if auth.AccountStatus == consts.AccountStatusDisabled {
		return nil, errs.New(errs.CodeAccountDisabled)
	}
	if uint32(auth.TokenVersion) != claims.TokenVersion {
		return nil, errs.New(errs.CodeRefreshTokenInvalid)
	}

	// 生成新会话 sid 与新 tokenPair。
	meta := extractRiskMeta(ctx, req.GetRisk())
	newSID, err := s.newSID()
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	pair, refreshHash, _, refreshExp, err := s.issueTokenPair(auth.UserId, newSID, uint32(auth.TokenVersion))
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	nowTime := gtime.NewFromTime(time.Now().UTC())
	err = dao.IamUserAuth.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		sessCols := dao.IamRefreshSession.Columns()
		// 事务步骤 1：撤销旧 sid，并记录替代 sid（审计可追溯）。
		if _, err = tx.Model(dao.IamRefreshSession.Table()).
			Where(sessCols.Sid, session.Sid).
			Data(do.IamRefreshSession{RevokedAt: nowTime, ReplacedBySid: newSID}).
			Update(); err != nil {
			return err
		}

		// 事务步骤 2：插入新 session（只存 refresh hash，不存明文 token）。
		_, err = tx.Model(dao.IamRefreshSession.Table()).Data(do.IamRefreshSession{
			Sid:              newSID,
			UserId:           auth.UserId,
			RefreshTokenHash: refreshHash,
			UaHash:           sha256Hex(meta.UserAgent),
			Ip:               meta.ClientIP,
			Geo:              "",
			Fingerprint:      meta.Fingerprint,
			ExpiresAt:        gtime.NewFromTime(refreshExp.UTC()),
		}).Insert()
		return err
	})
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	return &v1.RefreshTokenRes{TokenPair: pair}, nil
}

// Logout 让 refresh session 失效。
// 支持两种模式：
// 1) allDevices=true：注销用户所有未撤销会话。
// 2) allDevices=false：仅注销当前 sid（优先取 refresh token，否则从 access token 解析）。
func (s *Service) Logout(ctx context.Context, req *v1.LogoutReq) (*emptypb.Empty, error) {
	nowTime := gtime.NewFromTime(time.Now().UTC())

	if req.GetAllDevices() {
		accessToken := extractAccessTokenFromMetadata(ctx)
		claims, err := s.parseAccessToken(accessToken)
		if err != nil {
			return nil, err
		}
		cols := dao.IamRefreshSession.Columns()
		_, err = dao.IamRefreshSession.Ctx(ctx).
			Where(cols.UserId, claims.UserID).
			WhereNull(cols.RevokedAt).
			Data(do.IamRefreshSession{RevokedAt: nowTime}).
			Update()
		if err != nil {
			return nil, errs.Wrap(errs.CodeInternalError, err)
		}
		return &emptypb.Empty{}, nil
	}

	var sid string
	if token := strings.TrimSpace(req.GetRefreshToken()); token != "" {
		claims, err := s.parseRefreshToken(token)
		if err != nil {
			return nil, err
		}
		sid = claims.SID
	} else {
		accessToken := extractAccessTokenFromMetadata(ctx)
		claims, err := s.parseAccessToken(accessToken)
		if err != nil {
			return nil, err
		}
		sid = claims.SID
	}

	if strings.TrimSpace(sid) == "" {
		return nil, errs.New(errs.CodeInvalidParam)
	}

	cols := dao.IamRefreshSession.Columns()
	_, err := dao.IamRefreshSession.Ctx(ctx).
		Where(cols.Sid, sid).
		WhereNull(cols.RevokedAt).
		Data(do.IamRefreshSession{RevokedAt: nowTime}).
		Update()
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	return &emptypb.Empty{}, nil
}

// GetMySession 返回当前 access token 对应用户的会话摘要。
func (s *Service) GetMySession(ctx context.Context, req *emptypb.Empty) (*v1.GetMySessionRes, error) {
	accessToken := extractAccessTokenFromMetadata(ctx)
	claims, err := s.parseAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	auth, err := s.findAuthByUserID(ctx, claims.UserID)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	if auth == nil {
		return nil, errs.New(errs.CodeInvalidCredential)
	}
	// tokenVersion 对齐用于识别“全局失效”场景。
	if uint32(auth.TokenVersion) != claims.TokenVersion {
		return nil, errs.New(errs.CodeAccessTokenInvalid)
	}

	session, err := s.buildSessionSummary(ctx, auth)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	return &v1.GetMySessionRes{Session: session}, nil
}

// ChangePassword 登录态修改密码。
// 事务内同时完成：
// 1) 更新密码哈希。
// 2) tokenVersion +1（让历史 token 全部失效）。
// 3) 撤销全部 refresh session（强制重新登录）。
func (s *Service) ChangePassword(ctx context.Context, req *v1.ChangePasswordReq) (*v1.ChangePasswordRes, error) {
	var (
		oldPwd = req.GetOldPassword()
		newPwd = req.GetNewPassword()
	)
	if strings.TrimSpace(oldPwd) == "" || strings.TrimSpace(newPwd) == "" {
		return nil, errs.New(errs.CodeInvalidParam)
	}
	if err := validateNewPassword(newPwd); err != nil {
		return nil, errs.New(errs.CodeInvalidParam, err.Error())
	}

	accessToken := extractAccessTokenFromMetadata(ctx)
	claims, err := s.parseAccessToken(accessToken)
	if err != nil {
		return nil, err
	}
	auth, err := s.findAuthByUserID(ctx, claims.UserID)
	if err != nil || auth == nil {
		return nil, errs.New(errs.CodeInvalidCredential)
	}
	if !verifyPassword(oldPwd, auth.PasswordHash, auth.PasswordSalt, auth.PasswordAlgo) {
		return nil, errs.New(errs.CodeInvalidCredential)
	}

	hash, salt, algo, ver, err := makePasswordHash(newPwd)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	err = dao.IamUserAuth.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		authCols := dao.IamUserAuth.Columns()
		if _, err = tx.Model(dao.IamUserAuth.Table()).
			Where(authCols.UserId, auth.UserId).
			Data(do.IamUserAuth{
				PasswordHash:     hash,
				PasswordSalt:     salt,
				PasswordAlgo:     algo,
				PasswordVer:      ver,
				TokenVersion:     gdb.Raw(authCols.TokenVersion + " + 1"),
				AccountStatus:    consts.AccountStatusActive,
				FailedLoginCount: 0,
			}).
			Update(); err != nil {
			return err
		}

		sessCols := dao.IamRefreshSession.Columns()
		_, err = tx.Model(dao.IamRefreshSession.Table()).
			Where(sessCols.UserId, auth.UserId).
			WhereNull(sessCols.RevokedAt).
			Data(do.IamRefreshSession{RevokedAt: gtime.Now()}).
			Update()
		return err
	})
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	return &v1.ChangePasswordRes{Updated: true}, nil
}

// ResetPasswordBySms 通过短信码重置密码。
// 与改密策略一致：成功后递增 tokenVersion 并撤销全部 refresh session。
func (s *Service) ResetPasswordBySms(ctx context.Context, req *v1.ResetPasswordBySmsReq) (*v1.ResetPasswordBySmsRes, error) {
	var (
		phone   = normalizePhone(req.GetPhone())
		smsCode = strings.TrimSpace(req.GetSmsCode())
		newPwd  = req.GetNewPassword()
	)
	if phone == "" || smsCode == "" || strings.TrimSpace(newPwd) == "" {
		return nil, errs.New(errs.CodeInvalidParam)
	}
	if err := validateNewPassword(newPwd); err != nil {
		return nil, errs.New(errs.CodeInvalidParam, err.Error())
	}

	if err := s.verifySmsCode(ctx, v1.SmsScene_SMS_SCENE_RESET_PASSWORD, phone, smsCode); err != nil {
		return nil, err
	}

	auth, err := s.findAuthByPhone(ctx, phone)
	if err != nil || auth == nil {
		return nil, errs.New(errs.CodeInvalidCredential)
	}

	hash, salt, algo, ver, err := makePasswordHash(newPwd)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	err = dao.IamUserAuth.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		authCols := dao.IamUserAuth.Columns()
		if _, err = tx.Model(dao.IamUserAuth.Table()).
			Where(authCols.UserId, auth.UserId).
			Data(do.IamUserAuth{
				PasswordHash:     hash,
				PasswordSalt:     salt,
				PasswordAlgo:     algo,
				PasswordVer:      ver,
				TokenVersion:     gdb.Raw(authCols.TokenVersion + " + 1"),
				AccountStatus:    consts.AccountStatusActive,
				FailedLoginCount: 0,
			}).
			Update(); err != nil {
			return err
		}

		sessCols := dao.IamRefreshSession.Columns()
		_, err = tx.Model(dao.IamRefreshSession.Table()).
			Where(sessCols.UserId, auth.UserId).
			WhereNull(sessCols.RevokedAt).
			Data(do.IamRefreshSession{RevokedAt: gtime.Now()}).
			Update()
		return err
	})
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}

	return &v1.ResetPasswordBySmsRes{Updated: true}, nil
}

// LoginByOAuth 第三方登录暂未实现。
func (s *Service) LoginByOAuth(ctx context.Context, req *v1.LoginByOAuthReq) (*v1.LoginByOAuthRes, error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "oauth login is TODO")
}

// BindOAuth 第三方账号绑定暂未实现。
func (s *Service) BindOAuth(ctx context.Context, req *v1.BindOAuthReq) (*v1.BindOAuthRes, error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "oauth bind is TODO")
}

// UnbindOAuth 第三方账号解绑暂未实现。
func (s *Service) UnbindOAuth(ctx context.Context, req *v1.UnbindOAuthReq) (*v1.UnbindOAuthRes, error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "oauth unbind is TODO")
}
