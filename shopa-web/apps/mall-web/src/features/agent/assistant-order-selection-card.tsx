import React from "react";
import { formatAssistantCandidateUpdateTime } from "./assistant-display";
import {
  AssistantOrderSelectionCard,
  AssistantRecentOrderCandidate,
} from "./types";

type AssistantOrderSelectionCardViewProps = {
  card: AssistantOrderSelectionCard;
  isZh: boolean;
  onSelect: (candidate: AssistantRecentOrderCandidate) => void;
};

export function AssistantOrderSelectionCardView({
  card,
  isZh,
  onSelect,
}: AssistantOrderSelectionCardViewProps) {
  return (
    <article className="tb-agent-structured-card">
      <strong>
        {card.titleText || (isZh ? "请选择一个订单" : "Choose an order")}
      </strong>
      {card.helperText ? <p>{card.helperText}</p> : null}
      <div className="tb-agent-order-choice-list">
        {card.candidates.map((candidate) => (
          <button
            key={`${candidate.orderNo}-${candidate.subOrderNo}`}
            type="button"
            className="tb-agent-order-choice"
            onClick={() => onSelect(candidate)}
          >
            <div>
              <strong>{candidate.displayTitle || candidate.orderNo}</strong>
              <span>{candidate.orderNo}</span>
            </div>
            <div>
              <span>
                {candidate.selectionHint ||
                  candidate.fulfillmentStatus ||
                  candidate.logisticsStatus ||
                  "--"}
              </span>
              <small>
                {formatAssistantCandidateUpdateTime(
                  candidate.latestUpdateTime,
                  isZh ? "zh-CN" : "en-US",
                )}
              </small>
            </div>
          </button>
        ))}
      </div>
    </article>
  );
}
