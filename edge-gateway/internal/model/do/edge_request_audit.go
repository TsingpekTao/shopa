// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EdgeRequestAudit is the golang structure of table edge_request_audit for DAO operations like Where/Data.
type EdgeRequestAudit struct {
	g.Meta             `orm:"table:edge_request_audit, do:true"`
	Id                 any         // 涓婚敭ID
	RequestId          any         // 璇锋眰ID锛堝叏灞?敮涓?級
	RouteCode          any         // 鍛戒腑璺?敱缂栫爜
	UserId             any         // 璋冪敤鐢ㄦ埛ID
	Method             any         // HTTP鏂规硶
	Path               any         // 璇锋眰璺?緞
	UpstreamService    any         // 涓婃父鏈嶅姟鏍囪瘑
	StatusCode         any         // HTTP鐘舵?鐮
	LatencyMs          any         // 璇锋眰鑰楁椂姣??
	Partial            any         // 鏄?惁闄嶇骇杩斿洖
	DegradedFieldsJson any         // 闄嶇骇瀛楁?鍒楄〃
	ErrorCode          any         // 涓氬姟閿欒?鐮
	ClientIp           any         // 瀹㈡埛绔疘P
	UserAgent          any         // 瀹㈡埛绔疷A
	CreatedAt          *gtime.Time // 鍒涘缓鏃堕棿
}
