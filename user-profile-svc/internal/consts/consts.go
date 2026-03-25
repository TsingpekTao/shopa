package consts

const (
	// user_address 表的地址状态值。
	AddressStatusActive   = 1
	AddressStatusDeleted  = 2
	AddressStatusReplaced = 3
)

const (
	// 列表查询的默认分页值。
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

const (
	// user_profile 表的展示名来源值。
	DisplayNameSourceSystemInit = 1
	DisplayNameSourceUserSet    = 2
)

const (
	// register init 事件的消费者名称。
	ConsumerRegisterInit = "user-profile.register-init.v1"
)
