// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentIdempotency is the golang structure for table agent_idempotency.
type AgentIdempotency struct {
	Id             uint64      `json:"id"             orm:"id"              ` //
	IdempotencyKey string      `json:"idempotencyKey" orm:"idempotency_key" ` //
	ScopeCode      string      `json:"scopeCode"      orm:"scope_code"      ` //
	ScopeId        string      `json:"scopeId"        orm:"scope_id"        ` //
	ResponseJson   string      `json:"responseJson"   orm:"response_json"   ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      ` //
}
