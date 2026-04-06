"use client";

import { FormEvent, useEffect, useMemo, useRef, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Avatar, message } from "antd";
import { useI18n } from "@shopa/ui";
import { createOrGetConversation, listBuyerConversations, listMessages, markConversationRead, sendBuyerMessage } from "@/features/chat/api";
import { ChatConversation, ChatMessage } from "@/features/chat/types";

const SHELL_UNREAD_QUERY_KEY = ["shell", "unread"] as const;
const BUYER_CONVERSATIONS_QUERY_KEY = ["chat-buyer-conversations"] as const;

function readMallAccessToken(): string {
  if (typeof window === "undefined") {
    return "";
  }
  const direct = window.localStorage.getItem("shopa_mall_access_token")?.trim() ?? "";
  if (direct) {
    return direct;
  }
  try {
    const raw = window.localStorage.getItem("shopa-mall-auth");
    if (!raw) {
      return "";
    }
    const parsed = JSON.parse(raw) as {
      state?: { tokenPair?: { accessToken?: string } };
    };
    return parsed?.state?.tokenPair?.accessToken?.trim() ?? "";
  } catch {
    return "";
  }
}

function formatListTime(raw: string, locale: string): string {
  if (!raw) {
    return "";
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

function conversationTitle(conversation: ChatConversation | undefined, isZh: boolean): string {
  return conversation?.shopName?.trim() || (isZh ? "\u5e97\u94fa\u5ba2\u670d" : "Shop Support");
}

function conversationSubtitle(conversation: ChatConversation, isZh: boolean): string {
  if (conversation.lastMessagePreview.trim()) {
    return conversation.lastMessagePreview.trim();
  }
  return isZh
    ? "\u70b9\u51fb\u8fdb\u5165\u540e\u5c31\u53ef\u4ee5\u76f4\u63a5\u5f00\u59cb\u54a8\u8be2\u5546\u54c1\u7ec6\u8282\u3002"
    : "Open the thread to start chatting about this product.";
}

function senderTitle(messageRow: ChatMessage, conversation: ChatConversation | undefined, isZh: boolean): string {
  if (messageRow.senderDisplayName.trim()) {
    return messageRow.senderDisplayName.trim();
  }
  if (messageRow.senderType === "BUYER") {
    return isZh ? "\u6211" : "Me";
  }
  if (messageRow.senderType === "SELLER") {
    return conversationTitle(conversation, isZh);
  }
  return isZh ? "\u7cfb\u7edf\u901a\u77e5" : "System";
}

function messageAvatarLabel(messageRow: ChatMessage, conversation: ChatConversation | undefined, isZh: boolean): string {
  const title = senderTitle(messageRow, conversation, isZh);
  return title.slice(0, 1).toUpperCase();
}

function peerReadText(messageRow: ChatMessage, isZh: boolean): string {
  if (messageRow.senderType !== "BUYER") {
    return "";
  }
  if (messageRow.peerRead) {
    return isZh ? "\u5546\u5bb6\u5df2\u8bfb" : "Read by seller";
  }
  return isZh ? "\u672a\u8bfb" : "Unread";
}

export default function MallMessagesPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const router = useRouter();
  const queryClient = useQueryClient();
  const [messageApi, contextHolder] = message.useMessage();
  const searchParams = useSearchParams();
  const lastMarkedReadRef = useRef("");
  const autoCreateKeyRef = useRef("");

  const defaultShopNo = searchParams.get("shop_no") ?? "";
  const defaultSpuNo = searchParams.get("spu_no") ?? "";
  const defaultSkuNo = searchParams.get("sku_no") ?? "";
  const fromOrderNo = searchParams.get("order_no") ?? "";
  const mallAccessToken = readMallAccessToken();

  const [currentConversationNo, setCurrentConversationNo] = useState("");
  const [draft, setDraft] = useState("");

  const autoCreateKey = defaultShopNo ? [defaultShopNo, defaultSpuNo, defaultSkuNo].join("::") : "";

  const conversationsQuery = useQuery({
    queryKey: BUYER_CONVERSATIONS_QUERY_KEY,
    queryFn: () => listBuyerConversations(30, ""),
    enabled: Boolean(mallAccessToken),
    staleTime: 15_000,
    refetchInterval: mallAccessToken ? 12_000 : false,
    refetchOnWindowFocus: false
  });

  const createConversationMutation = useMutation({
    mutationFn: () =>
      createOrGetConversation({
        shopNo: defaultShopNo,
        anchorSpuNo: defaultSpuNo,
        anchorSkuNo: defaultSkuNo
      }),
    onSuccess: (conversation) => {
      if (conversation.conversationNo) {
        setCurrentConversationNo(conversation.conversationNo);
      }
      void queryClient.invalidateQueries({ queryKey: BUYER_CONVERSATIONS_QUERY_KEY });
      void queryClient.invalidateQueries({ queryKey: SHELL_UNREAD_QUERY_KEY });
    },
    onError: () => {
      messageApi.error(isZh ? "\u521b\u5efa\u4f1a\u8bdd\u5931\u8d25\uff0c\u8bf7\u7a0d\u540e\u91cd\u8bd5" : "Failed to create conversation");
    }
  });

  useEffect(() => {
    if (!defaultShopNo) {
      return;
    }
    if (!mallAccessToken) {
      messageApi.info(isZh ? "\u8bf7\u5148\u767b\u5f55\u540e\u518d\u54a8\u8be2\u5546\u5bb6" : "Please sign in before chatting with the seller");
      router.replace("/login");
      return;
    }
    if (createConversationMutation.isPending) {
      return;
    }
    if (autoCreateKeyRef.current === autoCreateKey) {
      return;
    }
    autoCreateKeyRef.current = autoCreateKey;
    setCurrentConversationNo("");
    createConversationMutation.reset();
    createConversationMutation.mutate();
  }, [autoCreateKey, createConversationMutation, defaultShopNo, isZh, mallAccessToken, messageApi, router]);

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

  useEffect(() => {
    lastMarkedReadRef.current = "";
  }, [currentConversationNo]);

  const currentConversation = useMemo<ChatConversation | undefined>(() => {
    if (!currentConversationNo) {
      return undefined;
    }
    return conversations.find((item) => item.conversationNo === currentConversationNo);
  }, [conversations, currentConversationNo]);

  const messagesQuery = useQuery({
    queryKey: ["chat-buyer-messages", currentConversationNo],
    queryFn: () => listMessages(currentConversationNo, 100, ""),
    enabled: Boolean(mallAccessToken) && currentConversationNo.length > 0,
    staleTime: 8_000,
    refetchInterval: currentConversationNo ? 8_000 : false
  });

  const conversationMessages = sortMessages(messagesQuery.data?.list ?? []);

  useEffect(() => {
    const last = conversationMessages[conversationMessages.length - 1];
    if (!currentConversationNo || !last?.messageNo) {
      return;
    }
    const nextReadKey = `${currentConversationNo}:${last.messageNo}`;
    if (lastMarkedReadRef.current === nextReadKey) {
      return;
    }
    lastMarkedReadRef.current = nextReadKey;

    markConversationRead({
      conversationNo: currentConversationNo,
      readToMessageNo: last.messageNo
    })
      .then(() => {
        void queryClient.invalidateQueries({ queryKey: BUYER_CONVERSATIONS_QUERY_KEY });
        void queryClient.invalidateQueries({ queryKey: SHELL_UNREAD_QUERY_KEY });
      })
      .catch(() => {
        lastMarkedReadRef.current = "";
      });
  }, [conversationMessages, currentConversationNo, queryClient]);

  const sendMutation = useMutation({
    mutationFn: async () =>
      sendBuyerMessage({
        conversationNo: currentConversationNo,
        contentText: draft.trim()
      }),
    onSuccess: async () => {
      setDraft("");
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ["chat-buyer-messages", currentConversationNo] }),
        queryClient.invalidateQueries({ queryKey: BUYER_CONVERSATIONS_QUERY_KEY }),
        queryClient.invalidateQueries({ queryKey: SHELL_UNREAD_QUERY_KEY })
      ]);
    },
    onError: () => {
      messageApi.error(isZh ? "\u53d1\u9001\u5931\u8d25\uff0c\u8bf7\u7a0d\u540e\u91cd\u8bd5" : "Send failed");
    }
  });

  const emptyHint = defaultShopNo
    ? isZh
      ? "\u6b63\u5728\u4e3a\u4f60\u8fde\u63a5\u5546\u5bb6\u4f1a\u8bdd..."
      : "Connecting you to the shop conversation\u2026"
    : isZh
      ? "\u5f53\u524d\u8fd8\u6ca1\u6709\u4f1a\u8bdd\uff0c\u53bb\u5546\u54c1\u8be6\u60c5\u9875\u70b9\u51fb\u201c\u54a8\u8be2\u5546\u5bb6\u201d\u5c31\u80fd\u53d1\u8d77\u65b0\u5bf9\u8bdd\u3002"
      : "No conversation yet. Open a product page and click 'Consult Seller'.";

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    if (!currentConversationNo) {
      messageApi.warning(isZh ? "\u8bf7\u5148\u9009\u62e9\u4e00\u4e2a\u4f1a\u8bdd" : "Please select a conversation first");
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
          <h1>{isZh ? "\u6d88\u606f\u4e2d\u5fc3" : "Message Center"}</h1>
          <p>
            {fromOrderNo
              ? isZh
                ? `\u5f53\u524d\u6b63\u5728\u5904\u7406\u8ba2\u5355 ${fromOrderNo} \u7684\u6c9f\u901a\u6d88\u606f`
                : `Currently chatting about order ${fromOrderNo}`
              : isZh
                ? "\u50cf\u804a\u5929\u5de5\u5177\u4e00\u6837\u76f4\u63a5\u548c\u5e97\u94fa\u6c9f\u901a\uff0c\u672a\u8bfb\u6570\u548c\u5df2\u8bfb\u72b6\u6001\u4f1a\u5b9e\u65f6\u540c\u6b65\u3002"
                : "Chat with shops in a messenger-style layout with live unread and read states."}
          </p>
        </div>
        <Link href="/">{isZh ? "\u8fd4\u56de\u9996\u9875" : "Back Home"}</Link>
      </header>

      <div className="tb-chat-layout">
        <aside className="tb-chat-sidebar">
          <h2>{isZh ? "\u4f1a\u8bdd\u5217\u8868" : "Conversations"}</h2>
          {conversations.length === 0 && !conversationsQuery.isLoading ? <p className="tb-chat-empty">{emptyHint}</p> : null}

          <div className="tb-chat-conv-list">
            {conversations.map((conversation) => {
              const active = conversation.conversationNo === currentConversationNo;
              const title = conversationTitle(conversation, isZh);
              return (
                <button
                  type="button"
                  key={conversation.conversationNo}
                  className={`tb-chat-conv-item ${active ? "is-active" : ""}`}
                  onClick={() => setCurrentConversationNo(conversation.conversationNo)}
                >
                  <div className="tb-chat-conv-avatar">
                    <Avatar src={conversation.shopAvatarUrl || undefined} size={42}>
                      {title.slice(0, 1).toUpperCase()}
                    </Avatar>
                  </div>

                  <div className="tb-chat-conv-body">
                    <div className="tb-chat-conv-head">
                      <strong>{title}</strong>
                      <time>{formatListTime(conversation.lastMessageAt, locale)}</time>
                    </div>

                    <div className="tb-chat-conv-meta">
                      <p>{conversationSubtitle(conversation, isZh)}</p>
                      {conversation.unreadCount > 0 ? <em>{conversation.unreadCount > 99 ? "99+" : conversation.unreadCount}</em> : null}
                    </div>
                  </div>
                </button>
              );
            })}
          </div>
        </aside>

        <div className="tb-chat-main">
          <div className="tb-chat-main-head">
            <div className="tb-chat-main-head-title">
              <Avatar src={currentConversation?.shopAvatarUrl || undefined} size={42}>
                {conversationTitle(currentConversation, isZh).slice(0, 1).toUpperCase()}
              </Avatar>
              <div>
                <h2>{currentConversation ? conversationTitle(currentConversation, isZh) : isZh ? "\u8bf7\u9009\u62e9\u4f1a\u8bdd" : "Select Conversation"}</h2>
                <span>{fromOrderNo ? (isZh ? "\u8ba2\u5355\u6c9f\u901a\u4e2d" : "Order conversation") : isZh ? "\u5546\u54c1\u54a8\u8be2\u4f1a\u8bdd" : "Product consultation"}</span>
              </div>
            </div>
          </div>

          <div className="tb-chat-messages">
            {conversationMessages.map((msg) => {
              const self = msg.senderType === "BUYER";
              const avatarSrc = self ? msg.senderAvatarUrl : currentConversation?.shopAvatarUrl || msg.senderAvatarUrl;
              return (
                <article key={msg.messageNo} className={`tb-chat-msg ${self ? "is-self" : ""}`}>
                  {!self ? (
                    <Avatar src={avatarSrc || undefined} size={36} className="tb-chat-msg-avatar">
                      {messageAvatarLabel(msg, currentConversation, isZh)}
                    </Avatar>
                  ) : null}

                  <div className="tb-chat-bubble">
                    <div className="tb-chat-bubble-head">{senderTitle(msg, currentConversation, isZh)}</div>
                    <p>{msg.contentText || (isZh ? "\u6682\u4e0d\u652f\u6301\u8be5\u6d88\u606f\u7c7b\u578b\u9884\u89c8" : "Unsupported message type")}</p>
                    <small>
                      {formatListTime(msg.sentAt, locale)}
                      {peerReadText(msg, isZh) ? ` · ${peerReadText(msg, isZh)}` : ""}
                    </small>
                  </div>

                  {self ? (
                    <Avatar src={avatarSrc || undefined} size={36} className="tb-chat-msg-avatar tb-chat-msg-avatar-self">
                      {messageAvatarLabel(msg, currentConversation, isZh)}
                    </Avatar>
                  ) : null}
                </article>
              );
            })}

            {messagesQuery.isLoading ? <p className="tb-chat-empty">{isZh ? "\u6b63\u5728\u52a0\u8f7d\u6d88\u606f..." : "Loading messages\u2026"}</p> : null}
          </div>

          <form className="tb-chat-composer" onSubmit={handleSubmit}>
            <textarea
              value={draft}
              onChange={(event) => setDraft(event.target.value)}
              placeholder={isZh ? "\u8f93\u5165\u4f60\u60f3\u54a8\u8be2\u7684\u95ee\u9898..." : "Type your message\u2026"}
              rows={3}
            />
            <button type="submit" disabled={!currentConversationNo || sendMutation.isPending || !draft.trim()}>
              {sendMutation.isPending ? (isZh ? "\u53d1\u9001\u4e2d..." : "Sending\u2026") : isZh ? "\u53d1\u9001" : "Send"}
            </button>
          </form>
        </div>
      </div>
    </section>
  );
}

