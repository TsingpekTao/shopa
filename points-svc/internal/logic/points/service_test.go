package points

import (
	"context"
	"testing"
	"time"

	httpv1 "github.com/TsingpekTao/shopa/points-svc/api/points/v1"
	"github.com/TsingpekTao/shopa/points-svc/internal/dao"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func TestGrantOrderHTTPAllowsEmptyExtraJSON(t *testing.T) {
	ctx := context.Background()
	svc := New()
	userID := uint64(time.Now().UnixNano())

	if _, _, err := svc.initAccountIfAbsent(ctx, userID); err != nil {
		t.Fatalf("init account failed: %v", err)
	}

	orderNo := "TEST_GRANT_EMPTY_EXTRA_" + time.Now().Format("20060102150405.000000000")
	_, err := svc.grantOrderHTTP(ctx, &httpv1.GrantOrderReq{
		OrderNo:        orderNo,
		UserID:         userID,
		PaidAmount:     100,
		IdempotencyKey: "test-" + orderNo,
		RequestSource:  "unit-test",
	})
	if err != nil {
		t.Fatalf("grantOrderHTTP returned error: %v", err)
	}

	_, _ = dao.PointsGrantDetail.Ctx(ctx).
		Where(dao.PointsGrantDetail.Columns().OrderNo, orderNo).
		Delete()
	_, _ = dao.PointsLedger.Ctx(ctx).
		Where(dao.PointsLedger.Columns().BizType, "ORDER_GRANT").
		Where(dao.PointsLedger.Columns().BizNo, orderNo).
		Delete()
}
