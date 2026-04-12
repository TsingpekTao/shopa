const PROTOBUF_MAP_SECONDS_PATTERN = /seconds:([0-9eE+.-]+)/;
const PROTOBUF_MAP_NANOS_PATTERN = /nanos:([0-9eE+.-]+)/;

function parseAssistantTimestamp(raw: string): Date | null {
  if (!raw) {
    return null;
  }

  const normalized = raw.trim();
  if (!normalized) {
    return null;
  }

  const directInput = /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/.test(normalized)
    ? normalized.replace(" ", "T")
    : normalized;
  const directDate = new Date(directInput);
  if (!Number.isNaN(directDate.getTime())) {
    return directDate;
  }

  const secondsMatch = normalized.match(PROTOBUF_MAP_SECONDS_PATTERN);
  if (!secondsMatch) {
    return null;
  }

  const seconds = Number(secondsMatch[1]);
  const nanosMatch = normalized.match(PROTOBUF_MAP_NANOS_PATTERN);
  const nanos = nanosMatch ? Number(nanosMatch[1]) : 0;
  if (!Number.isFinite(seconds) || !Number.isFinite(nanos)) {
    return null;
  }

  const candidate = new Date(seconds * 1000 + nanos / 1_000_000);
  if (Number.isNaN(candidate.getTime())) {
    return null;
  }
  return candidate;
}

export function formatAssistantCandidateUpdateTime(
  raw: string,
  locale: string,
): string {
  const parsed = parseAssistantTimestamp(raw);
  if (!parsed) {
    return "--";
  }

  return parsed.toLocaleString(locale === "zh-CN" ? "zh-CN" : "en-US", {
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function shouldRenderSidebarStructuredCards(input: {
  hasOrderSnapshotCard: boolean;
  hasDecisionCard: boolean;
  hasLogisticsCard?: boolean;
  hasProductRecommendationCard?: boolean;
}): boolean {
  return Boolean(
    input.hasOrderSnapshotCard ||
    input.hasDecisionCard ||
    input.hasLogisticsCard ||
    input.hasProductRecommendationCard,
  );
}
