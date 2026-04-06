package auth

import (
	"context"

	authv1 "github.com/TsingpekTao/shopa/iam-svc/api/auth/v1"
	iamv1 "github.com/TsingpekTao/shopa/iam-svc/api/v1"
	"github.com/TsingpekTao/shopa/iam-svc/internal/errs"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SendSmsCode 闂?HTTP 闁荤姴娲弨閬嶆儑閻楀牊濮滄い鎺嶇鎼村﹤鈽夐幘铏儓闁哥偞鎸抽弻?gRPC 闁荤姴娲弨閬嶆儑娴煎宓侀柟顖滃缁侇噣鏌涘▎鎰仴闁糕晛鐭傚鐣岀矙鐠恒劌顩梻浣规綑閻°劌顭囬崼銉ュ珘鐎广儱鎳庨～銈夋煏?
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

// RegisterByPassword 婵犮垼娉涚€氼噣骞冩繝鍕ㄥ亾闂堟稒顥犻柣鏍ㄧ矊閳绘棃濡搁妷銉︽澑闂佺绻堥崕杈亹濞戙垺鏅柛顐ｇ箖閸庢挾鐥褍鏋︽い鎾瑰吹閳ь剟娼уΛ娑㈡偉濠婂牆鍐€闁跨喓濮峰畷锝夋煥濞戞﹩妲堕柍?
func (c *ControllerV1) RegisterByPassword(ctx context.Context, req *authv1.RegisterByPasswordReq) (*authv1.RegisterByPasswordRes, error) {
	if req.Password != req.ConfirmPassword {
		return nil, errs.New(errs.CodeInvalidParam, "\u4e24\u6b21\u5bc6\u7801\u8f93\u5165\u4e0d\u4e00\u81f4")
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

// LoginByPassword 闁诲酣娼уΛ娑㈡偉濠婂牊鍎岄悹鍥皺缁夊潡鏌?
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

// LoginBySms 闂佹椿鍙庨崢鐑樼┍婵犲洦鍎岄悹鍥皺缁夊潡鏌?
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

// VerifyMfa 闁诲海鎳撻張顒勫垂?MFA 婵炲瓨绮岄張顒勵敃閸忕⒈娈界€光偓閸愵亝顫嶉梺?
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

// RefreshToken 闂佸憡甯￠弨閬嶅蓟婵犲嫭濯奸柛褎顨嗛敍鏍归悩顔尖偓妤佹櫠濠靛违?
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

// Logout 濠电偛顦崝鎴﹀绩閵忥絻浜归柟鎯у暱椤ゅ懎霉閸忛棿浜㈤柣锔跨矙楠炲寮借瀵潡鎮规担鐟板姢妞わ富鍓氱€电厧顫濆畷鍥ㄢ枔闂?
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

// GetMySession 闂佸搫琚崕鎾敋濡ゅ嫨浜归柟鎯у暱椤ゅ懘鏌ｈ椤曆呯礊瀹ュ棗顕辨慨妯虹－濡牓鏌熼懞銉劸妞も晪绠撴俊?
func (c *ControllerV1) GetMySession(ctx context.Context, _ *authv1.GetMySessionReq) (*authv1.GetMySessionRes, error) {
	resp, err := c.auth.GetMySession(withRequestMetadata(ctx), &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return &authv1.GetMySessionRes{Session: toHTTPSession(resp.Session)}, nil
}

// ChangePassword 闂佽皫鍡╁殭缂傚秴绉归獮鈧ù锝夋敱閸欏繘鏌￠埀顒勬偡閹殿喗顫氶梺娲诲枙闂勫秹鍩€?
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

// ResetPasswordBySms 闂佹椿鍙庨崢鐑樼┍婵犲洦鐓傜€广儱娲ㄩ弸鍌炴倵闂堟稒顥犻柣鏍ㄧ矒婵?
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

// LoginByOAuth 缂備焦顨忛崗娑氱箔娓氣偓瀵剟宕烽鐕佷划閻熸粎澧楀ú鏍х暦閻楀牊濯寸€广儱鎳庡鎶芥煕濞嗘瑧绁烽柍?
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

// BindOAuth 缂備焦顨忛崗娑氱箔娓氣偓瀵剛娑甸崨顓фП闂佸憡鐟╅ˉ鎾跺垝閿旈敮鍋撶憴鍕鐎规洜澧楅幏鍛吋閸涱厼姹查梺鍛婄懕缁插鍩€?
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

// UnbindOAuth 缂備焦顨忛崗娑氱箔娓氣偓瀵剛娑甸崨顓фП闂佸憡鐟ч弻澶庮暰缂傚倷鐒﹂崹闈涚暦閻楀牊濯寸€广儱鎳庡鎶芥煕濞嗘瑧绁烽柍?
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
