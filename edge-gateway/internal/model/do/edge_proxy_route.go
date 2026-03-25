// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EdgeProxyRoute 是 edge_proxy_route 表用于 DAO 操作（如 Where/Data）的 Go 结构体。
type EdgeProxyRoute struct {
	g.Meta               `orm:"table:edge_proxy_route, do:true"`
	Id                   any         // 涓婚敭ID
	RouteCode            any         // 璺?敱缂栫爜锛堝叏灞?敮涓?級
	Method               any         // HTTP鏂规硶
	PathPattern          any         // 缃戝叧璺?敱鍖归厤璺?緞
	UpstreamService      any         // 鐩?爣涓婃父鏈嶅姟鏍囪瘑
	UpstreamPathTemplate any         // 鐩?爣涓婃父璺?緞妯℃澘
	AuthRequired         any         // 鏄?惁闇??閴存潈
	InjectUserContext    any         // 鏄?惁娉ㄥ叆鐢ㄦ埛涓婁笅鏂囧ご
	BodyMode             any         // 璇锋眰浣撳?鐞嗘ā寮
	TimeoutMs            any         // 涓婃父璋冪敤瓒呮椂鏃堕棿锛堟?绉掞級
	RateLimitRps         any         // 姣忕?闄愭祦闃堝?锛?琛ㄧず涓嶉檺娴
	Status               any         // 鐘舵?锛?鍚?敤 0鍋滅敤
	Remark               any         // 澶囨敞
	CreatedAt            *gtime.Time // 鍒涘缓鏃堕棿
	UpdatedAt            *gtime.Time // 鏇存柊鏃堕棿
	DeletedAt            *gtime.Time // 杞?垹闄ゆ椂闂
}
