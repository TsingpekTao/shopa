// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/cart-svc/api/v1"
)

type (
	ICart interface {
		// AddItem 买家加购，购物车以 Redis Hash 为主存，并刷新 30 天滑动 TTL。
		AddItem(ctx context.Context, req *v1.AddItemReq) (*v1.AddItemRes, error)
		// UpdateItemQty 更新购物车中某个 SKU 的购买数量。
		UpdateItemQty(ctx context.Context, req *v1.UpdateItemQtyReq) (*v1.UpdateItemQtyRes, error)
		// ToggleItemChecked 设置单个商品勾选状态。
		ToggleItemChecked(ctx context.Context, req *v1.ToggleItemCheckedReq) (*v1.ToggleItemCheckedRes, error)
		// BatchToggleItems 批量设置勾选状态，常用于“全选/取消全选”。
		BatchToggleItems(ctx context.Context, req *v1.BatchToggleItemsReq) (*v1.BatchToggleItemsRes, error)
		// RemoveItems 从购物车中删除指定商品。
		RemoveItems(ctx context.Context, req *v1.RemoveItemsReq) (*v1.RemoveItemsRes, error)
		// ClearInvalidItems 清空失效商品。
		ClearInvalidItems(ctx context.Context, req *v1.ClearInvalidItemsReq) (*v1.ClearInvalidItemsRes, error)
		// GetMyCart 获取当前用户购物车；Redis miss 时支持从 MySQL 备份冷恢复。
		GetMyCart(ctx context.Context, req *v1.GetMyCartReq) (*v1.GetMyCartRes, error)
		// PrepareCheckout 生成结算快照 token（默认 5 分钟有效），订单侧必须使用该快照下单。
		PrepareCheckout(ctx context.Context, req *v1.PrepareCheckoutReq) (*v1.PrepareCheckoutRes, error)
		// ConsumeCheckoutToken 通过 Lua 原子脚本消费 token，避免并发重放下单。
		ConsumeCheckoutToken(ctx context.Context, req *v1.ConsumeCheckoutTokenReq) (*v1.ConsumeCheckoutTokenRes, error)
		// MarkItemsOrdered 订单创建后，清理购物车中已下单的 SKU。
		MarkItemsOrdered(ctx context.Context, req *v1.MarkItemsOrderedReq) (*v1.MarkItemsOrderedRes, error)
		// BatchUpsertSkuProjection 批量回写 SKU 展示投影（标题、价格、可售状态等）。
		BatchUpsertSkuProjection(ctx context.Context, req *v1.BatchUpsertSkuProjectionReq) (*v1.BatchUpsertSkuProjectionRes, error)
	}
)

var (
	localCart ICart
)

func Cart() ICart {
	if localCart == nil {
		panic("implement not found for interface ICart, forgot register?")
	}
	return localCart
}

func RegisterCart(i ICart) {
	localCart = i
}
