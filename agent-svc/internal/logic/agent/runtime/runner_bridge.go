package runtime

import "context"

func (r *einoRunner) assistantIntentRouter(ctx context.Context, state *runState) (*runState, error) {
	// 先标记当前节点，保证上层能拿到最新图进度。
	state.CurrentNodeCode = nodeIntentRouter
	// 继续兼容已有的安全拒绝逻辑，避免新流程放松防护边界。
	if containsAny(state.NormalizedQuery, "ignore previous", "system prompt", "忽略之前", "免费拿", "卡 bug", "赔付保证") {
		state.IntentCode = "SAFETY_REFUSAL"
		state.PromptInjection = true
		state.RiskDecisionCode = "REFUSE"
		return state, state.emitCheckpoint(ctx)
	}
	// 对当前用户消息执行 Guard 分类，把结果写回运行时状态。
	guard := detectGuardIntent(state.Message)
	state.IntentCode = guard.IntentCode
	state.GuardResultCode = guard.GuardResultCode
	// 新的售后助手主链路不依赖旧 knowledge/tool 开关，因此这里统一关闭旧分支。
	state.NeedKnowledge = false
	state.NeedTool = false
	// 返回新的 Guard 状态，供后续节点继续消费。
	return state, state.emitCheckpoint(ctx)
}

func (r *einoRunner) assistantGenerateAnswer(ctx context.Context, state *runState) (*runState, error) {
	// 先标记当前生成节点，并切到生成中状态。
	state.CurrentNodeCode = nodeGenerateAnswer
	state.RunStatusCode = runStatusGenerating
	// 在正式回复前先打检查点，确保运行中状态可持久化。
	if err := state.emitCheckpoint(ctx); err != nil {
		return nil, err
	}
	// 如果命中安全拒绝，就直接输出安全回复并结束。
	if state.PromptInjection {
		state.AnswerText = "这类请求我不能协助。我可以继续帮你处理订单、物流和售后问题。"
		state.ReplyPayload = ReplyPayload{
			ReplyText:       state.AnswerText,
			IntentCode:      state.IntentCode,
			GuardResultCode: state.GuardResultCode,
			Confidence:      0.99,
		}
		state.RunStatusCode = runStatusSuccess
		return state, state.emitCheckpoint(ctx)
	}
	// 走新的售后任务规划逻辑，把多轮任务会话和结构化回复一次性算出来。
	planned, err := planAfterSaleTurn(ctx, AfterSaleTurnInput{
		Message:         state.Message,
		Conversation:    state.Conversation,
		PreviousSession: state.TaskSession,
		Security:        state.Security,
		OrderRepository: r.orderRepository,
	})
	// 如果新的任务规划失败，就继续向上抛出统一错误。
	if err != nil {
		return nil, err
	}
	// 把规划后的会话状态回写到运行时，供 graph_state_json 和上层状态查询复用。
	state.TaskSession = planned.Session
	// 把结构化回复载荷写回运行时，供 SendAssistantMessage/GetAssistantRunStatus 返回。
	state.ReplyPayload = planned.Reply
	// 让历史字符串字段继续与结构化回复保持一致，兼容旧调用点。
	state.AnswerText = planned.Reply.ReplyText
	state.IntentCode = planned.Reply.IntentCode
	state.GuardResultCode = planned.Reply.GuardResultCode
	// 如果本轮建议转人工，就把运行状态标记为正式升级出口。
	if planned.Reply.HandoffRecommended {
		state.RunStatusCode = runStatusEscalated
	} else {
		// 其余正常完成的情况都落为成功状态。
		state.RunStatusCode = runStatusSuccess
	}
	// 返回带结构化 payload 的新状态，并落最终检查点。
	return state, state.emitCheckpoint(ctx)
}
