// =================================================================================
// 鏈枃浠跺畾涔夌敤鎴?API (v1) 鐨?HTTP 鎺у埗鍣ㄦ帴鍙ｃ€?
// 淇濇寔涓?GoFrame 鐢熸垚椋庢牸涓€鑷达紝纭繚鎺у埗鍣ㄧ粨鏋勭粺涓€銆?
// =================================================================================

package user

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/user-profile-svc/api/user/v1"
)

// IUserV1 鏄綋鍓嶇敤鎴风敾鍍?鍦板潃绠＄悊鐨?HTTP API 鎺ュ彛銆?
// 娉ㄦ剰锛氳璇?鎺堟潈鐢?middleware/iam-svc 璐熻矗锛屾湰鏈嶅姟浠呭埄鐢ㄤ笂涓嬫枃涓殑 user_id銆?
type IUserV1 interface {
	// GetMyProfile 杩斿洖褰撳墠鐢ㄦ埛鐢诲儚锛屽彲閫夊寘鍚湴鍧€鍒楄〃銆?
	GetMyProfile(ctx context.Context, req *v1.GetMyProfileReq) (res *v1.GetMyProfileRes, err error)
	// UpdateMyProfile 鎸夋洿鏂版帺鐮佽ˉ涓佸綋鍓嶇敤鎴风敾鍍忋€?
	UpdateMyProfile(ctx context.Context, req *v1.UpdateMyProfileReq) (res *v1.UpdateMyProfileRes, err error)

	// ListMyAddresses 鍒楀嚭褰撳墠鐢ㄦ埛鐨勬墍鏈夊湴鍧€銆?
	ListMyAddresses(ctx context.Context, req *v1.ListMyAddressesReq) (res *v1.ListMyAddressesRes, err error)
	// CreateMyAddress 鍒涘缓涓€鏉＄敤鎴峰湴鍧€銆?
	CreateMyAddress(ctx context.Context, req *v1.CreateMyAddressReq) (res *v1.CreateMyAddressRes, err error)
	// ReplaceMyAddress 閫氳繃鏂板缓琛屾浛鎹㈡寚瀹氬湴鍧€銆?
	ReplaceMyAddress(ctx context.Context, req *v1.ReplaceMyAddressReq) (res *v1.ReplaceMyAddressRes, err error)
	// UpdateMyAddress 琛ヤ竵鎸囧畾鍦板潃銆?
	UpdateMyAddress(ctx context.Context, req *v1.UpdateMyAddressReq) (res *v1.UpdateMyAddressRes, err error)
	// DeleteMyAddress 杞垹闄ゆ寚瀹氬湴鍧€銆?
	DeleteMyAddress(ctx context.Context, req *v1.DeleteMyAddressReq) (res *v1.DeleteMyAddressRes, err error)
	// SetMyDefaultAddress 璁剧疆鎴栨竻绌洪粯璁ゅ湴鍧€銆?
	SetMyDefaultAddress(ctx context.Context, req *v1.SetMyDefaultAddressReq) (res *v1.SetMyDefaultAddressRes, err error)
}

