package auth

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/TsingpekTao/shopa/iam-svc/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

// sAuth 是 GoFrame service 层适配器，仅负责协议层转发逻辑。
// 具体状态管理和业务判定都由 core Service 承担。
type sAuth struct {
	core *Service
}

// newAuthLogic 创建 service 适配器并绑定 core 单例。
func newAuthLogic() *sAuth {
	return &sAuth{core: New()}
}

func init() {
	// 在包初始化阶段注册 IAuth 实现，供 service.Auth() 全局获取。
	service.RegisterAuth(newAuthLogic())
}

// StartBackgroundWorkers 启动 IAM outbox 相关的后台 worker。
func (s *sAuth) StartBackgroundWorkers(ctx context.Context) {
	s.core.StartBackgroundWorkers(ctx)
}

// SendSmsCode 发送短信验证码。
func (s *sAuth) SendSmsCode(ctx context.Context, req *v1.SendSmsCodeReq) (*v1.SendSmsCodeRes, error) {
	return s.core.SendSmsCode(ctx, req)
}

// RegisterByPassword 使用手机号、短信码和密码完成注册。
func (s *sAuth) RegisterByPassword(ctx context.Context, req *v1.RegisterByPasswordReq) (*v1.RegisterByPasswordRes, error) {
	return s.core.RegisterByPassword(ctx, req)
}

// LoginByPassword 使用账号密码登录。
func (s *sAuth) LoginByPassword(ctx context.Context, req *v1.LoginByPasswordReq) (*v1.LoginByPasswordRes, error) {
	return s.core.LoginByPassword(ctx, req)
}

// LoginBySms 使用手机号短信码登录。
func (s *sAuth) LoginBySms(ctx context.Context, req *v1.LoginBySmsReq) (*v1.LoginBySmsRes, error) {
	return s.core.LoginBySms(ctx, req)
}

// VerifyMfaChallenge 完成 MFA 二次校验。
func (s *sAuth) VerifyMfaChallenge(ctx context.Context, req *v1.VerifyMfaChallengeReq) (*v1.VerifyMfaChallengeRes, error) {
	return s.core.VerifyMfaChallenge(ctx, req)
}

// RefreshToken 刷新令牌对。
func (s *sAuth) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (*v1.RefreshTokenRes, error) {
	return s.core.RefreshToken(ctx, req)
}

// Logout 注销当前会话或全部会话。
func (s *sAuth) Logout(ctx context.Context, req *v1.LogoutReq) (*emptypb.Empty, error) {
	return s.core.Logout(ctx, req)
}

// GetMySession 查询当前登录态会话摘要。
func (s *sAuth) GetMySession(ctx context.Context, req *emptypb.Empty) (*v1.GetMySessionRes, error) {
	return s.core.GetMySession(ctx, req)
}

// ChangePassword 在登录态下修改密码。
func (s *sAuth) ChangePassword(ctx context.Context, req *v1.ChangePasswordReq) (*v1.ChangePasswordRes, error) {
	return s.core.ChangePassword(ctx, req)
}

// ResetPasswordBySms 通过短信码重置密码。
func (s *sAuth) ResetPasswordBySms(ctx context.Context, req *v1.ResetPasswordBySmsReq) (*v1.ResetPasswordBySmsRes, error) {
	return s.core.ResetPasswordBySms(ctx, req)
}

// LoginByOAuth 第三方登录占位接口。
func (s *sAuth) LoginByOAuth(ctx context.Context, req *v1.LoginByOAuthReq) (*v1.LoginByOAuthRes, error) {
	return s.core.LoginByOAuth(ctx, req)
}

// BindOAuth 第三方账号绑定占位接口。
func (s *sAuth) BindOAuth(ctx context.Context, req *v1.BindOAuthReq) (*v1.BindOAuthRes, error) {
	return s.core.BindOAuth(ctx, req)
}

// UnbindOAuth 第三方账号解绑占位接口。
func (s *sAuth) UnbindOAuth(ctx context.Context, req *v1.UnbindOAuthReq) (*v1.UnbindOAuthRes, error) {
	return s.core.UnbindOAuth(ctx, req)
}

// GetAuthUserById 供内部服务按 user_id 查询认证信息。
func (s *sAuth) GetAuthUserById(ctx context.Context, req *v1.GetAuthUserByIdReq) (*v1.GetAuthUserByIdRes, error) {
	return s.core.GetAuthUserById(ctx, req)
}

// BatchGetAuthUsers 供内部服务批量查询认证信息。
func (s *sAuth) BatchGetAuthUsers(ctx context.Context, req *v1.BatchGetAuthUsersReq) (*v1.BatchGetAuthUsersRes, error) {
	return s.core.BatchGetAuthUsers(ctx, req)
}

// VerifyAccessToken 供内部服务验证 access token 并返回鉴权上下文。
func (s *sAuth) VerifyAccessToken(ctx context.Context, req *v1.VerifyAccessTokenReq) (*v1.VerifyAccessTokenRes, error) {
	return s.core.VerifyAccessToken(ctx, req)
}

// HasShopRole 判断用户是否拥有指定店铺作用域下的卖家角色。
func (s *sAuth) HasShopRole(ctx context.Context, req *v1.HasShopRoleReq) (*v1.HasShopRoleRes, error) {
	return s.core.HasShopRole(ctx, req)
}

// AssignShopSellerRole 为用户授予指定店铺的卖家角色。
func (s *sAuth) AssignShopSellerRole(ctx context.Context, req *v1.AssignShopSellerRoleReq) (*v1.AssignShopSellerRoleRes, error) {
	return s.core.AssignShopSellerRole(ctx, req)
}

// RevokeShopSellerRoleAndBumpToken 回收卖家角色并提升 token_version。
func (s *sAuth) RevokeShopSellerRoleAndBumpToken(ctx context.Context, req *v1.RevokeShopSellerRoleAndBumpTokenReq) (*v1.RevokeShopSellerRoleAndBumpTokenRes, error) {
	return s.core.RevokeShopSellerRoleAndBumpToken(ctx, req)
}
