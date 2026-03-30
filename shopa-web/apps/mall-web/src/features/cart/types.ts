export type CartItem = {
  userId: number;
  skuNo: string;
  spuNo: string;
  shopNo: string;
  qty: number;
  checked: boolean;
  status: number;
  invalidReasonCode: string;
  spuTitle: string;
  skuName: string;
  skuImageAssetId: string;
  salePrice: number;
  marketPrice: number;
  saleAttrsJson: string;
};

export type CartSummary = {
  totalItemCount: number;
  checkedItemCount: number;
  checkedGoodsAmount: number;
  checkedPayableAmount: number;
};

export type MyCart = {
  items: CartItem[];
  summary: CartSummary;
  loadedFromBackup: boolean;
};

export type CheckoutSnapshotItem = {
  skuNo: string;
  spuNo: string;
  shopNo: string;
  qty: number;
  settlePrice: number;
  marketPrice: number;
  spuTitle: string;
  skuName: string;
  skuImageAssetId: string;
  saleAttrsJson: string;
};

export type CheckoutSnapshot = {
  checkoutToken: string;
  userId: number;
  items: CheckoutSnapshotItem[];
  goodsAmount: number;
  freightAmount: number;
  payableAmount: number;
  snapshotDigest: string;
  expireAt?: string;
  createdAt?: string;
};

