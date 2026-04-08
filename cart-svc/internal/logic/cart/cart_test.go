package cart

import (
	"context"
	"database/sql"
	"testing"

	catalogv1 "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	"google.golang.org/grpc"
)

func TestEnrichCheckoutItemsUsesCatalogSnapshots(t *testing.T) {
	oldFactory := newCatalogSnapshotClient
	defer func() {
		newCatalogSnapshotClient = oldFactory
	}()

	newCatalogSnapshotClient = func(context.Context) (catalogSnapshotClient, error) {
		return catalogSnapshotClientFunc(func(ctx context.Context, req *catalogv1.GetSkuSnapshotForOrderReq, _ ...grpc.CallOption) (*catalogv1.GetSkuSnapshotForOrderRes, error) {
			if len(req.GetSkuNos()) != 1 || req.GetSkuNos()[0] != "SKU-1" {
				t.Fatalf("expected checkout enrichment to request selected sku, got %#v", req.GetSkuNos())
			}
			return &catalogv1.GetSkuSnapshotForOrderRes{
				Snapshots: []*catalogv1.SkuOrderSnapshot{
					{
						SkuNo:           "SKU-1",
						SpuNo:           "SPU-1",
						ShopNo:          "SHOP-1",
						SpuTitle:        "测试商品",
						SkuName:         "默认规格",
						SkuImageAssetId: 1001,
						SalePrice:       1999,
						MarketPrice:     2999,
						SaleAttrs: []*catalogv1.SkuSaleAttr{
							{AttrName: "颜色", Value: "黑色"},
						},
					},
				},
			}, nil
		}), nil
	}

	items := []*redisCartItem{
		{SkuNo: "SKU-1", Qty: 2},
	}

	if err := (&sCart{}).enrichCheckoutItems(context.Background(), items); err != nil {
		t.Fatalf("expected enrichCheckoutItems to succeed, got error: %v", err)
	}
	if items[0].SalePrice != 1999 {
		t.Fatalf("expected sale price to be hydrated, got %d", items[0].SalePrice)
	}
	if items[0].ShopNo != "SHOP-1" {
		t.Fatalf("expected shop_no to be hydrated, got %q", items[0].ShopNo)
	}
	if items[0].SaleAttrsJSON == "" || items[0].SaleAttrsJSON == "[]" {
		t.Fatalf("expected sale attrs json to be hydrated, got %q", items[0].SaleAttrsJSON)
	}
}

func TestSyncCheckpointMissIsAllowed(t *testing.T) {
	if !isCartSyncCheckpointMissError(sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows to be treated as missing cart sync checkpoint")
	}
}

func TestBackfillCartItemsFromCatalogSnapshotsHydratesMissingImageAssetID(t *testing.T) {
	oldFactory := newCatalogSnapshotClient
	defer func() {
		newCatalogSnapshotClient = oldFactory
	}()

	newCatalogSnapshotClient = func(context.Context) (catalogSnapshotClient, error) {
		return catalogSnapshotClientFunc(func(ctx context.Context, req *catalogv1.GetSkuSnapshotForOrderReq, _ ...grpc.CallOption) (*catalogv1.GetSkuSnapshotForOrderRes, error) {
			return &catalogv1.GetSkuSnapshotForOrderRes{
				Snapshots: []*catalogv1.SkuOrderSnapshot{
					{
						SkuNo:           "SKU-2",
						SpuNo:           "SPU-2",
						ShopNo:          "SHOP-2",
						SpuTitle:        "测试商品二",
						SkuName:         "默认规格二",
						SkuImageAssetId: 2002,
						SalePrice:       2999,
						MarketPrice:     3999,
						SaleAttrs: []*catalogv1.SkuSaleAttr{
							{AttrName: "颜色", Value: "白色"},
						},
					},
				},
			}, nil
		}), nil
	}

	items := map[string]*redisCartItem{
		"SKU-2": {
			SkuNo:           "SKU-2",
			SpuNo:           "SPU-2",
			ShopNo:          "SHOP-2",
			Qty:             1,
			SkuImageAssetID: 0,
			SaleAttrsJSON:   "[]",
		},
	}

	changed, err := (&sCart{}).backfillCartItemsFromCatalogSnapshots(context.Background(), items)
	if err != nil {
		t.Fatalf("expected backfillCartItemsFromCatalogSnapshots to succeed, got error: %v", err)
	}
	if !changed {
		t.Fatalf("expected cart item metadata backfill to report changes")
	}
	if items["SKU-2"].SkuImageAssetID != 2002 {
		t.Fatalf("expected sku image asset id to be hydrated, got %d", items["SKU-2"].SkuImageAssetID)
	}
	if items["SKU-2"].SaleAttrsJSON == "" || items["SKU-2"].SaleAttrsJSON == "[]" {
		t.Fatalf("expected sale attrs json to be hydrated, got %q", items["SKU-2"].SaleAttrsJSON)
	}
}

type catalogSnapshotClientFunc func(ctx context.Context, req *catalogv1.GetSkuSnapshotForOrderReq, opts ...grpc.CallOption) (*catalogv1.GetSkuSnapshotForOrderRes, error)

func (f catalogSnapshotClientFunc) GetSkuSnapshotForOrder(ctx context.Context, req *catalogv1.GetSkuSnapshotForOrderReq, opts ...grpc.CallOption) (*catalogv1.GetSkuSnapshotForOrderRes, error) {
	return f(ctx, req, opts...)
}
