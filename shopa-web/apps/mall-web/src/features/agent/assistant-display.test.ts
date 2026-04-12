import assert from "node:assert/strict";
import { test } from "vitest";
import {
  formatAssistantCandidateUpdateTime,
  shouldRenderSidebarStructuredCards,
} from "./assistant-display.js";

test("sidebar structured cards ignore order-selection-only payloads", () => {
  assert.equal(
    shouldRenderSidebarStructuredCards({
      hasOrderSnapshotCard: false,
      hasDecisionCard: false,
    }),
    false,
  );

  assert.equal(
    shouldRenderSidebarStructuredCards({
      hasOrderSnapshotCard: true,
      hasDecisionCard: false,
    }),
    true,
  );

  assert.equal(
    shouldRenderSidebarStructuredCards({
      hasOrderSnapshotCard: false,
      hasDecisionCard: true,
    }),
    true,
  );
});

test("candidate update time formats protobuf-like map timestamp strings", () => {
  const formatted = formatAssistantCandidateUpdateTime(
    "map[nanos:6.47e+08 seconds:1.775644071e+09]",
    "zh-CN",
  );

  assert.doesNotMatch(formatted, /map\[/);
  assert.notEqual(formatted, "--");
});

test("candidate update time falls back safely for invalid values", () => {
  assert.equal(formatAssistantCandidateUpdateTime("", "zh-CN"), "--");
  assert.equal(formatAssistantCandidateUpdateTime("nonsense", "zh-CN"), "--");
});
