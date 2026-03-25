package consts

const (
	// Address status values in table user_address.
	AddressStatusActive   = 1
	AddressStatusDeleted  = 2
	AddressStatusReplaced = 3
)

const (
	// Default pagination values for list queries.
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

const (
	// Display name source values in table user_profile.
	DisplayNameSourceSystemInit = 1
	DisplayNameSourceUserSet    = 2
)

const (
	// Consumer name for register init events.
	ConsumerRegisterInit = "user-profile.register-init.v1"
)
