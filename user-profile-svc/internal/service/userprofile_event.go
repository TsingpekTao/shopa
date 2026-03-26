package service

import "time"

// RegisterInitEvent 是“注册初始化资料事件”的最小消费载荷。
// 约束如下：EventID 保障幂等，EventVersion 预留给事件演进，OccurredAt 用于排序乱序事件。
type RegisterInitEvent struct {
	EventID         string
	EventVersion    string
	UserID          uint64
	InitDisplayName string
	OccurredAt      time.Time
}
