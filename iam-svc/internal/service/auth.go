// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	IAuth interface {
		StartBackgroundWorkers(ctx context.Context)
		SendSmsCode(ctx context.Context, req *v1.SendSmsCodeReq) (*v1.SendSmsCodeRes, error)
		RegisterByPassword(ctx context.Context, req *v1.RegisterByPasswordReq) (*v1.RegisterByPasswordRes, error)
		LoginByPassword(ctx context.Context, req *v1.LoginByPasswordReq) (*v1.LoginByPasswordRes, error)
		LoginBySms(ctx context.Context, req *v1.LoginBySmsReq) (*v1.LoginBySmsRes, error)
		VerifyMfaChallenge(ctx context.Context, req *v1.VerifyMfaChallengeReq) (*v1.VerifyMfaChallengeRes, error)
		RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (*v1.RefreshTokenRes, error)
		Logout(ctx context.Context, req *v1.LogoutReq) (*emptypb.Empty, error)
		GetMySession(ctx context.Context, req *emptypb.Empty) (*v1.GetMySessionRes, error)
		ChangePassword(ctx context.Context, req *v1.ChangePasswordReq) (*v1.ChangePasswordRes, error)
		ResetPasswordBySms(ctx context.Context, req *v1.ResetPasswordBySmsReq) (*v1.ResetPasswordBySmsRes, error)
		LoginByOAuth(ctx context.Context, req *v1.LoginByOAuthReq) (*v1.LoginByOAuthRes, error)
		BindOAuth(ctx context.Context, req *v1.BindOAuthReq) (*v1.BindOAuthRes, error)
		UnbindOAuth(ctx context.Context, req *v1.UnbindOAuthReq) (*v1.UnbindOAuthRes, error)
		GetAuthUserById(ctx context.Context, req *v1.GetAuthUserByIdReq) (*v1.GetAuthUserByIdRes, error)
		BatchGetAuthUsers(ctx context.Context, req *v1.BatchGetAuthUsersReq) (*v1.BatchGetAuthUsersRes, error)
		VerifyAccessToken(ctx context.Context, req *v1.VerifyAccessTokenReq) (*v1.VerifyAccessTokenRes, error)
	}
)

var (
	localAuth IAuth
)

// Auth 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func Auth() IAuth {
	if localAuth == nil {
		panic("implement not found for interface IAuth, forgot register?")
	}
	return localAuth
}

// RegisterAuth 澶勭悊娉ㄥ唽涓绘祦绋嬪強鍒濆鍖栧姩浣溿€
func RegisterAuth(i IAuth) {
	localAuth = i
}
