package inventory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
	"github.com/TsingpekTao/shopa/inventory-svc/internal/dao"
	"github.com/TsingpekTao/shopa/inventory-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/inventory-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/inventory-svc/internal/service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sInventory struct{}

// New 创建 Inventory 领域服务实例。
func New() *sInventory {
	return &sInventory{}
}

// init 在包加载时注册 Inventory 服务实现到 service 层。
func init() {
	service.RegisterInventory(New())
}

// BatchAdjustMySkuStock 卖家批量调库存，先校验 SKU 归属，再执行库存调整事务。
func (s *sInventory) BatchAdjustMySkuStock(ctx context.Context, req *v1.BatchAdjustMySkuStockReq) (*v1.BatchAdjustMySkuStockRes, error) {
	if strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	if len(req.GetItems()) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "items is required")
	}

	results := make([]*v1.AdjustResultItem, 0, len(req.GetItems()))
	err := dao.InventoryStock.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, item := range req.GetItems() {
			// 每个请求项都必须显式指定 sku_no，避免把增减量写到错误的库存行。
			if strings.TrimSpace(item.GetSkuNo()) == "" {
				return gerror.NewCode(gcode.CodeInvalidParameter, "sku_no is required")
			}
			// 先读取 SKU 上下文（spu/shop 归属），这是库存写入前的第一道鉴权防线。
			ctxRow, err := s.getSkuContextTx(ctx, tx, item.GetSkuNo())
			if err != nil {
				return err
			}
			if ctxRow == nil {
				if strings.TrimSpace(item.GetSpuNo()) == "" {
					return gerror.NewCodef(gcode.CodeNotFound, "sku context not found for %s", item.GetSkuNo())
				}
				if err := s.upsertSkuContextItemTx(ctx, tx, &v1.SkuContextItem{
					SkuNo:   item.GetSkuNo(),
					SpuNo:   item.GetSpuNo(),
					ShopNo:  req.GetShopNo(),
					Enabled: true,
				}); err != nil {
					return err
				}
				ctxRow, err = s.getSkuContextTx(ctx, tx, item.GetSkuNo())
				if err != nil {
					return err
				}
				if ctxRow == nil {
					return gerror.NewCodef(gcode.CodeNotFound, "sku context not found for %s", item.GetSkuNo())
				}
			}
			if ctxRow.ShopNo != req.GetShopNo() {
				return gerror.NewCodef(gcode.CodeInvalidParameter, "sku %s does not belong to shop %s", item.GetSkuNo(), req.GetShopNo())
			}
			if item.GetSpuNo() != "" && item.GetSpuNo() != ctxRow.SpuNo {
				return gerror.NewCodef(gcode.CodeInvalidParameter, "sku %s spu_no mismatch", item.GetSkuNo())
			}
			// 复用统一的库存增减函数，确保 total/available/version 的维护逻辑一致。
			adjusted, err := s.adjustStockTx(ctx, tx, item.GetSkuNo(), ctxRow.SpuNo, ctxRow.ShopNo, item.GetDeltaTotalQty())
			if err != nil {
				return err
			}
			results = append(results, adjusted)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logCatalogProjectionSyncFailure(ctx, "seller stock adjust", s.syncCatalogProjectionByAdjustResults(ctx, results))
	return &v1.BatchAdjustMySkuStockRes{Results: results}, nil
}

// ReserveStock 订单预占库存，支持批量校验与失败明细返回（ALL_OR_NOTHING 模式可整体回滚）。
func (s *sInventory) ReserveStock(ctx context.Context, req *v1.ReserveStockReq) (*v1.ReserveStockRes, error) {
	if strings.TrimSpace(req.GetReservationNo()) == "" || strings.TrimSpace(req.GetOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "reservation_no and order_no are required")
	}
	if len(req.GetItems()) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "items is required")
	}

	var (
		// response 在事务内组装，事务提交后统一返回，确保返回内容与落库状态一致。
		response *v1.ReserveStockRes
		// 默认模式为“要么全部成功，要么全部失败并回滚”。
		defaultMode = v1.ReserveMode_RESERVE_MODE_ALL_OR_NOTHING
		// reservationMode 允许调用方按场景切换为“部分成功可接受”。
		reservationMode = req.GetMode()
	)
	if reservationMode == v1.ReserveMode_RESERVE_MODE_UNSPECIFIED {
		reservationMode = defaultMode
	}

	err := dao.InventoryStock.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 幂等保护：若 reservation_no 已存在，直接回放历史结果，避免重复扣减库存。
		existing, err := s.getReservationEntityTx(ctx, tx, req.GetReservationNo())
		if err != nil {
			return err
		}
		if existing != nil {
			record, err := s.loadReservationRecordTx(ctx, tx, req.GetReservationNo())
			if err != nil {
				return err
			}
			response = &v1.ReserveStockRes{
				Success:     record.GetStatus() == v1.ReservationStatus_RESERVATION_STATUS_RESERVED,
				Reservation: record,
				RolledBack:  false,
				FailedItems: nil,
			}
			return nil
		}

		type checkedStock struct {
			// item 是请求输入项，stock 是数据库中同一 SKU 的当前库存快照。
			item  *v1.ReserveItem
			stock *entity.InventoryStock
		}
		// checked 保存通过校验、可进入预占扣减的项目。
		checked := make([]checkedStock, 0, len(req.GetItems()))
		// failed 保存失败明细，直接透传给前端用于逐项提示。
		failed := make([]*v1.ReserveFailedItem, 0)
		for _, item := range req.GetItems() {
			// 先做最轻量参数校验，尽早失败并生成精确错误码。
			if strings.TrimSpace(item.GetSkuNo()) == "" || item.GetQty() == 0 {
				failed = append(failed, &v1.ReserveFailedItem{
					SkuNo:        item.GetSkuNo(),
					SpuNo:        item.GetSpuNo(),
					ShopNo:       item.GetShopNo(),
					RequestedQty: item.GetQty(),
					ErrorCode:    "INVALID_ITEM",
					ErrorMessage: "sku_no and qty are required",
				})
				continue
			}
			// 读取库存行后做上下文核对（spu/shop）+ 数量校验，防止跨店/跨商品误预占。
			stock, err := s.getStockEntityTx(ctx, tx, item.GetSkuNo())
			if err != nil {
				return err
			}
			if stock == nil {
				failed = append(failed, &v1.ReserveFailedItem{
					SkuNo:        item.GetSkuNo(),
					SpuNo:        item.GetSpuNo(),
					ShopNo:       item.GetShopNo(),
					RequestedQty: item.GetQty(),
					ErrorCode:    "SKU_NOT_FOUND",
					ErrorMessage: "sku stock not found",
				})
				continue
			}
			if item.GetSpuNo() != "" && item.GetSpuNo() != stock.SpuNo {
				failed = append(failed, &v1.ReserveFailedItem{
					SkuNo:        item.GetSkuNo(),
					SpuNo:        item.GetSpuNo(),
					ShopNo:       item.GetShopNo(),
					RequestedQty: item.GetQty(),
					AvailableQty: stock.AvailableQty,
					ErrorCode:    "SPU_MISMATCH",
					ErrorMessage: "spu_no mismatch",
				})
				continue
			}
			if item.GetShopNo() != "" && item.GetShopNo() != stock.ShopNo {
				failed = append(failed, &v1.ReserveFailedItem{
					SkuNo:        item.GetSkuNo(),
					SpuNo:        item.GetSpuNo(),
					ShopNo:       item.GetShopNo(),
					RequestedQty: item.GetQty(),
					AvailableQty: stock.AvailableQty,
					ErrorCode:    "SHOP_MISMATCH",
					ErrorMessage: "shop_no mismatch",
				})
				continue
			}
			if stock.AvailableQty < uint64(item.GetQty()) {
				failed = append(failed, &v1.ReserveFailedItem{
					SkuNo:        item.GetSkuNo(),
					SpuNo:        stock.SpuNo,
					ShopNo:       stock.ShopNo,
					RequestedQty: item.GetQty(),
					AvailableQty: stock.AvailableQty,
					ErrorCode:    "INSUFFICIENT_STOCK",
					ErrorMessage: "available qty not enough",
				})
				continue
			}
			checked = append(checked, checkedStock{item: item, stock: stock})
		}

		// 全有或全无模式：只要有任意失败项，就不落库 reservation，直接返回失败明细。
		if len(failed) > 0 && reservationMode == v1.ReserveMode_RESERVE_MODE_ALL_OR_NOTHING {
			response = &v1.ReserveStockRes{
				Success:     false,
				Reservation: nil,
				FailedItems: failed,
				RolledBack:  true,
			}
			return nil
		}

		// 预占过期时间兜底：调用方未传时默认 15 分钟，配合延迟任务释放“幽灵锁”。
		expiredAt := protoTsToGTime(req.GetExpiredAt())
		if expiredAt == nil {
			expiredAt = gtime.NewFromTime(time.Now().Add(15 * time.Minute))
		}
		// 先写预占主记录，再写每个 SKU 的预占明细，保证链路可追踪。
		_, err = tx.Model(dao.InventoryReservation.Table()).Data(do.InventoryReservation{
			ReservationNo:     req.GetReservationNo(),
			OrderNo:           req.GetOrderNo(),
			UserId:            req.GetUserId(),
			ReservationStatus: uint(v1.ReservationStatus_RESERVATION_STATUS_RESERVED),
			ReserveMode:       uint(reservationMode),
			ExpiredAt:         expiredAt,
		}).Insert()
		if err != nil {
			return gerror.Wrap(err, "insert reservation failed")
		}

		for _, row := range checked {
			stock := row.stock
			// qty 是本次预占量；newLocked/newAvailable/newVersion 是写回库存的新值。
			qty := uint64(row.item.GetQty())
			newLocked := stock.LockedQty + qty
			newAvailable := stock.AvailableQty - qty
			newVersion := stock.StockVersion + 1
			status := uint(v1.StockStatus_STOCK_STATUS_OUT_OF_STOCK)
			if newAvailable > 0 {
				status = uint(v1.StockStatus_STOCK_STATUS_IN_STOCK)
			}

			_, err = tx.Model(dao.InventoryStock.Table()).
				Where(dao.InventoryStock.Columns().SkuNo, stock.SkuNo).
				Data(do.InventoryStock{
					LockedQty:    newLocked,
					AvailableQty: newAvailable,
					StockVersion: newVersion,
					StockStatus:  status,
					RowVersion:   stock.RowVersion + 1,
				}).
				Update()
			if err != nil {
				return gerror.Wrapf(err, "reserve stock for %s failed", stock.SkuNo)
			}
			// 明细记录固化当时的 stock_version，便于后续审计与问题追溯。
			_, err = tx.Model(dao.InventoryReservationItem.Table()).Data(do.InventoryReservationItem{
				ReservationNo: req.GetReservationNo(),
				OrderNo:       req.GetOrderNo(),
				SkuNo:         stock.SkuNo,
				SpuNo:         stock.SpuNo,
				ShopNo:        stock.ShopNo,
				Qty:           uint(row.item.GetQty()),
				StockVersion:  newVersion,
			}).Insert()
			if err != nil {
				return gerror.Wrapf(err, "insert reservation item for %s failed", stock.SkuNo)
			}
		}

		record, err := s.loadReservationRecordTx(ctx, tx, req.GetReservationNo())
		if err != nil {
			return err
		}
		// 即使是“部分成功模式”，也会把 failed_items 一并返回给调用方处理 UI 提示。
		response = &v1.ReserveStockRes{
			Success:     true,
			Reservation: record,
			FailedItems: failed,
			RolledBack:  false,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if response != nil && response.GetSuccess() {
		s.logCatalogProjectionSyncFailure(ctx, "reserve stock", s.syncCatalogProjectionByReservation(ctx, response.GetReservation()))
	}
	return response, nil
}

// ConfirmReservation 在支付成功后确认预占，执行真实库存扣减（total/locked 联动更新）。
func (s *sInventory) ConfirmReservation(ctx context.Context, req *v1.ConfirmReservationReq) (*v1.ConfirmReservationRes, error) {
	if strings.TrimSpace(req.GetReservationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "reservation_no is required")
	}
	var response *v1.ConfirmReservationRes
	err := dao.InventoryStock.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 读取并校验预占单状态，确保只有 RESERVED 才能进入确认扣减流程。
		reservation, err := s.getReservationEntityTx(ctx, tx, req.GetReservationNo())
		if err != nil {
			return err
		}
		if reservation == nil {
			return gerror.NewCode(gcode.CodeNotFound, "reservation not found")
		}
		if req.GetOrderNo() != "" && reservation.OrderNo != req.GetOrderNo() {
			return gerror.NewCode(gcode.CodeInvalidParameter, "order_no mismatch")
		}
		if v1.ReservationStatus(reservation.ReservationStatus) == v1.ReservationStatus_RESERVATION_STATUS_CONFIRMED {
			// 幂等语义：重复确认直接返回已确认记录，不重复扣减。
			record, err := s.loadReservationRecordTx(ctx, tx, req.GetReservationNo())
			if err != nil {
				return err
			}
			response = &v1.ConfirmReservationRes{Reservation: record}
			return nil
		}
		if v1.ReservationStatus(reservation.ReservationStatus) != v1.ReservationStatus_RESERVATION_STATUS_RESERVED {
			return gerror.NewCode(gcode.CodeInvalidParameter, "reservation status is not reservable")
		}

		items, err := s.listReservationItemsTx(ctx, tx, req.GetReservationNo())
		if err != nil {
			return err
		}
		for _, item := range items {
			// 逐条把 locked 转为真实扣减：total 减、locked 减、available 保持不变。
			stock, err := s.getStockEntityTx(ctx, tx, item.SkuNo)
			if err != nil {
				return err
			}
			if stock == nil {
				return gerror.NewCodef(gcode.CodeNotFound, "stock not found for sku %s", item.SkuNo)
			}
			if stock.LockedQty < uint64(item.Qty) || stock.TotalQty < uint64(item.Qty) {
				return gerror.NewCodef(gcode.CodeInvalidParameter, "locked qty not enough for sku %s", item.SkuNo)
			}
			newLocked := stock.LockedQty - uint64(item.Qty)
			newTotal := stock.TotalQty - uint64(item.Qty)
			// 预占时 available 已扣过，这里确认阶段不再二次扣减 available。
			newAvailable := stock.AvailableQty
			newVersion := stock.StockVersion + 1
			status := uint(v1.StockStatus_STOCK_STATUS_OUT_OF_STOCK)
			if newAvailable > 0 {
				status = uint(v1.StockStatus_STOCK_STATUS_IN_STOCK)
			}
			_, err = tx.Model(dao.InventoryStock.Table()).
				Where(dao.InventoryStock.Columns().SkuNo, stock.SkuNo).
				Data(do.InventoryStock{
					TotalQty:     newTotal,
					LockedQty:    newLocked,
					AvailableQty: newAvailable,
					StockVersion: newVersion,
					StockStatus:  status,
					RowVersion:   stock.RowVersion + 1,
				}).Update()
			if err != nil {
				return gerror.Wrapf(err, "confirm reservation stock update failed for %s", stock.SkuNo)
			}
		}
		_, err = tx.Model(dao.InventoryReservation.Table()).
			Where(dao.InventoryReservation.Columns().ReservationNo, req.GetReservationNo()).
			Data(do.InventoryReservation{
				ReservationStatus: uint(v1.ReservationStatus_RESERVATION_STATUS_CONFIRMED),
				ConfirmedAt:       gtime.Now(),
			}).
			Update()
		if err != nil {
			return gerror.Wrap(err, "update reservation to confirmed failed")
		}
		// 回读统一结构，保证返回给调用方的是数据库最终确认态。
		record, err := s.loadReservationRecordTx(ctx, tx, req.GetReservationNo())
		if err != nil {
			return err
		}
		response = &v1.ConfirmReservationRes{Reservation: record}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if response != nil {
		s.logCatalogProjectionSyncFailure(ctx, "confirm reservation", s.syncCatalogProjectionByReservation(ctx, response.GetReservation()))
	}
	return response, nil
}

// CancelReservation 取消或超时释放预占库存，回补 available 并降低 locked。
func (s *sInventory) CancelReservation(ctx context.Context, req *v1.CancelReservationReq) (*v1.CancelReservationRes, error) {
	if strings.TrimSpace(req.GetReservationNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "reservation_no is required")
	}
	var response *v1.CancelReservationRes
	err := dao.InventoryStock.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 取消只允许作用于 RESERVED；已取消/已过期请求按幂等直接回放结果。
		reservation, err := s.getReservationEntityTx(ctx, tx, req.GetReservationNo())
		if err != nil {
			return err
		}
		if reservation == nil {
			return gerror.NewCode(gcode.CodeNotFound, "reservation not found")
		}
		if req.GetOrderNo() != "" && reservation.OrderNo != req.GetOrderNo() {
			return gerror.NewCode(gcode.CodeInvalidParameter, "order_no mismatch")
		}
		status := v1.ReservationStatus(reservation.ReservationStatus)
		if status == v1.ReservationStatus_RESERVATION_STATUS_CANCELED || status == v1.ReservationStatus_RESERVATION_STATUS_EXPIRED {
			record, err := s.loadReservationRecordTx(ctx, tx, req.GetReservationNo())
			if err != nil {
				return err
			}
			response = &v1.CancelReservationRes{Reservation: record}
			return nil
		}
		if status != v1.ReservationStatus_RESERVATION_STATUS_RESERVED {
			return gerror.NewCode(gcode.CodeInvalidParameter, "reservation status cannot be canceled")
		}

		items, err := s.listReservationItemsTx(ctx, tx, req.GetReservationNo())
		if err != nil {
			return err
		}
		for _, item := range items {
			// 释放逻辑与预占相反：locked 回退，available 回补。
			stock, err := s.getStockEntityTx(ctx, tx, item.SkuNo)
			if err != nil {
				return err
			}
			if stock == nil {
				return gerror.NewCodef(gcode.CodeNotFound, "stock not found for sku %s", item.SkuNo)
			}
			if stock.LockedQty < uint64(item.Qty) {
				return gerror.NewCodef(gcode.CodeInvalidParameter, "locked qty not enough for sku %s", item.SkuNo)
			}
			newLocked := stock.LockedQty - uint64(item.Qty)
			newAvailable := stock.AvailableQty + uint64(item.Qty)
			newVersion := stock.StockVersion + 1
			status := uint(v1.StockStatus_STOCK_STATUS_OUT_OF_STOCK)
			if newAvailable > 0 {
				status = uint(v1.StockStatus_STOCK_STATUS_IN_STOCK)
			}
			_, err = tx.Model(dao.InventoryStock.Table()).
				Where(dao.InventoryStock.Columns().SkuNo, stock.SkuNo).
				Data(do.InventoryStock{
					LockedQty:    newLocked,
					AvailableQty: newAvailable,
					StockVersion: newVersion,
					StockStatus:  status,
					RowVersion:   stock.RowVersion + 1,
				}).Update()
			if err != nil {
				return gerror.Wrapf(err, "cancel reservation stock update failed for %s", stock.SkuNo)
			}
		}
		finalStatus := v1.ReservationStatus_RESERVATION_STATUS_CANCELED
		// reason_code=EXPIRED 时落 EXPIRED，便于区分用户主动取消与系统超时释放。
		if strings.EqualFold(req.GetReasonCode(), "EXPIRED") {
			finalStatus = v1.ReservationStatus_RESERVATION_STATUS_EXPIRED
		}
		_, err = tx.Model(dao.InventoryReservation.Table()).
			Where(dao.InventoryReservation.Columns().ReservationNo, req.GetReservationNo()).
			Data(do.InventoryReservation{
				ReservationStatus: uint(finalStatus),
				CanceledAt:        gtime.Now(),
				CancelReasonCode:  req.GetReasonCode(),
			}).Update()
		if err != nil {
			return gerror.Wrap(err, "update reservation canceled failed")
		}
		record, err := s.loadReservationRecordTx(ctx, tx, req.GetReservationNo())
		if err != nil {
			return err
		}
		response = &v1.CancelReservationRes{Reservation: record}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if response != nil {
		s.logCatalogProjectionSyncFailure(ctx, "cancel reservation", s.syncCatalogProjectionByReservation(ctx, response.GetReservation()))
	}
	return response, nil
}

// BatchAdjustStockByAdmin 管理员批量调库存，可同时修正缺失的 SKU 上下文。
func (s *sInventory) BatchAdjustStockByAdmin(ctx context.Context, req *v1.BatchAdjustStockByAdminReq) (*v1.BatchAdjustStockByAdminRes, error) {
	if len(req.GetItems()) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "items is required")
	}
	results := make([]*v1.AdjustResultItem, 0, len(req.GetItems()))
	err := dao.InventoryStock.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, item := range req.GetItems() {
			// 管理员入口要求显式携带 sku/spu/shop 三元组，防止上下文不完整。
			if strings.TrimSpace(item.GetSkuNo()) == "" || strings.TrimSpace(item.GetSpuNo()) == "" || strings.TrimSpace(item.GetShopNo()) == "" {
				return gerror.NewCode(gcode.CodeInvalidParameter, "sku_no, spu_no and shop_no are required")
			}
			// 管理员可顺手补齐上下文，避免后续库存事件丢失 spu/shop 信息。
			if err := s.upsertSkuContextItemTx(ctx, tx, &v1.SkuContextItem{
				SkuNo:   item.GetSkuNo(),
				SpuNo:   item.GetSpuNo(),
				ShopNo:  item.GetShopNo(),
				Enabled: true,
			}); err != nil {
				return err
			}
			adjusted, err := s.adjustStockTx(ctx, tx, item.GetSkuNo(), item.GetSpuNo(), item.GetShopNo(), item.GetDeltaTotalQty())
			if err != nil {
				return err
			}
			results = append(results, adjusted)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logCatalogProjectionSyncFailure(ctx, "admin stock adjust", s.syncCatalogProjectionByAdjustResults(ctx, results))
	return &v1.BatchAdjustStockByAdminRes{Results: results}, nil
}

// SetHotSku 设置或取消热点 SKU 标记，便于后续热点治理与监控。
func (s *sInventory) SetHotSku(ctx context.Context, req *v1.SetHotSkuReq) (*v1.SetHotSkuRes, error) {
	if strings.TrimSpace(req.GetSkuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sku_no is required")
	}
	_, err := dao.InventoryStock.Ctx(ctx).
		Where(dao.InventoryStock.Columns().SkuNo, req.GetSkuNo()).
		Data(do.InventoryStock{
			IsHot: boolToInt(req.GetIsHot()),
		}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "set hot sku failed")
	}
	return &v1.SetHotSkuRes{
		SkuNo: req.GetSkuNo(),
		IsHot: req.GetIsHot(),
	}, nil
}

// GetSkuInventory 查询单个 SKU 的库存快照。
func (s *sInventory) GetSkuInventory(ctx context.Context, req *v1.GetSkuInventoryReq) (*v1.GetSkuInventoryRes, error) {
	if strings.TrimSpace(req.GetSkuNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sku_no is required")
	}
	row, err := s.getStockEntity(ctx, req.GetSkuNo())
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "sku not found")
	}
	return &v1.GetSkuInventoryRes{Stock: toProtoStock(row)}, nil
}

// BatchGetSkuInventory 批量查询 SKU 库存快照。
func (s *sInventory) BatchGetSkuInventory(ctx context.Context, req *v1.BatchGetSkuInventoryReq) (*v1.BatchGetSkuInventoryRes, error) {
	if len(req.GetSkuNos()) == 0 {
		return &v1.BatchGetSkuInventoryRes{Stocks: []*v1.InventoryStock{}}, nil
	}
	var rows []*entity.InventoryStock
	err := dao.InventoryStock.Ctx(ctx).
		WhereIn(dao.InventoryStock.Columns().SkuNo, req.GetSkuNos()).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "batch get sku inventory failed")
	}
	out := make([]*v1.InventoryStock, 0, len(rows))
	for _, row := range rows {
		out = append(out, toProtoStock(row))
	}
	return &v1.BatchGetSkuInventoryRes{Stocks: out}, nil
}

// UpsertSkuContext 同步 SKU 与 SPU/店铺归属关系，为库存鉴权和事件补全提供上下文。
func (s *sInventory) UpsertSkuContext(ctx context.Context, req *v1.UpsertSkuContextReq) (*v1.UpsertSkuContextRes, error) {
	var updated uint64
	err := dao.InventorySkuContext.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, item := range req.GetItems() {
			if err := s.upsertSkuContextItemTx(ctx, tx, item); err != nil {
				return err
			}
			updated++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpsertSkuContextRes{UpdatedRows: updated}, nil
}

// adjustStockTx 在事务内执行库存增减，统一维护 total/available/status/version。
func (s *sInventory) adjustStockTx(ctx context.Context, tx gdb.TX, skuNo, spuNo, shopNo string, delta int64) (*v1.AdjustResultItem, error) {
	stock, err := s.getStockEntityTx(ctx, tx, skuNo)
	if err != nil {
		return nil, err
	}
	if stock == nil {
		// 首次建库存仅允许正向加库存，防止“无记录直接扣减”造成脏数据。
		if delta < 0 {
			return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "stock for sku %s not found", skuNo)
		}
		total := uint64(delta)
		status := uint(v1.StockStatus_STOCK_STATUS_OUT_OF_STOCK)
		if total > 0 {
			status = uint(v1.StockStatus_STOCK_STATUS_IN_STOCK)
		}
		_, err = tx.Model(dao.InventoryStock.Table()).Data(do.InventoryStock{
			SkuNo:        skuNo,
			SpuNo:        spuNo,
			ShopNo:       shopNo,
			TotalQty:     total,
			LockedQty:    uint64(0),
			AvailableQty: total,
			StockVersion: uint64(1),
			StockStatus:  status,
			IsHot:        0,
			RowVersion:   uint64(1),
		}).Insert()
		if err != nil {
			return nil, gerror.Wrap(err, "insert stock failed")
		}
		// 首次插入后直接返回当前快照，供调用方展示最新库存结果。
		return &v1.AdjustResultItem{
			SkuNo:        skuNo,
			SpuNo:        spuNo,
			ShopNo:       shopNo,
			TotalQty:     total,
			LockedQty:    0,
			AvailableQty: total,
			StockVersion: 1,
			StockStatus:  v1.StockStatus(status),
		}, nil
	}
	if spuNo != "" && stock.SpuNo != "" && stock.SpuNo != spuNo {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "sku %s spu mismatch", skuNo)
	}
	if shopNo != "" && stock.ShopNo != "" && stock.ShopNo != shopNo {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "sku %s shop mismatch", skuNo)
	}

	// newTotalInt 先用有符号整型计算，便于处理 delta 可能为负数的场景。
	newTotalInt := int64(stock.TotalQty) + delta
	if newTotalInt < 0 {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "total qty cannot be negative for sku %s", skuNo)
	}
	// total 不能小于 locked，否则说明还有订单占用库存，不能被管理员调没。
	if uint64(newTotalInt) < stock.LockedQty {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "total qty cannot be less than locked qty for sku %s", skuNo)
	}
	newTotal := uint64(newTotalInt)
	newAvailable := newTotal - stock.LockedQty
	newVersion := stock.StockVersion + 1
	status := uint(v1.StockStatus_STOCK_STATUS_OUT_OF_STOCK)
	if newAvailable > 0 {
		status = uint(v1.StockStatus_STOCK_STATUS_IN_STOCK)
	}
	_, err = tx.Model(dao.InventoryStock.Table()).
		Where(dao.InventoryStock.Columns().SkuNo, skuNo).
		Data(do.InventoryStock{
			TotalQty:     newTotal,
			AvailableQty: newAvailable,
			// stock_version 用于事件有序投影；row_version 用于数据库行更新轨迹。
			StockVersion: newVersion,
			StockStatus:  status,
			RowVersion:   stock.RowVersion + 1,
		}).Update()
	if err != nil {
		return nil, gerror.Wrap(err, "update stock failed")
	}
	return &v1.AdjustResultItem{
		SkuNo:        skuNo,
		SpuNo:        stock.SpuNo,
		ShopNo:       stock.ShopNo,
		TotalQty:     newTotal,
		LockedQty:    stock.LockedQty,
		AvailableQty: newAvailable,
		StockVersion: newVersion,
		StockStatus:  v1.StockStatus(status),
	}, nil
}

// upsertSkuContextItemTx 在事务内写入或更新单条 SKU 上下文记录。
func (s *sInventory) upsertSkuContextItemTx(ctx context.Context, tx gdb.TX, item *v1.SkuContextItem) error {
	if item == nil || strings.TrimSpace(item.GetSkuNo()) == "" {
		return gerror.NewCode(gcode.CodeInvalidParameter, "sku context item is invalid")
	}
	exists, err := s.getSkuContextTx(ctx, tx, item.GetSkuNo())
	if err != nil {
		return err
	}
	if exists == nil {
		// 上下文不存在时直接插入，建立 sku -> spu/shop 的稳定映射关系。
		_, err = tx.Model(dao.InventorySkuContext.Table()).Data(do.InventorySkuContext{
			SkuNo:   item.GetSkuNo(),
			SpuNo:   item.GetSpuNo(),
			ShopNo:  item.GetShopNo(),
			Enabled: boolToInt(item.GetEnabled()),
		}).Insert()
		if err != nil {
			return gerror.Wrap(err, "insert sku context failed")
		}
		return nil
	}
	// 已存在则覆盖更新，支持商家迁移或后台修正上下文。
	_, err = tx.Model(dao.InventorySkuContext.Table()).
		Where(dao.InventorySkuContext.Columns().SkuNo, item.GetSkuNo()).
		Data(do.InventorySkuContext{
			SpuNo:   item.GetSpuNo(),
			ShopNo:  item.GetShopNo(),
			Enabled: boolToInt(item.GetEnabled()),
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "update sku context failed")
	}
	return nil
}

// getStockEntity 查询单条库存实体（非事务版本）。
func (s *sInventory) getStockEntity(ctx context.Context, skuNo string) (*entity.InventoryStock, error) {
	var row entity.InventoryStock
	err := dao.InventoryStock.Ctx(ctx).
		Where(dao.InventoryStock.Columns().SkuNo, skuNo).
		Scan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, gerror.Wrap(err, "query stock failed")
	}
	if row.Id == 0 {
		return nil, nil
	}
	return &row, nil
}

// getStockEntityTx 在事务内查询单条库存实体。
func (s *sInventory) getStockEntityTx(ctx context.Context, tx gdb.TX, skuNo string) (*entity.InventoryStock, error) {
	var row entity.InventoryStock
	err := tx.Model(dao.InventoryStock.Table()).
		Where(dao.InventoryStock.Columns().SkuNo, skuNo).
		Scan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, gerror.Wrap(err, "query stock failed")
	}
	if row.Id == 0 {
		return nil, nil
	}
	return &row, nil
}

// getSkuContextTx 在事务内查询 SKU 上下文。
func (s *sInventory) getSkuContextTx(ctx context.Context, tx gdb.TX, skuNo string) (*entity.InventorySkuContext, error) {
	var row entity.InventorySkuContext
	err := tx.Model(dao.InventorySkuContext.Table()).
		Where(dao.InventorySkuContext.Columns().SkuNo, skuNo).
		Scan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, gerror.Wrap(err, "query sku context failed")
	}
	if row.Id == 0 {
		return nil, nil
	}
	return &row, nil
}

// getReservationEntityTx 在事务内查询预占主记录。
func (s *sInventory) getReservationEntityTx(ctx context.Context, tx gdb.TX, reservationNo string) (*entity.InventoryReservation, error) {
	var row entity.InventoryReservation
	err := tx.Model(dao.InventoryReservation.Table()).
		Where(dao.InventoryReservation.Columns().ReservationNo, reservationNo).
		Scan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, gerror.Wrap(err, "query reservation failed")
	}
	if row.Id == 0 {
		return nil, nil
	}
	return &row, nil
}

// listReservationItemsTx 在事务内查询预占明细列表。
func (s *sInventory) listReservationItemsTx(ctx context.Context, tx gdb.TX, reservationNo string) ([]*entity.InventoryReservationItem, error) {
	var rows []*entity.InventoryReservationItem
	err := tx.Model(dao.InventoryReservationItem.Table()).
		Where(dao.InventoryReservationItem.Columns().ReservationNo, reservationNo).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "query reservation items failed")
	}
	return rows, nil
}

// loadReservationRecordTx 在事务内组装完整的 ReservationRecord 返回值。
func (s *sInventory) loadReservationRecordTx(ctx context.Context, tx gdb.TX, reservationNo string) (*v1.ReservationRecord, error) {
	reservation, err := s.getReservationEntityTx(ctx, tx, reservationNo)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "reservation not found")
	}
	items, err := s.listReservationItemsTx(ctx, tx, reservationNo)
	if err != nil {
		return nil, err
	}
	reservedItems := make([]*v1.ReservedItem, 0, len(items))
	for _, row := range items {
		reservedItems = append(reservedItems, &v1.ReservedItem{
			SkuNo:        row.SkuNo,
			SpuNo:        row.SpuNo,
			ShopNo:       row.ShopNo,
			Qty:          uint32(row.Qty),
			StockVersion: row.StockVersion,
		})
	}
	return &v1.ReservationRecord{
		ReservationNo: reservation.ReservationNo,
		OrderNo:       reservation.OrderNo,
		UserId:        reservation.UserId,
		Status:        v1.ReservationStatus(reservation.ReservationStatus),
		ExpiredAt:     toProtoTs(reservation.ExpiredAt),
		Items:         reservedItems,
		CreatedAt:     toProtoTs(reservation.CreatedAt),
		UpdatedAt:     toProtoTs(reservation.UpdatedAt),
	}, nil
}

// toProtoStock 将库存实体转换为 Proto 返回结构。
func toProtoStock(row *entity.InventoryStock) *v1.InventoryStock {
	return &v1.InventoryStock{
		SkuNo:        row.SkuNo,
		SpuNo:        row.SpuNo,
		ShopNo:       row.ShopNo,
		TotalQty:     row.TotalQty,
		LockedQty:    row.LockedQty,
		AvailableQty: row.AvailableQty,
		StockVersion: row.StockVersion,
		StockStatus:  v1.StockStatus(row.StockStatus),
		IsHot:        row.IsHot == 1,
		UpdatedAt:    toProtoTs(row.UpdatedAt),
	}
}

// boolToInt 将布尔值转换为数据库常用的 1/0 标记。
func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

// toProtoTs 将 gtime 时间转换为 protobuf Timestamp。
func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(t.Time)
}

// protoTsToGTime 将 protobuf Timestamp 转换为 gtime 时间。
func protoTsToGTime(ts *timestamppb.Timestamp) *gtime.Time {
	if ts == nil {
		return nil
	}
	return gtime.NewFromTime(ts.AsTime())
}

// buildMismatchError 构造字段不匹配错误，便于定位调用方参数问题。
func buildMismatchError(field, expected, actual string) error {
	return gerror.NewCodef(gcode.CodeInvalidParameter, "%s mismatch, expected=%s actual=%s", field, expected, actual)
}

// formatSkuError 构造带 sku 维度上下文的业务错误信息。
func formatSkuError(skuNo string, msg string) error {
	return gerror.NewCode(gcode.CodeInvalidParameter, fmt.Sprintf("sku=%s %s", skuNo, msg))
}
