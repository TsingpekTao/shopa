export type SellerConversationStatus = "ACTIVE" | "CLOSED" | "BLOCKED" | "UNKNOWN";
export type SellerConversationSceneCode = "" | "PRE_SALE" | "AFTER_SALE";

export type SellerChatConversation = {
  conversationNo: string;
  buyerId: number;
  shopNo: string;
  shopName: string;
  buyerDisplayName: string;
  buyerAvatarUrl: string;
  shopAvatarUrl: string;
  sceneCode: SellerConversationSceneCode;
  orderNo: string;
  subOrderNo: string;
  anchorSpuNo: string;
  anchorSkuNo: string;
  unreadCount: number;
  lastMessagePreview: string;
  lastMessageAt: string;
  status: SellerConversationStatus;
  buyerReadToMessageNo: string;
  sellerReadToMessageNo: string;
  latestMessageReadByPeer: boolean;
  latestMessagePeerReadAt: string;
};

export type SellerChatProductCard = {
  title: string;
  subtitle?: string;
  coverImageUrl?: string;
  price?: string;
  priceLabel?: string;
  skuNo?: string;
  spuNo?: string;
  tags?: string[];
  detail?: string;
  link?: string;
};

export type SellerChatOrderItem = {
  title?: string;
  quantity?: string;
  price?: string;
  meta?: string;
};

export type SellerChatOrderCard = {
  orderNo?: string;
  subOrderNo?: string;
  shopNo?: string;
  status?: string;
  sceneCode?: string;
  totalAmount?: string;
  currency?: string;
  summary?: string;
  meta?: string[];
  items?: SellerChatOrderItem[];
  link?: string;
};

export type SellerChatMessageCard =
  | { kind: "PRODUCT_CARD"; product: SellerChatProductCard }
  | { kind: "ORDER_CARD"; order: SellerChatOrderCard };

export type SellerChatMessageType = "TEXT" | "IMAGE" | "PRODUCT_CARD" | "SYSTEM_NOTICE" | "UNKNOWN";

export type SellerChatMessage = {
  messageNo: string;
  conversationNo: string;
  senderType: "BUYER" | "SELLER" | "SYSTEM" | "UNKNOWN";
  senderUserId: number;
  messageType: SellerChatMessageType;
  contentText: string;
  mediaAssetId: number;
  extJson: string;
  card?: SellerChatMessageCard;
  sentAt: string;
  senderDisplayName: string;
  senderAvatarUrl: string;
  peerRead: boolean;
  peerReadAt: string;
};
