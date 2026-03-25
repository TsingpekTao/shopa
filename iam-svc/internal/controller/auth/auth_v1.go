package auth

import (
	"context"
	"strings"

	authv1 "github.com/TsingpekTao/shopa/iam-svc/api/auth/v1"
	iamv1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/TsingpekTao/shopa/iam-svc/internal/errs"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SendSmsCode 把 HTTP 请求转换为内部 gRPC 请求并转发到核心鉴权服务。
func (c *ControllerV1) SendSmsCode(ctx context.Context, req *authv1.SendSmsCodeReq) (*authv1.SendSmsCodeRes, error) {
	resp, err := c.auth.SendSmsCode(withRequestMetadata(ctx), &iamv1.SendSmsCodeReq{
		Scene: iamv1.SmsScene(req.Scene),
		Phone: req.Phone,
		Risk:  toProtoRisk(req.Risk),
	})
	if err != nil {
		return nil, err
	}
	return &authv1.SendSmsCodeRes{ResendAfterSeconds: resp.ResendAfterSeconds}, nil
}

// RegisterByPassword 处理密码注册入口（含确认密码校验）。
func (c *ControllerV1) RegisterByPassword(ctx context.Context, req *authv1.RegisterByPasswordReq) (*authv1.RegisterByPasswordRes, error) {
	if strings.TrimSpace(req.Password) != strings.TrimSpace(req.ConfirmPassword) {
		return nil, errs.New(errs.CodeInvalidParam, "两次密码输入不一致")
	}

	resp, err := c.auth.RegisterByPassword(withRequestMetadata(ctx), &iamv1.RegisterByPasswordReq{
		Phone:    req.Phone,
		SmsCode:  req.SmsCode,
		Password: req.Password,
		Risk:     toProtoRisk(req.Risk),
	})
	if err != nil {
		return nil, err
	}

	return &authv1.RegisterByPasswordRes{
		UserID:          resp.UserId,
		InitDisplayName: resp.InitDisplayName,
		Auth:            toHTTPAuthResult(resp.Auth),
		Session:         toHTTPSession(resp.Session),
	}, nil
}

// LoginByPassword 密码登录。
func (c *ControllerV1) LoginByPassword(ctx context.Context, req *authv1.LoginByPasswordReq) (*authv1.LoginByPasswordRes, error) {
	resp, err := c.auth.LoginByPassword(withRequestMetadata(ctx), &iamv1.LoginByPasswordReq{
		Identifier: req.Identifier,
		Password:   req.Password,
		Risk:       toProtoRisk(req.Risk),
	})
	if err != nil {
		return nil, err
	}
	return &authv1.LoginByPasswordRes{
		Channel: int32(resp.Channel),
		Auth:    toHTTPAuthResult(resp.Auth),
		Session: toHTTPSession(resp.Session),
	}, nil
}

// LoginBySms 短信登录。
func (c *ControllerV1) LoginBySms(ctx context.Context, req *authv1.LoginBySmsReq) (*authv1.LoginBySmsRes, error) {
	resp, err := c.auth.LoginBySms(withRequestMetadata(ctx), &iamv1.LoginBySmsReq{
		Phone:   req.Phone,
		SmsCode: req.SmsCode,
		Risk:    toProtoRisk(req.Risk),
	})
	if err != nil {
		return nil, err
	}
	return &authv1.LoginBySmsRes{
		Channel: int32(resp.Channel),
		Auth:    toHTTPAuthResult(resp.Auth),
		Session: toHTTPSession(resp.Session),
	}, nil
}

// VerifyMfa 完成 MFA 二次验证。
func (c *ControllerV1) VerifyMfa(ctx context.Context, req *authv1.VerifyMfaReq) (*authv1.VerifyMfaRes, error) {
	resp, err := c.auth.VerifyMfaChallenge(withRequestMetadata(ctx), &iamv1.VerifyMfaChallengeReq{
		ChallengeId: req.ChallengeID,
		SmsCode:     req.SmsCode,
		Risk:        toProtoRisk(req.Risk),
	})
	if err != nil {
		return nil, err
	}
	return &authv1.VerifyMfaRes{
		TokenPair: toHTTPTokenPair(resp.TokenPair),
		Session:   toHTTPSession(resp.Session),
	}, nil
}

// RefreshToken 刷新访问令牌。
func (c *ControllerV1) RefreshToken(ctx context.Context, req *authv1.RefreshTokenReq) (*authv1.RefreshTokenRes, error) {
	resp, err := c.auth.RefreshToken(withRequestMetadata(ctx), &iamv1.RefreshTokenReq{
		RefreshToken: req.RefreshToken,
		Risk:         toProtoRisk(req.Risk),
	})
	if err != nil {
		return nil, err
	}
	return &authv1.RefreshTokenRes{TokenPair: toHTTPTokenPair(resp.TokenPair)}, nil
}

// Logout 注销当前会话或全设备会话。
func (c *ControllerV1) Logout(ctx context.Context, req *authv1.LogoutReq) (*authv1.LogoutRes, error) {
	_, err := c.auth.Logout(withRequestMetadata(ctx), &iamv1.LogoutReq{
		RefreshToken: req.RefreshToken,
		AllDevices:   req.AllDevices,
	})
	if err != nil {
		return nil, err
	}
	return &authv1.LogoutRes{}, nil
}

// GetMySession 查询当前登录会话摘要。
func (c *ControllerV1) GetMySession(ctx context.Context, _ *authv1.GetMySessionReq) (*authv1.GetMySessionRes, error) {
	resp, err := c.auth.GetMySession(withRequestMetadata(ctx), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return &authv1.GetMySessionRes{Session: toHTTPSession(resp.Session)}, nil
}

// ChangePassword 登录态修改密码。
func (c *ControllerV1) ChangePassword(ctx context.Context, req *authv1.ChangePasswordReq) (*authv1.ChangePasswordRes, error) {
	resp, err := c.auth.ChangePassword(withRequestMetadata(ctx), &iamv1.ChangePasswordReq{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
		Risk:        toProtoRisk(req.Risk),
	})
	if err != nil {
		return nil, err
	}
	return &authv1.ChangePasswordRes{Updated: resp.Updated}, nil
}

// ResetPasswordBySms 短信重置密码。
func (c *ControllerV1) ResetPasswordBySms(ctx context.Context, req *authv1.ResetPasswordBySmsReq) (*authv1.ResetPasswordBySmsRes, error) {
	resp, err := c.auth.ResetPasswordBySms(withRequestMetadata(ctx), &iamv1.ResetPasswordBySmsReq{
		Phone:       req.Phone,
		SmsCode:     req.SmsCode,
		NewPassword: req.NewPassword,
		Risk:        toProtoRisk(req.Risk),
	})
	if err != nil {
		return nil, err
	}
	return &authv1.ResetPasswordBySmsRes{Updated: resp.Updated}, nil
}

// LoginByOAuth 第三方登录占位入口。
func (c *ControllerV1) LoginByOAuth(ctx context.Context, req *authv1.LoginByOAuthReq) (*authv1.LoginByOAuthRes, error) {
	resp, err := c.auth.LoginByOAuth(withRequestMetadata(ctx), &iamv1.LoginByOAuthReq{
		Provider: iamv1.OAuthProvider(req.Provider),
		Code:     req.Code,
		State:    req.State,
		Risk:     toProtoRisk(req.Risk),
	})
	if err != nil {
		return nil, err
	}
	return &authv1.LoginByOAuthRes{
		Channel: int32(resp.Channel),
		Auth:    toHTTPAuthResult(resp.Auth),
		Session: toHTTPSession(resp.Session),
	}, nil
}

// BindOAuth 第三方账号绑定占位入口。
func (c *ControllerV1) BindOAuth(ctx context.Context, req *authv1.BindOAuthReq) (*authv1.BindOAuthRes, error) {
	resp, err := c.auth.BindOAuth(withRequestMetadata(ctx), &iamv1.BindOAuthReq{
		Provider: iamv1.OAuthProvider(req.Provider),
		Code:     req.Code,
		State:    req.State,
		Risk:     toProtoRisk(req.Risk),
	})
	if err != nil {
		return nil, err
	}
	return &authv1.BindOAuthRes{Bound: resp.Bound}, nil
}

// UnbindOAuth 第三方账号解绑占位入口。
func (c *ControllerV1) UnbindOAuth(ctx context.Context, req *authv1.UnbindOAuthReq) (*authv1.UnbindOAuthRes, error) {
	resp, err := c.auth.UnbindOAuth(withRequestMetadata(ctx), &iamv1.UnbindOAuthReq{
		Provider:    iamv1.OAuthProvider(req.Provider),
		ProviderUid: req.ProviderUID,
		Risk:        toProtoRisk(req.Risk),
	})
	if err != nil {
		return nil, err
	}
	return &authv1.UnbindOAuthRes{Unbound: resp.Unbound}, nil
}
