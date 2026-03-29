export const Permissions = {
  AdminOverview: "admin:me:overview",
  DashboardMallView: "dashboard:mall:view",
  DashboardShopView: "dashboard:shop:view",
  MerchantReviewView: "merchant:review:view",
  MerchantReviewApprove: "merchant:review:approve",
  MerchantReviewReject: "merchant:review:reject",
  ProductReviewView: "product:review:view",
  ProductReviewApprove: "product:review:approve",
  ProductReviewReject: "product:review:reject",
  ProductReviewFreeze: "product:review:freeze",
  ProductReviewUnfreeze: "product:review:unfreeze",
  ProductReviewForceOffShelf: "product:review:force_off_shelf",
  CSConversationView: "cs:conversation:view"
} as const;

export type PermissionKey = (typeof Permissions)[keyof typeof Permissions];
