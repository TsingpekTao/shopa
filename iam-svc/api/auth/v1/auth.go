package v1

import "github.com/gogf/gf/v2/frame/g"

// RiskContext 閹佃儻娴囬崜宥囶伂闁插洭娉﹂崚鎵畱鐠佹儳顦搴㈠付娣団剝浼呴妴
type RiskContext struct {
	Fingerprint  string `json:"fingerprint" dc:"鐠佹儳顦幐鍥╂睏"`
	DeviceID     string `json:"deviceId" dc:"鐠佹儳顦琁D"`
	CaptchaToken string `json:"captchaToken" dc:"娴滅儤婧€妤犲矁鐦塼oken"`
}

// TokenPair 閹诲繗鍫?access/refresh 娴犮倗澧濈€靛箍鈧
type TokenPair struct {
	TokenType        string `json:"tokenType" dc:"娴犮倗澧濈猾璇茬€烽敍宀勨偓姘埗娑?Bearer"`
	AccessToken      string `json:"accessToken" dc:"鐠佸潡妫舵禒銈囧"`
	AccessExpiresIn  uint32 `json:"accessExpiresIn" dc:"鐠佸潡妫舵禒銈囧鏉╁洦婀＄粔鎺撴殶"`
	RefreshToken     string `json:"refreshToken,omitempty" dc:"閸掗攱鏌婃禒銈囧"`
	RefreshExpiresIn uint32 `json:"refreshExpiresIn,omitempty" dc:"閸掗攱鏌婃禒銈囧鏉╁洦婀＄粔鎺撴殶"`
	Sid              string `json:"sid" dc:"娴兼俺鐦絀D"`
}

// MfaChallenge 鐞涖劎銇氶棁鈧憰浣风癌濞嗏剝鐗庢宀€娈戦幐鎴炲灛娣団剝浼呴妴
type MfaChallenge struct {
	ChallengeID string `json:"challengeId" dc:"MFA閹告垶鍨琁D"`
	ExpireAt    string `json:"expireAt,omitempty" dc:"RFC3339鏉╁洦婀￠弮鍫曟？"`
}

// AuthResult 閻ц缍嶇紒鎾寸亯閺勵垯绨╅柅澶夌閿涙氨娲块幒銉ㄧ箲閸?token閿涘本鍨ㄦ潻鏂挎礀 MFA challenge閵
type AuthResult struct {
	TokenPair    *TokenPair    `json:"tokenPair,omitempty"`
	MfaChallenge *MfaChallenge `json:"mfaChallenge,omitempty"`
}

type RoleItem struct {
	RoleCode  int32  `json:"roleCode" dc:"1妞ゆ儳顓?2鎼存ぞ瀵?3缁狅紕鎮婇崨?4鐎广垺婀?`
	ScopeType int32  `json:"scopeType" dc:"1閸忋劌鐪?2鎼存鎽?`
	ScopeID   uint64 `json:"scopeId" dc:"娴ｆ粎鏁ら崺鐑瓺"`
}

type MembershipSummary struct {
	LevelCode string `json:"levelCode" dc:"娴兼艾鎲崇粵澶岄獓缂傛牜鐖?`
	Points    uint64 `json:"points" dc:"娴兼艾鎲崇粔顖氬瀻"`
	ExpireAt  string `json:"expireAt,omitempty" dc:"RFC3339鏉╁洦婀￠弮鍫曟？"`
}

type SessionSummary struct {
	UserID        uint64             `json:"userId"`
	AccountStatus int32              `json:"accountStatus" dc:"1濮濓絽鐖?2闁夸礁鐣?3缁備胶鏁?4鐎光剝鐗虫稉?`
	Roles         []*RoleItem        `json:"roles"`
	Membership    *MembershipSummary `json:"membership"`
	LastLoginAt   string             `json:"lastLoginAt,omitempty" dc:"RFC3339閻ц缍嶉弮鍫曟？"`
	LastLoginIP   string             `json:"lastLoginIp,omitempty"`
}

type SendSmsCodeReq struct {
	g.Meta `path:"/v1/auth/sms/send" method:"post" tags:"Auth" summary:"閸欐垿鈧胶鐓穱锟犵崣鐠囦胶鐖?`
	Scene  int32        `json:"scene" v:"required|in:1,2,3,4" dc:"1濞夈劌鍞?2閻ц缍?3MFA 4闁插秶鐤嗙€靛棛鐖?`
	Phone  string       `json:"phone" v:"required|phone" dc:"閹靛婧€閸?`
	Risk   *RiskContext `json:"risk" dc:"妞嬪孩甯舵稉濠佺瑓閺?`
}

type SendSmsCodeRes struct {
	ResendAfterSeconds uint32 `json:"resendAfterSeconds"`
}

type RegisterByPasswordReq struct {
	g.Meta  `path:"/v1/auth/register/password" method:"post" tags:"Auth" summary:"閹靛婧€閸?妤犲矁鐦夐惍?鐎靛棛鐖滃▔銊ュ斀"`
	Phone   string `json:"phone" v:"required|phone"`
	SmsCode string `json:"smsCode" v:"required|length:4,8"`
	// HTTP 鐏炲倸顤冮崝鐘碘€樼拋銈呯槕閻礁鐡у▓纰夌礉閸忓牆婀幒褍鍩楅崳銊﹀閹搭亙绗夋稉鈧懛纾嬵嚞濮瑰偊绱?	// gRPC 閸愬懘鍎存總鎴犲娴犲秳绻氶幐浣稿礋 password閿涘矂浼╅崗宥呭閸濆秵婀囬崝锟犳？鐠嬪啰鏁ら妴
	Password        string       `json:"password" v:"required|length:8,20"`
	ConfirmPassword string       `json:"confirmPassword" v:"required|length:8,20"`
	Risk            *RiskContext `json:"risk"`
}

type RegisterByPasswordRes struct {
	UserID          uint64          `json:"userId"`
	InitDisplayName string          `json:"initDisplayName"`
	Auth            *AuthResult     `json:"auth"`
	Session         *SessionSummary `json:"session"`
}

type LoginByPasswordReq struct {
	g.Meta     `path:"/v1/auth/login/password" method:"post" tags:"Auth" summary:"鐠愶箑褰跨€靛棛鐖滈惂璇茬秿"`
	Identifier string       `json:"identifier" v:"required|length:3,128" dc:"閹靛婧€閸欓攱鍨ㄩ柇顔绢唸"`
	Password   string       `json:"password" v:"required|length:1,64"`
	Risk       *RiskContext `json:"risk"`
}

type LoginByPasswordRes struct {
	Channel int32           `json:"channel"`
	Auth    *AuthResult     `json:"auth"`
	Session *SessionSummary `json:"session"`
}

type LoginBySmsReq struct {
	g.Meta  `path:"/v1/auth/login/sms" method:"post" tags:"Auth" summary:"閻厺淇婃宀冪槈閻胶娅ヨぐ?`
	Phone   string       `json:"phone" v:"required|phone"`
	SmsCode string       `json:"smsCode" v:"required|length:4,8"`
	Risk    *RiskContext `json:"risk"`
}

type LoginBySmsRes struct {
	Channel int32           `json:"channel"`
	Auth    *AuthResult     `json:"auth"`
	Session *SessionSummary `json:"session"`
}

type VerifyMfaReq struct {
	g.Meta      `path:"/v1/me/mfa/verify" method:"post" tags:"Me" summary:"妤犲矁鐦塎FA閹告垶鍨?`
	ChallengeID string       `json:"challengeId" v:"required|length:6,64"`
	SmsCode     string       `json:"smsCode" v:"required|length:4,8"`
	Risk        *RiskContext `json:"risk"`
}

type VerifyMfaRes struct {
	TokenPair *TokenPair      `json:"tokenPair"`
	Session   *SessionSummary `json:"session"`
}

type RefreshTokenReq struct {
	g.Meta       `path:"/v1/auth/token/refresh" method:"post" tags:"Auth" summary:"閸掗攱鏌婄拋鍧楁６娴犮倗澧?`
	RefreshToken string       `json:"refreshToken" v:"required"`
	Risk         *RiskContext `json:"risk"`
}

type RefreshTokenRes struct {
	TokenPair *TokenPair `json:"tokenPair"`
}

type LogoutReq struct {
	g.Meta       `path:"/v1/auth/logout" method:"post" tags:"Auth" summary:"闁偓閸戣櫣娅ヨぐ?`
	RefreshToken string `json:"refreshToken" dc:"閸欘垶鈧绱辨稉铏光敄閺冭埖瀵滆ぐ鎾冲access token娴兼俺鐦介柅鈧崙?`
	AllDevices   bool   `json:"allDevices"`
}

type LogoutRes struct{}

type GetMySessionReq struct {
	g.Meta `path:"/v1/me/session" method:"get" tags:"Me" summary:"閼惧嘲褰囪ぐ鎾冲閻ц缍嶆导姘崇樈"`
}

type GetMySessionRes struct {
	Session *SessionSummary `json:"session"`
}

type ChangePasswordReq struct {
	g.Meta      `path:"/v1/auth/password/change" method:"post" tags:"Auth" summary:"閻ц缍嶉幀浣锋叏閺€鐟扮槕閻?`
	OldPassword string       `json:"oldPassword" v:"required|length:1,64"`
	NewPassword string       `json:"newPassword" v:"required|length:8,20"`
	Risk        *RiskContext `json:"risk"`
}

type ChangePasswordRes struct {
	Updated bool `json:"updated"`
}

type ResetPasswordBySmsReq struct {
	g.Meta      `path:"/v1/auth/password/reset/sms" method:"post" tags:"Auth" summary:"閻厺淇婃宀冪槈閻線鍣哥純顔肩槕閻?`
	Phone       string       `json:"phone" v:"required|phone"`
	SmsCode     string       `json:"smsCode" v:"required|length:4,8"`
	NewPassword string       `json:"newPassword" v:"required|length:8,20"`
	Risk        *RiskContext `json:"risk"`
}

type ResetPasswordBySmsRes struct {
	Updated bool `json:"updated"`
}

type LoginByOAuthReq struct {
	g.Meta   `path:"/v1/auth/login/oauth" method:"post" tags:"Auth" summary:"缁楊兛绗侀弬鍦瑜版洩绱欓崡鐘辩秴閿?`
	Provider int32        `json:"provider" v:"required|in:1,2" dc:"1瀵邦喕淇?2閺€顖欑帛鐎?`
	Code     string       `json:"code" v:"required"`
	State    string       `json:"state"`
	Risk     *RiskContext `json:"risk"`
}

type LoginByOAuthRes struct {
	Channel int32           `json:"channel"`
	Auth    *AuthResult     `json:"auth"`
	Session *SessionSummary `json:"session"`
}

type BindOAuthReq struct {
	g.Meta   `path:"/v1/auth/oauth/bind" method:"post" tags:"Auth" summary:"缂佹垵鐣剧粭顑跨瑏閺傜澶勯崣鍑ょ礄閸楃姳缍呴敍?`
	Provider int32        `json:"provider" v:"required|in:1,2"`
	Code     string       `json:"code" v:"required"`
	State    string       `json:"state"`
	Risk     *RiskContext `json:"risk"`
}

type BindOAuthRes struct {
	Bound bool `json:"bound"`
}

type UnbindOAuthReq struct {
	g.Meta      `path:"/v1/auth/oauth/unbind" method:"post" tags:"Auth" summary:"鐟欙絿绮︾粭顑跨瑏閺傜澶勯崣鍑ょ礄閸楃姳缍呴敍?`
	Provider    int32        `json:"provider" v:"required|in:1,2"`
	ProviderUID string       `json:"providerUid" v:"required"`
	Risk        *RiskContext `json:"risk"`
}

type UnbindOAuthRes struct {
	Unbound bool `json:"unbound"`
}
