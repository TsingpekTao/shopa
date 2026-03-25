package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/TsingpekTao/shopa/iam-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Controller 封装 Auth/ Internal gRPC 服务，透传基础 auth 逻辑。
type Controller struct {
	v1.UnimplementedAuthServiceServer
	v1.UnimplementedInternalServiceServer
	auth service.IAuth
}

// Register 向 gRPC 服务器注册 auth 与 internal 服务。
func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{auth: service.Auth()}
	v1.RegisterAuthServiceServer(s.Server, ctrl)
	v1.RegisterInternalServiceServer(s.Server, ctrl)
}

// SendSmsCode 通过 auth 逻辑发送短信验证码。
func (c *Controller) SendSmsCode(ctx context.Context, req *v1.SendSmsCodeReq) (*v1.SendSmsCodeRes, error) {
	return c.auth.SendSmsCode(ctx, req)
}

// Register 向 gRPC 服务器注册 auth 与 internal 服务。
func (c *Controller) RegisterByPassword(ctx context.Context, req *v1.RegisterByPasswordReq) (*v1.RegisterByPasswordRes, error) {
	return c.auth.RegisterByPassword(ctx, req)
}

// LoginByPassword 使用密码登录授权用户。
func (c *Controller) LoginByPassword(ctx context.Context, req *v1.LoginByPasswordReq) (*v1.LoginByPasswordRes, error) {
	return c.auth.LoginByPassword(ctx, req)
}

// LoginBySms 使用短信验证码登录。
func (c *Controller) LoginBySms(ctx context.Context, req *v1.LoginBySmsReq) (*v1.LoginBySmsRes, error) {
	return c.auth.LoginBySms(ctx, req)
}

// VerifyMfaChallenge 校验多因素挑战响应。
func (c *Controller) VerifyMfaChallenge(ctx context.Context, req *v1.VerifyMfaChallengeReq) (*v1.VerifyMfaChallengeRes, error) {
	return c.auth.VerifyMfaChallenge(ctx, req)
}

// RefreshToken 刷新访问令牌，延长会话。
func (c *Controller) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (*v1.RefreshTokenRes, error) {
	return c.auth.RefreshToken(ctx, req)
}

// Logout 注销会话并释放登录状态。
func (c *Controller) Logout(ctx context.Context, req *v1.LogoutReq) (*emptypb.Empty, error) {
	return c.auth.Logout(ctx, req)
}

// GetMySession 查询当前会话信息。
func (c *Controller) GetMySession(ctx context.Context, req *emptypb.Empty) (*v1.GetMySessionRes, error) {
	return c.auth.GetMySession(ctx, req)
}

// ChangePassword 修改用户登录密码。
func (c *Controller) ChangePassword(ctx context.Context, req *v1.ChangePasswordReq) (*v1.ChangePasswordRes, error) {
	return c.auth.ChangePassword(ctx, req)
}

// ResetPasswordBySms 通过短信重置密码。
func (c *Controller) ResetPasswordBySms(ctx context.Context, req *v1.ResetPasswordBySmsReq) (*v1.ResetPasswordBySmsRes, error) {
	return c.auth.ResetPasswordBySms(ctx, req)
}

// LoginByOAuth 通过 OAuth 提供者登录。
func (c *Controller) LoginByOAuth(ctx context.Context, req *v1.LoginByOAuthReq) (*v1.LoginByOAuthRes, error) {
	return c.auth.LoginByOAuth(ctx, req)
}

// BindOAuth 绑定第三方 OAuth 账号。
func (c *Controller) BindOAuth(ctx context.Context, req *v1.BindOAuthReq) (*v1.BindOAuthRes, error) {
	return c.auth.BindOAuth(ctx, req)
}

// UnbindOAuth 解除绑定的第三方 OAuth 账号。
func (c *Controller) UnbindOAuth(ctx context.Context, req *v1.UnbindOAuthReq) (*v1.UnbindOAuthRes, error) {
	return c.auth.UnbindOAuth(ctx, req)
}

// GetAuthUserById 根据 ID 获取授权用户详情。
func (c *Controller) GetAuthUserById(ctx context.Context, req *v1.GetAuthUserByIdReq) (*v1.GetAuthUserByIdRes, error) {
	return c.auth.GetAuthUserById(ctx, req)
}

// BatchGetAuthUsers 批量查询授权用户信息。
func (c *Controller) BatchGetAuthUsers(ctx context.Context, req *v1.BatchGetAuthUsersReq) (*v1.BatchGetAuthUsersRes, error) {
	return c.auth.BatchGetAuthUsers(ctx, req)
}

// VerifyAccessToken 校验访问令牌的有效性。
func (c *Controller) VerifyAccessToken(ctx context.Context, req *v1.VerifyAccessTokenReq) (*v1.VerifyAccessTokenRes, error) {
	return c.auth.VerifyAccessToken(ctx, req)
}
