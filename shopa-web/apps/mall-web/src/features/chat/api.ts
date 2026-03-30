import { apiClient } from "@/lib/api-client";
import { ChatConversation, ChatMessage } from "./types";

type RawConversation = {
  conversation_no?: string;
  conversationNo?: string;
  buyer_id?: number | string;
  buyerId?: number | string;
  shop_no?: string;
  shopNo?: string;
  anchor_spu_no?: string;
  anchorSpuNo?: string;
  anchor_sku_no?: string;
  anchorSkuNo?: string;
  unread_count?: number | string;
  unreadCount?: number | string;
  last_message_preview?: string;
  lastMessagePreview?: string;
  last_message_at?: string;
  lastMessageAt?: string;
  conversation_status?: number | string;
  conversationStatus?: number | string;
};

type RawMessage = {
  message_no?: string;
  messageNo?: string;
  conversation_no?: string;
  conversationNo?: string;
  sender_type?: number | string;
  senderType?: number | string;
  sender_user_id?: number | string;
  senderUserId?: number | string;
  message_type?: number | string;
  messageType?: number | string;
  content_text?: string;
  contentText?: string;
  media_asset_id?: number | string;
  mediaAssetId?: number | string;
  ext_json?: string;
  extJson?: string;
  sent_at?: string;
  sentAt?: string;
};

type RawCreateConversationRes = {
  conversation?: RawConversation;
};

type RawListConversationRes = {
  list?: RawConversation[];
  has_more?: boolean;
  hasMore?: boolean;
  next_cursor?: string;
  nextCursor?: string;
};

type RawListMessagesRes = {
  list?: RawMessage[];
  has_more?: boolean;
  hasMore?: boolean;
  next_cursor?: string;
  nextCursor?: string;
};

type RawSendMessageRes = {
  conversation?: RawConversation;
  message?: RawMessage;
  idempotent_replay?: boolean;
  idempotentReplay?: boolean;
};

type RawUnreadSummaryRes = {
  total_unread_conversations?: number | string;
  totalUnreadConversations?: number | string;
  total_unread_messages?: number | string;
  totalUnreadMessages?: number | string;
};

type PageResult<T> = {
  list: T[];
  hasMore: boolean;
  nextCursor: string;
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

function mapStatus(value: number): ChatConversation["status"] {
  switch (value) {
    case 1:
      return "ACTIVE";
    case 2:
      return "CLOSED";
    case 3:
      return "BLOCKED";
    default:
      return "UNKNOWN";
  }
}

function mapSenderType(value: number): ChatMessage["senderType"] {
  switch (value) {
    case 1:
      return "BUYER";
    case 2:
      return "SELLER";
    case 3:
      return "SYSTEM";
    default:
      return "UNKNOWN";
  }
}

function mapMessageType(value: number): ChatMessage["messageType"] {
  switch (value) {
    case 1:
      return "TEXT";
    case 2:
      return "IMAGE";
    case 3:
      return "PRODUCT_CARD";
    case 4:
      return "SYSTEM_NOTICE";
    default:
      return "UNKNOWN";
  }
}

function normalizeConversation(raw?: RawConversation): ChatConversation {
  const statusCode = toNumber(raw?.conversation_status ?? raw?.conversationStatus);
  return {
    conversationNo: toString(raw?.conversation_no ?? raw?.conversationNo),
    buyerId: toNumber(raw?.buyer_id ?? raw?.buyerId),
    shopNo: toString(raw?.shop_no ?? raw?.shopNo),
    anchorSpuNo: toString(raw?.anchor_spu_no ?? raw?.anchorSpuNo),
    anchorSkuNo: toString(raw?.anchor_sku_no ?? raw?.anchorSkuNo),
    unreadCount: toNumber(raw?.unread_count ?? raw?.unreadCount),
    lastMessagePreview: toString(raw?.last_message_preview ?? raw?.lastMessagePreview),
    lastMessageAt: toString(raw?.last_message_at ?? raw?.lastMessageAt),
    status: mapStatus(statusCode)
  };
}

function normalizeMessage(raw?: RawMessage): ChatMessage {
  return {
    messageNo: toString(raw?.message_no ?? raw?.messageNo),
    conversationNo: toString(raw?.conversation_no ?? raw?.conversationNo),
    senderType: mapSenderType(toNumber(raw?.sender_type ?? raw?.senderType)),
    senderUserId: toNumber(raw?.sender_user_id ?? raw?.senderUserId),
    messageType: mapMessageType(toNumber(raw?.message_type ?? raw?.messageType)),
    contentText: toString(raw?.content_text ?? raw?.contentText),
    mediaAssetId: toNumber(raw?.media_asset_id ?? raw?.mediaAssetId),
    extJson: toString(raw?.ext_json ?? raw?.extJson),
    sentAt: toString(raw?.sent_at ?? raw?.sentAt)
  };
}

export async function createOrGetConversation(payload: {
  shopNo: string;
  anchorSpuNo?: string;
  anchorSkuNo?: string;
}): Promise<ChatConversation> {
  const response = await apiClient.post<RawCreateConversationRes>("/v1/chat/buyer/conversations:get-or-create", {
    shop_no: payload.shopNo,
    anchor_spu_no: payload.anchorSpuNo ?? "",
    anchor_sku_no: payload.anchorSkuNo ?? ""
  });
  return normalizeConversation(response?.conversation);
}

export async function listBuyerConversations(pageSize = 20, nextCursor = ""): Promise<PageResult<ChatConversation>> {
  const response = await apiClient.get<RawListConversationRes>("/v1/chat/buyer/conversations", {
    params: {
      page_size: pageSize,
      next_cursor: nextCursor
    }
  });
  return {
    list: (response?.list ?? []).map((item) => normalizeConversation(item)),
    hasMore: Boolean(response?.has_more ?? response?.hasMore),
    nextCursor: toString(response?.next_cursor ?? response?.nextCursor)
  };
}

export async function listMessages(conversationNo: string, pageSize = 40, nextCursor = ""): Promise<PageResult<ChatMessage>> {
  const response = await apiClient.get<RawListMessagesRes>(`/v1/chat/buyer/conversations/${encodeURIComponent(conversationNo)}/messages`, {
    params: {
      page_size: pageSize,
      next_cursor: nextCursor
    }
  });
  return {
    list: (response?.list ?? []).map((item) => normalizeMessage(item)),
    hasMore: Boolean(response?.has_more ?? response?.hasMore),
    nextCursor: toString(response?.next_cursor ?? response?.nextCursor)
  };
}

export async function sendBuyerMessage(payload: {
  conversationNo: string;
  contentText: string;
  clientMessageNo?: string;
}): Promise<{ conversation: ChatConversation; message: ChatMessage; idempotentReplay: boolean }> {
  const response = await apiClient.post<RawSendMessageRes>("/v1/chat/buyer/messages:send", {
    conversation_no: payload.conversationNo,
    client_message_no: payload.clientMessageNo ?? `m_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
    message_type: 1,
    content_text: payload.contentText
  });
  return {
    conversation: normalizeConversation(response?.conversation),
    message: normalizeMessage(response?.message),
    idempotentReplay: Boolean(response?.idempotent_replay ?? response?.idempotentReplay)
  };
}

export async function markConversationRead(payload: { conversationNo: string; readToMessageNo: string }): Promise<void> {
  await apiClient.post("/v1/chat/buyer/conversations:mark-read", {
    conversation_no: payload.conversationNo,
    read_to_message_no: payload.readToMessageNo
  });
}

export async function getUnreadSummary(): Promise<{ totalUnreadConversations: number; totalUnreadMessages: number }> {
  const response = await apiClient.get<RawUnreadSummaryRes>("/v1/chat/buyer/unread-summary");
  return {
    totalUnreadConversations: toNumber(response?.total_unread_conversations ?? response?.totalUnreadConversations),
    totalUnreadMessages: toNumber(response?.total_unread_messages ?? response?.totalUnreadMessages)
  };
}

