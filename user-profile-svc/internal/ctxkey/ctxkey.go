package ctxkey

// UserIDKey 是用于存储认证用户 ID 的上下文键类型。
// 使用专属类型可避免跨包键名冲突。
type UserIDKey struct{}
