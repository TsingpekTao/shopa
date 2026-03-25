// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EdgeRequestAudit is the golang structure for table edge_request_audit.
type EdgeRequestAudit struct {
	Id                 uint64      `json:"id"                 orm:"id"                   description:"涓婚敭ID"`          // 涓婚敭ID
	RequestId          string      `json:"requestId"          orm:"request_id"           description:"璇锋眰ID锛堝叏灞?敮涓?級"` // 璇锋眰ID锛堝叏灞?敮涓?級
	RouteCode          string      `json:"routeCode"          orm:"route_code"           description:"鍛戒腑璺?敱缂栫爜"`      // 鍛戒腑璺?敱缂栫爜
	UserId             uint64      `json:"userId"             orm:"user_id"              description:"璋冪敤鐢ㄦ埛ID"`       // 璋冪敤鐢ㄦ埛ID
	Method             string      `json:"method"             orm:"method"               description:"HTTP鏂规硶"`        // HTTP鏂规硶
	Path               string      `json:"path"               orm:"path"                 description:"璇锋眰璺?緞"`         // 璇锋眰璺?緞
	UpstreamService    string      `json:"upstreamService"    orm:"upstream_service"     description:"涓婃父鏈嶅姟鏍囪瘑"`      // 涓婃父鏈嶅姟鏍囪瘑
	StatusCode         int         `json:"statusCode"         orm:"status_code"          description:"HTTP鐘舵?鐮"`       // HTTP鐘舵?鐮
	LatencyMs          uint        `json:"latencyMs"          orm:"latency_ms"           description:"璇锋眰鑰楁椂姣??"`      // 璇锋眰鑰楁椂姣??
	Partial            int         `json:"partial"            orm:"partial"              description:"鏄?惁闄嶇骇杩斿洖"`      // 鏄?惁闄嶇骇杩斿洖
	DegradedFieldsJson string      `json:"degradedFieldsJson" orm:"degraded_fields_json" description:"闄嶇骇瀛楁?鍒楄〃"`      // 闄嶇骇瀛楁?鍒楄〃
	ErrorCode          string      `json:"errorCode"          orm:"error_code"           description:"涓氬姟閿欒?鐮"`        // 涓氬姟閿欒?鐮
	ClientIp           string      `json:"clientIp"           orm:"client_ip"            description:"瀹㈡埛绔疘P"`         // 瀹㈡埛绔疘P
	UserAgent          string      `json:"userAgent"          orm:"user_agent"           description:"瀹㈡埛绔疷A"`         // 瀹㈡埛绔疷A
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           description:"鍒涘缓鏃堕棿"`         // 鍒涘缓鏃堕棿
}
