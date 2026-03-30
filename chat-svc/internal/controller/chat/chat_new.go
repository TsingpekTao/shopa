package chat

import (
	chatapi "github.com/TsingpekTao/shopa/chat-svc/api/chat"
	"github.com/TsingpekTao/shopa/chat-svc/internal/service"
)

type ControllerV1 struct {
	chat service.IChat
}

func NewV1() chatapi.IChatV1 {
	return &ControllerV1{chat: service.Chat()}
}
