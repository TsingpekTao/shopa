package errs

import (
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// TestToBiz_PreservesBusinessCode 验证业务错误码不会被覆盖。
func TestToBiz_PreservesBusinessCode(t *testing.T) {
	err := New(CodePhoneRegistered)
	code, msg := ToBiz(err)
	if code != 101001 {
		t.Fatalf("unexpected code: %d", code)
	}
	if msg == "" {
		t.Fatalf("empty message")
	}
}

// TestToBiz_MapsInvalidParameter 验证无效参数错误映射到业务码。
func TestToBiz_MapsInvalidParameter(t *testing.T) {
	err := gerror.NewCode(gcode.CodeInvalidParameter, "x")
	code, _ := ToBiz(err)
	if code != 100001 {
		t.Fatalf("unexpected code: %d", code)
	}
}

// TestToBiz_MapsNotAuthorized 验证未授权错误映射到 access token 失效码。
func TestToBiz_MapsNotAuthorized(t *testing.T) {
	err := gerror.NewCode(gcode.CodeNotAuthorized, "x")
	code, _ := ToBiz(err)
	if code != 101013 {
		t.Fatalf("unexpected code: %d", code)
	}
}

// TestToBiz_DefaultsToInternal 验证未知错误默认映射到内部错误码。
func TestToBiz_DefaultsToInternal(t *testing.T) {
	err := gerror.New("boom")
	code, _ := ToBiz(err)
	if code != 100002 {
		t.Fatalf("unexpected code: %d", code)
	}
}
