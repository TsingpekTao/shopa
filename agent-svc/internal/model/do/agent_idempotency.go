// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentIdempotency is the golang structure of table agent_idempotency for DAO operations like Where/Data.
type AgentIdempotency struct {
	g.Meta         `orm:"table:agent_idempotency, do:true"`
	Id             any         //
	IdempotencyKey any         //
	ScopeCode      any         //
	ScopeId        any         //
	ResponseJson   any         //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
}
