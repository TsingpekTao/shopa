// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CartSyncCheckpoint is the golang structure of table cart_sync_checkpoint for DAO operations like Where/Data.
type CartSyncCheckpoint struct {
	g.Meta          `orm:"table:cart_sync_checkpoint, do:true"`
	Id              any         //
	UserId          any         //
	LastSyncedAt    *gtime.Time //
	LastSyncVersion any         //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}
