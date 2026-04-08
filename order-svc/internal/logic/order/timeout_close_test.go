package order

import (
	"context"
	"errors"
	"testing"
)

func TestRunExpiredUnpaidOrderBatchContinuesAfterCloseError(t *testing.T) {
	var called []string

	processed, closed := runExpiredUnpaidOrderBatch(context.Background(), []string{"", "ORD-1", "ORD-2", "  "}, func(ctx context.Context, orderNo string) (bool, error) {
		called = append(called, orderNo)
		if orderNo == "ORD-1" {
			return false, errors.New("close failed")
		}
		return true, nil
	})

	if processed != 2 {
		t.Fatalf("expected 2 non-empty orders processed, got %d", processed)
	}
	if closed != 1 {
		t.Fatalf("expected 1 order closed, got %d", closed)
	}
	if len(called) != 2 || called[0] != "ORD-1" || called[1] != "ORD-2" {
		t.Fatalf("expected both order numbers to be attempted in order, got %#v", called)
	}
}

func TestNormalizeUnpaidCloseBatchSizeFallsBackToWorkerDefault(t *testing.T) {
	if got := normalizeUnpaidCloseBatchSize(0, unpaidCloseWorkerConf{BatchSize: 25}); got != 25 {
		t.Fatalf("expected zero limit to fall back to worker batch size 25, got %d", got)
	}
	if got := normalizeUnpaidCloseBatchSize(-1, unpaidCloseWorkerConf{BatchSize: 25}); got != 25 {
		t.Fatalf("expected negative limit to fall back to worker batch size 25, got %d", got)
	}
	if got := normalizeUnpaidCloseBatchSize(8, unpaidCloseWorkerConf{BatchSize: 25}); got != 8 {
		t.Fatalf("expected explicit limit 8 to win, got %d", got)
	}
}
