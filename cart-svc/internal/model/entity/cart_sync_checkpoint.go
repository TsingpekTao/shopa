// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CartSyncCheckpoint is the golang structure for table cart_sync_checkpoint.
type CartSyncCheckpoint struct {
	Id              uint64      `json:"id"              orm:"id"                description:""` //
	UserId          uint64      `json:"userId"          orm:"user_id"           description:""` //
	LastSyncedAt    *gtime.Time `json:"lastSyncedAt"    orm:"last_synced_at"    description:""` //
	LastSyncVersion uint64      `json:"lastSyncVersion" orm:"last_sync_version" description:""` //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"        description:""` //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"        description:""` //
}
