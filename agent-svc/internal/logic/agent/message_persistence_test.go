package agent

import (
	"testing"

	"github.com/gogf/gf/v2/os/gtime"
)

func TestBuildInternalMessageDOUsesMessageNoAsClientMessageNo(t *testing.T) {
	t.Parallel()

	now := gtime.Now()
	data := buildInternalMessageDO(
		"AMSG202604101700000001",
		"ACV20260410145439326500",
		"ARN20260410162223920700",
		3,
		"AGENT",
		"TEXT",
		"我先把你最近的订单拉出来，点一个我马上继续查物流和发货进度。",
		"",
		0,
		now,
	)

	if got := data.ClientMessageNo; got != "AMSG202604101700000001" {
		t.Fatalf("expected assistant client_message_no to reuse message_no, got %#v", got)
	}
	if got := data.MessageNo; got != "AMSG202604101700000001" {
		t.Fatalf("expected message_no to be preserved, got %#v", got)
	}
	if got := data.ConversationNo; got != "ACV20260410145439326500" {
		t.Fatalf("expected conversation_no to be preserved, got %#v", got)
	}
	if got := data.RunNo; got != "ARN20260410162223920700" {
		t.Fatalf("expected run_no to be preserved, got %#v", got)
	}
	if got := data.ReplyToTurnNo; got != uint64(3) {
		t.Fatalf("expected reply_to_turn_no to be preserved, got %#v", got)
	}
	if got := data.ExtJson; got != "{}" {
		t.Fatalf("expected empty ext_json to normalize to empty object, got %#v", got)
	}
	if data.SentAt != now {
		t.Fatalf("expected sent_at pointer to be reused")
	}
}

func TestBuildInternalMessageDOSupportsSystemQueueNotice(t *testing.T) {
	t.Parallel()

	now := gtime.Now()
	data := buildInternalMessageDO(
		"AMSG202604101700000002",
		"ACV20260410145439326500",
		"",
		0,
		"SYSTEM",
		"QUEUE_NOTICE",
		"当前人工客服繁忙，请稍后查看排队进度。",
		`{"queue_blocked":true}`,
		0,
		now,
	)

	if got := data.ClientMessageNo; got != "AMSG202604101700000002" {
		t.Fatalf("expected system client_message_no to reuse message_no, got %#v", got)
	}
	if got := data.SenderTypeCode; got != "SYSTEM" {
		t.Fatalf("expected sender type SYSTEM, got %#v", got)
	}
	if got := data.MessageTypeCode; got != "QUEUE_NOTICE" {
		t.Fatalf("expected message type QUEUE_NOTICE, got %#v", got)
	}
	if got := data.ExtJson; got != `{"queue_blocked":true}` {
		t.Fatalf("expected queue notice ext_json to be preserved, got %#v", got)
	}
	if data.CreatedAt != now || data.UpdatedAt != now {
		t.Fatalf("expected audit timestamps to reuse sent_at")
	}
}
