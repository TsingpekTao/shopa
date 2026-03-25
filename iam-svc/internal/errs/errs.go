package errs

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// IAM 娑撴艾濮熼柨娆掝嚖閻緤绱濈紓鏍垳鐟欏嫬鍨敍娆癆BBBB閵?// AA 鐞涖劎銇氶張宥呭閸欓攱顔岄敍鍦汚M=10閿涘绱滲BBB 鐞涖劎銇氶崷鐑樻珯缂佸棗鍨庨妴
var (
	CodeOK = gcode.New(0, "OK", nil)

	// 闁氨鏁ら柨娆掝嚖閵
	CodeInvalidParam  = gcode.New(100001, "Invalid parameter", nil)
	CodeInternalError = gcode.New(100002, "Internal error", nil)

	// 濞夈劌鍞?閻ц缍嶇€瑰鍙忛惄绋垮彠闁挎瑨顕ら妴
	CodePhoneRegistered     = gcode.New(101001, "Phone already registered", nil)
	CodeEmailRegistered     = gcode.New(101002, "Email already registered", nil)
	CodeInvalidCredential   = gcode.New(101003, "Account or credential invalid", nil)
	CodeAccountLocked       = gcode.New(101004, "Account locked", nil)
	CodeSmsCodeInvalid      = gcode.New(101005, "Verification code invalid", nil)
	CodeSmsCodeExpired      = gcode.New(101006, "Verification code expired", nil)
	CodeSmsTooFrequent      = gcode.New(101007, "Verification code requested too frequently", nil)
	CodeRegisterTooFrequent = gcode.New(101008, "Register attempts too frequent", nil)
	CodeCaptchaFailed       = gcode.New(101009, "Captcha verification failed", nil)
	CodeRiskRejected        = gcode.New(101010, "Risk control rejected", nil)
	CodeMFARequired         = gcode.New(101011, "Secondary verification required", nil)
	CodeMFAFailed           = gcode.New(101012, "Secondary verification failed", nil)
	CodeAccessTokenInvalid  = gcode.New(101013, "Access token invalid", nil)
	CodeAccessTokenExpired  = gcode.New(101014, "Access token expired", nil)
	CodeRefreshTokenInvalid = gcode.New(101015, "Refresh token invalid", nil)
	CodeAccountDisabled     = gcode.New(101016, "Account disabled", nil)
	CodeThirdPartyBound     = gcode.New(101017, "Third-party account already bound", nil)
)

// New 鏋勯€犲嚱鏁帮細鍒涘缓骞惰繑鍥炴湇鍔″疄渚嬨€
func New(code gcode.Code, msg ...string) error {
	if len(msg) > 0 && msg[0] != "" {
		return gerror.NewCode(code, msg[0])
	}
	return gerror.NewCode(code, code.Message())
}

// Wrap 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func Wrap(code gcode.Code, cause error, msg ...string) error {
	if cause == nil {
		return New(code, msg...)
	}
	if len(msg) > 0 && msg[0] != "" {
		return gerror.WrapCode(code, cause, msg[0])
	}
	return gerror.WrapCode(code, cause, code.Message())
}

// ToBiz 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func ToBiz(err error) (code int, message string) {
	if err == nil {
		return CodeOK.Code(), CodeOK.Message()
	}
	c := gerror.Code(err)
	if c != gcode.CodeNil && c.Code() >= 100000 {
		return c.Code(), c.Message()
	}

	switch c {
	case gcode.CodeInvalidParameter, gcode.CodeMissingParameter, gcode.CodeValidationFailed, gcode.CodeInvalidRequest:
		return CodeInvalidParam.Code(), CodeInvalidParam.Message()
	case gcode.CodeNotAuthorized, gcode.CodeSecurityReason:
		return CodeAccessTokenInvalid.Code(), CodeAccessTokenInvalid.Message()
	case gcode.CodeNotFound:
		return CodeInvalidParam.Code(), CodeInvalidParam.Message()
	}
	return CodeInternalError.Code(), CodeInternalError.Message()
}
