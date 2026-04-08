"use client";

import { FormEvent, useEffect, useMemo, useRef, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Avatar, message } from "antd";
import { useI18n } from "@shopa/ui";
import { getBuyerProductDetail, listBuyerProductImages } from "@/features/catalog/api";
import { createOrGetConversation, listBuyerConversations, listMessages, markConversationRead, sendBuyerMessage } from "@/features/chat/api";
import { ChatCardPayload, ChatConversation, ChatMessage } from "@/features/chat/types";
import { getMyOrderDetail } from "@/features/order/api";
import { formatCnyFromCents } from "@/lib/price";
import { pickProductImageBySpuNo } from "@/lib/product-images";

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
    const parsed = JSON.parse(raw) as { state?: { tokenPair?: { accessToken?: string } } };
    return parsed?.state?.tokenPair?.accessToken?.trim() ?? "";
  } catch {
    return "";
  }
}

function formatListTime(raw: string, locale: string): string {
  if (!raw) {
    return "";
  }
  const date = new Date(raw);
  if (Number.isNaN(date.getTime())) {
    return raw;
  }
  return date.toLocaleString(locale === "zh-CN" ? "zh-CN" : "en-US", {
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
  return conversation?.shopName?.trim() || (isZh ? "店铺客服" : "Shop Support");
}

function conversationSceneText(conversation: ChatConversation | undefined, isZh: boolean, routeOrderNo: string): string {
  if (routeOrderNo) {
    return isZh ? `订单 ${routeOrderNo}` : `Order ${routeOrderNo}`;
  }
  if (!conversation) {
    return isZh ? "请选择会话" : "Select a conversation";
  }
  if (conversation.sceneCode === "AFTER_SALE" && conversation.orderNo) {
    return isZh ? `售后订单 ${conversation.orderNo}` : `After-sale order ${conversation.orderNo}`;
  }
  if (conversation.anchorSpuNo) {
    return isZh ? "商品咨询" : "Product consultation";
  }
  return isZh ? "店铺会话" : "Shop conversation";
}

function conversationSubtitle(conversation: ChatConversation, isZh: boolean): string {
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

function senderTitle(messageRow: ChatMessage, conversation: ChatConversation | undefined, isZh: boolean): string {
  if (messageRow.senderDisplayName.trim()) {
    return messageRow.senderDisplayName.trim();
  }
  if (messageRow.senderType === "BUYER") {
    return isZh ? "我" : "Me";
  }
  if (messageRow.senderType === "SELLER") {
    return conversationTitle(conversation, isZh);
  }
  return isZh ? "系统通知" : "System";
}

function messageAvatarLabel(messageRow: ChatMessage, conversation: ChatConversation | undefined, isZh: boolean): string {
  return senderTitle(messageRow, conversation, isZh).slice(0, 1).toUpperCase();
}

function peerReadText(messageRow: ChatMessage, isZh: boolean): string {
  if (messageRow.senderType !== "BUYER") {
    return "";
  }
  return messageRow.peerRead ? (isZh ? "商家已读" : "Read by seller") : isZh ? "未读" : "Unread";
}

function parseCardPayload(messageRow: ChatMessage): ChatCardPayload | null {
  if (messageRow.messageType !== "PRODUCT_CARD" || !messageRow.extJson.trim()) {
    return null;
  }
  try {
    const parsed = JSON.parse(messageRow.extJson) as ChatCardPayload;
    if (parsed && (parsed.kind === "product" || parsed.kind === "order")) {
      return parsed;
    }
  } catch {
    return null;
  }
  return null;
}

function renderCardPayload(card: ChatCardPayload, isZh: boolean) {
  if (card.kind === "product") {
    return (
      <div className="tb-chat-card">
        <div className="tb-chat-card-head">{isZh ? "商品详情" : "Product Detail"}</div>
        <div className="tb-chat-card-body">
          <div className="tb-chat-card-thumb" style={card.imageUrl ? { backgroundImage: `url("${card.imageUrl}")`, backgroundSize: "cover" } : undefined} />
          <div className="tb-chat-card-main">
            <strong>{card.title}</strong>
            <span>{card.skuName || card.skuNo}</span>
            <small>{`CNY ${formatCnyFromCents(card.priceCents)}`}</small>
            <Link href={card.href}>{isZh ? "查看商品" : "View Product"}</Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="tb-chat-card">
      <div className="tb-chat-card-head">{isZh ? "订单详情" : "Order Detail"}</div>
      <div className="tb-chat-card-body">
        <div className="tb-chat-card-main">
          <strong>{card.title}</strong>
          <span>
            {isZh ? "订单号：" : "Order No: "}
            {card.orderNo}
          </span>
          <span>{card.statusText}</span>
          <small>{`${isZh ? "金额：" : "Amount: "}CNY ${formatCnyFromCents(card.totalAmountCents)}`}</small>
          <Link href={card.href}>{isZh ? "查看订单" : "View Order"}</Link>
        </div>
      </div>
    </div>
  );
}

function mergedConversationList(conversations: ChatConversation[]): ChatConversation[] {
  const groups = new Map<string, ChatConversation[]>();
  conversations.forEach((conversation) => {
    const key = conversation.shopNo || conversation.conversationNo;
    const bucket = groups.get(key) ?? [];
    bucket.push(conversation);
    groups.set(key, bucket);
  });

  return Array.from(groups.values())
    .map((bucket) => {
      const sorted = [...bucket].sort((a, b) => {
        const av = a.lastMessageAt ? new Date(a.lastMessageAt).getTime() : 0;
        const bv = b.lastMessageAt ? new Date(b.lastMessageAt).getTime() : 0;
        return bv - av;
      });
      const latest = sorted[0];
      return {
        ...latest,
        unreadCount: bucket.reduce((sum, item) => sum + item.unreadCount, 0)
      };
    })
    .sort((a, b) => {
      const av = a.lastMessageAt ? new Date(a.lastMessageAt).getTime() : 0;
      const bv = b.lastMessageAt ? new Date(b.lastMessageAt).getTime() : 0;
      return bv - av;
    });
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
  const defaultOrderNo = searchParams.get("order_no") ?? "";
  const defaultSubOrderNo = searchParams.get("sub_order_no") ?? "";
  const mallAccessToken = readMallAccessToken();

  const [currentConversationNo, setCurrentConversationNo] = useState("");
  const [draft, setDraft] = useState("");
  const [plusMenuOpen, setPlusMenuOpen] = useState(false);

  const autoCreateKey = defaultShopNo ? [defaultShopNo, defaultOrderNo, defaultSubOrderNo, defaultSpuNo, defaultSkuNo].join("::") : "";

  const conversationsQuery = useQuery({
    queryKey: BUYER_CONVERSATIONS_QUERY_KEY,
    queryFn: () => listBuyerConversations(30, ""),
    enabled: Boolean(mallAccessToken),
    staleTime: 15_000,
    refetchInterval: mallAccessToken ? 12_000 : false,
    refetchOnWindowFocus: false
  });

  const productDetailQuery = useQuery({
    queryKey: ["buyer-chat-product", defaultSpuNo],
    queryFn: () => getBuyerProductDetail(defaultSpuNo),
    enabled: Boolean(defaultSpuNo),
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const productImagesQuery = useQuery({
    queryKey: ["buyer-chat-product-images", defaultSpuNo],
    queryFn: () => listBuyerProductImages([defaultSpuNo]),
    enabled: Boolean(defaultSpuNo),
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const orderDetailQuery = useQuery({
    queryKey: ["buyer-chat-order", defaultOrderNo],
    queryFn: () => getMyOrderDetail(defaultOrderNo),
    enabled: Boolean(defaultOrderNo),
    staleTime: 30_000,
    refetchOnWindowFocus: false
  });

  const createConversationMutation = useMutation({
    mutationFn: () =>
      createOrGetConversation({
        shopNo: defaultShopNo,
        sceneCode: defaultOrderNo ? "AFTER_SALE" : "PRE_SALE",
        orderNo: defaultOrderNo,
        subOrderNo: defaultSubOrderNo,
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
      messageApi.error(isZh ? "创建会话失败，请稍后重试" : "Failed to create conversation");
    }
  });

  const rawConversations = useMemo(() => conversationsQuery.data?.list ?? [], [conversationsQuery.data?.list]);
  const conversations = useMemo(() => mergedConversationList(rawConversations), [rawConversations]);

  useEffect(() => {
    if (!defaultShopNo) {
      return;
    }
    if (!mallAccessToken) {
      messageApi.info(isZh ? "请先登录后再联系商家" : "Please sign in before chatting with the seller");
      router.replace("/login");
      return;
    }

    const existingShopConversation = conversations.find((item) => item.shopNo === defaultShopNo);
    if (existingShopConversation) {
      setCurrentConversationNo(existingShopConversation.conversationNo);
      autoCreateKeyRef.current = autoCreateKey;
      return;
    }

    if (createConversationMutation.isLoading || autoCreateKeyRef.current === autoCreateKey) {
      return;
    }

    autoCreateKeyRef.current = autoCreateKey;
    setCurrentConversationNo("");
    createConversationMutation.reset();
    createConversationMutation.mutate();
  }, [autoCreateKey, conversations, createConversationMutation, defaultShopNo, isZh, mallAccessToken, messageApi, router]);

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
    setPlusMenuOpen(false);
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
    mutationFn: sendBuyerMessage,
    onSuccess: async () => {
      setDraft("");
      setPlusMenuOpen(false);
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ["chat-buyer-messages", currentConversationNo] }),
        queryClient.invalidateQueries({ queryKey: BUYER_CONVERSATIONS_QUERY_KEY }),
        queryClient.invalidateQueries({ queryKey: SHELL_UNREAD_QUERY_KEY })
      ]);
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
      ? "还没有会话，去商品页发起咨询吧。"
      : "No conversation yet. Open a product page and click 'Consult Seller'.";

  const quickProductCard = useMemo<Extract<ChatCardPayload, { kind: "product" }> | null>(() => {
    if (!defaultShopNo || !defaultSpuNo) {
      return null;
    }
    const detail = productDetailQuery.data;
    const sku = detail?.skus?.find((item) => item.skuNo === defaultSkuNo) ?? detail?.skus?.[0];
    const imageUrl = productImagesQuery.data?.[defaultSpuNo] || pickProductImageBySpuNo(defaultSpuNo);
    return {
      kind: "product",
      shopNo: defaultShopNo,
      spuNo: defaultSpuNo,
      skuNo: defaultSkuNo || sku?.skuNo || "",
      title: detail?.title || (isZh ? `商品 ${defaultSpuNo}` : `Product ${defaultSpuNo}`),
      skuName: sku?.skuName || defaultSkuNo,
      imageUrl,
      priceCents: Number(sku?.salePrice ?? detail?.minSalePrice ?? 0),
      href: `/item/${encodeURIComponent(defaultSpuNo)}`
    };
  }, [defaultShopNo, defaultSkuNo, defaultSpuNo, isZh, productDetailQuery.data, productImagesQuery.data]);

  const quickOrderCard = useMemo<Extract<ChatCardPayload, { kind: "order" }> | null>(() => {
    if (!defaultOrderNo) {
      return null;
    }
    const order = orderDetailQuery.data;
    if (!order) {
      return {
        kind: "order",
        orderNo: defaultOrderNo,
        subOrderNo: defaultSubOrderNo,
        shopNo: defaultShopNo,
        statusText: isZh ? "售后沟通中" : "After-sale conversation",
        totalAmountCents: 0,
        itemCount: 0,
        title: isZh ? `订单 ${defaultOrderNo}` : `Order ${defaultOrderNo}`,
        href: `/me/orders?order_no=${encodeURIComponent(defaultOrderNo)}`
      };
    }
    const items = order.subOrders.flatMap((sub) => sub.items);
    const firstItem = items[0];
    return {
      kind: "order",
      orderNo: order.orderNo,
      subOrderNo: defaultSubOrderNo || firstItem?.subOrderNo || "",
      shopNo: firstItem?.shopNo || defaultShopNo,
      statusText: isZh ? "售后订单沟通" : "After-sale order conversation",
      totalAmountCents: order.amount.payableAmount,
      itemCount: items.reduce((sum, item) => sum + item.qty, 0),
      title: firstItem?.spuTitle || (isZh ? `订单 ${order.orderNo}` : `Order ${order.orderNo}`),
      href: `/me/orders?order_no=${encodeURIComponent(order.orderNo)}`
    };
  }, [defaultOrderNo, defaultShopNo, defaultSubOrderNo, isZh, orderDetailQuery.data]);

  function sendCard(card: ChatCardPayload) {
    if (!currentConversationNo) {
      messageApi.warning(isZh ? "请先选择一个会话" : "Please select a conversation first");
      return;
    }
    const contentText =
      card.kind === "product"
        ? isZh
          ? `商品详情：${card.title}`
          : `Product detail: ${card.title}`
        : isZh
          ? `订单详情：${card.orderNo}`
          : `Order detail: ${card.orderNo}`;
    sendMutation.mutate({
      conversationNo: currentConversationNo,
      contentText,
      messageType: "PRODUCT_CARD",
      extJson: JSON.stringify(card)
    });
  }

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    if (!currentConversationNo) {
      messageApi.warning(isZh ? "请先选择一个会话" : "Please select a conversation first");
      return;
    }
    if (!draft.trim()) {
      return;
    }
    sendMutation.mutate({
      conversationNo: currentConversationNo,
      contentText: draft.trim()
    });
  }

  return (
    <section className="tb-chat-page">
      {contextHolder}

      <header className="tb-chat-header">
        <div>
          <h1>{isZh ? "消息中心" : "Message Center"}</h1>
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
                <h2>{currentConversation ? conversationTitle(currentConversation, isZh) : isZh ? "请选择会话" : "Select Conversation"}</h2>
                <span>{conversationSceneText(currentConversation, isZh, defaultOrderNo)}</span>
              </div>
            </div>
          </div>

          <div className="tb-chat-messages">
            {conversationMessages.map((msg) => {
              const self = msg.senderType === "BUYER";
              const avatarSrc = self ? msg.senderAvatarUrl : currentConversation?.shopAvatarUrl || msg.senderAvatarUrl;
              const cardPayload = parseCardPayload(msg);
              return (
                <article key={msg.messageNo} className={`tb-chat-msg ${self ? "is-self" : ""}`}>
                  {!self ? (
                    <Avatar src={avatarSrc || undefined} size={36} className="tb-chat-msg-avatar">
                      {messageAvatarLabel(msg, currentConversation, isZh)}
                    </Avatar>
                  ) : null}

                  <div className="tb-chat-bubble">
                    <div className="tb-chat-bubble-head">{senderTitle(msg, currentConversation, isZh)}</div>
                    {cardPayload ? renderCardPayload(cardPayload, isZh) : <p>{msg.contentText || (isZh ? "暂不支持该消息类型" : "Unsupported message type")}</p>}
                    <small>
                      {formatListTime(msg.sentAt, locale)}
                      {peerReadText(msg, isZh) ? ` / ${peerReadText(msg, isZh)}` : ""}
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

            {messagesQuery.isLoading ? <p className="tb-chat-empty">{isZh ? "正在加载消息..." : "Loading messages..."}</p> : null}
          </div>

          <form className="tb-chat-composer" onSubmit={handleSubmit}>
            <textarea value={draft} onChange={(event) => setDraft(event.target.value)} placeholder={isZh ? "输入消息..." : "Type your message..."} rows={3} />

            <div className="tb-chat-composer-actions">
              <div className="tb-chat-plus-wrap">
                <button
                  type="button"
                  className="tb-chat-plus-btn"
                  onClick={() => setPlusMenuOpen((value) => !value)}
                  disabled={!currentConversationNo || sendMutation.isLoading || (!quickProductCard && !quickOrderCard)}
                >
                  +
                </button>
                {plusMenuOpen ? (
                  <div className="tb-chat-plus-menu">
                    {quickProductCard ? (
                      <button type="button" onClick={() => sendCard(quickProductCard)}>
                        {isZh ? "发送商品详情" : "Send Product Detail"}
                      </button>
                    ) : null}
                    {quickOrderCard ? (
                      <button type="button" onClick={() => sendCard(quickOrderCard)}>
                        {isZh ? "发送订单详情" : "Send Order Detail"}
                      </button>
                    ) : null}
                  </div>
                ) : null}
              </div>

              <button type="submit" disabled={!currentConversationNo || sendMutation.isLoading || !draft.trim()}>
                {sendMutation.isLoading ? (isZh ? "发送中..." : "Sending...") : isZh ? "发送" : "Send"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </section>
  );
}
