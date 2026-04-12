package runtime

import (
	"encoding/json"
	"fmt"
	"strings"
)

const maxPlannerRepairAttempts = 1

type rawPlannerDecision struct {
	Thought            string               `json:"thought,omitempty"`
	Intent             string               `json:"intent,omitempty"`
	UserGoal           string               `json:"user_goal,omitempty"`
	Confidence         float64              `json:"confidence,omitempty"`
	SlotValues         map[string]any       `json:"slot_values,omitempty"`
	MissingSlots       []string             `json:"missing_slots,omitempty"`
	NextAction         rawPlannerNextAction `json:"next_action,omitempty"`
	HandoffRecommended bool                 `json:"handoff_recommended,omitempty"`
	HandoffReasonCode  string               `json:"handoff_reason_code,omitempty"`
}

type rawPlannerNextAction struct {
	Type     string         `json:"type,omitempty"`
	ToolName string         `json:"tool_name,omitempty"`
	ToolArgs map[string]any `json:"tool_args,omitempty"`
	Reason   string         `json:"reason,omitempty"`
}

type plannerValidationError struct {
	ErrorCode    string
	ErrorMessage string
	FieldErrors  []string
	Retryable    bool
}

func decodePlannerDecision(raw string) (PlannerDecision, *plannerValidationError) {
	var decoded rawPlannerDecision
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &decoded); err != nil {
		return PlannerDecision{}, &plannerValidationError{
			ErrorCode:    "INVALID_JSON",
			ErrorMessage: err.Error(),
			Retryable:    true,
		}
	}

	decision := PlannerDecision{
		Thought:            decoded.Thought,
		Intent:             decoded.Intent,
		UserGoal:           decoded.UserGoal,
		Confidence:         decoded.Confidence,
		MissingSlots:       append([]string(nil), decoded.MissingSlots...),
		HandoffRecommended: decoded.HandoffRecommended,
		HandoffReasonCode:  decoded.HandoffReasonCode,
		SlotValues:         make(map[string]string),
		NextAction: PlannerNextAction{
			Type:     decoded.NextAction.Type,
			ToolName: decoded.NextAction.ToolName,
			ToolArgs: make(map[string]string),
			Reason:   decoded.NextAction.Reason,
		},
	}

	for key, value := range decoded.SlotValues {
		text, ok := plannerArgumentAsString(value)
		if !ok {
			return PlannerDecision{}, &plannerValidationError{
				ErrorCode:    "INVALID_SLOT_VALUE",
				ErrorMessage: fmt.Sprintf("slot_values.%s must be a string", key),
				FieldErrors:  []string{fmt.Sprintf("slot_values.%s", key)},
				Retryable:    true,
			}
		}
		decision.SlotValues[key] = text
	}
	for key, value := range decoded.NextAction.ToolArgs {
		text, ok := plannerArgumentAsString(value)
		if !ok {
			return PlannerDecision{}, &plannerValidationError{
				ErrorCode:    "INVALID_TOOL_ARGUMENTS",
				ErrorMessage: fmt.Sprintf("tool_args.%s must be a string", key),
				FieldErrors:  []string{fmt.Sprintf("tool_args.%s", key)},
				Retryable:    true,
			}
		}
		decision.NextAction.ToolArgs[key] = text
	}
	if validationErr := validatePlannerDecision(decision); validationErr != nil {
		return PlannerDecision{}, validationErr
	}
	return normalizePlannerDecision(decision), nil
}

func plannerArgumentAsString(value any) (string, bool) {
	switch typed := value.(type) {
	case nil:
		return "", true
	case string:
		return strings.TrimSpace(typed), true
	default:
		return "", false
	}
}

func validatePlannerDecision(decision PlannerDecision) *plannerValidationError {
	decision = normalizePlannerDecision(decision)
	if decision.Intent == "" {
		return &plannerValidationError{
			ErrorCode:    "MISSING_INTENT",
			ErrorMessage: "intent is required",
			FieldErrors:  []string{"intent"},
			Retryable:    true,
		}
	}
	switch decision.NextAction.Type {
	case officialActionToolCall:
		if !isAllowedPlannerTool(decision.NextAction.ToolName) {
			return &plannerValidationError{
				ErrorCode:    "UNKNOWN_TOOL",
				ErrorMessage: fmt.Sprintf("tool %q is not allowed", decision.NextAction.ToolName),
				FieldErrors:  []string{"next_action.tool_name"},
				Retryable:    false,
			}
		}
		if fieldErrors := validatePlannerToolArguments(decision.NextAction.ToolName, decision.NextAction.ToolArgs); len(fieldErrors) > 0 {
			return &plannerValidationError{
				ErrorCode:    "INVALID_TOOL_ARGUMENTS",
				ErrorMessage: "tool arguments are invalid",
				FieldErrors:  fieldErrors,
				Retryable:    true,
			}
		}
	case officialActionClarification, officialActionDirectReply:
		return nil
	default:
		return &plannerValidationError{
			ErrorCode:    "INVALID_ACTION",
			ErrorMessage: fmt.Sprintf("next_action.type %q is invalid", decision.NextAction.Type),
			FieldErrors:  []string{"next_action.type"},
			Retryable:    true,
		}
	}
	return nil
}

func validatePlannerToolArguments(toolName string, toolArgs map[string]string) []string {
	trimmedToolName := strings.TrimSpace(toolName)
	trimmedArgs := make(map[string]string, len(toolArgs))
	for key, value := range toolArgs {
		trimmedArgs[key] = strings.TrimSpace(value)
	}

	switch trimmedToolName {
	case "search_knowledge_chunks", "search_products", "batch_get_spu_cards":
		if trimmedArgs["query"] == "" {
			return []string{"tool_args.query"}
		}
	case "get_order_snapshot", "query_logistics":
		if orderNo := trimmedArgs["order_no"]; orderNo != "" && !strings.HasPrefix(strings.ToUpper(orderNo), "ORD") {
			return []string{"tool_args.order_no"}
		}
		if subOrderNo := trimmedArgs["sub_order_no"]; subOrderNo != "" && !strings.HasPrefix(strings.ToUpper(subOrderNo), "SUB") {
			return []string{"tool_args.sub_order_no"}
		}
	}
	return nil
}

func buildPlannerRepairUserPrompt(originalPrompt string, rawResponse string, validationErr *plannerValidationError) string {
	if validationErr == nil {
		return strings.TrimSpace(originalPrompt)
	}

	var builder strings.Builder
	builder.WriteString(strings.TrimSpace(originalPrompt))
	builder.WriteString("\n\n")
	builder.WriteString("Your previous planner JSON was rejected.\n")
	builder.WriteString("error_code: ")
	builder.WriteString(validationErr.ErrorCode)
	builder.WriteString("\nerror_message: ")
	builder.WriteString(validationErr.ErrorMessage)
	if len(validationErr.FieldErrors) > 0 {
		builder.WriteString("\nfield_errors: ")
		builder.WriteString(strings.Join(validationErr.FieldErrors, ", "))
	}
	builder.WriteString("\nprevious_response: ")
	builder.WriteString(strings.TrimSpace(rawResponse))
	builder.WriteString("\nReturn corrected JSON only. Keep the same tool intent unless the error makes it impossible.")
	return builder.String()
}
