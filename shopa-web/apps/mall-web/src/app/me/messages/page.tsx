"use client";

import { FormEvent, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import { createOrGetConversation, listBuyerConversations, listMessages, markConversationRead, sendBuyerMessage } from "@/features/chat/api";
import { ChatConversation, ChatMessage } from "@/features/chat/types";

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

function sortMessages(list: ChatMessage[]): ChatMessage[] {
  return [...list].sort((a, b) => {
    const av = a.sentAt ? new Date(a.sentAt).getTime() : 0;
    const bv = b.sentAt ? new Date(b.sentAt).getTime() : 0;
    return av - bv;
  });
}

export default function MallMessagesPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const queryClient = useQueryClient();
  const [messageApi, contextHolder] = message.useMessage();
  const searchParams = useSearchParams();

  const defaultShopNo = searchParams.get("shop_no") ?? "";
  const defaultSpuNo = searchParams.get("spu_no") ?? "";
  const fromOrderNo = searchParams.get("order_no") ?? "";

  const [currentConversationNo, setCurrentConversationNo] = useState("");
  const [draft, setDraft] = useState("");

  const conversationsQuery = useQuery({
    queryKey: ["chat-buyer-conversations"],
    queryFn: () => listBuyerConversations(30, ""),
    staleTime: 15_000,
    refetchOnWindowFocus: false
  });

  const createConversationMutation = useMutation({
    mutationFn: () =>
      createOrGetConversation({
        shopNo: defaultShopNo,
        anchorSpuNo: defaultSpuNo
      }),
    onSuccess: (conversation) => {
      if (conversation.conversationNo) {
        setCurrentConversationNo(conversation.conversationNo);
      }
      queryClient.invalidateQueries({ queryKey: ["chat-buyer-conversations"] });
    },
    onError: () => {
      messageApi.error(isZh ? "创建会话失败，请稍后重试" : "Failed to create conversation");
    }
  });

  useEffect(() => {
    if (!defaultShopNo) {
      return;
    }
    if (createConversationMutation.isLoading || createConversationMutation.isSuccess) {
      return;
    }
    createConversationMutation.mutate();
  }, [createConversationMutation, defaultShopNo]);

  const conversations = conversationsQuery.data?.list ?? [];

  useEffect(() => {
    if (currentConversationNo) {
      return;
    }
    const first = conversations[0];
    if (first?.conversationNo) {
      setCurrentConversationNo(first.conversationNo);
    }
  }, [conversations, currentConversationNo]);

  const currentConversation = useMemo<ChatConversation | undefined>(() => {
    if (!currentConversationNo) {
      return undefined;
    }
    return conversations.find((item) => item.conversationNo === currentConversationNo);
  }, [conversations, currentConversationNo]);

  const messagesQuery = useQuery({
    queryKey: ["chat-buyer-messages", currentConversationNo],
    queryFn: () => listMessages(currentConversationNo, 100, ""),
    enabled: currentConversationNo.length > 0,
    staleTime: 8_000,
    refetchInterval: 8_000
  });

  const conversationMessages = sortMessages(messagesQuery.data?.list ?? []);

  useEffect(() => {
    const last = conversationMessages[conversationMessages.length - 1];
    if (!currentConversationNo || !last?.messageNo) {
      return;
    }
    markConversationRead({ conversationNo: currentConversationNo, readToMessageNo: last.messageNo }).catch(() => {
      // ignore read marker failures
    });
  }, [conversationMessages, currentConversationNo]);

  const sendMutation = useMutation({
    mutationFn: async () =>
      sendBuyerMessage({
        conversationNo: currentConversationNo,
        contentText: draft.trim()
      }),
    onSuccess: () => {
      setDraft("");
      queryClient.invalidateQueries({ queryKey: ["chat-buyer-messages", currentConversationNo] });
      queryClient.invalidateQueries({ queryKey: ["chat-buyer-conversations"] });
    },
    onError: () => {
      messageApi.error(isZh ? "发送失败，请稍后重试" : "Send failed");
    }
  });

  const emptyHint = defaultShopNo
    ? isZh
      ? "正在为你连接商家会话..."
      : "Connecting you to the shop conversation..."
    : isZh
      ? "当前没有会话，去商品详情页点击“咨询商家”即可发起。"
      : "No conversation yet. Open a product page and click 'Consult Seller'.";

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    if (!currentConversationNo) {
      messageApi.warning(isZh ? "请先选择会话" : "Please select a conversation first");
      return;
    }
    if (!draft.trim()) {
      return;
    }
    sendMutation.mutate();
  }

  return (
    <section className="tb-chat-page">
      {contextHolder}
      <header className="tb-chat-header">
        <div>
          <h1>{isZh ? "商家咨询" : "Shop Chat"}</h1>
          <p>
            {fromOrderNo
              ? isZh
                ? `售后沟通（订单 ${fromOrderNo}）`
                : `After-sale conversation (Order ${fromOrderNo})`
              : isZh
                ? "售前咨询与售后沟通统一入口"
                : "Unified pre-sale and after-sale conversation entry"}
          </p>
        </div>
        <Link href="/">{isZh ? "返回首页" : "Back Home"}</Link>
      </header>

      <div className="tb-chat-layout">
        <aside className="tb-chat-sidebar">
          <h2>{isZh ? "会话列表" : "Conversations"}</h2>
          {conversations.length === 0 && !conversationsQuery.isLoading ? <p className="tb-chat-empty">{emptyHint}</p> : null}
          <div className="tb-chat-conv-list">
            {conversations.map((conversation) => {
              const active = conversation.conversationNo === currentConversationNo;
              return (
                <button
                  type="button"
                  key={conversation.conversationNo}
                  className={`tb-chat-conv-item ${active ? "is-active" : ""}`}
                  onClick={() => setCurrentConversationNo(conversation.conversationNo)}
                >
                  <div>
                    <strong>{conversation.shopNo || "SHOP"}</strong>
                    {conversation.unreadCount > 0 ? <em>{conversation.unreadCount}</em> : null}
                  </div>
                  <p>{conversation.lastMessagePreview || (isZh ? "暂无消息" : "No messages yet")}</p>
                </button>
              );
            })}
          </div>
        </aside>

        <div className="tb-chat-main">
          <div className="tb-chat-main-head">
            <h2>{currentConversation?.shopNo || (isZh ? "请选择会话" : "Select Conversation")}</h2>
            <span>{currentConversation?.anchorSpuNo ? `${isZh ? "关联商品" : "Product"}: ${currentConversation.anchorSpuNo}` : ""}</span>
          </div>
          <div className="tb-chat-messages">
            {conversationMessages.map((msg) => (
              <article key={msg.messageNo} className={`tb-chat-msg ${msg.senderType === "BUYER" ? "is-self" : ""}`}>
                <div className="tb-chat-bubble">
                  <p>{msg.contentText || (isZh ? "（暂不支持该消息类型渲染）" : "(Unsupported message type)")}</p>
                  <small>{formatSentAt(msg.sentAt, locale)}</small>
                </div>
              </article>
            ))}
            {messagesQuery.isLoading ? <p className="tb-chat-empty">{isZh ? "正在加载消息..." : "Loading messages..."}</p> : null}
          </div>

          <form className="tb-chat-composer" onSubmit={handleSubmit}>
            <textarea
              value={draft}
              onChange={(event) => setDraft(event.target.value)}
              placeholder={isZh ? "输入你想咨询的问题..." : "Type your message..."}
              rows={3}
            />
            <button type="submit" disabled={!currentConversationNo || sendMutation.isLoading || !draft.trim()}>
              {sendMutation.isLoading ? (isZh ? "发送中..." : "Sending...") : isZh ? "发送" : "Send"}
            </button>
          </form>
        </div>
      </div>
    </section>
  );
}

