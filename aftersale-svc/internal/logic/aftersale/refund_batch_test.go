package aftersale

import (
	"database/sql"
	"testing"

	v1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
)

func TestDeriveRefundBatchStatusPendingWins(t *testing.T) {
	status := deriveRefundBatchStatus([]v1.AfterSaleStatus{
		v1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW,
		v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUNDED,
	})
	if status != v1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW {
		t.Fatalf("expected pending seller review, got %v", status)
	}
}

func TestDeriveRefundBatchStatusAllRefunded(t *testing.T) {
	status := deriveRefundBatchStatus([]v1.AfterSaleStatus{
		v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUNDED,
		v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUNDED,
	})
	if status != v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUNDED {
		t.Fatalf("expected refunded, got %v", status)
	}
}

func TestDeriveRefundBatchStatusMixedTerminalFallsBackToClosed(t *testing.T) {
	status := deriveRefundBatchStatus([]v1.AfterSaleStatus{
		v1.AfterSaleStatus_AFTER_SALE_STATUS_REFUNDED,
		v1.AfterSaleStatus_AFTER_SALE_STATUS_SELLER_REJECTED,
	})
	if status != v1.AfterSaleStatus_AFTER_SALE_STATUS_CLOSED {
		t.Fatalf("expected closed for mixed terminal batch, got %v", status)
	}
}

func TestIsRefundBatchIdempotencyMissErrorTreatsSQLNoRowsAsMiss(t *testing.T) {
	if !isRefundBatchIdempotencyMissError(sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows to be treated as idempotency miss")
	}
}
