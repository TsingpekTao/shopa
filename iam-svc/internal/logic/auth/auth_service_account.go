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

// RefreshToken 接收 refresh token，并换发一组新的 token 对。
func (s *Service) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (*v1.RefreshTokenRes, error) {
	refreshToken := strings.TrimSpace(req.GetRefreshToken())
	if refreshToken == "" {
		return nil, errs.New(errs.CodeInvalidParam)
	}

	// 第一步：解析 refresh token，完成签名与基础时效校验。
	claims, err := s.parseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 第二步：读取 refresh_session，并校验会话状态与持久化哈希。
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

	// 第三步：校验账号状态与 tokenVersion，确保旧批次 token 无法刷新。
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
		// 事务步骤 1：撤销旧 sid，并记录被新的 sid 替代。
		if _, err = tx.Model(dao.IamRefreshSession.Table()).
			Where(sessCols.Sid, session.Sid).
			Data(do.IamRefreshSession{RevokedAt: nowTime, ReplacedBySid: newSID}).
			Update(); err != nil {
			return err
		}

		// 事务步骤 2：插入新 session，仅保存 refresh token 哈希。
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

// Logout 注销当前设备或全部设备的 refresh session。
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

// GetMySession 通过 access token 查询当前会话摘要。
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
	// tokenVersion 不匹配说明 access token 已被刷新或失效。
	if uint32(auth.TokenVersion) != claims.TokenVersion {
		return nil, errs.New(errs.CodeAccessTokenInvalid)
	}

	session, err := s.buildSessionSummary(ctx, auth)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, err)
	}
	return &v1.GetMySessionRes{Session: session}, nil
}

// ChangePassword 校验旧密码并更新哈希，同时让历史会话失效。
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

// ResetPasswordBySms 校验短信验证码后重置密码，并清空旧会话。
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

// LoginByOAuth 是第三方登录占位接口。
func (s *Service) LoginByOAuth(ctx context.Context, req *v1.LoginByOAuthReq) (*v1.LoginByOAuthRes, error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "oauth login is TODO")
}

// BindOAuth 是第三方账号绑定占位接口。
func (s *Service) BindOAuth(ctx context.Context, req *v1.BindOAuthReq) (*v1.BindOAuthRes, error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "oauth bind is TODO")
}

// UnbindOAuth 是第三方账号解绑占位接口。
func (s *Service) UnbindOAuth(ctx context.Context, req *v1.UnbindOAuthReq) (*v1.UnbindOAuthRes, error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented, "oauth unbind is TODO")
}
