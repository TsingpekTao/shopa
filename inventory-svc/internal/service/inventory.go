// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
)

type (
	IInventory interface {
		// BatchAdjustMySkuStock 卖家侧批量调整库存，要求 SKU 必须归属当前 shop_no。
		BatchAdjustMySkuStock(ctx context.Context, req *v1.BatchAdjustMySkuStockReq) (*v1.BatchAdjustMySkuStockRes, error)
		// ReserveStock 订单侧预占库存，支持返回失败明细，便于前端定位哪个 SKU 无货。
		ReserveStock(ctx context.Context, req *v1.ReserveStockReq) (*v1.ReserveStockRes, error)
		// ConfirmReservation 订单支付成功后确认预占，正式扣减 total 并释放 locked。
		ConfirmReservation(ctx context.Context, req *v1.ConfirmReservationReq) (*v1.ConfirmReservationRes, error)
		// CancelReservation 订单取消或超时后释放预占库存，回补 available。
		CancelReservation(ctx context.Context, req *v1.CancelReservationReq) (*v1.CancelReservationRes, error)
		// BatchAdjustStockByAdmin 管理员批量调整库存，可补齐或修正 SKU 上下文。
		BatchAdjustStockByAdmin(ctx context.Context, req *v1.BatchAdjustStockByAdminReq) (*v1.BatchAdjustStockByAdminRes, error)
		// SetHotSku 设置或取消热点 SKU 标记，为热点治理和流量策略提供依据。
		SetHotSku(ctx context.Context, req *v1.SetHotSkuReq) (*v1.SetHotSkuRes, error)
		// GetSkuInventory 查询单个 SKU 库存快照，常用于详情页或后台库存看板。
		GetSkuInventory(ctx context.Context, req *v1.GetSkuInventoryReq) (*v1.GetSkuInventoryRes, error)
		// BatchGetSkuInventory 批量查询 SKU 库存快照，适合购物车/结算页批量拉取。
		BatchGetSkuInventory(ctx context.Context, req *v1.BatchGetSkuInventoryReq) (*v1.BatchGetSkuInventoryRes, error)
		// UpsertSkuContext 同步 SKU 与 SPU/店铺上下文，用于库存鉴权与事件补全字段。
		UpsertSkuContext(ctx context.Context, req *v1.UpsertSkuContextReq) (*v1.UpsertSkuContextRes, error)
	}
)

var (
	localInventory IInventory
)

func Inventory() IInventory {
	if localInventory == nil {
		panic("implement not found for interface IInventory, forgot register?")
	}
	return localInventory
}

// RegisterInventory 注册 Inventory 领域服务实现，供 controller 层通过 service.Inventory() 统一访问。
func RegisterInventory(i IInventory) {
	localInventory = i
}
