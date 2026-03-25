// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserProfile is the golang structure of table user_profile for DAO operations like Where/Data.
type UserProfile struct {
	g.Meta             `orm:"table:user_profile, do:true"`
	UserId             any         // User ID from iam-svc
	DisplayName        any         // Display name
	AvatarAssetId      any         // Avatar asset id from media-svc
	AvatarUrl          any         // Avatar URL cache
	Gender             any         // 0=unspecified,1=male,2=female,3=other
	Birthday           *gtime.Time // Birthday date
	Locale             any         // Locale
	Timezone           any         // Timezone
	MarketingOptIn     any         // Marketing opt in
	DefaultAddressId   any         // Deprecated compatibility field
	ProfileVersion     any         // Optimistic version of profile aggregate
	AddressBookVersion any         // CAS version of address-book aggregate
	DisplayNameSource  any         // 1 SYSTEM_INIT,2 USER_SET
	LastInitEventAt    *gtime.Time // Last register-init event time
	LastInitEventId    any         // Last register-init event id
	Ext                any         // Extension payload
	CreatedAt          *gtime.Time //
	UpdatedAt          *gtime.Time //
}
