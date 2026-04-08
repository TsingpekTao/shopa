package bff

import (
	"testing"
	"time"

	sellerv1 "github.com/TsingpekTao/shopa/edge-gateway/api/seller/v1"
)

func TestBuildSalesBucketsDailyRange(t *testing.T) {
	loc := time.FixedZone("CST", 8*60*60)
	now := time.Date(2026, 4, 7, 15, 30, 0, 0, loc)

	buckets, err := buildSalesBuckets(now, "7D", loc)
	if err != nil {
		t.Fatalf("buildSalesBuckets returned error: %v", err)
	}
	if len(buckets) != 7 {
		t.Fatalf("expected 7 daily buckets, got %d", len(buckets))
	}
	if buckets[0].Key != "2026-04-01" {
		t.Fatalf("unexpected first bucket key: %s", buckets[0].Key)
	}
	if buckets[6].Key != "2026-04-07" {
		t.Fatalf("unexpected last bucket key: %s", buckets[6].Key)
	}
	if buckets[0].Label != "04-01" {
		t.Fatalf("unexpected first bucket label: %s", buckets[0].Label)
	}
}

func TestBuildSalesBucketsWeeklyRange(t *testing.T) {
	loc := time.FixedZone("CST", 8*60*60)
	now := time.Date(2026, 4, 7, 15, 30, 0, 0, loc)

	buckets, err := buildSalesBuckets(now, "8W", loc)
	if err != nil {
		t.Fatalf("buildSalesBuckets returned error: %v", err)
	}
	if len(buckets) != 8 {
		t.Fatalf("expected 8 weekly buckets, got %d", len(buckets))
	}
	if buckets[0].Start.Weekday() != time.Monday {
		t.Fatalf("expected weekly bucket to start on Monday, got %s", buckets[0].Start.Weekday())
	}
	if buckets[7].Key != "2026-04-06" {
		t.Fatalf("unexpected latest weekly bucket key: %s", buckets[7].Key)
	}
}

func TestAggregateSalesAnalyticsUsesPaidAndRefundSemantics(t *testing.T) {
	loc := time.FixedZone("CST", 8*60*60)
	buckets, err := buildSalesBuckets(time.Date(2026, 4, 7, 15, 30, 0, 0, loc), "7D", loc)
	if err != nil {
		t.Fatalf("buildSalesBuckets returned error: %v", err)
	}

	summary, series := aggregateSalesAnalytics(
		buckets,
		[]salesOrderRecord{
			{
				OrderNo:      "ORD-1",
				UserID:       101,
				PaidAt:       time.Date(2026, 4, 6, 10, 0, 0, 0, loc),
				PaidAmount:   10000,
				SourceShopNo: "SHOP-1",
			},
			{
				OrderNo:      "ORD-1",
				UserID:       101,
				PaidAt:       time.Date(2026, 4, 6, 10, 1, 0, 0, loc),
				PaidAmount:   2000,
				SourceShopNo: "SHOP-1",
			},
			{
				OrderNo:      "ORD-2",
				UserID:       101,
				PaidAt:       time.Date(2026, 4, 7, 9, 0, 0, 0, loc),
				PaidAmount:   8000,
				SourceShopNo: "SHOP-1",
			},
			{
				OrderNo:      "ORD-3",
				UserID:       202,
				PaidAt:       time.Date(2026, 4, 7, 12, 0, 0, 0, loc),
				PaidAmount:   4000,
				SourceShopNo: "SHOP-2",
			},
		},
		[]salesRefundRecord{
			{
				AfterSaleNo:   "AS-1",
				ShopNo:        "SHOP-1",
				OrderNo:       "ORD-2",
				AppliedAt:     time.Date(2026, 4, 7, 13, 0, 0, 0, loc),
				RefundAmount:  1000,
				AfterSaleOpen: true,
			},
			{
				AfterSaleNo:   "AS-CANCELED",
				ShopNo:        "SHOP-1",
				OrderNo:       "ORD-2",
				AppliedAt:     time.Date(2026, 4, 7, 14, 0, 0, 0, loc),
				RefundAmount:  999,
				AfterSaleOpen: false,
			},
		},
	)

	if summary.Gmv != 24000 {
		t.Fatalf("expected GMV 24000, got %d", summary.Gmv)
	}
	if summary.PaidOrderCount != 3 {
		t.Fatalf("expected paid order count 3, got %d", summary.PaidOrderCount)
	}
	if summary.PaidBuyerCount != 2 {
		t.Fatalf("expected paid buyer count 2, got %d", summary.PaidBuyerCount)
	}
	if summary.RefundAmount != 1000 {
		t.Fatalf("expected refund amount 1000, got %d", summary.RefundAmount)
	}
	if summary.AvgOrderValue != 8000 {
		t.Fatalf("expected avg order value 8000, got %d", summary.AvgOrderValue)
	}
	if !nearlyEqual(summary.RefundRate, 1000.0/24000.0) {
		t.Fatalf("expected refund rate 0.041666..., got %v", summary.RefundRate)
	}

	var bucket0407 sellerv1.SellerSalesBucket
	for _, bucket := range series {
		if bucket.BucketKey == "2026-04-07" {
			bucket0407 = bucket
			break
		}
	}
	if bucket0407.BucketKey == "" {
		t.Fatalf("expected to find bucket for 2026-04-07")
	}
	if bucket0407.PaidOrderCount != 2 {
		t.Fatalf("expected bucket paid order count 2, got %d", bucket0407.PaidOrderCount)
	}
	if bucket0407.PaidBuyerCount != 2 {
		t.Fatalf("expected bucket paid buyer count 2, got %d", bucket0407.PaidBuyerCount)
	}
	if bucket0407.RefundAmount != 1000 {
		t.Fatalf("expected bucket refund amount 1000, got %d", bucket0407.RefundAmount)
	}
}

func TestAggregateSalesAnalyticsRefundRateFallsBackToZeroWhenNoGMV(t *testing.T) {
	loc := time.FixedZone("CST", 8*60*60)
	buckets, err := buildSalesBuckets(time.Date(2026, 4, 7, 15, 30, 0, 0, loc), "7D", loc)
	if err != nil {
		t.Fatalf("buildSalesBuckets returned error: %v", err)
	}

	summary, series := aggregateSalesAnalytics(
		buckets,
		nil,
		[]salesRefundRecord{
			{
				AfterSaleNo:   "AS-1",
				ShopNo:        "SHOP-1",
				OrderNo:       "ORD-1",
				AppliedAt:     time.Date(2026, 4, 7, 13, 0, 0, 0, loc),
				RefundAmount:  1000,
				AfterSaleOpen: true,
			},
		},
	)

	if summary.RefundRate != 0 {
		t.Fatalf("expected summary refund rate 0, got %v", summary.RefundRate)
	}

	for _, bucket := range series {
		if bucket.BucketKey == "2026-04-07" && bucket.RefundRate != 0 {
			t.Fatalf("expected bucket refund rate 0, got %v", bucket.RefundRate)
		}
	}
}

func nearlyEqual(a, b float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < 0.0000001
}
