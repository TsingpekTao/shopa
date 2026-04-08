package points

import "testing"

func TestCalculatePreviewDeductionCapsAtFivePercent(t *testing.T) {
	pointsUsed, cashDiscount := calculatePointsDeduction(29900, 5000, 0, 500, 1)
	if pointsUsed != 1495 {
		t.Fatalf("expected points used to cap at 1495, got %d", pointsUsed)
	}
	if cashDiscount != 1495 {
		t.Fatalf("expected cash discount to cap at 1495 cents, got %d", cashDiscount)
	}
}

func TestCalculatePreviewDeductionHonorsRequestedPoints(t *testing.T) {
	pointsUsed, cashDiscount := calculatePointsDeduction(29900, 5000, 800, 500, 1)
	if pointsUsed != 800 {
		t.Fatalf("expected requested points to be honored, got %d", pointsUsed)
	}
	if cashDiscount != 800 {
		t.Fatalf("expected cash discount to equal requested points cents, got %d", cashDiscount)
	}
}

func TestCalculateGrantedPointsRoundsDownByYuan(t *testing.T) {
	granted := calculateGrantedPoints(29950, 1)
	if granted != 299 {
		t.Fatalf("expected 299 points for 299.50 yuan paid, got %d", granted)
	}
}
