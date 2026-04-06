export type SellerConversationStatus = "ACTIVE" | "CLOSED" | "BLOCKED" | "UNKNOWN";

export type SellerChatConversation = {
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
  status: SellerConversationStatus;
  buyerReadToMessageNo: string;
  sellerReadToMessageNo: string;
  latestMessageReadByPeer: boolean;
  latestMessagePeerReadAt: string;
};

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
  sentAt: string;
  senderDisplayName: string;
  senderAvatarUrl: string;
  peerRead: boolean;
  peerReadAt: string;
};
