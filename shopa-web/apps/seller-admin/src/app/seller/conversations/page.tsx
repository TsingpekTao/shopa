"use client";

import { CSSProperties, FormEvent, KeyboardEvent, ReactNode, useEffect, useMemo, useRef, useState } from "react";
import Link from "next/link";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Alert, Avatar, Empty, Select, Skeleton, message } from "antd";
import { ReloadOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import { listSellerConversations, listSellerMessages, markSellerConversationRead, sendSellerMessage } from "@/features/chat/api";
import { buildSellerChatShopOptions } from "@/features/chat/shop-selection";
import type { SellerChatConversation, SellerChatMessage, SellerChatMessageCard } from "@/features/chat/types";
import { listSellerProducts } from "@/features/catalog/api";
import { getSellerCurrentShopNo, setSellerCurrentShopNo } from "@/features/navigation/nav-badges";
import { fetchSellerWorkbench } from "@/features/seller-shop/api";

const MALL_BASE_URL = "http://127.0.0.1:3100";
const SELLER_CHAT_LAST_SHOP_KEY = "seller-chat-last-shop-no";

type LocalReadState = {
  conversationNo: string;
  messageNo: string;
  lastMessageAt: string;
};

function formatListTime(value: string, locale: string): string {
  if (!value) {
    return "";
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleString(locale === "zh-CN" ? "zh-CN" : "en-US", {
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit"
  });
}

function sortMessages(list: SellerChatMessage[]): SellerChatMessage[] {
  return [...list].sort((a, b) => {
    const av = a.sentAt ? new Date(a.sentAt).getTime() : 0;
    const bv = b.sentAt ? new Date(b.sentAt).getTime() : 0;
    return av - bv;
  });
}

function fallbackPreview(text: string, isZh: boolean): string {
  if (text.trim()) {
    return text;
  }
  return isZh ? "暂无新消息" : "No messages yet";
}

function conversationTitle(conversation: SellerChatConversation | undefined, isZh: boolean): string {
  return conversation?.buyerDisplayName?.trim() || (isZh ? "买家" : "Buyer");
}

function conversationSceneLabel(conversation: SellerChatConversation, isZh: boolean): string {
  if (conversation.sceneCode === "AFTER_SALE") {
    return isZh ? "售后" : "After-sale";
  }
  if (conversation.orderNo) {
    return isZh ? "订单" : "Order";
  }
  return isZh ? "咨询" : "Consult";
}

function conversationSubtitle(conversation: SellerChatConversation, isZh: boolean): string {
  if (conversation.lastMessagePreview.trim()) {
    return conversation.lastMessagePreview.trim();
  }
  if (conversation.orderNo) {
    return isZh ? `订单 ${conversation.orderNo}` : `Order ${conversation.orderNo}`;
  }
  if (conversation.anchorSpuNo) {
    return isZh ? `商品 ${conversation.anchorSpuNo}` : `Product ${conversation.anchorSpuNo}`;
  }
  return isZh ? "继续沟通" : "Continue chatting";
}

function threadIntroText(conversation: SellerChatConversation | undefined, isZh: boolean): string {
  if (!conversation) {
    return isZh ? "请选择会话" : "Select a conversation";
  }
  if (conversation.orderNo) {
    return isZh ? `订单 ${conversation.orderNo}` : `Order ${conversation.orderNo}`;
  }
  if (conversation.anchorSpuNo) {
    return isZh ? `商品 ${conversation.anchorSpuNo}` : `Product ${conversation.anchorSpuNo}`;
  }
  return isZh ? "买家会话" : "Buyer conversation";
}

function senderTitle(messageRow: SellerChatMessage, conversation: SellerChatConversation | undefined, shopName: string, isZh: boolean): string {
  if (messageRow.senderDisplayName.trim()) {
    return messageRow.senderDisplayName.trim();
  }
  if (messageRow.senderType === "SELLER") {
    return shopName || (isZh ? "商家" : "Seller");
  }
  if (messageRow.senderType === "BUYER") {
    return conversationTitle(conversation, isZh);
  }
  return isZh ? "系统通知" : "System";
}

function readStateText(messageRow: SellerChatMessage, isZh: boolean): string {
  if (messageRow.senderType !== "SELLER") {
    return "";
  }
  return messageRow.peerRead ? (isZh ? "买家已读" : "Read by buyer") : isZh ? "未读" : "Unread";
}

function resolveMallLink(rawLink?: string): string {
  const link = rawLink?.trim() ?? "";
  if (!link) {
    return "";
  }
  if (/^https?:\/\//i.test(link)) {
    return link;
  }
  if (!link.startsWith("/")) {
    return `${MALL_BASE_URL}/${link}`;
  }
  return `${MALL_BASE_URL}${link}`;
}

function resolveSellerOrderLink(cardSubOrderNo?: string, conversationSubOrderNo?: string): string {
  const subOrderNo = cardSubOrderNo?.trim() || conversationSubOrderNo?.trim() || "";
  return subOrderNo ? `/seller/orders/${encodeURIComponent(subOrderNo)}` : "";
}

function normalizeRemoteAssetUrl(raw?: string): string {
  const value = raw?.trim() ?? "";
  if (!value) {
    return "";
  }
  try {
    const parsed = new URL(value);
    parsed.pathname = parsed.pathname
      .split("/")
      .map((segment, index) => (index === 0 ? segment : normalizeRemoteAssetSegment(segment)))
      .join("/");
    return parsed.toString();
  } catch {
    return encodeURI(value);
  }
}

function normalizeRemoteAssetSegment(segment: string): string {
  if (!segment) {
    return segment;
  }
  try {
    return encodeURIComponent(decodeURIComponent(segment));
  } catch {
    return encodeURIComponent(segment);
  }
}

function renderSellerCard(
  card: SellerChatMessageCard,
  contentText: string,
  isZh: boolean,
  conversation?: SellerChatConversation
): ReactNode {
  const note = contentText.trim();

  if (card.kind === "ORDER_CARD") {
    const order = card.order;
    const href = resolveSellerOrderLink(order.subOrderNo, conversation?.subOrderNo);
    const body = (
      <div className="tb-chat-card">
        <div className="tb-chat-card-head">{isZh ? "订单详情" : "Order Detail"}</div>
        <div className="tb-chat-card-body">
          <div className="tb-chat-card-main">
            <strong>{order.orderNo ? `${isZh ? "订单号：" : "Order No: "}${order.orderNo}` : isZh ? "订单卡片" : "Order card"}</strong>
            {order.subOrderNo ? <span>{`Sub Order No: ${order.subOrderNo}`}</span> : null}
            {order.shopNo ? <span>{`Shop: ${order.shopNo}`}</span> : null}
            {order.status ? <span>{order.status}</span> : null}
            {order.summary ? <span>{order.summary}</span> : null}
            {order.totalAmount ? <small>{order.currency ? `${order.currency} ${order.totalAmount}` : order.totalAmount}</small> : null}
            {order.items?.length ? (
              <ul className="tb-chat-card-list">
                {order.items.map((item, index) => (
                  <li key={`${order.orderNo || "order"}-${index}`}>
                    <strong>{item.title || (isZh ? "商品" : "Item")}</strong>
                    <span>{[item.quantity ? `${isZh ? "数量" : "Qty"} ${item.quantity}` : "", item.price ? `${isZh ? "价格" : "Price"} ${item.price}` : "", item.meta || ""].filter(Boolean).join(" / ")}</span>
                  </li>
                ))}
              </ul>
            ) : null}
            {note ? <small className="tb-chat-card-note">{note}</small> : null}
            {href ? <a href={href} target="_blank" rel="noreferrer">{isZh ? "查看订单" : "View Order"}</a> : null}
          </div>
        </div>
      </div>
    );
    return body;
  }

  const product = card.product;
  const href = resolveMallLink(product.link);
  const normalizedCoverImageUrl = normalizeRemoteAssetUrl(product.coverImageUrl);
  const thumbStyle: CSSProperties | undefined = product.coverImageUrl
    ? {
        backgroundImage: `url("${normalizedCoverImageUrl}")`,
        backgroundSize: "cover",
        backgroundPosition: "center"
      }
    : undefined;

  const body = (
    <div className="tb-chat-card">
      <div className="tb-chat-card-head">{isZh ? "商品详情" : "Product Detail"}</div>
      <div className="tb-chat-card-body">
        <div className="tb-chat-card-thumb" style={thumbStyle} />
        <div className="tb-chat-card-main">
          <strong>{product.title || (isZh ? "商品卡片" : "Product card")}</strong>
          <span>{product.subtitle || product.skuNo || product.spuNo || (isZh ? "商品信息" : "Product info")}</span>
          {product.price ? <small>{product.priceLabel ? `${product.priceLabel} ${product.price}` : product.price}</small> : null}
          {note ? <small className="tb-chat-card-note">{note}</small> : null}
          {href ? <a href={href} target="_blank" rel="noreferrer">{isZh ? "查看商品" : "View Product"}</a> : null}
        </div>
      </div>
    </div>
  );

  return href ? (
    <a href={href} target="_blank" rel="noreferrer" className="tb-chat-card-link">
      {body}
    </a>
  ) : (
    body
  );
}

export default function SellerConversationsPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const queryClient = useQueryClient();
  const [messageApi, contextHolder] = message.useMessage();
  const threadRef = useRef<HTMLDivElement | null>(null);
  const lastReadRef = useRef("");

  const [selectedShopNo, setSelectedShopNo] = useState("");
  const [currentConversationNo, setCurrentConversationNo] = useState("");
  const [draft, setDraft] = useState("");
  const [persistedShopNo, setPersistedShopNo] = useState("");
  const [locallyReadState, setLocallyReadState] = useState<LocalReadState>({
    conversationNo: "",
    messageNo: "",
    lastMessageAt: ""
  });

  const workbenchQuery = useQuery({
    queryKey: ["seller-chat-workbench"],
    queryFn: fetchSellerWorkbench,
    staleTime: 30_000
  });

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    setPersistedShopNo(getSellerCurrentShopNo() || window.localStorage.getItem(SELLER_CHAT_LAST_SHOP_KEY)?.trim() || "");
  }, []);

  const workbenchShops = useMemo(() => {
    return (workbenchQuery.data?.shops ?? [])
      .map((shop) => ({
        shopNo: String(shop.shopNo ?? "").trim(),
        shopName: String(shop.shopDisplayName || shop.shopName || shop.shopNo || "").trim(),
        shopStatusCode: String(shop.shopStatusCode ?? "").trim()
      }))
      .filter((shop) => shop.shopNo);
  }, [workbenchQuery.data?.shops]);

  const fallbackCatalogShopsQuery = useQuery({
    queryKey: ["seller-chat-fallback-shops"],
    queryFn: async () => {
      const result = await listSellerProducts({ page: 1, pageSize: 100 });
      return Array.from(new Set(result.products.map((product) => product.shopNo.trim()).filter(Boolean)));
    },
    enabled: workbenchShops.length === 0 || Boolean(workbenchQuery.data?.partial),
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const shops = useMemo(() => {
    return buildSellerChatShopOptions({
      workbenchShops,
      catalogShopNos: fallbackCatalogShopsQuery.data ?? [],
      persistedShopNo
    });
  }, [fallbackCatalogShopsQuery.data, persistedShopNo, workbenchShops]);

  const normalizedSelectedShopNo = selectedShopNo.trim();

  useEffect(() => {
    if (!shops.length) {
      setSelectedShopNo("");
      return;
    }
    setSelectedShopNo((current) => {
      if (current && shops.some((shop) => shop.shopNo === current)) {
        return current;
      }
      return shops[0].shopNo;
    });
  }, [shops]);

  useEffect(() => {
    if (!normalizedSelectedShopNo || typeof window === "undefined") {
      return;
    }
    setSellerCurrentShopNo(normalizedSelectedShopNo);
    window.localStorage.setItem(SELLER_CHAT_LAST_SHOP_KEY, normalizedSelectedShopNo);
    setPersistedShopNo(normalizedSelectedShopNo);
  }, [normalizedSelectedShopNo]);

  const selectedShop = useMemo(() => shops.find((shop) => shop.shopNo === normalizedSelectedShopNo), [normalizedSelectedShopNo, shops]);

  const conversationsQuery = useQuery({
    queryKey: ["seller-chat-conversations", normalizedSelectedShopNo],
    queryFn: () => listSellerConversations(normalizedSelectedShopNo, 60, ""),
    enabled: Boolean(normalizedSelectedShopNo),
    staleTime: 10_000,
    refetchInterval: 12_000
  });

  const conversations = conversationsQuery.data?.list ?? [];

  useEffect(() => {
    setCurrentConversationNo("");
    setLocallyReadState({ conversationNo: "", messageNo: "", lastMessageAt: "" });
    lastReadRef.current = "";
  }, [normalizedSelectedShopNo]);

  const displayConversations = useMemo(() => {
    return conversations.map((conversation) => {
      const shouldHideUnread =
        conversation.conversationNo === locallyReadState.conversationNo &&
        locallyReadState.lastMessageAt &&
        conversation.lastMessageAt === locallyReadState.lastMessageAt;
      return shouldHideUnread && conversation.unreadCount > 0 ? { ...conversation, unreadCount: 0 } : conversation;
    });
  }, [conversations, locallyReadState]);

  useEffect(() => {
    if (!displayConversations.length) {
      setCurrentConversationNo("");
      return;
    }
    setCurrentConversationNo((current) => {
      if (current && displayConversations.some((item) => item.conversationNo === current)) {
        return current;
      }
      return displayConversations[0].conversationNo;
    });
  }, [displayConversations]);

  const normalizedConversationNo = currentConversationNo.trim();

  const currentConversation = useMemo<SellerChatConversation | undefined>(() => {
    return displayConversations.find((item) => item.conversationNo === normalizedConversationNo);
  }, [displayConversations, normalizedConversationNo]);

  const messagesQuery = useQuery({
    queryKey: ["seller-chat-messages", normalizedSelectedShopNo, normalizedConversationNo],
    queryFn: () => listSellerMessages(normalizedSelectedShopNo, normalizedConversationNo, 100, ""),
    enabled: Boolean(normalizedSelectedShopNo) && Boolean(normalizedConversationNo),
    staleTime: 5_000,
    refetchInterval: normalizedConversationNo ? 8_000 : false
  });

  const conversationMessages = useMemo(() => sortMessages(messagesQuery.data?.list ?? []), [messagesQuery.data?.list]);

  useEffect(() => {
    if (!threadRef.current) {
      return;
    }
    threadRef.current.scrollTop = threadRef.current.scrollHeight;
  }, [conversationMessages]);

  useEffect(() => {
    const lastMessage = conversationMessages[conversationMessages.length - 1];
    if (!normalizedConversationNo || !lastMessage?.messageNo) {
      return;
    }
    const readKey = `${normalizedConversationNo}:${lastMessage.messageNo}`;
    if (lastReadRef.current === readKey) {
      return;
    }
    lastReadRef.current = readKey;
    const snapshotLastMessageAt = currentConversation?.lastMessageAt || lastMessage.sentAt;
    setLocallyReadState({
      conversationNo: normalizedConversationNo,
      messageNo: lastMessage.messageNo,
      lastMessageAt: snapshotLastMessageAt
    });
    markSellerConversationRead({
      conversationNo: normalizedConversationNo,
      readToMessageNo: lastMessage.messageNo
    })
      .then(() => queryClient.invalidateQueries({ queryKey: ["seller-chat-conversations", normalizedSelectedShopNo] }))
      .catch(() => {
        lastReadRef.current = "";
        setLocallyReadState((current) =>
          current.conversationNo === normalizedConversationNo ? { conversationNo: "", messageNo: "", lastMessageAt: "" } : current
        );
      });
  }, [conversationMessages, currentConversation, normalizedConversationNo, normalizedSelectedShopNo, queryClient]);

  const sendMutation = useMutation({
    mutationFn: async () =>
      sendSellerMessage({
        conversationNo: normalizedConversationNo,
        contentText: draft.trim()
      }),
    onSuccess: async () => {
      setDraft("");
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ["seller-chat-messages", normalizedSelectedShopNo, normalizedConversationNo] }),
        queryClient.invalidateQueries({ queryKey: ["seller-chat-conversations", normalizedSelectedShopNo] })
      ]);
    },
    onError: (error) => {
      const errorMessage = error instanceof Error ? error.message : isZh ? "发送消息失败" : "Send message failed";
      messageApi.error(errorMessage);
    }
  });

  const unreadTotal = useMemo(() => displayConversations.reduce((sum, item) => sum + item.unreadCount, 0), [displayConversations]);

  function submitMessage() {
    if (!normalizedConversationNo) {
      messageApi.warning(isZh ? "请先选择一个会话" : "Please select a conversation first");
      return;
    }
    if (!draft.trim()) {
      return;
    }
    sendMutation.mutate();
  }

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    submitMessage();
  }

  function handleTextareaEnter(event: KeyboardEvent<HTMLTextAreaElement>) {
    if (event.shiftKey) {
      return;
    }
    event.preventDefault();
    submitMessage();
  }

  if (workbenchQuery.isLoading) {
    return <Skeleton active paragraph={{ rows: 10 }} />;
  }

  return (
    <section className="tb-chat-page tb-chat-page-seller">
      {contextHolder}

      <header className="tb-chat-header">
        <div>
          <h1>{isZh ? "消息中心" : "Message Center"}</h1>
        </div>
        <div className="tb-chat-header-actions">
          <Select
            value={selectedShopNo || undefined}
            onChange={(value) => setSelectedShopNo(value.trim())}
            className="tb-chat-shop-select"
            placeholder={isZh ? "选择店铺" : "Select shop"}
            options={shops.map((shop) => ({
              value: shop.shopNo,
              label: shop.shopName
            }))}
          />
          <button
            type="button"
            className="tb-chat-header-btn"
            onClick={() => {
              void Promise.all([workbenchQuery.refetch(), conversationsQuery.refetch(), messagesQuery.refetch()]);
            }}
          >
            <ReloadOutlined />
            <span>{isZh ? "刷新" : "Refresh"}</span>
          </button>
          <Link href="/seller/workbench">{isZh ? "返回工作台" : "Back"}</Link>
        </div>
      </header>

      {workbenchQuery.data?.partial ? (
        <Alert
          type="warning"
          showIcon
          message={isZh ? "部分店铺数据暂时不可用" : "Some shop data is temporarily unavailable"}
          description={(workbenchQuery.data.degradedFields ?? []).join(", ") || undefined}
        />
      ) : null}

      {!shops.length ? (
        <div className="tb-chat-empty-panel">
          <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={isZh ? "当前账号下还没有可接待消息的店铺" : "No shop is available for seller chat yet"} />
        </div>
      ) : (
        <div className="tb-chat-layout">
          <aside className="tb-chat-sidebar">
            <div className="tb-chat-sidebar-top">
              <h2>{isZh ? "会话列表" : "Conversations"}</h2>
              {unreadTotal > 0 ? <span className="tb-chat-unread-pill">{unreadTotal > 99 ? "99+" : unreadTotal}</span> : null}
            </div>

            <div className="tb-chat-sidebar-status">
              <strong>{selectedShop?.shopName || (isZh ? "当前店铺" : "Current shop")}</strong>
              <span>{selectedShop?.shopStatusCode || "UNKNOWN"}</span>
            </div>

            {displayConversations.length === 0 && !conversationsQuery.isLoading ? (
              <p className="tb-chat-empty">{isZh ? "这个店铺暂时还没有买家会话" : "No buyer conversation for this shop yet"}</p>
            ) : null}

            <div className="tb-chat-conv-list">
              {conversationsQuery.isLoading ? (
                <Skeleton active paragraph={{ rows: 6 }} title={false} />
              ) : (
                displayConversations.map((conversation) => {
                  const active = conversation.conversationNo === normalizedConversationNo;
                  const title = conversationTitle(conversation, isZh);
                  return (
                    <button
                      type="button"
                      key={conversation.conversationNo}
                      className={`tb-chat-conv-item ${active ? "is-active" : ""}`}
                      onClick={() => {
                        lastReadRef.current = "";
                        setCurrentConversationNo(conversation.conversationNo.trim());
                      }}
                    >
                      <div className="tb-chat-conv-avatar">
                        <Avatar src={conversation.buyerAvatarUrl || undefined} size={42}>
                          {title.slice(0, 1).toUpperCase()}
                        </Avatar>
                      </div>

                      <div className="tb-chat-conv-body">
                        <div className="tb-chat-conv-head">
                          <strong>{title}</strong>
                          <time>{formatListTime(conversation.lastMessageAt, locale)}</time>
                        </div>

                        <div className="tb-chat-conv-meta">
                          <p>{fallbackPreview(conversationSubtitle(conversation, isZh), isZh)}</p>
                          {conversation.unreadCount > 0 ? <em>{conversation.unreadCount > 99 ? "99+" : conversation.unreadCount}</em> : null}
                        </div>
                      </div>
                    </button>
                  );
                })
              )}
            </div>
          </aside>

          <div className="tb-chat-main">
            <div className="tb-chat-main-head">
              <div className="tb-chat-main-head-title">
                <Avatar src={currentConversation?.buyerAvatarUrl || undefined} size={42}>
                  {conversationTitle(currentConversation, isZh).slice(0, 1).toUpperCase()}
                </Avatar>
                <div>
                  <h2>{currentConversation ? conversationTitle(currentConversation, isZh) : isZh ? "请选择会话" : "Select Conversation"}</h2>
                  <span>{threadIntroText(currentConversation, isZh)}</span>
                </div>
              </div>
            </div>

            <div className="tb-chat-messages" ref={threadRef}>
              {messagesQuery.isLoading ? (
                <Skeleton active paragraph={{ rows: 8 }} title={false} />
              ) : !currentConversation ? (
                <p className="tb-chat-empty">{isZh ? "请选择左侧会话" : "Select a conversation"}</p>
              ) : conversationMessages.length === 0 ? (
                <p className="tb-chat-empty">{isZh ? "当前会话还没有消息" : "This conversation has no messages yet"}</p>
              ) : (
                conversationMessages.map((item) => {
                  const self = item.senderType === "SELLER";
                  const avatarLabel = senderTitle(item, currentConversation, selectedShop?.shopName || "", isZh).slice(0, 1).toUpperCase();
                  const stateText = readStateText(item, isZh);
                  return (
                    <article key={item.messageNo} className={`tb-chat-msg ${self ? "is-self" : ""}`}>
                      {!self ? (
                        <Avatar src={item.senderAvatarUrl || currentConversation?.buyerAvatarUrl || undefined} size={36} className="tb-chat-msg-avatar">
                          {avatarLabel}
                        </Avatar>
                      ) : null}

                      <div className="tb-chat-bubble">
                        <div className="tb-chat-bubble-head">{senderTitle(item, currentConversation, selectedShop?.shopName || "", isZh)}</div>
                        {item.card ? renderSellerCard(item.card, item.contentText, isZh, currentConversation) : <p>{item.contentText || (isZh ? "暂不支持该消息类型" : "Unsupported message type")}</p>}
                        <small>
                          {formatListTime(item.sentAt, locale)}
                          {stateText ? ` / ${stateText}` : ""}
                        </small>
                      </div>

                      {self ? (
                        <Avatar src={item.senderAvatarUrl || currentConversation?.shopAvatarUrl || undefined} size={36} className="tb-chat-msg-avatar tb-chat-msg-avatar-self">
                          {avatarLabel}
                        </Avatar>
                      ) : null}
                    </article>
                  );
                })
              )}
            </div>

            <form className="tb-chat-composer" onSubmit={handleSubmit}>
              <textarea
                value={draft}
                onChange={(event) => setDraft(event.target.value)}
                onKeyDown={handleTextareaEnter}
                placeholder={isZh ? "输入消息..." : "Type a message..."}
                rows={3}
                disabled={!currentConversation || !normalizedSelectedShopNo}
              />
              <div className="tb-chat-composer-actions">
                <span className="tb-chat-composer-tip">{currentConversation ? (isZh ? "Shift + Enter 换行" : "Shift + Enter for a new line") : isZh ? "请先选择会话" : "Choose a conversation first"}</span>
                <button type="submit" disabled={!currentConversation || !normalizedSelectedShopNo || !draft.trim() || sendMutation.isLoading}>
                  {isZh ? "发送" : "Send"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </section>
  );
}
