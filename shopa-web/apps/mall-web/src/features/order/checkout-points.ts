export const POINTS_PER_YUAN = 100;
export const POINTS_PER_CENT = 1;
export const MAX_POINTS_DEDUCTION_RATE_BPS = 500;

export type CheckoutPointsSummary = {
  availablePoints: number;
  maxDiscountAmount: number;
  maxUsablePoints: number;
  selectedPoints: number;
  selectedDiscountAmount: number;
  finalPayableAmount: number;
  canToggle: boolean;
};

export type OrderAmountBreakdown = {
  goodsAmount: number;
  freightAmount: number;
  discountAmount: number;
  pointsDiscountAmount: number;
  payableAmount: number;
  paidAmount: number;
  originalPayableAmount: number;
  totalDiscountAmount: number;
};

export function buildCheckoutPointsSummary(input: {
  payableAmount: number;
  availablePoints: number;
  usePoints: boolean;
}): CheckoutPointsSummary {
  const payableAmount = Math.max(0, Math.floor(input.payableAmount));
  const availablePoints = Math.max(0, Math.floor(input.availablePoints));
  const maxDiscountAmount = Math.floor((payableAmount * MAX_POINTS_DEDUCTION_RATE_BPS) / 10000);
  const maxUsablePoints = Math.min(availablePoints, maxDiscountAmount * POINTS_PER_CENT);
  const selectedPoints = input.usePoints ? maxUsablePoints : 0;
  const selectedDiscountAmount = input.usePoints ? Math.floor(selectedPoints / POINTS_PER_CENT) : 0;
  return {
    availablePoints,
    maxDiscountAmount,
    maxUsablePoints,
    selectedPoints,
    selectedDiscountAmount,
    finalPayableAmount: Math.max(0, payableAmount - selectedDiscountAmount),
    canToggle: maxUsablePoints > 0
  };
}

export function buildOrderAmountBreakdown(input: {
  goodsAmount: number;
  freightAmount: number;
  discountAmount: number;
  pointsDiscountAmount: number;
  payableAmount: number;
  paidAmount: number;
}): OrderAmountBreakdown {
  const goodsAmount = Math.max(0, Math.floor(input.goodsAmount));
  const freightAmount = Math.max(0, Math.floor(input.freightAmount));
  const discountAmount = Math.max(0, Math.floor(input.discountAmount));
  const pointsDiscountAmount = Math.max(0, Math.floor(input.pointsDiscountAmount));
  const payableAmount = Math.max(0, Math.floor(input.payableAmount));
  const paidAmount = Math.max(0, Math.floor(input.paidAmount));
  const totalDiscountAmount = discountAmount + pointsDiscountAmount;
  return {
    goodsAmount,
    freightAmount,
    discountAmount,
    pointsDiscountAmount,
    payableAmount,
    paidAmount,
    originalPayableAmount: Math.max(payableAmount + totalDiscountAmount, goodsAmount + freightAmount),
    totalDiscountAmount
  };
}
