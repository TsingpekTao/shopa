package service

import "time"

// RegisterInitEvent is the minimal payload for register-init event handling.
type RegisterInitEvent struct {
	EventID         string
	EventVersion    string
	UserID          uint64
	InitDisplayName string
	OccurredAt      time.Time
}
