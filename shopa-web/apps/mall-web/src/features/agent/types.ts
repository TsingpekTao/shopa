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
