export type ConversationStatus = "ACTIVE" | "CLOSED" | "BLOCKED" | "UNKNOWN";
export type ConversationSceneCode = "" | "PRE_SALE" | "AFTER_SALE";

export type ChatConversation = {
  conversationNo: string;
  buyerId: number;
  shopNo: string;
  shopName: string;
  buyerDisplayName: string;
  buyerAvatarUrl: string;
  shopAvatarUrl: string;
  sceneCode: ConversationSceneCode;
  orderNo: string;
  subOrderNo: string;
  anchorSpuNo: string;
  anchorSkuNo: string;
  unreadCount: number;
  lastMessagePreview: string;
  lastMessageAt: string;
  status: ConversationStatus;
  buyerReadToMessageNo: string;
  sellerReadToMessageNo: string;
  latestMessageReadByPeer: boolean;
  latestMessagePeerReadAt: string;
};

export type ChatMessageType = "TEXT" | "IMAGE" | "PRODUCT_CARD" | "SYSTEM_NOTICE" | "UNKNOWN";

export type ChatProductCardPayload = {
  kind: "product";
  shopNo: string;
  spuNo: string;
  skuNo: string;
  title: string;
  skuName: string;
  imageUrl: string;
  priceCents: number;
  href: string;
};

export type ChatOrderCardPayload = {
  kind: "order";
  orderNo: string;
  subOrderNo: string;
  shopNo: string;
  statusText: string;
  totalAmountCents: number;
  itemCount: number;
  title: string;
  href: string;
};

export type ChatCardPayload = ChatProductCardPayload | ChatOrderCardPayload;

export type ChatMessage = {
  messageNo: string;
  conversationNo: string;
  senderType: "BUYER" | "SELLER" | "SYSTEM" | "UNKNOWN";
  senderUserId: number;
  messageType: ChatMessageType;
  contentText: string;
  mediaAssetId: number;
  extJson: string;
  sentAt: string;
  senderDisplayName: string;
  senderAvatarUrl: string;
  peerRead: boolean;
  peerReadAt: string;
};
