package agent

import (
	agentapi "github.com/TsingpekTao/shopa/agent-svc/api/agent"
	"github.com/TsingpekTao/shopa/agent-svc/internal/service"
)

type ControllerV1 struct {
	svc service.IAgent
}

// NewV1 创建 agent HTTP 控制器。
func NewV1() agentapi.IAgentV1 {
	return &ControllerV1{svc: service.Agent()}
}
