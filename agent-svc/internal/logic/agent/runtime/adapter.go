package runtime

import (
	"context"
)

type AdapterRequest struct {
	IntentCode   string
	AdapterCode  string
	ToolName     string
	ToolScope    string
	UserQuery    string
	RequestID    string
	Conversation string
}

type AdapterResponse struct {
	Status             string
	Message            string
	DataJSON           string
	DegradedReasonCode string
}

type McpAdapter interface {
	Code() string
	Call(ctx context.Context, req AdapterRequest) (*AdapterResponse, error)
}

type AdapterRegistry struct {
	adapters map[string]McpAdapter
}

func NewAdapterRegistry(adapters ...McpAdapter) *AdapterRegistry {
	// 初始化按业务域聚合的 MCP Adapter 注册表，隔离上层图编排与下游服务接入实现。
	registry := &AdapterRegistry{adapters: make(map[string]McpAdapter)}
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for _, adapter := range adapters {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		if adapter == nil {
			// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
			continue
		}
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		registry.adapters[adapter.Code()] = adapter
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return registry
}

func (r *AdapterRegistry) Call(ctx context.Context, req AdapterRequest) (*AdapterResponse, bool, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if r == nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, false, nil
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	adapter, ok := r.adapters[req.AdapterCode]
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if !ok {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, false, nil
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	resp, err := adapter.Call(ctx, req)
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return resp, true, err
}

type staticAdapter struct {
	code     string
	response func(req AdapterRequest) *AdapterResponse
}

func (a *staticAdapter) Code() string {
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return a.code
}

func (a *staticAdapter) Call(ctx context.Context, req AdapterRequest) (*AdapterResponse, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_ = ctx
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return a.response(req), nil
}

func newDefaultAdapterRegistry() *AdapterRegistry {
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return NewAdapterRegistry(
		&staticAdapter{
			code: "assistant-order-mcp",
			response: func(req AdapterRequest) *AdapterResponse {
				// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
				return &AdapterResponse{
					Status:             toolResultDegraded,
					Message:            "订单相关系统当前繁忙，我先给你流程建议；如果需要人工继续跟进，也可以直接发起转人工。",
					DegradedReasonCode: "DOWNSTREAM_TEMPORARY_UNAVAILABLE",
				}
			},
		},
		&staticAdapter{
			code: "assistant-account-mcp",
			response: func(req AdapterRequest) *AdapterResponse {
				// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
				return &AdapterResponse{
					Status:             toolResultDegraded,
					Message:            "账号相关系统当前繁忙，我先给你通用说明；如果需要进一步核验，可以稍后再试或转人工。",
					DegradedReasonCode: "DOWNSTREAM_TEMPORARY_UNAVAILABLE",
				}
			},
		},
		&staticAdapter{
			code: "assistant-policy-mcp",
			response: func(req AdapterRequest) *AdapterResponse {
				// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
				return &AdapterResponse{
					Status:   toolResultSuccess,
					Message:  "已命中平台规则知识，可继续按规则答复。",
					DataJSON: `{"scope":"policy","confidence":"medium"}`,
				}
			},
		},
		&staticAdapter{
			code: "assistant-catalog-mcp",
			response: func(req AdapterRequest) *AdapterResponse {
				// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
				return &AdapterResponse{
					Status:   toolResultSuccess,
					Message:  "已命中商品导购上下文，可继续补充预算、品牌或规格。",
					DataJSON: `{"scope":"catalog","confidence":"medium"}`,
				}
			},
		},
		&staticAdapter{
			code: "assistant-risk-mcp",
			response: func(req AdapterRequest) *AdapterResponse {
				// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
				return &AdapterResponse{
					Status:   toolResultSuccess,
					Message:  "风险校验已完成。",
					DataJSON: `{"scope":"risk","decision":"pass"}`,
				}
			},
		},
	)
}
