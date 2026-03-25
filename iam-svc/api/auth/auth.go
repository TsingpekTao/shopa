package auth

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/auth/v1"
)

// IAuthV1 鐎规矮绠?IAM 鐎电懓顦?HTTP 閹恒儱褰涢懗钘夊閵
type IAuthV1 interface {
	SendSmsCode(ctx context.Context, req *v1.SendSmsCodeReq) (res *v1.SendSmsCodeRes, err error)
	RegisterByPassword(ctx context.Context, req *v1.RegisterByPasswordReq) (res *v1.RegisterByPasswordRes, err error)
	LoginByPassword(ctx context.Context, req *v1.LoginByPasswordReq) (res *v1.LoginByPasswordRes, err error)
	LoginBySms(ctx context.Context, req *v1.LoginBySmsReq) (res *v1.LoginBySmsRes, err error)
	VerifyMfa(ctx context.Context, req *v1.VerifyMfaReq) (res *v1.VerifyMfaRes, err error)
	RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error)
	Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error)
	GetMySession(ctx context.Context, req *v1.GetMySessionReq) (res *v1.GetMySessionRes, err error)
	ChangePassword(ctx context.Context, req *v1.ChangePasswordReq) (res *v1.ChangePasswordRes, err error)
	ResetPasswordBySms(ctx context.Context, req *v1.ResetPasswordBySmsReq) (res *v1.ResetPasswordBySmsRes, err error)
	LoginByOAuth(ctx context.Context, req *v1.LoginByOAuthReq) (res *v1.LoginByOAuthRes, err error)
	BindOAuth(ctx context.Context, req *v1.BindOAuthReq) (res *v1.BindOAuthRes, err error)
	UnbindOAuth(ctx context.Context, req *v1.UnbindOAuthReq) (res *v1.UnbindOAuthRes, err error)
}
