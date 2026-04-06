import { apiClient } from "@/lib/api-client";
import { AnswerSource, AssistantConversation, AssistantMessage, AssistantRun, HandoffTicket } from "./types";

type RawConversation = Record<string, unknown>;
type RawRun = Record<string, unknown>;
type RawMessage = Record<string, unknown>;
type RawSource = Record<string, unknown>;
type RawTicket = Record<string, unknown>;

type RawCreateConversationRes = {
  conversation?: RawConversation;
  latest_run?: RawRun;
  latestRun?: RawRun;
};

type RawSendAssistantMessageRes = {
  conversation?: RawConversation;
  run?: RawRun;
  user_message?: RawMessage;
  userMessage?: RawMessage;
  idempotent_replay?: boolean;
  idempotentReplay?: boolean;
};

type RawGetRunStatusRes = {
  conversation?: RawConversation;
  run?: RawRun;
  latest_assistant_message?: RawMessage;
  latestAssistantMessage?: RawMessage;
  answer_sources?: RawSource[];
  answerSources?: RawSource[];
  pending_message_count?: number | string;
  pendingMessageCount?: number | string;
};

type RawListMessagesRes = {
  list?: RawMessage[];
  next_cursor?: string;
  nextCursor?: string;
  has_more?: boolean;
  hasMore?: boolean;
};

type RawEscalateRes = {
  ticket?: RawTicket;
  conversation?: RawConversation;
};

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  return String(value);
}

function toNumber(value: unknown): number {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : 0;
}

function toBoolean(value: unknown): boolean {
  return Boolean(value);
}

function normalizeConversation(raw?: RawConversation): AssistantConversation {
  return {
    conversationNo: toString(raw?.conversation_no ?? raw?.conversationNo),
    botCode: toString(raw?.bot_code ?? raw?.botCode),
    sceneCode: toString(raw?.scene_code ?? raw?.sceneCode),
    shopNo: toString(raw?.shop_no ?? raw?.shopNo),
    orderNo: toString(raw?.order_no ?? raw?.orderNo),
    anchorSpuNo: toString(raw?.anchor_spu_no ?? raw?.anchorSpuNo),
    anchorSkuNo: toString(raw?.anchor_sku_no ?? raw?.anchorSkuNo),
    conversationStatusCode: toString(raw?.conversation_status_code ?? raw?.conversationStatusCode),
    lastRunStatusCode: toString(raw?.last_run_status_code ?? raw?.lastRunStatusCode),
    isHumanHandover: toBoolean(raw?.is_human_handover ?? raw?.isHumanHandover),
    pendingMessageCount: toNumber(raw?.pending_message_count ?? raw?.pendingMessageCount),
    sessionSummary: toString(raw?.session_summary ?? raw?.sessionSummary)
  };
}

function normalizeRun(raw?: RawRun): AssistantRun {
  return {
    runNo: toString(raw?.run_no ?? raw?.runNo),
    conversationNo: toString(raw?.conversation_no ?? raw?.conversationNo),
    turnNo: toNumber(raw?.turn_no ?? raw?.turnNo),
    runStatusCode: toString(raw?.run_status_code ?? raw?.runStatusCode),
    riskDecisionCode: toString(raw?.risk_decision_code ?? raw?.riskDecisionCode),
    promptInjectionFlag: toBoolean(raw?.prompt_injection_flag ?? raw?.promptInjectionFlag),
    replyInterrupted: toBoolean(raw?.reply_interrupted ?? raw?.replyInterrupted),
    interruptReasonCode: toString(raw?.interrupt_reason_code ?? raw?.interruptReasonCode),
    errorCode: toString(raw?.error_code ?? raw?.errorCode),
    errorMessage: toString(raw?.error_message ?? raw?.errorMessage)
  };
}

function normalizeMessage(raw?: RawMessage): AssistantMessage {
  const assetIds = raw?.asset_ids ?? raw?.assetIds;
  return {
    messageNo: toString(raw?.message_no ?? raw?.messageNo),
    conversationNo: toString(raw?.conversation_no ?? raw?.conversationNo),
    runNo: toString(raw?.run_no ?? raw?.runNo),
    senderTypeCode: toString(raw?.sender_type_code ?? raw?.senderTypeCode),
    messageTypeCode: toString(raw?.message_type_code ?? raw?.messageTypeCode),
    contentText: toString(raw?.content_text ?? raw?.contentText),
    assetIds: Array.isArray(assetIds) ? assetIds.map((item) => toNumber(item)) : [],
    extJson: toString(raw?.ext_json ?? raw?.extJson),
    interrupted: toBoolean(raw?.interrupted),
    sentAt: toString(raw?.sent_at ?? raw?.sentAt)
  };
}

function normalizeSource(raw?: RawSource): AnswerSource {
  return {
    sourceTypeCode: toString(raw?.source_type_code ?? raw?.sourceTypeCode),
    sourceId: toString(raw?.source_id ?? raw?.sourceId),
    sourceVersion: toNumber(raw?.source_version ?? raw?.sourceVersion),
    title: toString(raw?.title),
    snippet: toString(raw?.snippet)
  };
}

function normalizeTicket(raw?: RawTicket): HandoffTicket {
  return {
    ticketNo: toString(raw?.ticket_no ?? raw?.ticketNo),
    conversationNo: toString(raw?.conversation_no ?? raw?.conversationNo),
    escalationReasonCode: toString(raw?.escalation_reason_code ?? raw?.escalationReasonCode),
    statusCode: toString(raw?.status_code ?? raw?.statusCode),
    handoffSummary: toString(raw?.handoff_summary ?? raw?.handoffSummary)
  };
}

export async function createOrGetAssistantConversation(payload: {
  shopNo?: string;
  orderNo?: string;
  anchorSpuNo?: string;
  anchorSkuNo?: string;
  forceNew?: boolean;
}): Promise<{ conversation: AssistantConversation; latestRun?: AssistantRun }> {
  const response = await apiClient.post<RawCreateConversationRes>("/v1/agent/buyer/conversations:get-or-create", {
    scene_code: "BUYER_ASSISTANT",
    bot_code: "BUYER_ASSISTANT",
    shop_no: payload.shopNo ?? "",
    order_no: payload.orderNo ?? "",
    anchor_spu_no: payload.anchorSpuNo ?? "",
    anchor_sku_no: payload.anchorSkuNo ?? "",
    force_new: Boolean(payload.forceNew)
  });
  const latestRun = response?.latest_run ?? response?.latestRun;
  return {
    conversation: normalizeConversation(response?.conversation),
    latestRun: latestRun ? normalizeRun(latestRun) : undefined
  };
}

export async function sendAssistantMessage(payload: {
  conversationNo: string;
  contentText: string;
  clientMessageNo?: string;
}): Promise<{ conversation: AssistantConversation; run: AssistantRun; userMessage: AssistantMessage; idempotentReplay: boolean }> {
  const response = await apiClient.post<RawSendAssistantMessageRes>("/v1/agent/buyer/messages:send", {
    conversation_no: payload.conversationNo,
    client_message_no: payload.clientMessageNo ?? `agent_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
    message_type_code: "TEXT",
    content_text: payload.contentText,
    asset_ids: []
  });
  return {
    conversation: normalizeConversation(response?.conversation),
    run: normalizeRun(response?.run),
    userMessage: normalizeMessage(response?.user_message ?? response?.userMessage),
    idempotentReplay: toBoolean(response?.idempotent_replay ?? response?.idempotentReplay)
  };
}

export async function getAssistantRunStatus(conversationNo: string, runNo = ""): Promise<{ conversation: AssistantConversation; run: AssistantRun; latestAssistantMessage?: AssistantMessage; answerSources: AnswerSource[]; pendingMessageCount: number }> {
  const response = await apiClient.get<RawGetRunStatusRes>(`/v1/agent/buyer/conversations/${encodeURIComponent(conversationNo)}/run-status`, {
    params: { run_no: runNo }
  });
  const latestAssistantMessage = response?.latest_assistant_message ?? response?.latestAssistantMessage;
  return {
    conversation: normalizeConversation(response?.conversation),
    run: normalizeRun(response?.run),
    latestAssistantMessage: latestAssistantMessage ? normalizeMessage(latestAssistantMessage) : undefined,
    answerSources: (response?.answer_sources ?? response?.answerSources ?? []).map((item) => normalizeSource(item)),
    pendingMessageCount: toNumber(response?.pending_message_count ?? response?.pendingMessageCount)
  };
}

export async function listAssistantMessages(conversationNo: string, pageSize = 100, nextCursor = ""): Promise<{ list: AssistantMessage[]; nextCursor: string; hasMore: boolean }> {
  const response = await apiClient.get<RawListMessagesRes>(`/v1/agent/buyer/conversations/${encodeURIComponent(conversationNo)}/messages`, {
    params: { page_size: pageSize, next_cursor: nextCursor }
  });
  return {
    list: (response?.list ?? []).map((item) => normalizeMessage(item)),
    nextCursor: toString(response?.next_cursor ?? response?.nextCursor),
    hasMore: toBoolean(response?.has_more ?? response?.hasMore)
  };
}

export async function escalateAssistantConversation(payload: { conversationNo: string; escalationReasonCode: string; remark?: string }): Promise<{ ticket: HandoffTicket; conversation: AssistantConversation }> {
  const response = await apiClient.post<RawEscalateRes>(`/v1/agent/buyer/conversations/${encodeURIComponent(payload.conversationNo)}:escalate`, {
    escalation_reason_code: payload.escalationReasonCode,
    remark: payload.remark ?? ""
  });
  return {
    ticket: normalizeTicket(response?.ticket),
    conversation: normalizeConversation(response?.conversation)
  };
}

export async function submitAssistantFeedback(payload: { conversationNo: string; runNo: string; feedbackCode: "HELPFUL" | "NOT_HELPFUL"; resolved?: boolean; comment?: string }): Promise<void> {
  await apiClient.post("/v1/agent/buyer/feedback:submit", {
    conversation_no: payload.conversationNo,
    run_no: payload.runNo,
    feedback_code: payload.feedbackCode,
    resolved: Boolean(payload.resolved),
    comment: payload.comment ?? ""
  });
}
