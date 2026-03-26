// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// 你可以删除这些注释，并按需手动维护该接口文件。
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
)

type (
	IInventory interface {
		// BatchAdjustMySkuStock 商家批量调整自己店铺下 SKU 的库存。
		BatchAdjustMySkuStock(ctx context.Context, req *v1.BatchAdjustMySkuStockReq) (*v1.BatchAdjustMySkuStockRes, error)
		// ReserveStock 订单创建阶段预占库存。
		ReserveStock(ctx context.Context, req *v1.ReserveStockReq) (*v1.ReserveStockRes, error)
		// ConfirmReservation 支付成功后确认预占，扣减锁定库存。
		ConfirmReservation(ctx context.Context, req *v1.ConfirmReservationReq) (*v1.ConfirmReservationRes, error)
		// CancelReservation 取消订单或超时未支付时释放预占库存。
		CancelReservation(ctx context.Context, req *v1.CancelReservationReq) (*v1.CancelReservationRes, error)
		// BatchAdjustStockByAdmin 管理员批量调整库存。
		BatchAdjustStockByAdmin(ctx context.Context, req *v1.BatchAdjustStockByAdminReq) (*v1.BatchAdjustStockByAdminRes, error)
		// SetHotSku 标记或取消热点 SKU。
		SetHotSku(ctx context.Context, req *v1.SetHotSkuReq) (*v1.SetHotSkuRes, error)
		// GetSkuInventory 查询单个 SKU 库存详情。
		GetSkuInventory(ctx context.Context, req *v1.GetSkuInventoryReq) (*v1.GetSkuInventoryRes, error)
		// BatchGetSkuInventory 批量查询 SKU 库存详情。
		BatchGetSkuInventory(ctx context.Context, req *v1.BatchGetSkuInventoryReq) (*v1.BatchGetSkuInventoryRes, error)
		// UpsertSkuContext 写入或更新 SKU 所属上下文（SPU、店铺等）。
		UpsertSkuContext(ctx context.Context, req *v1.UpsertSkuContextReq) (*v1.UpsertSkuContextRes, error)
	}
)

var (
	localInventory IInventory
)

// Inventory 返回已注册的库存服务实现。
func Inventory() IInventory {
	if localInventory == nil {
		panic("implement not found for interface IInventory, forgot register?")
	}
	return localInventory
}

// RegisterInventory 注册库存服务实现，通常在 internal/logic 的 init 中调用。
func RegisterInventory(i IInventory) {
	localInventory = i
}
