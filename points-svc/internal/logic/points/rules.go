package points

func calculatePointsDeduction(payableAmountCent, availablePoints, requestedPoints uint64, maxDeductionRateBps uint32, deductPointsPerCent uint64) (uint64, uint64) {
	if payableAmountCent == 0 || availablePoints == 0 {
		return 0, 0
	}

	pointsPerCent := deductPointsPerCent
	if pointsPerCent == 0 {
		pointsPerCent = 1
	}

	requested := requestedPoints
	if requested == 0 || requested > availablePoints {
		requested = availablePoints
	}

	maxCashByRule := payableAmountCent * uint64(maxDeductionRateBps) / 10000
	if maxCashByRule == 0 || maxCashByRule > payableAmountCent {
		maxCashByRule = payableAmountCent
	}

	maxPointsByRule := maxCashByRule * pointsPerCent
	pointsUsed := requested
	if pointsUsed > maxPointsByRule {
		pointsUsed = maxPointsByRule
	}
	if pointsUsed == 0 {
		return 0, 0
	}

	return pointsUsed, pointsUsed / pointsPerCent
}

func calculateGrantedPoints(paidAmountCent, grantPointsPerYuan uint64) uint64 {
	if paidAmountCent == 0 || grantPointsPerYuan == 0 {
		return 0
	}
	return paidAmountCent * grantPointsPerYuan / 100
}
