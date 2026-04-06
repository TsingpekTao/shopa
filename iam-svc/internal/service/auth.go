// ================================================================================
// 该文件由 GoFrame CLI 生成，承载 service 层对外的鉴权接口契约。
// logic 层通过 RegisterAuth 注册实现，其他业务模块则统一通过 Auth() 获取实例。
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	// IAuth 约束 IAM 认证域能力边界。
	// 约定：每个方法都应为“幂等读”或“显式写入”，并返回统一错误码语义。
	IAuth interface {
		// StartBackgroundWorkers 启动 outbox 投递、归档、清理与指标后台任务。
		StartBackgroundWorkers(ctx context.Context)
		// SendSmsCode 发送短信验证码（注册/登录/MFA/重置密码场景）。
		SendSmsCode(ctx context.Context, req *v1.SendSmsCodeReq) (*v1.SendSmsCodeRes, error)
		// RegisterByPassword 使用手机号+短信验证码+密码完成注册。
		RegisterByPassword(ctx context.Context, req *v1.RegisterByPasswordReq) (*v1.RegisterByPasswordRes, error)
		// LoginByPassword 账号密码登录（可按策略触发 MFA）。
		LoginByPassword(ctx context.Context, req *v1.LoginByPasswordReq) (*v1.LoginByPasswordRes, error)
		// LoginBySms 手机号短信码登录。
		LoginBySms(ctx context.Context, req *v1.LoginBySmsReq) (*v1.LoginBySmsRes, error)
		// VerifyMfaChallenge 校验二次验证码并签发最终登录 token。
		VerifyMfaChallenge(ctx context.Context, req *v1.VerifyMfaChallengeReq) (*v1.VerifyMfaChallengeRes, error)
		// RefreshToken 使用 refresh token 轮换并换发新的 token 对。
		RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (*v1.RefreshTokenRes, error)
		// Logout 注销会话（当前设备或全部设备）。
		Logout(ctx context.Context, req *v1.LogoutReq) (*emptypb.Empty, error)
		// GetMySession 获取当前 access token 对应会话摘要。
		GetMySession(ctx context.Context, req *emptypb.Empty) (*v1.GetMySessionRes, error)
		// ChangePassword 登录态修改密码，并使历史会话失效。
		ChangePassword(ctx context.Context, req *v1.ChangePasswordReq) (*v1.ChangePasswordRes, error)
		// ResetPasswordBySms 通过短信码重置密码并失效历史会话。
		ResetPasswordBySms(ctx context.Context, req *v1.ResetPasswordBySmsReq) (*v1.ResetPasswordBySmsRes, error)
		// LoginByOAuth 第三方登录（当前为预留接口）。
		LoginByOAuth(ctx context.Context, req *v1.LoginByOAuthReq) (*v1.LoginByOAuthRes, error)
		// BindOAuth 绑定第三方账号（当前为预留接口）。
		BindOAuth(ctx context.Context, req *v1.BindOAuthReq) (*v1.BindOAuthRes, error)
		// UnbindOAuth 解绑第三方账号（当前为预留接口）。
		UnbindOAuth(ctx context.Context, req *v1.UnbindOAuthReq) (*v1.UnbindOAuthRes, error)
		// GetAuthUserById 内部接口：按 user_id 查询认证聚合信息。
		GetAuthUserById(ctx context.Context, req *v1.GetAuthUserByIdReq) (*v1.GetAuthUserByIdRes, error)
		// BatchGetAuthUsers 内部接口：批量查询认证聚合信息。
		BatchGetAuthUsers(ctx context.Context, req *v1.BatchGetAuthUsersReq) (*v1.BatchGetAuthUsersRes, error)
		// VerifyAccessToken 内部接口：校验 access token 并返回鉴权上下文。
		VerifyAccessToken(ctx context.Context, req *v1.VerifyAccessTokenReq) (*v1.VerifyAccessTokenRes, error)
		// HasShopRole 内部接口：判断用户是否已具备指定店铺的卖家角色。
		HasShopRole(ctx context.Context, req *v1.HasShopRoleReq) (*v1.HasShopRoleRes, error)
		// AssignShopSellerRole 内部接口：为用户授予指定店铺的卖家角色。
		AssignShopSellerRole(ctx context.Context, req *v1.AssignShopSellerRoleReq) (*v1.AssignShopSellerRoleRes, error)
		// RevokeShopSellerRoleAndBumpToken 内部接口：回收店铺卖家角色并提升 token_version 使旧令牌失效。
		RevokeShopSellerRoleAndBumpToken(ctx context.Context, req *v1.RevokeShopSellerRoleAndBumpTokenReq) (*v1.RevokeShopSellerRoleAndBumpTokenRes, error)
	}
)

var (
	localAuth IAuth
)

// Auth 返回已注册的 IAuth 实现；若未注册会 panic，避免在“无实现”状态下静默运行。
func Auth() IAuth {
	if localAuth == nil {
		panic("implement not found for interface IAuth, forgot register?")
	}
	return localAuth
}

// RegisterAuth 在进程 init 期间注入 IAuth 实现，调用者要确保并发调用不会重复注册。
func RegisterAuth(i IAuth) {
	localAuth = i
}
