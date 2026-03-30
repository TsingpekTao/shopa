// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RiskHitLog is the golang structure of table risk_hit_log for DAO operations like Where/Data.
type RiskHitLog struct {
	g.Meta    `orm:"table:risk_hit_log, do:true"`
	Id        any         //
	HitNo     any         //
	UserId    any         //
	RuleCode  any         //
	BizType   any         //
	BizNo     any         //
	Decision  any         //
	CreatedAt *gtime.Time //
}
