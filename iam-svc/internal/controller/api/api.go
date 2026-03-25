package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/TsingpekTao/shopa/iam-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Controller 閹佃儻娴?gRPC 鐏炲倹甯撮崣锝忕礉娑撴艾濮熼柅鏄忕帆閸忋劑鍎存稉瀣焽閸?service閵
type Controller struct {
	v1.UnimplementedAuthServiceServer
	v1.UnimplementedInternalServiceServer
	auth service.IAuth
}

// Register 澶勭悊娉ㄥ唽涓绘祦绋嬪強鍒濆鍖栧姩浣溿€
func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{auth: service.Auth()}
	v1.RegisterAuthServiceServer(s.Server, ctrl)
	v1.RegisterInternalServiceServer(s.Server, ctrl)
}

// SendSmsCode 鍙戦€侀獙璇佺爜/娑堟伅骞惰繑鍥炵姸鎬併€
func (c *Controller) SendSmsCode(ctx context.Context, req *v1.SendSmsCodeReq) (*v1.SendSmsCodeRes, error) {
	return c.auth.SendSmsCode(ctx, req)
}

// RegisterByPassword 澶勭悊娉ㄥ唽涓绘祦绋嬪強鍒濆鍖栧姩浣溿€
func (c *Controller) RegisterByPassword(ctx context.Context, req *v1.RegisterByPasswordReq) (*v1.RegisterByPasswordRes, error) {
	return c.auth.RegisterByPassword(ctx, req)
}

// LoginByPassword 澶勭悊鐧诲綍璁よ瘉骞惰繑鍥炰細璇濈粨鏋溿€
func (c *Controller) LoginByPassword(ctx context.Context, req *v1.LoginByPasswordReq) (*v1.LoginByPasswordRes, error) {
	return c.auth.LoginByPassword(ctx, req)
}

// LoginBySms 澶勭悊鐧诲綍璁よ瘉骞惰繑鍥炰細璇濈粨鏋溿€
func (c *Controller) LoginBySms(ctx context.Context, req *v1.LoginBySmsReq) (*v1.LoginBySmsRes, error) {
	return c.auth.LoginBySms(ctx, req)
}

// VerifyMfaChallenge 鎵ц鏍￠獙閫昏緫骞惰繑鍥炴牎楠岀粨鏋溿€
func (c *Controller) VerifyMfaChallenge(ctx context.Context, req *v1.VerifyMfaChallengeReq) (*v1.VerifyMfaChallengeRes, error) {
	return c.auth.VerifyMfaChallenge(ctx, req)
}

// RefreshToken 鍒锋柊浠ょ墝骞惰繑鍥炴柊鍑瘉銆
func (c *Controller) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (*v1.RefreshTokenRes, error) {
	return c.auth.RefreshToken(ctx, req)
}

// Logout 鍚婇攢浼氳瘽骞舵竻鐞嗙櫥褰曠姸鎬併€
func (c *Controller) Logout(ctx context.Context, req *v1.LogoutReq) (*emptypb.Empty, error) {
	return c.auth.Logout(ctx, req)
}

// GetMySession 鎸夋潯浠惰鍙栧苟杩斿洖鍗曟潯缁撴灉銆
func (c *Controller) GetMySession(ctx context.Context, req *emptypb.Empty) (*v1.GetMySessionRes, error) {
	return c.auth.GetMySession(ctx, req)
}

// ChangePassword 鎵ц鍙樻洿鎿嶄綔骞惰繑鍥炵粨鏋溿€
func (c *Controller) ChangePassword(ctx context.Context, req *v1.ChangePasswordReq) (*v1.ChangePasswordRes, error) {
	return c.auth.ChangePassword(ctx, req)
}

// ResetPasswordBySms 鎵ц閲嶇疆娴佺▼骞惰繑鍥炲鐞嗙粨鏋溿€
func (c *Controller) ResetPasswordBySms(ctx context.Context, req *v1.ResetPasswordBySmsReq) (*v1.ResetPasswordBySmsRes, error) {
	return c.auth.ResetPasswordBySms(ctx, req)
}

// LoginByOAuth 澶勭悊鐧诲綍璁よ瘉骞惰繑鍥炰細璇濈粨鏋溿€
func (c *Controller) LoginByOAuth(ctx context.Context, req *v1.LoginByOAuthReq) (*v1.LoginByOAuthRes, error) {
	return c.auth.LoginByOAuth(ctx, req)
}

// BindOAuth 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (c *Controller) BindOAuth(ctx context.Context, req *v1.BindOAuthReq) (*v1.BindOAuthRes, error) {
	return c.auth.BindOAuth(ctx, req)
}

// UnbindOAuth 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (c *Controller) UnbindOAuth(ctx context.Context, req *v1.UnbindOAuthReq) (*v1.UnbindOAuthRes, error) {
	return c.auth.UnbindOAuth(ctx, req)
}

// GetAuthUserById 鎸夋潯浠惰鍙栧苟杩斿洖鍗曟潯缁撴灉銆
func (c *Controller) GetAuthUserById(ctx context.Context, req *v1.GetAuthUserByIdReq) (*v1.GetAuthUserByIdRes, error) {
	return c.auth.GetAuthUserById(ctx, req)
}

// BatchGetAuthUsers 鎵归噺澶勭悊璇锋眰锛屽噺灏戝線杩斿紑閿€銆
func (c *Controller) BatchGetAuthUsers(ctx context.Context, req *v1.BatchGetAuthUsersReq) (*v1.BatchGetAuthUsersRes, error) {
	return c.auth.BatchGetAuthUsers(ctx, req)
}

// VerifyAccessToken 鎵ц鏍￠獙閫昏緫骞惰繑鍥炴牎楠岀粨鏋溿€
func (c *Controller) VerifyAccessToken(ctx context.Context, req *v1.VerifyAccessTokenReq) (*v1.VerifyAccessTokenRes, error) {
	return c.auth.VerifyAccessToken(ctx, req)
}
