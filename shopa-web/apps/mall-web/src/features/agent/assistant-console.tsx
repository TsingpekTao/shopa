"use client";

import { FormEvent, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import {
  AssistantConsoleVariant,
  buildActionPrompt,
  buildOrderSelectionPrompt,
  getIntroMessageText,
  getVariantCopy,
} from "@/features/agent/assistant-copy";
import {
  formatAssistantCandidateUpdateTime,
  shouldRenderSidebarStructuredCards,
} from "@/features/agent/assistant-display";
import { AssistantOrderSelectionCardView } from "@/features/agent/assistant-order-selection-card";
import {
  createOrGetAssistantConversation,
  escalateAssistantConversation,
  getAssistantRunStatus,
  listAssistantMessages,
  sendAssistantMessage,
  submitAssistantFeedback,
} from "@/features/agent/api";
import {
  AssistantAfterSaleDecisionCard,
  AssistantHiddenAction,
  AssistantLogisticsTrackingCard,
  AssistantMessage,
  AssistantOrderSelectionCard,
  AssistantOrderSnapshotCard,
  AssistantProductRecommendationCard,
  AssistantRecentOrderCandidate,
  AssistantReplyPayload,
  AssistantSuggestedAction,
} from "@/features/agent/types";

type AssistantConsoleProps = {
  variant?: AssistantConsoleVariant;
  backHref?: string;
};

type RenderedMessage = {
  key: string;
  runNo: string;
  senderTypeCode: string;
  contentText: string;
  sentAt: string;
  virtual?: boolean;
};

type SendMessageInput = {
  contentText: string;
  hiddenAction?: AssistantHiddenAction;
};

function formatSentAt(raw: string, locale: string): string {
  if (!raw) {
    return "--";
  }
  const dt = new Date(raw);
  if (Number.isNaN(dt.getTime())) {
    return raw;
  }
  return dt.toLocaleString(locale === "zh-CN" ? "zh-CN" : "en-US", {
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function sortMessages(list: AssistantMessage[]): AssistantMessage[] {
  return [...list].sort((a, b) => {
    const av = a.sentAt ? new Date(a.sentAt).getTime() : 0;
    const bv = b.sentAt ? new Date(b.sentAt).getTime() : 0;
    return av - bv;
  });
}

function buildIntroMessage(
  variant: AssistantConsoleVariant,
  isZh: boolean,
): RenderedMessage {
  return {
    key: `welcome-${variant}`,
    runNo: "",
    senderTypeCode: "AGENT",
    contentText: getIntroMessageText(variant, isZh),
    sentAt: "",
    virtual: true,
  };
}

function findOrderSnapshotCard(
  replyPayload?: AssistantReplyPayload,
): AssistantOrderSnapshotCard | undefined {
  return replyPayload?.dataCards.find((item) => item.orderSnapshotCard)
    ?.orderSnapshotCard;
}

function findDecisionCard(
  replyPayload?: AssistantReplyPayload,
): AssistantAfterSaleDecisionCard | undefined {
  return replyPayload?.dataCards.find((item) => item.afterSaleDecisionCard)
    ?.afterSaleDecisionCard;
}

function findOrderSelectionCard(
  replyPayload?: AssistantReplyPayload,
): AssistantOrderSelectionCard | undefined {
  return replyPayload?.dataCards.find((item) => item.orderSelectionCard)
    ?.orderSelectionCard;
}

function findLogisticsTrackingCard(
  replyPayload?: AssistantReplyPayload,
): AssistantLogisticsTrackingCard | undefined {
  return replyPayload?.dataCards.find((item) => item.logisticsTrackingCard)
    ?.logisticsTrackingCard;
}

function findProductRecommendationCard(
  replyPayload?: AssistantReplyPayload,
): AssistantProductRecommendationCard | undefined {
  return replyPayload?.dataCards.find((item) => item.productRecommendationCard)
    ?.productRecommendationCard;
}

function getRunPollingInterval(runStatusCode: string): number | false {
  switch (runStatusCode) {
    case "PENDING":
    case "PROCESSING":
    case "WAITING_TOOL_CALL":
    case "GENERATING":
      return 2_000;
    default:
      return false;
  }
}

function renderAssistantName(
  variant: AssistantConsoleVariant,
  isZh: boolean,
  senderTypeCode: string,
): string {
  if (senderTypeCode === "BUYER") {
    return isZh ? "我" : "You";
  }
  if (variant === "official") {
    return isZh ? "官方客服助手" : "Official support";
  }
  return "Shopa AI";
}

export function AssistantConsole({
  variant = "shopping",
  backHref,
}: AssistantConsoleProps) {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const copy = useMemo(() => getVariantCopy(variant, isZh), [isZh, variant]);
  const [messageApi, contextHolder] = message.useMessage();
  const queryClient = useQueryClient();
  const searchParams = useSearchParams();

  const shopNo = searchParams.get("shop_no") ?? "";
  const orderNo = searchParams.get("order_no") ?? "";
  const anchorSpuNo = searchParams.get("spu_no") ?? "";
  const anchorSkuNo = searchParams.get("sku_no") ?? "";

  const [conversationNo, setConversationNo] = useState("");
  const [runNo, setRunNo] = useState("");
  const [draft, setDraft] = useState("");

  const bootstrapQuery = useQuery({
    queryKey: [
      "assistant-bootstrap",
      variant,
      shopNo,
      orderNo,
      anchorSpuNo,
      anchorSkuNo,
    ],
    queryFn: () =>
      createOrGetAssistantConversation({
        shopNo,
        orderNo,
        anchorSpuNo,
        anchorSkuNo,
      }),
    retry: false,
    staleTime: 60_000,
    refetchOnWindowFocus: false,
  });

  useEffect(() => {
    const result = bootstrapQuery.data;
    if (!result) {
      return;
    }
    if (conversationNo !== result.conversation.conversationNo) {
      setConversationNo(result.conversation.conversationNo);
    }
    if (result.latestRun?.runNo && runNo !== result.latestRun.runNo) {
      setRunNo(result.latestRun.runNo);
    }
  }, [bootstrapQuery.data, conversationNo, runNo]);

  useEffect(() => {
    if (!bootstrapQuery.isError) {
      return;
    }
    messageApi.error(
      isZh
        ? "智能客服暂时不可用，请稍后再试。"
        : "Assistant is unavailable right now.",
    );
  }, [bootstrapQuery.isError, bootstrapQuery.errorUpdatedAt, isZh, messageApi]);

  const messagesQuery = useQuery({
    queryKey: ["assistant-messages", conversationNo],
    queryFn: () => listAssistantMessages(conversationNo, 100, ""),
    enabled: conversationNo.length > 0,
    staleTime: 5_000,
    refetchInterval: 8_000,
  });

  const runStatusQuery = useQuery({
    queryKey: ["assistant-run-status", conversationNo, runNo],
    queryFn: () => getAssistantRunStatus(conversationNo, runNo),
    enabled: conversationNo.length > 0 && runNo.length > 0,
    staleTime: 2_000,
    refetchInterval: (data) =>
      getRunPollingInterval(data?.run?.runStatusCode ?? ""),
  });

  const sendMutation = useMutation({
    mutationFn: (payload: SendMessageInput) =>
      sendAssistantMessage({
        conversationNo,
        contentText: payload.contentText,
        hiddenAction: payload.hiddenAction,
      }),
    onSuccess: (result) => {
      setDraft("");
      setRunNo(result.run.runNo);
      queryClient.invalidateQueries({
        queryKey: ["assistant-messages", conversationNo],
      });
      queryClient.invalidateQueries({
        queryKey: ["assistant-run-status", conversationNo],
      });
    },
    onError: () => {
      messageApi.error(
        isZh ? "发送失败，请稍后重试。" : "Failed to send message.",
      );
    },
  });

  const escalateMutation = useMutation({
    mutationFn: () =>
      escalateAssistantConversation({
        conversationNo,
        escalationReasonCode: "USER_REQUEST",
      }),
    onSuccess: (result) => {
      messageApi.success(
        result.ticket.handoffSummary ||
          (isZh ? "已转人工客服，请稍等。" : "Escalated to human support."),
      );
      queryClient.invalidateQueries({
        queryKey: ["assistant-run-status", conversationNo],
      });
    },
    onError: () => {
      messageApi.error(
        isZh ? "转人工失败，请稍后重试。" : "Failed to escalate.",
      );
    },
  });

  const feedbackMutation = useMutation({
    mutationFn: (feedbackCode: "HELPFUL" | "NOT_HELPFUL") =>
      submitAssistantFeedback({
        conversationNo,
        runNo,
        feedbackCode,
        resolved: feedbackCode === "HELPFUL",
      }),
    onSuccess: () => {
      messageApi.success(
        isZh ? "感谢反馈，我们会继续优化回答。" : "Thanks for the feedback.",
      );
    },
    onError: () => {
      messageApi.error(
        isZh ? "反馈提交失败，请稍后重试。" : "Failed to submit feedback.",
      );
    },
  });

  const orderedMessages = sortMessages(messagesQuery.data?.list ?? []);
  const lastAssistantMessage = [...orderedMessages]
    .reverse()
    .find((item) => item.senderTypeCode === "AGENT");

  const displayMessages = useMemo(() => {
    const hasAssistantReply = orderedMessages.some(
      (item) => item.senderTypeCode === "AGENT",
    );
    const base = orderedMessages.map<RenderedMessage>((item) => ({
      key: item.messageNo,
      runNo: item.runNo,
      senderTypeCode: item.senderTypeCode,
      contentText: item.contentText,
      sentAt: item.sentAt,
    }));
    if (hasAssistantReply) {
      return base;
    }
    return [buildIntroMessage(variant, isZh), ...base];
  }, [isZh, orderedMessages, variant]);

  const currentConversation =
    runStatusQuery.data?.conversation ?? bootstrapQuery.data?.conversation;
  const currentRun = runStatusQuery.data?.run ?? bootstrapQuery.data?.latestRun;
  const replyPayload = currentRun?.replyPayload;
  const answerSources = runStatusQuery.data?.answerSources ?? [];

  const orderSnapshotCard = findOrderSnapshotCard(replyPayload);
  const decisionCard = findDecisionCard(replyPayload);
  const orderSelectionCard = findOrderSelectionCard(replyPayload);
  const logisticsTrackingCard = findLogisticsTrackingCard(replyPayload);
  const productRecommendationCard = findProductRecommendationCard(replyPayload);

  const showSidebarStructuredCards = shouldRenderSidebarStructuredCards({
    hasOrderSnapshotCard: Boolean(orderSnapshotCard),
    hasDecisionCard: Boolean(decisionCard),
    hasLogisticsCard: Boolean(logisticsTrackingCard),
    hasProductRecommendationCard: Boolean(productRecommendationCard),
  });

  function submitText(payload: SendMessageInput) {
    const normalized = payload.contentText.trim();
    if (!conversationNo || !normalized || sendMutation.isLoading) {
      return;
    }
    sendMutation.mutate({
      contentText: normalized,
      hiddenAction: payload.hiddenAction,
    });
  }

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    submitText({ contentText: draft });
  }

  function handleSuggestedAction(action: AssistantSuggestedAction) {
    if (!action.enabled) {
      if (action.reasonIfDisabled) {
        messageApi.info(action.reasonIfDisabled);
      }
      return;
    }
    if (action.actionCode === "escalate_to_human") {
      escalateMutation.mutate();
      return;
    }
    submitText({ contentText: buildActionPrompt(action, isZh) });
  }

  function handleQuickPrompt(prompt: string) {
    submitText({ contentText: prompt });
  }

  function handleSelectRecentOrder(candidate: AssistantRecentOrderCandidate) {
    if (!orderSelectionCard) {
      return;
    }
    submitText({
      contentText: buildOrderSelectionPrompt(
        orderSelectionCard,
        candidate,
        isZh,
      ),
      hiddenAction: {
        type: "SET_SLOT",
        key: "selected_order_no",
        value: candidate.orderNo,
      },
    });
  }

  return (
    <main className="tb-agent-page">
      {contextHolder}

      <section className="tb-agent-hero">
        <div>
          <span className="tb-chip">{copy.chip}</span>
          <h1>{copy.title}</h1>
          <p>{copy.description}</p>
        </div>
        <div className="tb-agent-hero-tags">
          <span>
            {shopNo || (isZh ? "官方客服场景" : "Official support scene")}
          </span>
          <span>
            {anchorSpuNo || (isZh ? "当前未锚定商品" : "No anchored product")}
          </span>
          <span>
            {orderNo || (isZh ? "当前未锚定订单" : "No anchored order")}
          </span>
        </div>
      </section>

      <section className="tb-agent-layout">
        <aside className="tb-agent-sidebar">
          <div className="tb-agent-panel">
            <h2>{isZh ? "快捷提问" : "Quick prompts"}</h2>
            <div className="tb-agent-prompt-list">
              {copy.quickPrompts.map((item) => (
                <button
                  key={item}
                  type="button"
                  onClick={() => handleQuickPrompt(item)}
                >
                  {item}
                </button>
              ))}
            </div>
          </div>

          <div className="tb-agent-panel">
            <h2>{isZh ? "会话摘要" : "Session summary"}</h2>
            <p>{currentConversation?.sessionSummary || copy.emptySummary}</p>
            <dl>
              <div>
                <dt>{isZh ? "状态" : "Status"}</dt>
                <dd>
                  {currentRun?.runStatusCode || (isZh ? "待开始" : "Idle")}
                </dd>
              </div>
              <div>
                <dt>{isZh ? "待处理消息" : "Pending"}</dt>
                <dd>{currentConversation?.pendingMessageCount ?? 0}</dd>
              </div>
            </dl>
            <button
              type="button"
              className="tb-agent-escalate"
              disabled={!conversationNo || escalateMutation.isLoading}
              onClick={() => escalateMutation.mutate()}
            >
              {escalateMutation.isLoading
                ? isZh
                  ? "转接中..."
                  : "Escalating..."
                : isZh
                  ? "转人工客服"
                  : "Escalate to human"}
            </button>
          </div>
        </aside>

        <div className="tb-agent-main">
          <div className="tb-agent-panel tb-agent-thread-panel">
            <div className="tb-agent-thread-head">
              <div>
                <strong>
                  {variant === "official"
                    ? isZh
                      ? "官方智能客服会话"
                      : "Official support conversation"
                    : isZh
                      ? "智能助手会话"
                      : "Assistant conversation"}
                </strong>
                <span>
                  {conversationNo ||
                    (isZh ? "正在创建会话..." : "Creating conversation...")}
                </span>
              </div>
              <Link
                href={
                  backHref ?? (variant === "official" ? "/help" : "/service")
                }
              >
                {copy.backLabel}
              </Link>
            </div>

            <div className="tb-agent-thread">
              {displayMessages.map((item) => (
                <article
                  key={item.key}
                  className={`tb-agent-msg ${item.senderTypeCode === "BUYER" ? "is-self" : ""}`}
                >
                  <div className="tb-agent-bubble">
                    <strong>
                      {renderAssistantName(variant, isZh, item.senderTypeCode)}
                    </strong>
                    <p>
                      {item.contentText ||
                        (isZh
                          ? "当前消息类型暂不支持展示。"
                          : "This message type is not supported yet.")}
                    </p>
                    {item.sentAt ? (
                      <small>{formatSentAt(item.sentAt, locale)}</small>
                    ) : item.virtual ? (
                      <small>
                        {isZh ? "已准备好为你服务" : "Ready to help"}
                      </small>
                    ) : null}
                    {item.senderTypeCode === "AGENT" &&
                    item.runNo &&
                    item.runNo === currentRun?.runNo &&
                    orderSelectionCard ? (
                      <div className="tb-agent-inline-selection">
                        <AssistantOrderSelectionCardView
                          card={orderSelectionCard}
                          isZh={isZh}
                          onSelect={handleSelectRecentOrder}
                        />
                      </div>
                    ) : null}
                  </div>
                </article>
              ))}
              {bootstrapQuery.isLoading || messagesQuery.isLoading ? (
                <p className="tb-agent-empty">
                  {isZh ? "正在连接智能客服..." : "Connecting assistant..."}
                </p>
              ) : null}
            </div>

            <form className="tb-agent-composer" onSubmit={handleSubmit}>
              <textarea
                value={draft}
                onChange={(event) => setDraft(event.target.value)}
                placeholder={copy.placeholder}
                rows={4}
              />
              <div className="tb-agent-composer-actions">
                <span>
                  {isZh
                    ? "直接描述你的问题，或者点击上方快捷提问，我会继续帮你处理。"
                    : "Describe your issue directly or choose a quick prompt above to continue."}
                </span>
                <button
                  type="submit"
                  disabled={
                    !conversationNo || !draft.trim() || sendMutation.isLoading
                  }
                >
                  {sendMutation.isLoading
                    ? isZh
                      ? "发送中..."
                      : "Sending..."
                    : isZh
                      ? "发送给客服"
                      : "Send"}
                </button>
              </div>
            </form>
          </div>

          <div className="tb-agent-panel tb-agent-source-panel">
            <div className="tb-agent-source-head">
              <strong>{isZh ? "本轮建议" : "Current guidance"}</strong>
              {replyPayload?.handoffRecommended ? (
                <span>
                  {isZh ? "当前建议转人工" : "Human handoff recommended"}
                </span>
              ) : null}
            </div>

            {showSidebarStructuredCards ? (
              <div className="tb-agent-structured-list">
                {orderSnapshotCard ? (
                  <article className="tb-agent-structured-card">
                    <strong>{isZh ? "订单事实" : "Order facts"}</strong>
                    <dl>
                      <div>
                        <dt>{isZh ? "订单号" : "Order"}</dt>
                        <dd>{orderSnapshotCard.orderNo || "--"}</dd>
                      </div>
                      <div>
                        <dt>{isZh ? "主状态" : "Main status"}</dt>
                        <dd>{orderSnapshotCard.mainStatus || "--"}</dd>
                      </div>
                      <div>
                        <dt>{isZh ? "履约状态" : "Fulfillment"}</dt>
                        <dd>{orderSnapshotCard.fulfillmentStatus || "--"}</dd>
                      </div>
                      <div>
                        <dt>{isZh ? "物流状态" : "Logistics"}</dt>
                        <dd>{orderSnapshotCard.logisticsStatus || "--"}</dd>
                      </div>
                    </dl>
                  </article>
                ) : null}

                {logisticsTrackingCard ? (
                  <article className="tb-agent-structured-card">
                    <strong>{isZh ? "物流进度" : "Logistics tracking"}</strong>
                    <p>{logisticsTrackingCard.latestTraceText || "--"}</p>
                    <small>
                      {`${isZh ? "订单号" : "Order"}: ${logisticsTrackingCard.orderNo || "--"}`}
                    </small>
                    <small>
                      {`${isZh ? "物流状态" : "Logistics"}: ${logisticsTrackingCard.logisticsStatus || "--"}`}
                    </small>
                    {logisticsTrackingCard.latestUpdateTime ? (
                      <small>
                        {`${isZh ? "最近更新" : "Updated"}: ${formatAssistantCandidateUpdateTime(
                          logisticsTrackingCard.latestUpdateTime,
                          locale,
                        )}`}
                      </small>
                    ) : null}
                    {logisticsTrackingCard.timelineSummary ? (
                      <small>{logisticsTrackingCard.timelineSummary}</small>
                    ) : null}
                  </article>
                ) : null}

                {decisionCard ? (
                  <article className="tb-agent-structured-card">
                    <strong>{isZh ? "系统判断" : "Decision"}</strong>
                    <p>{decisionCard.reasonText || "--"}</p>
                    <small>
                      {`${isZh ? "场景码" : "Scene"}: ${decisionCard.sceneCode || "--"}`}
                    </small>
                    {decisionCard.constraintText ? (
                      <small>{decisionCard.constraintText}</small>
                    ) : null}
                    {decisionCard.nextStepText ? (
                      <small>{decisionCard.nextStepText}</small>
                    ) : null}
                  </article>
                ) : null}

                {productRecommendationCard ? (
                  <article className="tb-agent-structured-card">
                    <strong>
                      {productRecommendationCard.titleText ||
                        (isZh ? "商品推荐" : "Product recommendation")}
                    </strong>
                    {productRecommendationCard.helperText ? (
                      <p>{productRecommendationCard.helperText}</p>
                    ) : null}
                    <div className="tb-agent-source-list">
                      {productRecommendationCard.items.map((item) => (
                        <article
                          key={`${item.spuNo}-${item.title}`}
                          className="tb-agent-product-card"
                        >
                          <strong>{item.title || item.spuNo || "--"}</strong>
                          {item.priceText ? <p>{item.priceText}</p> : null}
                          {item.reasonText ? (
                            <small>{item.reasonText}</small>
                          ) : null}
                          {item.shopName ? (
                            <small>{item.shopName}</small>
                          ) : null}
                        </article>
                      ))}
                    </div>
                  </article>
                ) : null}
              </div>
            ) : (
              <p className="tb-agent-empty">
                {isZh
                  ? "发来一条消息后，这里会展示订单事实、物流进度、处理建议和推荐动作。"
                  : "After a turn, order facts, logistics, decisions, and actions will appear here."}
              </p>
            )}

            {replyPayload?.missingSlots?.length ? (
              <div className="tb-agent-slot-list">
                <strong>{isZh ? "还需要你补充" : "Still needed"}</strong>
                {replyPayload.missingSlots.map((slot) => (
                  <p key={slot.slotCode}>{slot.promptText}</p>
                ))}
              </div>
            ) : null}

            {replyPayload?.suggestedActions?.length ? (
              <div className="tb-agent-action-list">
                {replyPayload.suggestedActions.map((action) => (
                  <button
                    key={`${action.actionCode}-${action.label}`}
                    type="button"
                    disabled={!action.enabled}
                    title={action.reasonIfDisabled || ""}
                    onClick={() => handleSuggestedAction(action)}
                  >
                    {action.label}
                  </button>
                ))}
              </div>
            ) : null}

            <div className="tb-agent-source-head">
              <strong>{isZh ? "本轮回答依据" : "Answer sources"}</strong>
              {currentRun?.promptInjectionFlag ? (
                <span>
                  {isZh ? "已触发安全防护" : "Safety guard triggered"}
                </span>
              ) : null}
            </div>

            {answerSources.length === 0 ? (
              <p className="tb-agent-empty">
                {isZh
                  ? "本轮暂无外部依据。"
                  : "No explicit sources for this turn."}
              </p>
            ) : null}

            <div className="tb-agent-source-list">
              {answerSources.map((item) => (
                <article
                  key={`${item.sourceTypeCode}-${item.sourceId}-${item.sourceVersion}`}
                >
                  <strong>{item.title || item.sourceTypeCode}</strong>
                  {item.snippet ? <p>{item.snippet}</p> : null}
                  <small>{`${item.sourceTypeCode} / v${item.sourceVersion}`}</small>
                </article>
              ))}
            </div>

            {lastAssistantMessage && runNo ? (
              <div className="tb-agent-feedback-row">
                <button
                  type="button"
                  onClick={() => feedbackMutation.mutate("HELPFUL")}
                  disabled={feedbackMutation.isLoading}
                >
                  {isZh ? "这条有帮助" : "Helpful"}
                </button>
                <button
                  type="button"
                  onClick={() => feedbackMutation.mutate("NOT_HELPFUL")}
                  disabled={feedbackMutation.isLoading}
                >
                  {isZh ? "这条没帮助" : "Not helpful"}
                </button>
              </div>
            ) : null}
          </div>
        </div>
      </section>
    </main>
  );
}
