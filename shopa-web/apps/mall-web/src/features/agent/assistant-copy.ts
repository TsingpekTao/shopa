import {
  AssistantOrderSelectionCard,
  AssistantRecentOrderCandidate,
  AssistantSuggestedAction,
} from "./types";

export type AssistantConsoleVariant = "official" | "shopping";

export type AssistantVariantCopy = {
  chip: string;
  title: string;
  description: string;
  quickPrompts: string[];
  placeholder: string;
  emptySummary: string;
  backLabel: string;
};

export function buildActionPrompt(
  action: AssistantSuggestedAction,
  isZh: boolean,
): string {
  switch (action.actionCode) {
    case "request_refund":
      return isZh
        ? "这个订单现在适合直接退款吗"
        : "Can this order be refunded now?";
    case "request_return_refund":
      return isZh
        ? "这个订单现在更适合退货退款吗"
        : "Should this order use return and refund now?";
    case "request_exchange":
      return isZh
        ? "这个订单现在可以申请换货吗"
        : "Can I request an exchange for this order?";
    case "urge_shipment":
      return isZh
        ? "帮我看看这个订单为什么还没发货"
        : "Why has this order not shipped yet?";
    case "view_logistics":
      return isZh
        ? "帮我查一下这个订单的物流进度"
        : "Check the logistics progress for this order";
    case "wait_for_update":
      return isZh
        ? "帮我同步一下这个订单的最新进度"
        : "Share the latest update for this order";
    case "escalate_to_human":
      return isZh ? "我要转人工客服" : "I want a human agent";
    default:
      return action.label;
  }
}

export function getIntroMessageText(
  variant: AssistantConsoleVariant,
  isZh: boolean,
): string {
  if (variant === "official") {
    return isZh
      ? "您好，请问我能帮助您什么？我可以帮您查询订单状态、物流进度、退款退货建议，也可以为您转接人工客服。"
      : "Hello, how can I help you? I can check order status, logistics, refund guidance, or help escalate to a human agent.";
  }

  return isZh
    ? "您好，请问我能帮助您什么？您可以继续问我商品、订单、售后或平台规则问题。"
    : "Hello, how can I help you? You can ask about products, orders, after-sale support, or platform policies.";
}

export function getVariantCopy(
  variant: AssistantConsoleVariant,
  isZh: boolean,
): AssistantVariantCopy {
  if (variant === "official") {
    return {
      chip: isZh ? "Shopa 官方客服" : "Shopa Official Support",
      title: isZh ? "官方客服" : "Official Support",
      description: isZh
        ? "订单、物流、退款退货问题都可以直接聊，我会先理解你的问题，再帮你查订单、看物流或给出下一步建议。"
        : "Ask directly about orders, shipping, refunds, or returns. I will understand the request first, then check facts and suggest the next step.",
      quickPrompts: [
        isZh ? "帮我查一下这个订单现在到哪了" : "Check where this order is now",
        isZh ? "这个订单为什么还没发货" : "Why has this order not shipped yet?",
        isZh ? "这个订单现在能退款吗" : "Can this order be refunded now?",
        isZh ? "我要转人工客服" : "I want a human agent",
      ],
      placeholder: isZh
        ? "输入订单、物流、退款、退货、换货或人工客服相关问题..."
        : "Ask about orders, shipping, refunds, returns, exchanges or human support...",
      emptySummary: isZh
        ? "发来一条消息后，我会逐步整理当前订单事实和处理建议。"
        : "Send a message and I will organize the current order facts and next-step guidance.",
      backLabel: isZh ? "返回帮助中心" : "Back to help center",
    };
  }

  return {
    chip: isZh ? "Shopa AI 导购" : "Shopa AI Concierge",
    title: isZh ? "智能助手" : "Shopping Assistant",
    description: isZh
      ? "你可以问商品、订单、售后或平台规则问题，我会结合当前上下文继续帮你处理。"
      : "Ask about products, orders, after-sale support, or platform rules and I will continue from the current context.",
    quickPrompts: [
      isZh
        ? "帮我推荐一款适合通勤的手机"
        : "Recommend a commuter-friendly phone",
      isZh ? "这件商品什么时候能发货" : "When can this product ship?",
      isZh ? "退货退款流程怎么走" : "How does the return flow work?",
      isZh ? "我想找一款红色的商品" : "I want to find a red product",
    ],
    placeholder: isZh
      ? "输入商品、订单、退款或平台规则相关问题..."
      : "Ask about products, orders, refunds or platform policies...",
    emptySummary: isZh
      ? "先发来一条消息，我会逐步整理会话摘要。"
      : "Send a message first and I will build the session summary.",
    backLabel: isZh ? "返回服务中心" : "Back to service center",
  };
}

export function buildOrderSelectionPrompt(
  card: AssistantOrderSelectionCard,
  candidate: AssistantRecentOrderCandidate,
  isZh: boolean,
): string {
  const originalQuery = card.originalQuery.trim();
  if (originalQuery) {
    return isZh
      ? `${originalQuery}，订单号 ${candidate.orderNo}`
      : `${originalQuery}. Order number ${candidate.orderNo}`;
  }

  switch (card.taskCode) {
    case "logistics_query":
      return isZh
        ? `帮我查一下订单 ${candidate.orderNo} 现在到哪了`
        : `Check where order ${candidate.orderNo} is now`;
    case "aftersale_decision":
      return isZh
        ? `帮我判断订单 ${candidate.orderNo} 现在更适合退款还是退货退款`
        : `Help me judge whether order ${candidate.orderNo} should be refunded or returned now`;
    default:
      return isZh
        ? `帮我继续处理订单 ${candidate.orderNo}`
        : `Help me continue with order ${candidate.orderNo}`;
  }
}
