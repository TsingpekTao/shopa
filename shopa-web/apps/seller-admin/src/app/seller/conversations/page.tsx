"use client";

import { FormEvent, KeyboardEvent, useEffect, useMemo, useRef, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Alert, Avatar, Badge, Button, Empty, Input, Select, Skeleton, Space, Tag, Typography, message } from "antd";
import { MessageOutlined, ReloadOutlined, SendOutlined, ShopOutlined, UserOutlined } from "@ant-design/icons";
import { useI18n } from "@shopa/ui";
import { listSellerConversations, listSellerMessages, markSellerConversationRead, sendSellerMessage } from "@/features/chat/api";
import { SellerChatConversation, SellerChatMessage } from "@/features/chat/types";
import { fetchSellerWorkbench } from "@/features/seller-shop/api";

const { Title, Text, Paragraph } = Typography;

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
  return isZh ? "暂无消息内容" : "No messages yet";
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

  const workbenchQuery = useQuery({
    queryKey: ["seller-chat-workbench"],
    queryFn: fetchSellerWorkbench,
    staleTime: 30_000
  });

  const shops = useMemo(() => {
    return (workbenchQuery.data?.shops ?? [])
      .map((shop) => ({
        shopNo: String(shop.shopNo ?? "").trim(),
        shopName: String(shop.shopDisplayName || shop.shopName || shop.shopNo || "").trim(),
        shopStatusCode: String(shop.shopStatusCode ?? "").trim()
      }))
      .filter((shop) => shop.shopNo);
  }, [workbenchQuery.data?.shops]);

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

  const selectedShop = useMemo(() => {
    return shops.find((shop) => shop.shopNo === selectedShopNo);
  }, [selectedShopNo, shops]);

  const conversationsQuery = useQuery({
    queryKey: ["seller-chat-conversations", selectedShopNo],
    queryFn: () => listSellerConversations(selectedShopNo, 60, ""),
    enabled: Boolean(selectedShopNo),
    staleTime: 10_000,
    refetchInterval: 12_000
  });

  const conversations = conversationsQuery.data?.list ?? [];

  useEffect(() => {
    setCurrentConversationNo("");
    lastReadRef.current = "";
  }, [selectedShopNo]);

  useEffect(() => {
    if (!conversations.length) {
      setCurrentConversationNo("");
      return;
    }
    setCurrentConversationNo((current) => {
      if (current && conversations.some((item) => item.conversationNo === current)) {
        return current;
      }
      return conversations[0].conversationNo;
    });
  }, [conversations]);

  const currentConversation = useMemo<SellerChatConversation | undefined>(() => {
    return conversations.find((item) => item.conversationNo === currentConversationNo);
  }, [conversations, currentConversationNo]);

  const messagesQuery = useQuery({
    queryKey: ["seller-chat-messages", selectedShopNo, currentConversationNo],
    queryFn: () => listSellerMessages(selectedShopNo, currentConversationNo, 100, ""),
    enabled: Boolean(selectedShopNo) && Boolean(currentConversationNo),
    staleTime: 5_000,
    refetchInterval: 8_000
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
    if (!currentConversationNo || !lastMessage?.messageNo) {
      return;
    }
    const readKey = `${currentConversationNo}:${lastMessage.messageNo}`;
    if (lastReadRef.current === readKey) {
      return;
    }
    lastReadRef.current = readKey;
    markSellerConversationRead({
      conversationNo: currentConversationNo,
      readToMessageNo: lastMessage.messageNo
    })
      .then(() => queryClient.invalidateQueries({ queryKey: ["seller-chat-conversations", selectedShopNo] }))
      .catch(() => {
        lastReadRef.current = "";
      });
  }, [conversationMessages, currentConversationNo, queryClient, selectedShopNo]);

  const sendMutation = useMutation({
    mutationFn: async () =>
      sendSellerMessage({
        conversationNo: currentConversationNo,
        contentText: draft.trim()
      }),
    onSuccess: async () => {
      setDraft("");
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ["seller-chat-messages", selectedShopNo, currentConversationNo] }),
        queryClient.invalidateQueries({ queryKey: ["seller-chat-conversations", selectedShopNo] })
      ]);
    },
    onError: (error) => {
      const errorMessage = error instanceof Error ? error.message : isZh ? "发送消息失败" : "Send message failed";
      messageApi.error(errorMessage);
    }
  });

  const unreadTotal = useMemo(() => conversations.reduce((sum, item) => sum + item.unreadCount, 0), [conversations]);

  function submitMessage() {
    if (!currentConversationNo) {
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
    <section className="seller-page seller-chat-page">
      {contextHolder}
      <header className="seller-page-head">
        <div>
          <Title level={3}>{isZh ? "会话中心" : "Conversations"}</Title>
          <Text type="secondary">
            {isZh
              ? "像微信工作台一样处理买家消息，左侧看会话，右侧直接接待。"
              : "Handle buyer chats in a WeChat-like seller workspace."}
          </Text>
        </div>
        <Space wrap>
          <Select
            value={selectedShopNo || undefined}
            onChange={(value) => setSelectedShopNo(value)}
            className="seller-chat-shop-select"
            placeholder={isZh ? "选择店铺" : "Select shop"}
            options={shops.map((shop) => ({
              value: shop.shopNo,
              label: shop.shopName
            }))}
          />
          <Button
            icon={<ReloadOutlined />}
            onClick={() => {
              void Promise.all([workbenchQuery.refetch(), conversationsQuery.refetch(), messagesQuery.refetch()]);
            }}
          >
            {isZh ? "刷新" : "Refresh"}
          </Button>
        </Space>
      </header>

      {workbenchQuery.data?.partial ? (
        <Alert
          type="warning"
          showIcon
          message={isZh ? "商家工作台数据部分降级" : "Workbench data is partially degraded"}
          description={(workbenchQuery.data.degradedFields ?? []).join(", ")}
        />
      ) : null}

      {!shops.length ? (
        <div className="seller-chat-empty-panel">
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={isZh ? "当前账号下还没有可接待消息的店铺" : "No shop is available for seller chat yet"}
          />
        </div>
      ) : (
        <div className="seller-chat-shell">
          <aside className="seller-chat-sidebar">
            <div className="seller-chat-sidebar-head">
              <div>
                <Text strong>{selectedShop?.shopName || (isZh ? "当前店铺" : "Current shop")}</Text>
                <Text type="secondary">{isZh ? "会话接待中" : "Conversation inbox"}</Text>
              </div>
              <Badge count={unreadTotal} overflowCount={99} color="#ff6a00">
                <Avatar size={40} icon={<ShopOutlined />} className="seller-chat-avatar seller-chat-avatar-shop" />
              </Badge>
            </div>

            <div className="seller-chat-sidebar-meta">
              <Tag color={selectedShop?.shopStatusCode?.toUpperCase() === "ACTIVE" ? "success" : "default"}>
                {selectedShop?.shopStatusCode || "UNKNOWN"}
              </Tag>
              <span>{isZh ? `${conversations.length} 个会话` : `${conversations.length} conversations`}</span>
            </div>

            <div className="seller-chat-conversation-list">
              {conversationsQuery.isLoading ? (
                <Skeleton active paragraph={{ rows: 6 }} title={false} />
              ) : conversations.length === 0 ? (
                <Empty
                  image={Empty.PRESENTED_IMAGE_SIMPLE}
                  description={isZh ? "这个店铺暂时还没有买家会话" : "No buyer conversation for this shop yet"}
                />
              ) : (
                conversations.map((conversation) => {
                  const active = conversation.conversationNo === currentConversationNo;
                  return (
                    <button
                      type="button"
                      key={conversation.conversationNo}
                      className={`seller-chat-conversation-item ${active ? "is-active" : ""}`}
                      onClick={() => {
                        lastReadRef.current = "";
                        setCurrentConversationNo(conversation.conversationNo);
                      }}
                    >
                      <div className="seller-chat-conversation-top">
                        <div className="seller-chat-conversation-buyer">
                          <Avatar size={38} src={conversation.buyerAvatarUrl || undefined} icon={<UserOutlined />} className="seller-chat-avatar" />
                          <div>
                            <strong>{conversation.buyerDisplayName || (isZh ? "买家" : "Buyer")}</strong>
                            <span>{conversation.unreadCount > 0 ? (isZh ? "等待回复中" : "Waiting for reply") : isZh ? "最近有新互动" : "Recent interaction"}</span>
                          </div>
                        </div>
                        <div className="seller-chat-conversation-side">
                          {conversation.lastMessageAt ? <time>{formatListTime(conversation.lastMessageAt, locale)}</time> : null}
                          {conversation.unreadCount > 0 ? (
                            <span className="seller-chat-unread-badge">{conversation.unreadCount > 99 ? "99+" : conversation.unreadCount}</span>
                          ) : null}
                        </div>
                      </div>
                      <Paragraph ellipsis={{ rows: 2 }} className="seller-chat-conversation-preview">
                        {fallbackPreview(conversation.lastMessagePreview, isZh)}
                      </Paragraph>
                    </button>
                  );
                })
              )}
            </div>
          </aside>

          <section className="seller-chat-thread-panel">
            <div className="seller-chat-thread-head">
              {currentConversation ? (
                <>
                  <div className="seller-chat-thread-title">
                    <Avatar
                      size={44}
                      src={currentConversation.buyerAvatarUrl || undefined}
                      icon={<MessageOutlined />}
                      className="seller-chat-avatar seller-chat-avatar-thread"
                    />
                    <div>
                      <strong>{currentConversation.buyerDisplayName || (isZh ? "买家" : "Buyer")}</strong>
                      <span>{isZh ? "和这位买家继续沟通需求与尺码问题" : "Continue helping this buyer with product questions."}</span>
                    </div>
                  </div>
                  <div className="seller-chat-thread-tags">
                    <Tag color="processing">{currentConversation.status}</Tag>
                    {currentConversation.unreadCount > 0 ? <Tag color="orange">{isZh ? `未读 ${currentConversation.unreadCount}` : `${currentConversation.unreadCount} unread`}</Tag> : null}
                  </div>
                </>
              ) : (
                <div className="seller-chat-thread-placeholder">
                  <strong>{isZh ? "请选择左侧会话" : "Select a conversation"}</strong>
                  <span>{isZh ? "选中后即可在右侧像微信一样查看和回复。" : "Open one from the left list to reply on the right."}</span>
                </div>
              )}
            </div>

            <div className="seller-chat-thread-body" ref={threadRef}>
              {messagesQuery.isLoading ? (
                <Skeleton active paragraph={{ rows: 8 }} title={false} />
              ) : !currentConversation ? (
                <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={isZh ? "还没有选中会话" : "No conversation selected"} />
              ) : conversationMessages.length === 0 ? (
                <Empty
                  image={Empty.PRESENTED_IMAGE_SIMPLE}
                  description={isZh ? "当前会话还没有消息" : "This conversation has no messages yet"}
                />
              ) : (
                conversationMessages.map((item) => {
                  const isSelf = item.senderType === "SELLER";
                  const senderLabel =
                    item.senderType === "SYSTEM"
                      ? isZh
                        ? "系统通知"
                        : "System"
                      : isSelf
                        ? selectedShop?.shopName || (isZh ? "我" : "Me")
                        : item.senderDisplayName || currentConversation.buyerDisplayName || (isZh ? "买家" : "Buyer");
                  const readStateText =
                    isSelf && item.senderType !== "SYSTEM"
                      ? item.peerRead
                        ? isZh
                          ? "买家已读"
                          : "Read by buyer"
                        : isZh
                          ? "未读"
                          : "Unread"
                      : "";
                  const peerReadTime = item.peerReadAt ? formatListTime(item.peerReadAt, locale) : "";
                  return (
                    <article
                      key={item.messageNo}
                      className={`seller-chat-message-row ${isSelf ? "is-self" : ""} ${item.senderType === "SYSTEM" ? "is-system" : ""}`}
                    >
                      <Avatar
                        size={36}
                        src={item.senderAvatarUrl || undefined}
                        icon={item.senderType === "SYSTEM" ? <MessageOutlined /> : <UserOutlined />}
                        className={`seller-chat-avatar ${isSelf ? "seller-chat-avatar-self" : ""}`}
                      />
                      <div className="seller-chat-message-main">
                        <div className="seller-chat-message-meta">
                          <span>{senderLabel}</span>
                          <time>
                            {formatListTime(item.sentAt, locale)}
                            {readStateText ? ` · ${readStateText}${peerReadTime ? ` · ${peerReadTime}` : ""}` : ""}
                          </time>
                        </div>
                        <div className={`seller-chat-bubble ${isSelf ? "is-self" : ""} ${item.senderType === "SYSTEM" ? "is-system" : ""}`}>
                          {item.contentText || (isZh ? "暂不支持该消息类型预览" : "Unsupported message type")}
                        </div>
                      </div>
                    </article>
                  );
                })
              )}
            </div>

            <form className="seller-chat-composer" onSubmit={handleSubmit}>
              <Input.TextArea
                value={draft}
                onChange={(event) => setDraft(event.target.value)}
                rows={3}
                placeholder={isZh ? "输入回复内容，回车发送，Shift + Enter 换行" : "Type your reply. Enter to send, Shift + Enter for newline"}
                onPressEnter={handleTextareaEnter}
                disabled={!currentConversation}
              />
              <div className="seller-chat-composer-actions">
                <Text type="secondary">
                  {currentConversation
                    ? isZh
                      ? "消息会自动标记为已读，右侧支持实时刷新。"
                      : "Messages are auto-marked as read and refreshed in place."
                    : isZh
                      ? "先从左侧选择一个会话。"
                      : "Choose a conversation from the left first."}
                </Text>
                <Button
                  type="primary"
                  htmlType="submit"
                  icon={<SendOutlined />}
                  loading={sendMutation.isLoading}
                  disabled={!currentConversation || !draft.trim()}
                >
                  {isZh ? "发送" : "Send"}
                </Button>
              </div>
            </form>
          </section>
        </div>
      )}
    </section>
  );
}
