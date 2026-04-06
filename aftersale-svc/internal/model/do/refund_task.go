// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RefundTask is the golang structure of table refund_task for DAO operations like Where/Data.
type RefundTask struct {
	g.Meta                 `orm:"table:refund_task, do:true"`
	Id                     any         //
	RefundTaskNo           any         //
	AfterSaleNo            any         //
	OrderNo                any         //
	SubOrderNo             any         //
	PayNo                  any         //
	RefundAmount           any         //
	PointsReturnAmount     any         // 鏈??閫??杩旇繕鐨勫凡娑堣垂绉?垎锛屽崟浣嶏細鍒
	PointsReverseAmount    any         // 鏈??閫??鍐插洖鐨勫凡璧犵Н鍒嗭紝鍗曚綅锛氬垎
	PointsCashOffsetAmount any         // 璧犲垎鍐插洖浣欓?涓嶈冻鏃剁殑鐜伴噾鎶垫墸閲戦?锛屽崟浣嶏細鍒
	FinalCashRefundAmount  any         // 瀹為檯鐜伴噾閫??閲戦?锛屽崟浣嶏細鍒
	AccountDebtAfter       any         // 閫??鎵ц?鍚庣Н鍒嗚处鎴锋瑺璐︾粷瀵瑰?锛屽崟浣嶏細鍒
	Status                 any         //
	RetryCount             any         //
	NextRetryAt            *gtime.Time //
	LastErrorCode          any         //
	LastErrorMessage       any         //
	Version                any         //
	CreatedAt              *gtime.Time //
	UpdatedAt              *gtime.Time //
	DeletedAt              *gtime.Time //
}
