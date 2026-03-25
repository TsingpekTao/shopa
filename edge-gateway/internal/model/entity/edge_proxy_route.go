// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EdgeProxyRoute is the golang structure for table edge_proxy_route.
type EdgeProxyRoute struct {
	Id                   uint64      `json:"id"                   orm:"id"                     description:"涓婚敭ID"`              // 涓婚敭ID
	RouteCode            string      `json:"routeCode"            orm:"route_code"             description:"璺?敱缂栫爜锛堝叏灞?敮涓?級"`    // 璺?敱缂栫爜锛堝叏灞?敮涓?級
	Method               string      `json:"method"               orm:"method"                 description:"HTTP鏂规硶"`            // HTTP鏂规硶
	PathPattern          string      `json:"pathPattern"          orm:"path_pattern"           description:"缃戝叧璺?敱鍖归厤璺?緞"`       // 缃戝叧璺?敱鍖归厤璺?緞
	UpstreamService      string      `json:"upstreamService"      orm:"upstream_service"       description:"鐩?爣涓婃父鏈嶅姟鏍囪瘑"`       // 鐩?爣涓婃父鏈嶅姟鏍囪瘑
	UpstreamPathTemplate string      `json:"upstreamPathTemplate" orm:"upstream_path_template" description:"鐩?爣涓婃父璺?緞妯℃澘"`       // 鐩?爣涓婃父璺?緞妯℃澘
	AuthRequired         int         `json:"authRequired"         orm:"auth_required"          description:"鏄?惁闇??閴存潈"`          // 鏄?惁闇??閴存潈
	InjectUserContext    int         `json:"injectUserContext"    orm:"inject_user_context"    description:"鏄?惁娉ㄥ叆鐢ㄦ埛涓婁笅鏂囧ご"`    // 鏄?惁娉ㄥ叆鐢ㄦ埛涓婁笅鏂囧ご
	BodyMode             string      `json:"bodyMode"             orm:"body_mode"              description:"璇锋眰浣撳?鐞嗘ā寮"`         // 璇锋眰浣撳?鐞嗘ā寮
	TimeoutMs            uint        `json:"timeoutMs"            orm:"timeout_ms"             description:"涓婃父璋冪敤瓒呮椂鏃堕棿锛堟?绉掞級"` // 涓婃父璋冪敤瓒呮椂鏃堕棿锛堟?绉掞級
	RateLimitRps         uint        `json:"rateLimitRps"         orm:"rate_limit_rps"         description:"姣忕?闄愭祦闃堝?锛?琛ㄧず涓嶉檺娴"` // 姣忕?闄愭祦闃堝?锛?琛ㄧず涓嶉檺娴
	Status               uint        `json:"status"               orm:"status"                 description:"鐘舵?锛?鍚?敤 0鍋滅敤"`      // 鐘舵?锛?鍚?敤 0鍋滅敤
	Remark               string      `json:"remark"               orm:"remark"                 description:"澶囨敞"`                // 澶囨敞
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"             description:"鍒涘缓鏃堕棿"`             // 鍒涘缓鏃堕棿
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"             description:"鏇存柊鏃堕棿"`             // 鏇存柊鏃堕棿
	DeletedAt            *gtime.Time `json:"deletedAt"            orm:"deleted_at"             description:"杞?垹闄ゆ椂闂"`            // 杞?垹闄ゆ椂闂
}
