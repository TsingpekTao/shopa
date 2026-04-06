// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RefundTask is the golang structure for table refund_task.
type RefundTask struct {
	Id                     uint64      `json:"id"                     orm:"id"                        ` //
	RefundTaskNo           string      `json:"refundTaskNo"           orm:"refund_task_no"            ` //
	AfterSaleNo            string      `json:"afterSaleNo"            orm:"after_sale_no"             ` //
	OrderNo                string      `json:"orderNo"                orm:"order_no"                  ` //
	SubOrderNo             string      `json:"subOrderNo"             orm:"sub_order_no"              ` //
	PayNo                  string      `json:"payNo"                  orm:"pay_no"                    ` //
	RefundAmount           uint64      `json:"refundAmount"           orm:"refund_amount"             ` //
	PointsReturnAmount     uint64      `json:"pointsReturnAmount"     orm:"points_return_amount"      ` // 鏈??閫??杩旇繕鐨勫凡娑堣垂绉?垎锛屽崟浣嶏細鍒
	PointsReverseAmount    uint64      `json:"pointsReverseAmount"    orm:"points_reverse_amount"     ` // 鏈??閫??鍐插洖鐨勫凡璧犵Н鍒嗭紝鍗曚綅锛氬垎
	PointsCashOffsetAmount uint64      `json:"pointsCashOffsetAmount" orm:"points_cash_offset_amount" ` // 璧犲垎鍐插洖浣欓?涓嶈冻鏃剁殑鐜伴噾鎶垫墸閲戦?锛屽崟浣嶏細鍒
	FinalCashRefundAmount  uint64      `json:"finalCashRefundAmount"  orm:"final_cash_refund_amount"  ` // 瀹為檯鐜伴噾閫??閲戦?锛屽崟浣嶏細鍒
	AccountDebtAfter       uint64      `json:"accountDebtAfter"       orm:"account_debt_after"        ` // 閫??鎵ц?鍚庣Н鍒嗚处鎴锋瑺璐︾粷瀵瑰?锛屽崟浣嶏細鍒
	Status                 uint        `json:"status"                 orm:"status"                    ` //
	RetryCount             uint        `json:"retryCount"             orm:"retry_count"               ` //
	NextRetryAt            *gtime.Time `json:"nextRetryAt"            orm:"next_retry_at"             ` //
	LastErrorCode          string      `json:"lastErrorCode"          orm:"last_error_code"           ` //
	LastErrorMessage       string      `json:"lastErrorMessage"       orm:"last_error_message"        ` //
	Version                uint64      `json:"version"                orm:"version"                   ` //
	CreatedAt              *gtime.Time `json:"createdAt"              orm:"created_at"                ` //
	UpdatedAt              *gtime.Time `json:"updatedAt"              orm:"updated_at"                ` //
	DeletedAt              *gtime.Time `json:"deletedAt"              orm:"deleted_at"                ` //
}
