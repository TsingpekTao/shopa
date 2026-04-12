package agent

import (
	"encoding/json"
	"strings"

	agentv1 "github.com/TsingpekTao/shopa/agent-svc/api/v1"
)

type assistantHiddenAction struct {
	Type  string `json:"type,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type sendAssistantMessageExtEnvelope struct {
	HiddenAction *assistantHiddenAction `json:"hidden_action,omitempty"`
}

func extractHiddenAction(req *agentv1.SendAssistantMessageReq) *assistantHiddenAction {
	// 先读取显式协议字段，让新前端可以直接走结构化设槽。
	if req != nil && req.GetHiddenAction() != nil {
		action := normalizeHiddenAction(&assistantHiddenAction{
			Type:  req.GetHiddenAction().GetType(),
			Key:   req.GetHiddenAction().GetKey(),
			Value: req.GetHiddenAction().GetValue(),
		})
		if action != nil {
			return action
		}
	}
	// 再兼容 ext_json 中的 hidden_action，保证旧调用方也能继续工作。
	if req == nil || strings.TrimSpace(req.GetExtJson()) == "" {
		return nil
	}

	return extractHiddenActionFromJSON(req.GetExtJson())
}

func extractHiddenActionFromJSON(raw string) *assistantHiddenAction {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	var envelope sendAssistantMessageExtEnvelope
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		return nil
	}
	return normalizeHiddenAction(envelope.HiddenAction)
}

func mergeHiddenActionIntoExtJSON(raw string, action *assistantHiddenAction) string {
	if action == nil {
		return normalizeJSON(raw)
	}

	var envelope map[string]any
	if strings.TrimSpace(raw) != "" {
		_ = json.Unmarshal([]byte(raw), &envelope)
	}
	if envelope == nil {
		envelope = make(map[string]any)
	}
	envelope["hidden_action"] = action

	encoded, err := json.Marshal(envelope)
	if err != nil {
		return normalizeJSON(raw)
	}
	return string(encoded)
}

func normalizeHiddenAction(action *assistantHiddenAction) *assistantHiddenAction {
	if action == nil {
		return nil
	}

	normalized := &assistantHiddenAction{
		Type:  strings.ToUpper(strings.TrimSpace(action.Type)),
		Key:   strings.ToLower(strings.TrimSpace(action.Key)),
		Value: strings.TrimSpace(action.Value),
	}
	if normalized.Type != "SET_SLOT" {
		return nil
	}
	if normalized.Key != "selected_order_no" && normalized.Key != "selected_sub_order_no" {
		return nil
	}
	if normalized.Value == "" {
		return nil
	}
	return normalized
}
