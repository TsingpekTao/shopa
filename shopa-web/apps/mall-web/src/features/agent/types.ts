export type AssistantConversation = {
  conversationNo: string;
  botCode: string;
  sceneCode: string;
  shopNo: string;
  orderNo: string;
  anchorSpuNo: string;
  anchorSkuNo: string;
  conversationStatusCode: string;
  lastRunStatusCode: string;
  isHumanHandover: boolean;
  pendingMessageCount: number;
  sessionSummary: string;
};

export type AssistantRun = {
  runNo: string;
  conversationNo: string;
  turnNo: number;
  runStatusCode: string;
  riskDecisionCode: string;
  promptInjectionFlag: boolean;
  replyInterrupted: boolean;
  interruptReasonCode: string;
  errorCode: string;
  errorMessage: string;
  replyPayload?: AssistantReplyPayload;
};

export type AssistantMissingSlot = {
  slotCode: string;
  promptText: string;
  required: boolean;
};

export type AssistantSuggestedAction = {
  actionCode: string;
  label: string;
  enabled: boolean;
  reasonIfDisabled: string;
};

export type AssistantHiddenAction = {
  type: "SET_SLOT";
  key: "selected_order_no" | "selected_sub_order_no";
  value: string;
};

export type AssistantOrderSnapshotCard = {
  orderNo: string;
  mainStatus: string;
  paymentStatus: string;
  fulfillmentStatus: string;
  logisticsStatus: string;
  afterSaleStatus: string;
  latestUpdateTime: string;
};

export type AssistantAfterSaleDecisionCard = {
  sceneCode: string;
  decisionPathCode: string;
  reasonText: string;
  constraintText: string;
  nextStepText: string;
};

export type AssistantRecentOrderCandidate = {
  orderNo: string;
  subOrderNo: string;
  displayTitle: string;
  mainStatus: string;
  paymentStatus: string;
  fulfillmentStatus: string;
  logisticsStatus: string;
  afterSaleStatus: string;
  latestUpdateTime: string;
  selectionHint: string;
};

export type AssistantOrderSelectionCard = {
  titleText: string;
  helperText: string;
  originalQuery: string;
  taskCode: string;
  candidates: AssistantRecentOrderCandidate[];
};

export type AssistantLogisticsTrackingCard = {
  orderNo: string;
  subOrderNo: string;
  fulfillmentStatus: string;
  logisticsStatus: string;
  latestUpdateTime: string;
  latestTraceText: string;
  timelineSummary: string;
};

export type AssistantProductRecommendationItem = {
  spuNo: string;
  title: string;
  coverUrl: string;
  priceText: string;
  shopName: string;
  reasonText: string;
};

export type AssistantProductRecommendationCard = {
  titleText: string;
  helperText: string;
  items: AssistantProductRecommendationItem[];
};

export type AssistantDataCard = {
  orderSnapshotCard?: AssistantOrderSnapshotCard;
  afterSaleDecisionCard?: AssistantAfterSaleDecisionCard;
  orderSelectionCard?: AssistantOrderSelectionCard;
  logisticsTrackingCard?: AssistantLogisticsTrackingCard;
  productRecommendationCard?: AssistantProductRecommendationCard;
};

export type AssistantReplyPayload = {
  replyText: string;
  intentCode: string;
  guardResultCode: string;
  confidence: number;
  dataCards: AssistantDataCard[];
  suggestedActions: AssistantSuggestedAction[];
  missingSlots: AssistantMissingSlot[];
  slotRetryCount: number;
  handoffRecommended: boolean;
  handoffReasonCode: string;
};

export type AssistantMessage = {
  messageNo: string;
  conversationNo: string;
  runNo: string;
  senderTypeCode: string;
  messageTypeCode: string;
  contentText: string;
  assetIds: number[];
  extJson: string;
  interrupted: boolean;
  sentAt: string;
};

export type AnswerSource = {
  sourceTypeCode: string;
  sourceId: string;
  sourceVersion: number;
  title: string;
  snippet: string;
};

export type HandoffTicket = {
  ticketNo: string;
  conversationNo: string;
  escalationReasonCode: string;
  statusCode: string;
  handoffSummary: string;
};
