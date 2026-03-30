package consts

const (
	CartKeyPrefix                 = "cart:user:"
	CheckoutTokenKeyPrefix        = "checkout:token:"
	SkuProjectionKeyPrefix        = "cart:sku_projection:"
	DirtyUsersZSetKey             = "cart:dirty_users"
	CartKeyTTLSeconds       int64 = 30 * 24 * 60 * 60 // 30d
	CheckoutTokenTTLSeconds       = 5 * 60            // 5m
)

const (
	CartItemStatusActive  uint = 1
	CartItemStatusInvalid uint = 2
)
