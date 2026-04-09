package runtime

import "encoding/json"

type graphStateSnapshot struct {
	ReplyPayload ReplyPayload     `json:"reply_payload"`
	TaskSession  TaskSessionState `json:"task_session"`
}

func ParseGraphReplyPayload(raw string) *ReplyPayload {
	// 先判空，避免对空 graph_state_json 做无意义反序列化。
	if raw == "" {
		return nil
	}
	// 反序列化运行时快照，提取结构化回复载荷。
	var snapshot graphStateSnapshot
	if err := json.Unmarshal([]byte(raw), &snapshot); err != nil {
		return nil
	}
	// 如果载荷里没有有效 reply_text 和 intent，就视为未生成结构化回复。
	if snapshot.ReplyPayload.ReplyText == "" && snapshot.ReplyPayload.IntentCode == "" {
		return nil
	}
	// 返回解析成功的结构化回复，供协议映射和状态查询复用。
	return &snapshot.ReplyPayload
}

func ParseGraphTaskSession(raw string) *TaskSessionState {
	// 先判空，避免对空 graph_state_json 做无意义反序列化。
	if raw == "" {
		return nil
	}
	// 反序列化运行时快照，提取任务会话状态。
	var snapshot graphStateSnapshot
	if err := json.Unmarshal([]byte(raw), &snapshot); err != nil {
		return nil
	}
	// 如果会话里既没有任务码也没有槽位信息，就视为无可复用状态。
	if snapshot.TaskSession.ActiveTaskCode == "" && len(snapshot.TaskSession.SlotValues) == 0 && len(snapshot.TaskSession.MissingSlots) == 0 {
		return nil
	}
	// 返回解析成功的任务会话，供下一轮运行继续复用。
	return &snapshot.TaskSession
}
