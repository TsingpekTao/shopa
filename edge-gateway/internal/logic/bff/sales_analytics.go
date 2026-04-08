package bff

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	aftersalev1 "github.com/TsingpekTao/shopa/aftersale-svc/api/v1"
	sellerv1 "github.com/TsingpekTao/shopa/edge-gateway/api/seller/v1"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/consts"
	"github.com/TsingpekTao/shopa/edge-gateway/internal/service"
	orderv1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	sellershopv1 "github.com/TsingpekTao/shopa/seller-shop-svc/api/v1"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type salesBucket struct {
	Key   string
	Label string
	Start time.Time
	End   time.Time
}

type salesOrderRecord struct {
	OrderNo      string
	UserID       uint64
	PaidAt       time.Time
	PaidAmount   uint64
	SourceShopNo string
}

type salesRefundRecord struct {
	AfterSaleNo   string
	ShopNo        string
	OrderNo       string
	AppliedAt     time.Time
	RefundAmount  uint64
	AfterSaleOpen bool
}

type salesBucketAccumulator struct {
	gmv          uint64
	refundAmount uint64
	orderNos     map[string]struct{}
	userIDs      map[uint64]struct{}
}

type salesOrderMeta struct {
	UserID        uint64
	PaidAt        time.Time
	PaymentStatus orderv1.PaymentStatus
}

func (s *sBff) BuildSellerSalesAnalytics(
	ctx context.Context,
	accessToken string,
	shopNo string,
	rangeCode string,
) (*sellerv1.GetSalesAnalyticsRes, error) {
	verified, err := service.Auth().VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	if err = s.ensureClients(ctx); err != nil {
		return nil, err
	}

	loc := time.Local
	if loc == nil {
		loc = time.FixedZone("UTC+8", 8*60*60)
	}
	buckets, err := buildSalesBuckets(time.Now().In(loc), rangeCode, loc)
	if err != nil {
		return nil, err
	}

	shopNos, err := s.resolveSellerAnalyticsShopScope(ctx, verified.UserID, shopNo)
	if err != nil {
		return nil, err
	}

	res := &sellerv1.GetSalesAnalyticsRes{
		Series: make([]sellerv1.SellerSalesBucket, 0, len(buckets)),
	}
	if len(shopNos) == 0 {
		_, res.Series = aggregateSalesAnalytics(buckets, nil, nil)
		return res, nil
	}

	var (
		mu             sync.Mutex
		orderRecords   []salesOrderRecord
		refundRecords  []salesRefundRecord
		degradedFields []string
	)

	eg, egCtx := errgroup.WithContext(ctx)
	floor := buckets[0].Start

	eg.Go(func() error {
		records, loadErr := s.collectPaidOrderRecords(egCtx, shopNos, floor, loc)
		mu.Lock()
		defer mu.Unlock()
		if loadErr != nil {
			degradedFields = append(degradedFields, "sales_orders")
			return nil
		}
		orderRecords = records
		return nil
	})

	eg.Go(func() error {
		records, loadErr := s.collectRefundRecords(egCtx, shopNos, floor, loc)
		mu.Lock()
		defer mu.Unlock()
		if loadErr != nil {
			degradedFields = append(degradedFields, "sales_refunds")
			return nil
		}
		refundRecords = records
		return nil
	})

	_ = eg.Wait()

	res.Summary, res.Series = aggregateSalesAnalytics(buckets, orderRecords, refundRecords)
	res.Partial = len(degradedFields) > 0
	res.DegradedFields = uniqueStrings(degradedFields)
	return res, nil
}

func (s *sBff) resolveSellerAnalyticsShopScope(ctx context.Context, userID uint64, shopNo string) ([]string, error) {
	normalizedShopNo := strings.TrimSpace(shopNo)
	if normalizedShopNo != "" {
		ownerRes, err := s.sellerInternal.IsUserShopOwner(ctx, &sellershopv1.IsUserShopOwnerReq{
			UserId: userID,
			ShopNo: normalizedShopNo,
		})
		if err != nil {
			return nil, gerror.Wrap(err, "check shop owner failed")
		}
		if !ownerRes.GetIsOwner() {
			return nil, gerror.NewCode(consts.CodeForbidden, "forbidden")
		}
		return []string{normalizedShopNo}, nil
	}

	res, err := s.sellerInternal.ListShopsByOwnerUserId(ctx, &sellershopv1.ListShopsByOwnerUserIdReq{
		OwnerUserId: userID,
	})
	if err != nil {
		return nil, gerror.Wrap(err, "list seller shops failed")
	}

	shopNos := make([]string, 0, len(res.GetShops()))
	for _, row := range res.GetShops() {
		if value := strings.TrimSpace(row.GetShopNo()); value != "" {
			shopNos = append(shopNos, value)
		}
	}
	return uniqueStrings(shopNos), nil
}

func (s *sBff) collectPaidOrderRecords(
	ctx context.Context,
	shopNos []string,
	floor time.Time,
	loc *time.Location,
) ([]salesOrderRecord, error) {
	statuses := []orderv1.SubOrderStatus{
		orderv1.SubOrderStatus_SUB_ORDER_STATUS_PAID,
		orderv1.SubOrderStatus_SUB_ORDER_STATUS_WAIT_SHIP,
		orderv1.SubOrderStatus_SUB_ORDER_STATUS_SHIPPED,
		orderv1.SubOrderStatus_SUB_ORDER_STATUS_COMPLETED,
		orderv1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDING,
		orderv1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDED,
	}

	records := make([]salesOrderRecord, 0, 64)
	detailCache := make(map[string]salesOrderMeta)

	for _, shopNo := range shopNos {
		nextCursor := ""
		for page := 0; page < 200; page++ {
			res, err := s.orderSeller.ListShopOrders(ctx, &orderv1.ListShopOrdersReq{
				ShopNo:     shopNo,
				PageSize:   100,
				NextCursor: nextCursor,
				Statuses:   statuses,
			})
			if err != nil {
				return nil, gerror.Wrapf(err, "list paid shop orders failed for shop %s", shopNo)
			}

			rows := res.GetOrders()
			if len(rows) == 0 {
				break
			}

			stopPaging := false
			for _, sub := range rows {
				createdAt := tsToTime(sub.GetCreatedAt(), loc)
				if !createdAt.IsZero() && createdAt.Before(floor) {
					stopPaging = true
					break
				}

				orderNo := strings.TrimSpace(sub.GetOrderNo())
				if orderNo == "" {
					continue
				}

				meta, ok := detailCache[orderNo]
				if !ok {
					detail, err := s.orderSeller.GetShopOrderDetail(ctx, &orderv1.GetShopOrderDetailReq{
						SubOrderNo: sub.GetSubOrderNo(),
					})
					if err != nil {
						return nil, gerror.Wrapf(err, "get shop order detail failed for sub order %s", sub.GetSubOrderNo())
					}
					order := detail.GetOrder()
					if order == nil {
						meta = salesOrderMeta{}
					} else {
						meta = salesOrderMeta{
							UserID:        order.GetUserId(),
							PaidAt:        tsToTime(order.GetPaidAt(), loc),
							PaymentStatus: order.GetPaymentStatus(),
						}
					}
					detailCache[orderNo] = meta
				}

				if meta.PaymentStatus != orderv1.PaymentStatus_PAYMENT_STATUS_PAID || meta.PaidAt.IsZero() || meta.PaidAt.Before(floor) {
					continue
				}

				amount := uint64(0)
				if sub.GetAmount() != nil {
					amount = sub.GetAmount().GetPaidAmount()
					if amount == 0 {
						amount = sub.GetAmount().GetPayableAmount()
					}
				}

				records = append(records, salesOrderRecord{
					OrderNo:      orderNo,
					UserID:       meta.UserID,
					PaidAt:       meta.PaidAt,
					PaidAmount:   amount,
					SourceShopNo: shopNo,
				})
			}

			if stopPaging || !res.GetHasMore() || strings.TrimSpace(res.GetNextCursor()) == "" {
				break
			}
			nextCursor = res.GetNextCursor()
		}
	}

	return records, nil
}

func (s *sBff) collectRefundRecords(
	ctx context.Context,
	shopNos []string,
	floor time.Time,
	loc *time.Location,
) ([]salesRefundRecord, error) {
	statuses := []aftersalev1.AfterSaleStatus{
		aftersalev1.AfterSaleStatus_AFTER_SALE_STATUS_PENDING_SELLER_REVIEW,
		aftersalev1.AfterSaleStatus_AFTER_SALE_STATUS_SELLER_REJECTED,
		aftersalev1.AfterSaleStatus_AFTER_SALE_STATUS_WAIT_REFUND_TASK,
		aftersalev1.AfterSaleStatus_AFTER_SALE_STATUS_REFUND_PROCESSING,
		aftersalev1.AfterSaleStatus_AFTER_SALE_STATUS_REFUNDED,
		aftersalev1.AfterSaleStatus_AFTER_SALE_STATUS_CLOSED,
	}

	records := make([]salesRefundRecord, 0, 32)

	for _, shopNo := range shopNos {
		nextCursor := ""
		for page := 0; page < 200; page++ {
			res, err := s.aftersaleSeller.ListShopAfterSales(ctx, &aftersalev1.ListShopAfterSalesReq{
				ShopNo:     shopNo,
				PageSize:   100,
				NextCursor: nextCursor,
				Statuses:   statuses,
			})
			if err != nil {
				return nil, gerror.Wrapf(err, "list shop aftersales failed for shop %s", shopNo)
			}

			rows := res.GetList()
			if len(rows) == 0 {
				break
			}

			stopPaging := false
			for _, row := range rows {
				appliedAt := tsToTime(row.GetCreatedAt(), loc)
				if !appliedAt.IsZero() && appliedAt.Before(floor) {
					stopPaging = true
					break
				}
				if !includeRefundCaseForAnalytics(row.GetAfterSaleStatus()) {
					continue
				}
				records = append(records, salesRefundRecord{
					AfterSaleNo:   row.GetAfterSaleNo(),
					ShopNo:        row.GetShopNo(),
					OrderNo:       row.GetOrderNo(),
					AppliedAt:     appliedAt,
					RefundAmount:  row.GetApplyRefundAmount(),
					AfterSaleOpen: true,
				})
			}

			if stopPaging || !res.GetHasMore() || strings.TrimSpace(res.GetNextCursor()) == "" {
				break
			}
			nextCursor = res.GetNextCursor()
		}
	}

	return records, nil
}

func buildSalesBuckets(now time.Time, rangeCode string, loc *time.Location) ([]salesBucket, error) {
	if loc == nil {
		loc = time.Local
	}
	code := strings.ToUpper(strings.TrimSpace(rangeCode))
	switch code {
	case "7D", "14D", "30D":
		days := 7
		if code == "14D" {
			days = 14
		}
		if code == "30D" {
			days = 30
		}
		return buildDailySalesBuckets(now.In(loc), days), nil
	case "8W":
		return buildWeeklySalesBuckets(now.In(loc), 8), nil
	default:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "range must be one of 7D,14D,30D,8W")
	}
}

func buildDailySalesBuckets(now time.Time, days int) []salesBucket {
	currentDayStart := startOfDay(now)
	firstDayStart := currentDayStart.AddDate(0, 0, -days+1)
	out := make([]salesBucket, 0, days)
	for i := 0; i < days; i++ {
		start := firstDayStart.AddDate(0, 0, i)
		out = append(out, salesBucket{
			Key:   start.Format("2006-01-02"),
			Label: start.Format("01-02"),
			Start: start,
			End:   start.AddDate(0, 0, 1),
		})
	}
	return out
}

func buildWeeklySalesBuckets(now time.Time, weeks int) []salesBucket {
	currentWeekStart := startOfWeek(startOfDay(now))
	firstWeekStart := currentWeekStart.AddDate(0, 0, -7*(weeks-1))
	out := make([]salesBucket, 0, weeks)
	for i := 0; i < weeks; i++ {
		start := firstWeekStart.AddDate(0, 0, i*7)
		end := start.AddDate(0, 0, 7)
		out = append(out, salesBucket{
			Key:   start.Format("2006-01-02"),
			Label: fmt.Sprintf("%s-%s", start.Format("01/02"), end.AddDate(0, 0, -1).Format("01/02")),
			Start: start,
			End:   end,
		})
	}
	return out
}

func aggregateSalesAnalytics(
	buckets []salesBucket,
	orderRecords []salesOrderRecord,
	refundRecords []salesRefundRecord,
) (sellerv1.SellerSalesSummary, []sellerv1.SellerSalesBucket) {
	series := make([]sellerv1.SellerSalesBucket, len(buckets))
	accumulators := make([]salesBucketAccumulator, len(buckets))
	for i, bucket := range buckets {
		series[i] = sellerv1.SellerSalesBucket{
			BucketKey:   bucket.Key,
			BucketLabel: bucket.Label,
		}
		accumulators[i] = salesBucketAccumulator{
			orderNos: make(map[string]struct{}),
			userIDs:  make(map[uint64]struct{}),
		}
	}

	summaryOrderNos := make(map[string]struct{})
	summaryUserIDs := make(map[uint64]struct{})
	seenRefunds := make(map[string]struct{})
	summary := sellerv1.SellerSalesSummary{}

	for _, record := range orderRecords {
		index := findSalesBucketIndex(buckets, record.PaidAt)
		if index < 0 {
			continue
		}
		series[index].Gmv += record.PaidAmount
		accumulators[index].gmv += record.PaidAmount
		summary.Gmv += record.PaidAmount

		if record.OrderNo != "" {
			accumulators[index].orderNos[record.OrderNo] = struct{}{}
			summaryOrderNos[record.OrderNo] = struct{}{}
		}
		if record.UserID > 0 {
			accumulators[index].userIDs[record.UserID] = struct{}{}
			summaryUserIDs[record.UserID] = struct{}{}
		}
	}

	for _, record := range refundRecords {
		if !record.AfterSaleOpen {
			continue
		}
		if afterSaleNo := strings.TrimSpace(record.AfterSaleNo); afterSaleNo != "" {
			if _, ok := seenRefunds[afterSaleNo]; ok {
				continue
			}
			seenRefunds[afterSaleNo] = struct{}{}
		}
		index := findSalesBucketIndex(buckets, record.AppliedAt)
		if index < 0 {
			continue
		}
		series[index].RefundAmount += record.RefundAmount
		accumulators[index].refundAmount += record.RefundAmount
		summary.RefundAmount += record.RefundAmount
	}

	summary.PaidOrderCount = uint64(len(summaryOrderNos))
	summary.PaidBuyerCount = uint64(len(summaryUserIDs))
	summary.AvgOrderValue = calculateAvgOrderValue(summary.Gmv, summary.PaidOrderCount)
	summary.RefundRate = calculateRefundRate(summary.RefundAmount, summary.Gmv)

	for i := range series {
		series[i].PaidOrderCount = uint64(len(accumulators[i].orderNos))
		series[i].PaidBuyerCount = uint64(len(accumulators[i].userIDs))
		series[i].RefundRate = calculateRefundRate(series[i].RefundAmount, series[i].Gmv)
	}

	return summary, series
}

func findSalesBucketIndex(buckets []salesBucket, ts time.Time) int {
	if ts.IsZero() {
		return -1
	}
	for i, bucket := range buckets {
		if !ts.Before(bucket.Start) && ts.Before(bucket.End) {
			return i
		}
	}
	return -1
}

func includeRefundCaseForAnalytics(status aftersalev1.AfterSaleStatus) bool {
	switch status {
	case aftersalev1.AfterSaleStatus_AFTER_SALE_STATUS_UNSPECIFIED,
		aftersalev1.AfterSaleStatus_AFTER_SALE_STATUS_CANCELED:
		return false
	default:
		return true
	}
}

func calculateRefundRate(refundAmount, gmv uint64) float64 {
	if gmv == 0 {
		return 0
	}
	return float64(refundAmount) / float64(gmv)
}

func calculateAvgOrderValue(gmv, paidOrderCount uint64) uint64 {
	if paidOrderCount == 0 {
		return 0
	}
	return gmv / paidOrderCount
}

func startOfDay(ts time.Time) time.Time {
	return time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
}

func startOfWeek(ts time.Time) time.Time {
	weekday := int(ts.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return ts.AddDate(0, 0, -(weekday - 1))
}

func tsToTime(ts *timestamppb.Timestamp, loc *time.Location) time.Time {
	if ts == nil {
		return time.Time{}
	}
	if loc == nil {
		loc = time.Local
	}
	return ts.AsTime().In(loc)
}
