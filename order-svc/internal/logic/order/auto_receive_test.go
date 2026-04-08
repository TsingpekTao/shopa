package order

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

func TestShouldAutoReceiveSubOrder(t *testing.T) {
	now := gtime.NewFromTime(time.Date(2026, 4, 7, 21, 0, 0, 0, time.FixedZone("CST", 8*3600)))
	window := 7 * 24 * time.Hour

	if shouldAutoReceiveSubOrder(nil, now, window) {
		t.Fatalf("expected nil shipped_at to be ignored")
	}
	if shouldAutoReceiveSubOrder(gtime.NewFromTime(now.Time.Add(-6*24*time.Hour)), now, window) {
		t.Fatalf("expected sub order shipped within 7 days not to auto receive")
	}
	if !shouldAutoReceiveSubOrder(gtime.NewFromTime(now.Time.Add(-7*24*time.Hour)), now, window) {
		t.Fatalf("expected sub order shipped exactly 7 days ago to auto receive")
	}
	if !shouldAutoReceiveSubOrder(gtime.NewFromTime(now.Time.Add(-9*24*time.Hour)), now, window) {
		t.Fatalf("expected sub order shipped earlier than 7 days to auto receive")
	}
}

func TestRunAutoReceiveBatch(t *testing.T) {
	var captured []string

	processed, completed := runAutoReceiveBatch(context.Background(), []string{"ORD1001", "", " ORD1002 "}, func(ctx context.Context, orderNo string) (bool, error) {
		captured = append(captured, orderNo)
		return orderNo == "ORD1002", nil
	})

	if processed != 2 {
		t.Fatalf("expected 2 processed orders, got %d", processed)
	}
	if completed != 1 {
		t.Fatalf("expected 1 completed order, got %d", completed)
	}
	if len(captured) != 2 || captured[0] != "ORD1001" || captured[1] != "ORD1002" {
		t.Fatalf("expected trimmed order numbers to be forwarded, got %#v", captured)
	}
}
