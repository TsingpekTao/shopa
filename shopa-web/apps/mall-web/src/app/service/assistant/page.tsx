"use client";

import { FormEvent, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import {
  createOrGetAssistantConversation,
  escalateAssistantConversation,
  getAssistantRunStatus,
  listAssistantMessages,
  sendAssistantMessage,
  submitAssistantFeedback
} from "@/features/agent/api";
import { AssistantMessage } from "@/features/agent/types";

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
    minute: "2-digit"
  });
}

function sortMessages(list: AssistantMessage[]): AssistantMessage[] {
  return [...list].sort((a, b) => {
    const av = a.sentAt ? new Date(a.sentAt).getTime() : 0;
    const bv = b.sentAt ? new Date(b.sentAt).getTime() : 0;
    return av - bv;
  });
}

export default function AssistantServicePage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
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

  const quickPrompts = useMemo(
    () => [
      isZh ? "帮我推荐适合通勤的手机" : "Recommend a commuter-friendly phone",
      isZh ? "这件商品什么时候能发货" : "When can this product ship?",
      isZh ? "售后退货流程怎么走" : "How does the return flow work?",
      isZh ? "我想找一款红色的商品" : "I want to find a red product"
    ],
    [isZh]
  );

  const bootstrapMutation = useMutation({
    mutationFn: () =>
      createOrGetAssistantConversation({
        shopNo,
        orderNo,
        anchorSpuNo,
        anchorSkuNo
      }),
    onSuccess: (result) => {
      setConversationNo(result.conversation.conversationNo);
      if (result.latestRun?.runNo) {
        setRunNo(result.latestRun.runNo);
      }
    },
    onError: () => {
      messageApi.error(isZh ? "智能客服暂时不可用，请稍后再试" : "Assistant is unavailable right now");
    }
  });

  useEffect(() => {
    if (bootstrapMutation.isLoading || bootstrapMutation.isSuccess || conversationNo) {
      return;
    }
    bootstrapMutation.mutate();
  }, [bootstrapMutation, conversationNo]);

  const messagesQuery = useQuery({
    queryKey: ["assistant-messages", conversationNo],
    queryFn: () => listAssistantMessages(conversationNo, 100, ""),
    enabled: conversationNo.length > 0,
    staleTime: 5_000,
    refetchInterval: 8_000
  });

  const runStatusQuery = useQuery({
    queryKey: ["assistant-run-status", conversationNo, runNo],
    queryFn: () => getAssistantRunStatus(conversationNo, runNo),
    enabled: conversationNo.length > 0,
    staleTime: 2_000,
    refetchInterval: (data) => {
      const status = data?.run?.runStatusCode ?? "";
      return status === "PROCESSING" ? 2_000 : false;
    }
  });

  const sendMutation = useMutation({
    mutationFn: () =>
      sendAssistantMessage({
        conversationNo,
        contentText: draft.trim()
      }),
    onSuccess: (result) => {
      setDraft("");
      setRunNo(result.run.runNo);
      queryClient.invalidateQueries({ queryKey: ["assistant-messages", conversationNo] });
      queryClient.invalidateQueries({ queryKey: ["assistant-run-status", conversationNo] });
    },
    onError: () => {
      messageApi.error(isZh ? "发送失败，请稍后重试" : "Failed to send message");
    }
  });

  const escalateMutation = useMutation({
    mutationFn: () =>
      escalateAssistantConversation({
        conversationNo,
        escalationReasonCode: "USER_REQUEST"
      }),
    onSuccess: (result) => {
      messageApi.success(result.ticket.handoffSummary || (isZh ? "已转人工客服，请稍候" : "Escalated to human support"));
      queryClient.invalidateQueries({ queryKey: ["assistant-run-status", conversationNo] });
    },
    onError: () => {
      messageApi.error(isZh ? "转人工失败，请稍后重试" : "Failed to escalate");
    }
  });

  const feedbackMutation = useMutation({
    mutationFn: (feedbackCode: "HELPFUL" | "NOT_HELPFUL") =>
      submitAssistantFeedback({
        conversationNo,
        runNo,
        feedbackCode,
        resolved: feedbackCode === "HELPFUL"
      }),
    onSuccess: () => {
      messageApi.success(isZh ? "感谢反馈，我们会继续优化回答" : "Thanks for the feedback");
    },
    onError: () => {
      messageApi.error(isZh ? "反馈提交失败" : "Failed to submit feedback");
    }
  });

  const orderedMessages = sortMessages(messagesQuery.data?.list ?? []);
  const lastAssistantMessage = [...orderedMessages].reverse().find((item) => item.senderTypeCode === "AGENT");
  const answerSources = runStatusQuery.data?.answerSources ?? [];
  const currentConversation = runStatusQuery.data?.conversation;
  const currentRun = runStatusQuery.data?.run;

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    if (!conversationNo || !draft.trim()) {
      return;
    }
    sendMutation.mutate();
  }

  return (
    <main className="tb-agent-page">
      {contextHolder}
      <section className="tb-agent-hero">
        <div>
          <span className="tb-chip">{isZh ? "Shopa AI 导购" : "Shopa AI Concierge"}</span>
          <h1>{isZh ? "边逛边问，像淘宝智能客服一样给你答复" : "Ask while you browse and get marketplace-style help"}</h1>
          <p>
            {isZh
              ? "它会结合商品、店铺和平台规则帮你答疑。当前版本先支持导购、订单流程、售后规则和转人工。"
              : "It blends product, shop and policy context to help with shopping, order flow, after-sale and human handoff."}
          </p>
        </div>
        <div className="tb-agent-hero-tags">
          <span>{shopNo || (isZh ? "全站通用场景" : "General shopping scene")}</span>
          <span>{anchorSpuNo || (isZh ? "未绑定商品" : "No anchored product")}</span>
          <span>{orderNo || (isZh ? "未绑定订单" : "No anchored order")}</span>
        </div>
      </section>

      <section className="tb-agent-layout">
        <aside className="tb-agent-sidebar">
          <div className="tb-agent-panel">
            <h2>{isZh ? "快捷提问" : "Quick prompts"}</h2>
            <div className="tb-agent-prompt-list">
              {quickPrompts.map((item) => (
                <button key={item} type="button" onClick={() => setDraft(item)}>
                  {item}
                </button>
              ))}
            </div>
          </div>

          <div className="tb-agent-panel">
            <h2>{isZh ? "会话摘要" : "Session summary"}</h2>
            <p>{currentConversation?.sessionSummary || (isZh ? "先发一条消息，助手会逐步形成摘要。" : "Send a message first and the assistant will build a summary.")}</p>
            <dl>
              <div>
                <dt>{isZh ? "状态" : "Status"}</dt>
                <dd>{currentRun?.runStatusCode || (isZh ? "待开始" : "Idle")}</dd>
              </div>
              <div>
                <dt>{isZh ? "待处理消息" : "Pending"}</dt>
                <dd>{currentConversation?.pendingMessageCount ?? 0}</dd>
              </div>
            </dl>
            <button type="button" className="tb-agent-escalate" disabled={!conversationNo || escalateMutation.isLoading} onClick={() => escalateMutation.mutate()}>
              {escalateMutation.isLoading ? (isZh ? "转接中..." : "Escalating...") : isZh ? "转人工客服" : "Escalate to human"}
            </button>
          </div>
        </aside>

        <div className="tb-agent-main">
          <div className="tb-agent-panel tb-agent-thread-panel">
            <div className="tb-agent-thread-head">
              <div>
                <strong>{isZh ? "智能客服会话" : "Assistant conversation"}</strong>
                <span>{conversationNo || (isZh ? "正在创建会话..." : "Creating conversation...")}</span>
              </div>
              <Link href="/service">{isZh ? "返回服务中心" : "Back to service center"}</Link>
            </div>

            <div className="tb-agent-thread">
              {orderedMessages.map((item) => (
                <article key={item.messageNo} className={`tb-agent-msg ${item.senderTypeCode === "BUYER" ? "is-self" : ""}`}>
                  <div className="tb-agent-bubble">
                    <strong>{item.senderTypeCode === "BUYER" ? (isZh ? "我" : "You") : "Shopa AI"}</strong>
                    <p>{item.contentText || (isZh ? "暂不支持该消息类型展示。" : "Unsupported message type.")}</p>
                    <small>{formatSentAt(item.sentAt, locale)}</small>
                  </div>
                </article>
              ))}
              {bootstrapMutation.isLoading || messagesQuery.isLoading ? <p className="tb-agent-empty">{isZh ? "正在连接智能客服..." : "Connecting assistant..."}</p> : null}
            </div>

            <form className="tb-agent-composer" onSubmit={handleSubmit}>
              <textarea
                value={draft}
                onChange={(event) => setDraft(event.target.value)}
                placeholder={isZh ? "输入商品、订单、退货或平台规则问题..." : "Ask about products, orders, returns or platform policies..."}
                rows={4}
              />
              <div className="tb-agent-composer-actions">
                <span>{isZh ? "图片上传入口下一步接 media-svc。" : "Image upload can be connected to media-svc next."}</span>
                <button type="submit" disabled={!conversationNo || !draft.trim() || sendMutation.isLoading}>
                  {sendMutation.isLoading ? (isZh ? "发送中..." : "Sending...") : isZh ? "发送给 AI" : "Send to AI"}
                </button>
              </div>
            </form>
          </div>

          <div className="tb-agent-panel tb-agent-source-panel">
            <div className="tb-agent-source-head">
              <strong>{isZh ? "本轮回答依据" : "Answer sources"}</strong>
              {currentRun?.promptInjectionFlag ? <span>{isZh ? "已触发安全防护" : "Safety guard triggered"}</span> : null}
            </div>
            {answerSources.length === 0 ? <p className="tb-agent-empty">{isZh ? "本轮暂无外部知识依据。" : "No explicit sources for this turn."}</p> : null}
            <div className="tb-agent-source-list">
              {answerSources.map((item) => (
                <article key={`${item.sourceTypeCode}-${item.sourceId}-${item.sourceVersion}`}>
                  <strong>{item.title || item.sourceTypeCode}</strong>
                  <p>{item.snippet}</p>
                  <small>{`${item.sourceTypeCode} / v${item.sourceVersion}`}</small>
                </article>
              ))}
            </div>
            {lastAssistantMessage && runNo ? (
              <div className="tb-agent-feedback-row">
                <button type="button" onClick={() => feedbackMutation.mutate("HELPFUL")} disabled={feedbackMutation.isLoading}>
                  {isZh ? "这条有帮助" : "Helpful"}
                </button>
                <button type="button" onClick={() => feedbackMutation.mutate("NOT_HELPFUL")} disabled={feedbackMutation.isLoading}>
                  {isZh ? "不太有帮助" : "Not helpful"}
                </button>
              </div>
            ) : null}
          </div>
        </div>
      </section>
    </main>
  );
}
