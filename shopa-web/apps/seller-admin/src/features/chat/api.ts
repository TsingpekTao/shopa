import { apiClient } from "@/lib/api-client";
import type {
  SellerChatConversation,
  SellerChatMessage,
  SellerChatOrderCard,
  SellerChatProductCard,
  SellerChatMessageCard,
  SellerChatMessageType,
  SellerConversationSceneCode
} from "./types";

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
  scene_code?: string;
  sceneCode?: string;
  order_no?: string;
  orderNo?: string;
  sub_order_no?: string;
  subOrderNo?: string;
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
  if (typeof value === "string") {
    return value;
  }
  if (typeof value === "number" || typeof value === "boolean" || typeof value === "bigint") {
    return String(value);
  }
  if (value instanceof Date) {
    return value.toISOString();
  }
  if (typeof value === "object" && !Array.isArray(value)) {
    const record = value as Record<string, unknown>;
    const seconds = Number(record.seconds ?? 0);
    const nanos = Number(record.nanos ?? 0);
    if (Number.isFinite(seconds) || Number.isFinite(nanos)) {
      const millis = seconds * 1000 + Math.floor(nanos / 1_000_000);
      if (millis > 0) {
        return new Date(millis).toISOString();
      }
    }
    const candidates = [record.text, record.content, record.title, record.name, record.message, record.value];
    for (const item of candidates) {
      if (typeof item === "string" && item.trim()) {
        return item.trim();
      }
    }
  }
  return "";
}

function toNumber(value: unknown): number {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : 0;
}

function readText(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  if (typeof value === "string") {
    return value.trim();
  }
  return toString(value).trim();
}

function parseExtObject(raw: string): Record<string, unknown> {
  const trimmed = raw.trim();
  if (!trimmed) {
    return {};
  }
  try {
    const parsed = JSON.parse(trimmed);
    if (parsed && typeof parsed === "object" && !Array.isArray(parsed)) {
      return parsed as Record<string, unknown>;
    }
  } catch {
    return {};
  }
  return {};
}

function parseStringList(value: unknown): string[] {
  if (!value) {
    return [];
  }
  if (Array.isArray(value)) {
    return value.map((item) => readText(item)).filter(Boolean);
  }
  const text = readText(value);
  if (!text) {
    return [];
  }
  return text
    .split(",")
    .map((segment) => segment.trim())
    .filter(Boolean);
}

type SellerOrderCardItem = NonNullable<SellerChatOrderCard["items"]>[number];

function parseOrderItems(value: unknown): SellerOrderCardItem[] {
  if (!Array.isArray(value)) {
    return [];
  }
  const items: SellerOrderCardItem[] = [];
  value.forEach((entry) => {
    if (!entry || typeof entry !== "object" || Array.isArray(entry)) {
      return;
    }
    const record = entry as Record<string, unknown>;
    const item: SellerOrderCardItem = {
      title: readText(record.title) || readText(record.name) || undefined,
      quantity: readText(record.quantity) || readText(record.qty) || undefined,
      price: readText(record.price) || readText(record.amount) || undefined,
      meta: readText(record.meta) || readText(record.detail) || undefined
    };
    if (item.title || item.quantity || item.price || item.meta) {
      items.push(item);
    }
  });
  return items;
}

function buildSellerProductCard(ext: Record<string, unknown>, fallbackTitle: string): SellerChatProductCard {
  const tags = parseStringList(ext.tags ?? ext.tagList ?? ext.labels);
  return {
    title: readText(ext.title) || readText(ext.productTitle) || readText(ext.product_title) || fallbackTitle,
    subtitle: readText(ext.subtitle) || readText(ext.subTitle) || readText(ext.sub_title) || undefined,
    coverImageUrl: readText(ext.imageUrl) || readText(ext.image_url) || undefined,
    price: readText(ext.priceText) || readText(ext.price_text) || readText(ext.price) || undefined,
    priceLabel: readText(ext.priceLabel) || readText(ext.price_label) || undefined,
    skuNo: readText(ext.skuNo) || readText(ext.sku_no) || undefined,
    spuNo: readText(ext.spuNo) || readText(ext.spu_no) || undefined,
    tags: tags.length ? tags : undefined,
    detail: readText(ext.detail) || readText(ext.description) || undefined,
    link: readText(ext.link) || readText(ext.href) || readText(ext.url) || undefined
  };
}

function buildSellerOrderCard(ext: Record<string, unknown>): SellerChatOrderCard {
  const meta = parseStringList(ext.meta ?? ext.orderMeta ?? ext.metadata ?? ext.order_metadata);
  const items = parseOrderItems(ext.items ?? ext.orderItems ?? ext.products);
  return {
    orderNo: readText(ext.orderNo) || readText(ext.order_no) || undefined,
    subOrderNo: readText(ext.subOrderNo) || readText(ext.sub_order_no) || undefined,
    shopNo: readText(ext.shopNo) || readText(ext.shop_no) || undefined,
    status: readText(ext.status) || readText(ext.orderStatus) || readText(ext.statusCode) || undefined,
    sceneCode: readText(ext.sceneCode) || readText(ext.scene_code) || undefined,
    totalAmount: readText(ext.totalAmount) || readText(ext.total_amount) || undefined,
    currency: readText(ext.currency) || undefined,
    summary: readText(ext.summary) || readText(ext.description) || undefined,
    meta: meta.length ? meta : undefined,
    items: items.length ? items : undefined,
    link: readText(ext.link) || readText(ext.href) || readText(ext.url) || undefined
  };
}

function buildSellerCard(raw: string, messageType: SellerChatMessageType, fallbackTitle: string): SellerChatMessageCard | undefined {
  if (messageType !== "PRODUCT_CARD") {
    return undefined;
  }
  const ext = parseExtObject(raw);
  const kindRaw = readText(ext.kind).toLowerCase();
  const orderNo = readText(ext.orderNo) || readText(ext.order_no);
  const subOrderNo = readText(ext.subOrderNo) || readText(ext.sub_order_no);
  const isOrderKind = kindRaw === "order" || Boolean(orderNo) || Boolean(subOrderNo);
  if (isOrderKind) {
    const order = buildSellerOrderCard(ext);
    const hasOrderInfo =
      Boolean(order.orderNo) ||
      Boolean(order.subOrderNo) ||
      Boolean(order.shopNo) ||
      Boolean(order.status) ||
      Boolean(order.sceneCode) ||
      Boolean(order.totalAmount) ||
      Boolean(order.summary) ||
      Boolean(order.items?.length) ||
      Boolean(order.meta?.length);
    return hasOrderInfo ? { kind: "ORDER_CARD", order } : undefined;
  }
  const product = buildSellerProductCard(ext, fallbackTitle);
  const hasProductInfo =
    Boolean(product.title) ||
    Boolean(product.subtitle) ||
    Boolean(product.coverImageUrl) ||
    Boolean(product.price) ||
    Boolean(product.priceLabel) ||
    Boolean(product.tags?.length) ||
    Boolean(product.detail);
  return hasProductInfo ? { kind: "PRODUCT_CARD", product } : undefined;
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

function mapMessageType(value: number): SellerChatMessageType {
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

function mapSceneCode(value: unknown): SellerConversationSceneCode {
  const normalized = toString(value).trim().toUpperCase();
  if (normalized === "PRE_SALE") {
    return "PRE_SALE";
  }
  if (normalized === "AFTER_SALE") {
    return "AFTER_SALE";
  }
  return "";
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
    sceneCode: mapSceneCode(raw?.scene_code ?? raw?.sceneCode),
    orderNo: toString(raw?.order_no ?? raw?.orderNo),
    subOrderNo: toString(raw?.sub_order_no ?? raw?.subOrderNo),
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
  const messageType = mapMessageType(toNumber(raw?.message_type ?? raw?.messageType));
  const contentText = toString(raw?.content_text ?? raw?.contentText);
  const extJson = toString(raw?.ext_json ?? raw?.extJson);
  return {
    messageNo: toString(raw?.message_no ?? raw?.messageNo),
    conversationNo: toString(raw?.conversation_no ?? raw?.conversationNo),
    senderType: mapSenderType(toNumber(raw?.sender_type ?? raw?.senderType)),
    senderUserId: toNumber(raw?.sender_user_id ?? raw?.senderUserId),
    messageType,
    contentText,
    mediaAssetId: toNumber(raw?.media_asset_id ?? raw?.mediaAssetId),
    extJson,
    card: buildSellerCard(extJson, messageType, contentText),
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
  const normalizedShopNo = shopNo.trim();
  if (!normalizedShopNo) {
    return {
      list: [],
      hasMore: false,
      nextCursor: ""
    };
  }
  const response = await apiClient.get<RawListConversationRes>(
    `/v1/chat/seller/shops/${encodeURIComponent(normalizedShopNo)}/conversations`,
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
  const normalizedShopNo = shopNo.trim();
  const normalizedConversationNo = conversationNo.trim();
  if (!normalizedShopNo || !normalizedConversationNo) {
    return {
      list: [],
      hasMore: false,
      nextCursor: ""
    };
  }
  const response = await apiClient.get<RawListMessagesRes>(
    `/v1/chat/seller/shops/${encodeURIComponent(normalizedShopNo)}/conversations/${encodeURIComponent(normalizedConversationNo)}/messages`,
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
  messageType?: SellerChatMessageType;
  extJson?: string;
  clientMessageNo?: string;
}): Promise<{ conversation: SellerChatConversation; message: SellerChatMessage; idempotentReplay: boolean }> {
  const messageType = payload.messageType === "PRODUCT_CARD" ? 3 : payload.messageType === "SYSTEM_NOTICE" ? 4 : payload.messageType === "IMAGE" ? 2 : 1;
  const response = await apiClient.post<RawSendMessageRes>("/v1/chat/seller/messages:send", {
    conversation_no: payload.conversationNo,
    client_message_no: payload.clientMessageNo ?? `seller_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
    message_type: messageType,
    content_text: payload.contentText,
    ext_json: payload.extJson ?? ""
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
