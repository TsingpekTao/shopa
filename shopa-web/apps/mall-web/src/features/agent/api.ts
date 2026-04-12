import { apiClient } from "@/lib/api-client";
import {
  AnswerSource,
  AssistantConversation,
  AssistantDataCard,
  AssistantHiddenAction,
  AssistantMessage,
  AssistantMissingSlot,
  AssistantReplyPayload,
  AssistantRun,
  AssistantSuggestedAction,
  HandoffTicket,
} from "./types";

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

function toText(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  if (typeof value === "string") {
    return value;
  }
  if (
    typeof value === "number" ||
    typeof value === "boolean" ||
    typeof value === "bigint"
  ) {
    return String(value);
  }
  if (Array.isArray(value)) {
    return value
      .map((item) => toText(item))
      .filter((item) => item.trim().length > 0)
      .join(" ")
      .trim();
  }
  return "";
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
    conversationNo: toText(raw?.conversation_no ?? raw?.conversationNo),
    botCode: toText(raw?.bot_code ?? raw?.botCode),
    sceneCode: toText(raw?.scene_code ?? raw?.sceneCode),
    shopNo: toText(raw?.shop_no ?? raw?.shopNo),
    orderNo: toText(raw?.order_no ?? raw?.orderNo),
    anchorSpuNo: toText(raw?.anchor_spu_no ?? raw?.anchorSpuNo),
    anchorSkuNo: toText(raw?.anchor_sku_no ?? raw?.anchorSkuNo),
    conversationStatusCode: toText(
      raw?.conversation_status_code ?? raw?.conversationStatusCode,
    ),
    lastRunStatusCode: toText(
      raw?.last_run_status_code ?? raw?.lastRunStatusCode,
    ),
    isHumanHandover: toBoolean(raw?.is_human_handover ?? raw?.isHumanHandover),
    pendingMessageCount: toNumber(
      raw?.pending_message_count ?? raw?.pendingMessageCount,
    ),
    sessionSummary: toText(raw?.session_summary ?? raw?.sessionSummary),
  };
}

function normalizeRun(raw?: RawRun): AssistantRun {
  return {
    runNo: toText(raw?.run_no ?? raw?.runNo),
    conversationNo: toText(raw?.conversation_no ?? raw?.conversationNo),
    turnNo: toNumber(raw?.turn_no ?? raw?.turnNo),
    runStatusCode: toText(raw?.run_status_code ?? raw?.runStatusCode),
    riskDecisionCode: toText(raw?.risk_decision_code ?? raw?.riskDecisionCode),
    promptInjectionFlag: toBoolean(
      raw?.prompt_injection_flag ?? raw?.promptInjectionFlag,
    ),
    replyInterrupted: toBoolean(
      raw?.reply_interrupted ?? raw?.replyInterrupted,
    ),
    interruptReasonCode: toText(
      raw?.interrupt_reason_code ?? raw?.interruptReasonCode,
    ),
    errorCode: toText(raw?.error_code ?? raw?.errorCode),
    errorMessage: toText(raw?.error_message ?? raw?.errorMessage),
    replyPayload: normalizeReplyPayload(
      (raw?.reply_payload ?? raw?.replyPayload) as RawConversation | undefined,
    ),
  };
}

function normalizeMessage(raw?: RawMessage): AssistantMessage {
  const assetIds = raw?.asset_ids ?? raw?.assetIds;
  return {
    messageNo: toText(raw?.message_no ?? raw?.messageNo),
    conversationNo: toText(raw?.conversation_no ?? raw?.conversationNo),
    runNo: toText(raw?.run_no ?? raw?.runNo),
    senderTypeCode: toText(raw?.sender_type_code ?? raw?.senderTypeCode),
    messageTypeCode: toText(raw?.message_type_code ?? raw?.messageTypeCode),
    contentText: toText(raw?.content_text ?? raw?.contentText),
    assetIds: Array.isArray(assetIds)
      ? assetIds.map((item) => toNumber(item))
      : [],
    extJson: toText(raw?.ext_json ?? raw?.extJson),
    interrupted: toBoolean(raw?.interrupted),
    sentAt: toText(raw?.sent_at ?? raw?.sentAt),
  };
}

function normalizeSource(raw?: RawSource): AnswerSource {
  return {
    sourceTypeCode: toText(raw?.source_type_code ?? raw?.sourceTypeCode),
    sourceId: toText(raw?.source_id ?? raw?.sourceId),
    sourceVersion: toNumber(raw?.source_version ?? raw?.sourceVersion),
    title: toText(raw?.title),
    snippet: toText(raw?.snippet),
  };
}

function normalizeTicket(raw?: RawTicket): HandoffTicket {
  return {
    ticketNo: toText(raw?.ticket_no ?? raw?.ticketNo),
    conversationNo: toText(raw?.conversation_no ?? raw?.conversationNo),
    escalationReasonCode: toText(
      raw?.escalation_reason_code ?? raw?.escalationReasonCode,
    ),
    statusCode: toText(raw?.status_code ?? raw?.statusCode),
    handoffSummary: toText(raw?.handoff_summary ?? raw?.handoffSummary),
  };
}

function normalizeMissingSlot(raw?: RawConversation): AssistantMissingSlot {
  return {
    slotCode: toText(raw?.slot_code ?? raw?.slotCode),
    promptText: toText(raw?.prompt_text ?? raw?.promptText),
    required: toBoolean(raw?.required),
  };
}

function normalizeSuggestedAction(
  raw?: RawConversation,
): AssistantSuggestedAction {
  return {
    actionCode: toText(raw?.action_code ?? raw?.actionCode),
    label: toText(raw?.label),
    enabled: toBoolean(raw?.enabled),
    reasonIfDisabled: toText(raw?.reason_if_disabled ?? raw?.reasonIfDisabled),
  };
}

function normalizeDataCard(raw?: RawConversation): AssistantDataCard {
  const orderSnapshotCard = (raw?.order_snapshot_card ??
    raw?.orderSnapshotCard) as RawConversation | undefined;
  const afterSaleDecisionCard = (raw?.after_sale_decision_card ??
    raw?.afterSaleDecisionCard) as RawConversation | undefined;
  const orderSelectionCard = (raw?.order_selection_card ??
    raw?.orderSelectionCard) as RawConversation | undefined;
  const logisticsTrackingCard = (raw?.logistics_tracking_card ??
    raw?.logisticsTrackingCard) as RawConversation | undefined;
  const productRecommendationCard = (raw?.product_recommendation_card ??
    raw?.productRecommendationCard) as RawConversation | undefined;
  return {
    orderSnapshotCard: orderSnapshotCard
      ? {
          orderNo: toText(
            orderSnapshotCard.order_no ?? orderSnapshotCard.orderNo,
          ),
          mainStatus: toText(
            orderSnapshotCard.main_status ?? orderSnapshotCard.mainStatus,
          ),
          paymentStatus: toText(
            orderSnapshotCard.payment_status ?? orderSnapshotCard.paymentStatus,
          ),
          fulfillmentStatus: toText(
            orderSnapshotCard.fulfillment_status ??
              orderSnapshotCard.fulfillmentStatus,
          ),
          logisticsStatus: toText(
            orderSnapshotCard.logistics_status ??
              orderSnapshotCard.logisticsStatus,
          ),
          afterSaleStatus: toText(
            orderSnapshotCard.after_sale_status ??
              orderSnapshotCard.afterSaleStatus,
          ),
          latestUpdateTime: toText(
            orderSnapshotCard.latest_update_time ??
              orderSnapshotCard.latestUpdateTime,
          ),
        }
      : undefined,
    afterSaleDecisionCard: afterSaleDecisionCard
      ? {
          sceneCode: toText(
            afterSaleDecisionCard.scene_code ?? afterSaleDecisionCard.sceneCode,
          ),
          decisionPathCode: toText(
            afterSaleDecisionCard.decision_path_code ??
              afterSaleDecisionCard.decisionPathCode,
          ),
          reasonText: toText(
            afterSaleDecisionCard.reason_text ??
              afterSaleDecisionCard.reasonText,
          ),
          constraintText: toText(
            afterSaleDecisionCard.constraint_text ??
              afterSaleDecisionCard.constraintText,
          ),
          nextStepText: toText(
            afterSaleDecisionCard.next_step_text ??
              afterSaleDecisionCard.nextStepText,
          ),
        }
      : undefined,
    orderSelectionCard: orderSelectionCard
      ? {
          titleText: toText(
            orderSelectionCard.title_text ?? orderSelectionCard.titleText,
          ),
          helperText: toText(
            orderSelectionCard.helper_text ?? orderSelectionCard.helperText,
          ),
          originalQuery: toText(
            orderSelectionCard.original_query ??
              orderSelectionCard.originalQuery,
          ),
          taskCode: toText(
            orderSelectionCard.task_code ?? orderSelectionCard.taskCode,
          ),
          candidates: Array.isArray(orderSelectionCard.candidates)
            ? (orderSelectionCard.candidates as RawConversation[]).map(
                (item) => ({
                  orderNo: toText(item.order_no ?? item.orderNo),
                  subOrderNo: toText(item.sub_order_no ?? item.subOrderNo),
                  displayTitle: toText(item.display_title ?? item.displayTitle),
                  mainStatus: toText(item.main_status ?? item.mainStatus),
                  paymentStatus: toText(
                    item.payment_status ?? item.paymentStatus,
                  ),
                  fulfillmentStatus: toText(
                    item.fulfillment_status ?? item.fulfillmentStatus,
                  ),
                  logisticsStatus: toText(
                    item.logistics_status ?? item.logisticsStatus,
                  ),
                  afterSaleStatus: toText(
                    item.after_sale_status ?? item.afterSaleStatus,
                  ),
                  latestUpdateTime: toText(
                    item.latest_update_time ?? item.latestUpdateTime,
                  ),
                  selectionHint: toText(
                    item.selection_hint ?? item.selectionHint,
                  ),
                }),
              )
            : [],
        }
      : undefined,
    logisticsTrackingCard: logisticsTrackingCard
      ? {
          orderNo: toText(
            logisticsTrackingCard.order_no ?? logisticsTrackingCard.orderNo,
          ),
          subOrderNo: toText(
            logisticsTrackingCard.sub_order_no ??
              logisticsTrackingCard.subOrderNo,
          ),
          fulfillmentStatus: toText(
            logisticsTrackingCard.fulfillment_status ??
              logisticsTrackingCard.fulfillmentStatus,
          ),
          logisticsStatus: toText(
            logisticsTrackingCard.logistics_status ??
              logisticsTrackingCard.logisticsStatus,
          ),
          latestUpdateTime: toText(
            logisticsTrackingCard.latest_update_time ??
              logisticsTrackingCard.latestUpdateTime,
          ),
          latestTraceText: toText(
            logisticsTrackingCard.latest_trace_text ??
              logisticsTrackingCard.latestTraceText,
          ),
          timelineSummary: toText(
            logisticsTrackingCard.timeline_summary ??
              logisticsTrackingCard.timelineSummary,
          ),
        }
      : undefined,
    productRecommendationCard: productRecommendationCard
      ? {
          titleText: toText(
            productRecommendationCard.title_text ??
              productRecommendationCard.titleText,
          ),
          helperText: toText(
            productRecommendationCard.helper_text ??
              productRecommendationCard.helperText,
          ),
          items: Array.isArray(productRecommendationCard.items)
            ? (productRecommendationCard.items as RawConversation[]).map(
                (item) => ({
                  spuNo: toText(item.spu_no ?? item.spuNo),
                  title: toText(item.title),
                  coverUrl: toText(item.cover_url ?? item.coverUrl),
                  priceText: toText(item.price_text ?? item.priceText),
                  shopName: toText(item.shop_name ?? item.shopName),
                  reasonText: toText(item.reason_text ?? item.reasonText),
                }),
              )
            : [],
        }
      : undefined,
  };
}

function normalizeReplyPayload(
  raw?: RawConversation,
): AssistantReplyPayload | undefined {
  if (!raw) {
    return undefined;
  }
  return {
    replyText: toText(raw.reply_text ?? raw.replyText),
    intentCode: toText(raw.intent_code ?? raw.intentCode),
    guardResultCode: toText(raw.guard_result_code ?? raw.guardResultCode),
    confidence: toNumber(raw.confidence),
    dataCards: Array.isArray(raw.data_cards ?? raw.dataCards)
      ? ((raw.data_cards ?? raw.dataCards) as RawConversation[]).map((item) =>
          normalizeDataCard(item),
        )
      : [],
    suggestedActions: Array.isArray(
      raw.suggested_actions ?? raw.suggestedActions,
    )
      ? (
          (raw.suggested_actions ?? raw.suggestedActions) as RawConversation[]
        ).map((item) => normalizeSuggestedAction(item))
      : [],
    missingSlots: Array.isArray(raw.missing_slots ?? raw.missingSlots)
      ? ((raw.missing_slots ?? raw.missingSlots) as RawConversation[]).map(
          (item) => normalizeMissingSlot(item),
        )
      : [],
    slotRetryCount: toNumber(raw.slot_retry_count ?? raw.slotRetryCount),
    handoffRecommended: toBoolean(
      raw.handoff_recommended ?? raw.handoffRecommended,
    ),
    handoffReasonCode: toText(raw.handoff_reason_code ?? raw.handoffReasonCode),
  };
}

export async function createOrGetAssistantConversation(payload: {
  shopNo?: string;
  orderNo?: string;
  anchorSpuNo?: string;
  anchorSkuNo?: string;
  forceNew?: boolean;
}): Promise<{ conversation: AssistantConversation; latestRun?: AssistantRun }> {
  const response = await apiClient.post<RawCreateConversationRes>(
    "/v1/agent/buyer/conversations:get-or-create",
    {
      scene_code: "BUYER_ASSISTANT",
      bot_code: "BUYER_ASSISTANT",
      shop_no: payload.shopNo ?? "",
      order_no: payload.orderNo ?? "",
      anchor_spu_no: payload.anchorSpuNo ?? "",
      anchor_sku_no: payload.anchorSkuNo ?? "",
      force_new: Boolean(payload.forceNew),
    },
  );
  const latestRun = response?.latest_run ?? response?.latestRun;
  return {
    conversation: normalizeConversation(response?.conversation),
    latestRun: latestRun ? normalizeRun(latestRun) : undefined,
  };
}

export async function sendAssistantMessage(payload: {
  conversationNo: string;
  contentText: string;
  clientMessageNo?: string;
  hiddenAction?: AssistantHiddenAction;
}): Promise<{
  conversation: AssistantConversation;
  run: AssistantRun;
  userMessage: AssistantMessage;
  idempotentReplay: boolean;
}> {
  const response = await apiClient.post<RawSendAssistantMessageRes>(
    "/v1/agent/buyer/messages:send",
    {
      conversation_no: payload.conversationNo,
      client_message_no:
        payload.clientMessageNo ??
        `agent_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
      message_type_code: "TEXT",
      content_text: payload.contentText,
      asset_ids: [],
      hidden_action: payload.hiddenAction,
    },
  );
  return {
    conversation: normalizeConversation(response?.conversation),
    run: normalizeRun(response?.run),
    userMessage: normalizeMessage(
      response?.user_message ?? response?.userMessage,
    ),
    idempotentReplay: toBoolean(
      response?.idempotent_replay ?? response?.idempotentReplay,
    ),
  };
}

export async function getAssistantRunStatus(
  conversationNo: string,
  runNo = "",
): Promise<{
  conversation: AssistantConversation;
  run: AssistantRun;
  latestAssistantMessage?: AssistantMessage;
  answerSources: AnswerSource[];
  pendingMessageCount: number;
}> {
  const response = await apiClient.get<RawGetRunStatusRes>(
    `/v1/agent/buyer/conversations/${encodeURIComponent(conversationNo)}/run-status`,
    {
      params: { run_no: runNo },
    },
  );
  const latestAssistantMessage =
    response?.latest_assistant_message ?? response?.latestAssistantMessage;
  return {
    conversation: normalizeConversation(response?.conversation),
    run: normalizeRun(response?.run),
    latestAssistantMessage: latestAssistantMessage
      ? normalizeMessage(latestAssistantMessage)
      : undefined,
    answerSources: (
      response?.answer_sources ??
      response?.answerSources ??
      []
    ).map((item) => normalizeSource(item)),
    pendingMessageCount: toNumber(
      response?.pending_message_count ?? response?.pendingMessageCount,
    ),
  };
}

export async function listAssistantMessages(
  conversationNo: string,
  pageSize = 100,
  nextCursor = "",
): Promise<{ list: AssistantMessage[]; nextCursor: string; hasMore: boolean }> {
  const response = await apiClient.get<RawListMessagesRes>(
    `/v1/agent/buyer/conversations/${encodeURIComponent(conversationNo)}/messages`,
    {
      params: { page_size: pageSize, next_cursor: nextCursor },
    },
  );
  return {
    list: (response?.list ?? []).map((item) => normalizeMessage(item)),
    nextCursor: toText(response?.next_cursor ?? response?.nextCursor),
    hasMore: toBoolean(response?.has_more ?? response?.hasMore),
  };
}

export async function escalateAssistantConversation(payload: {
  conversationNo: string;
  escalationReasonCode: string;
  remark?: string;
}): Promise<{ ticket: HandoffTicket; conversation: AssistantConversation }> {
  const response = await apiClient.post<RawEscalateRes>(
    `/v1/agent/buyer/conversations/${encodeURIComponent(payload.conversationNo)}:escalate`,
    {
      escalation_reason_code: payload.escalationReasonCode,
      remark: payload.remark ?? "",
    },
  );
  return {
    ticket: normalizeTicket(response?.ticket),
    conversation: normalizeConversation(response?.conversation),
  };
}

export async function submitAssistantFeedback(payload: {
  conversationNo: string;
  runNo: string;
  feedbackCode: "HELPFUL" | "NOT_HELPFUL";
  resolved?: boolean;
  comment?: string;
}): Promise<void> {
  await apiClient.post("/v1/agent/buyer/feedback:submit", {
    conversation_no: payload.conversationNo,
    run_no: payload.runNo,
    feedback_code: payload.feedbackCode,
    resolved: Boolean(payload.resolved),
    comment: payload.comment ?? "",
  });
}
