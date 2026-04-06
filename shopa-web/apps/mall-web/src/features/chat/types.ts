export type ConversationStatus = "ACTIVE" | "CLOSED" | "BLOCKED" | "UNKNOWN";

export type ChatConversation = {
  conversationNo: string;
  buyerId: number;
  shopNo: string;
  shopName: string;
  buyerDisplayName: string;
  buyerAvatarUrl: string;
  shopAvatarUrl: string;
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

