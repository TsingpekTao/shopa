import { apiClient } from "@/lib/api-client";
import { SellerChatConversation, SellerChatMessage } from "./types";

type RawConversation = {
  conversation_no?: string;
  conversationNo?: string;
  buyer_id?: number | string;
  buyerId?: number | string;
  shop_no?: string;
  shopNo?: string;
  shop_name?: string;
  shopName?: string;
  buyer_display_name?: string;
  buyerDisplayName?: string;
  buyer_avatar_url?: string;
  buyerAvatarUrl?: string;
  shop_avatar_url?: string;
  shopAvatarUrl?: string;
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
  buyer_read_to_message_no?: string;
  buyerReadToMessageNo?: string;
  seller_read_to_message_no?: string;
  sellerReadToMessageNo?: string;
  latest_message_read_by_peer?: boolean;
  latestMessageReadByPeer?: boolean;
  latest_message_peer_read_at?: string;
  latestMessagePeerReadAt?: string;
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
  sender_display_name?: string;
  senderDisplayName?: string;
  sender_avatar_url?: string;
  senderAvatarUrl?: string;
  peer_read?: boolean;
  peerRead?: boolean;
  peer_read_at?: string;
  peerReadAt?: string;
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

export type SellerChatPageResult<T> = {
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

function mapStatus(value: number): SellerChatConversation["status"] {
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

function mapSenderType(value: number): SellerChatMessage["senderType"] {
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

function mapMessageType(value: number): SellerChatMessage["messageType"] {
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

function normalizeConversation(raw?: RawConversation): SellerChatConversation {
  const statusCode = toNumber(raw?.conversation_status ?? raw?.conversationStatus);
  return {
    conversationNo: toString(raw?.conversation_no ?? raw?.conversationNo),
    buyerId: toNumber(raw?.buyer_id ?? raw?.buyerId),
    shopNo: toString(raw?.shop_no ?? raw?.shopNo),
    shopName: toString(raw?.shop_name ?? raw?.shopName),
    buyerDisplayName: toString(raw?.buyer_display_name ?? raw?.buyerDisplayName),
    buyerAvatarUrl: toString(raw?.buyer_avatar_url ?? raw?.buyerAvatarUrl),
    shopAvatarUrl: toString(raw?.shop_avatar_url ?? raw?.shopAvatarUrl),
    anchorSpuNo: toString(raw?.anchor_spu_no ?? raw?.anchorSpuNo),
    anchorSkuNo: toString(raw?.anchor_sku_no ?? raw?.anchorSkuNo),
    unreadCount: toNumber(raw?.unread_count ?? raw?.unreadCount),
    lastMessagePreview: toString(raw?.last_message_preview ?? raw?.lastMessagePreview),
    lastMessageAt: toString(raw?.last_message_at ?? raw?.lastMessageAt),
    status: mapStatus(statusCode),
    buyerReadToMessageNo: toString(raw?.buyer_read_to_message_no ?? raw?.buyerReadToMessageNo),
    sellerReadToMessageNo: toString(raw?.seller_read_to_message_no ?? raw?.sellerReadToMessageNo),
    latestMessageReadByPeer: Boolean(raw?.latest_message_read_by_peer ?? raw?.latestMessageReadByPeer),
    latestMessagePeerReadAt: toString(raw?.latest_message_peer_read_at ?? raw?.latestMessagePeerReadAt)
  };
}

function normalizeMessage(raw?: RawMessage): SellerChatMessage {
  return {
    messageNo: toString(raw?.message_no ?? raw?.messageNo),
    conversationNo: toString(raw?.conversation_no ?? raw?.conversationNo),
    senderType: mapSenderType(toNumber(raw?.sender_type ?? raw?.senderType)),
    senderUserId: toNumber(raw?.sender_user_id ?? raw?.senderUserId),
    messageType: mapMessageType(toNumber(raw?.message_type ?? raw?.messageType)),
    contentText: toString(raw?.content_text ?? raw?.contentText),
    mediaAssetId: toNumber(raw?.media_asset_id ?? raw?.mediaAssetId),
    extJson: toString(raw?.ext_json ?? raw?.extJson),
    sentAt: toString(raw?.sent_at ?? raw?.sentAt),
    senderDisplayName: toString(raw?.sender_display_name ?? raw?.senderDisplayName),
    senderAvatarUrl: toString(raw?.sender_avatar_url ?? raw?.senderAvatarUrl),
    peerRead: Boolean(raw?.peer_read ?? raw?.peerRead),
    peerReadAt: toString(raw?.peer_read_at ?? raw?.peerReadAt)
  };
}

export async function listSellerConversations(
  shopNo: string,
  pageSize = 20,
  nextCursor = ""
): Promise<SellerChatPageResult<SellerChatConversation>> {
  const response = await apiClient.get<RawListConversationRes>(
    `/v1/chat/seller/shops/${encodeURIComponent(shopNo)}/conversations`,
    {
      params: {
        page_size: pageSize,
        next_cursor: nextCursor
      }
    }
  );
  return {
    list: (response?.list ?? []).map((item) => normalizeConversation(item)),
    hasMore: Boolean(response?.has_more ?? response?.hasMore),
    nextCursor: toString(response?.next_cursor ?? response?.nextCursor)
  };
}

export async function listSellerMessages(
  shopNo: string,
  conversationNo: string,
  pageSize = 40,
  nextCursor = ""
): Promise<SellerChatPageResult<SellerChatMessage>> {
  const response = await apiClient.get<RawListMessagesRes>(
    `/v1/chat/seller/shops/${encodeURIComponent(shopNo)}/conversations/${encodeURIComponent(conversationNo)}/messages`,
    {
      params: {
        page_size: pageSize,
        next_cursor: nextCursor
      }
    }
  );
  return {
    list: (response?.list ?? []).map((item) => normalizeMessage(item)),
    hasMore: Boolean(response?.has_more ?? response?.hasMore),
    nextCursor: toString(response?.next_cursor ?? response?.nextCursor)
  };
}

export async function sendSellerMessage(payload: {
  conversationNo: string;
  contentText: string;
  clientMessageNo?: string;
}): Promise<{ conversation: SellerChatConversation; message: SellerChatMessage; idempotentReplay: boolean }> {
  const response = await apiClient.post<RawSendMessageRes>("/v1/chat/seller/messages:send", {
    conversation_no: payload.conversationNo,
    client_message_no: payload.clientMessageNo ?? `seller_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
    message_type: 1,
    content_text: payload.contentText
  });
  return {
    conversation: normalizeConversation(response?.conversation),
    message: normalizeMessage(response?.message),
    idempotentReplay: Boolean(response?.idempotent_replay ?? response?.idempotentReplay)
  };
}

export async function markSellerConversationRead(payload: {
  conversationNo: string;
  readToMessageNo: string;
}): Promise<void> {
  await apiClient.post("/v1/chat/seller/conversations:mark-read", {
    conversation_no: payload.conversationNo,
    read_to_message_no: payload.readToMessageNo
  });
}
