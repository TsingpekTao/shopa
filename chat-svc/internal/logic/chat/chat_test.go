package chat

import "testing"

func TestNormalizeReaderUserIDUsesSharedSellerScope(t *testing.T) {
	if got := normalizeReaderUserID(readerTypeSeller, 2040403725247320064); got != 0 {
		t.Fatalf("expected seller reader scope to use shared user id 0, got %d", got)
	}
}

func TestNormalizeReaderUserIDPreservesBuyerIdentity(t *testing.T) {
	const buyerID uint64 = 2040403725247320064
	if got := normalizeReaderUserID(readerTypeBuyer, buyerID); got != buyerID {
		t.Fatalf("expected buyer reader scope to keep user id %d, got %d", buyerID, got)
	}
}
