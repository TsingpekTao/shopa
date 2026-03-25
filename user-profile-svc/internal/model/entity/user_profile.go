// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserProfile is the golang structure for table user_profile.
type UserProfile struct {
	UserId             uint64      `json:"userId"             orm:"user_id"              description:"User ID from iam-svc"`                    // User ID from iam-svc
	DisplayName        string      `json:"displayName"        orm:"display_name"         description:"Display name"`                            // Display name
	AvatarAssetId      uint64      `json:"avatarAssetId"      orm:"avatar_asset_id"      description:"Avatar asset id from media-svc"`          // Avatar asset id from media-svc
	AvatarUrl          string      `json:"avatarUrl"          orm:"avatar_url"           description:"Avatar URL cache"`                        // Avatar URL cache
	Gender             uint        `json:"gender"             orm:"gender"               description:"0=unspecified,1=male,2=female,3=other"`   // 0=unspecified,1=male,2=female,3=other
	Birthday           *gtime.Time `json:"birthday"           orm:"birthday"             description:"Birthday date"`                           // Birthday date
	Locale             string      `json:"locale"             orm:"locale"               description:"Locale"`                                  // Locale
	Timezone           string      `json:"timezone"           orm:"timezone"             description:"Timezone"`                                // Timezone
	MarketingOptIn     int         `json:"marketingOptIn"     orm:"marketing_opt_in"     description:"Marketing opt in"`                        // Marketing opt in
	DefaultAddressId   uint64      `json:"defaultAddressId"   orm:"default_address_id"   description:"Deprecated compatibility field"`          // Deprecated compatibility field
	ProfileVersion     uint64      `json:"profileVersion"     orm:"profile_version"      description:"Optimistic version of profile aggregate"` // Optimistic version of profile aggregate
	AddressBookVersion uint64      `json:"addressBookVersion" orm:"address_book_version" description:"CAS version of address-book aggregate"`   // CAS version of address-book aggregate
	DisplayNameSource  uint        `json:"displayNameSource"  orm:"display_name_source"  description:"1 SYSTEM_INIT,2 USER_SET"`                // 1 SYSTEM_INIT,2 USER_SET
	LastInitEventAt    *gtime.Time `json:"lastInitEventAt"    orm:"last_init_event_at"   description:"Last register-init event time"`           // Last register-init event time
	LastInitEventId    string      `json:"lastInitEventId"    orm:"last_init_event_id"   description:"Last register-init event id"`             // Last register-init event id
	Ext                string      `json:"ext"                orm:"ext"                  description:"Extension payload"`                       // Extension payload
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           description:""`                                        //
	UpdatedAt          *gtime.Time `json:"updatedAt"          orm:"updated_at"           description:""`                                        //
}
