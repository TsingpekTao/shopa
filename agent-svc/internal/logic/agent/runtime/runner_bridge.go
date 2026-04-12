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
	replyPayload, taskSession, err := r.generateOfficialReply(ctx, state)
	if err != nil {
		return nil, err
	}
	taskSession = finalizeTaskSessionState(state.TaskSession, taskSession, state.Security.ConversationNo)
	state.TaskSession = taskSession
	state.ReplyPayload = *replyPayload
	state.AnswerText = replyPayload.ReplyText
	state.IntentCode = replyPayload.IntentCode
	state.GuardResultCode = replyPayload.GuardResultCode
	if replyPayload.HandoffRecommended {
		state.RunStatusCode = runStatusEscalated
	} else {
		state.RunStatusCode = runStatusSuccess
	}
	return state, state.emitCheckpoint(ctx)
}
