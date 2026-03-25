// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InventoryReservationDao 是 inventory_reservation 表的数据访问对象。
type InventoryReservationDao struct {
	table    string                      // table 是 DAO 的底层表名。
	group    string                      // group 是当前 DAO 使用的数据库配置组名。
	columns  InventoryReservationColumns // columns 持有该表所有字段名，便于直接调用。
	handlers []gdb.ModelHandler          // handlers 用于自定义模型处理器。
}

// InventoryReservationColumns 定义并存储 inventory_reservation 表的字段名。
type InventoryReservationColumns struct {
	Id                string //
	ReservationNo     string //
	OrderNo           string //
	UserId            string //
	ReservationStatus string //
	ReserveMode       string //
	ExpiredAt         string //
	ConfirmedAt       string //
	CanceledAt        string //
	CancelReasonCode  string //
	IdempotencyKey    string //
	RequestId         string //
	CreatedAt         string //
	UpdatedAt         string //
}

// inventoryReservationColumns 保持表 inventory_reservation 的字段信息。
var inventoryReservationColumns = InventoryReservationColumns{
	Id:                "id",
	ReservationNo:     "reservation_no",
	OrderNo:           "order_no",
	UserId:            "user_id",
	ReservationStatus: "reservation_status",
	ReserveMode:       "reserve_mode",
	ExpiredAt:         "expired_at",
	ConfirmedAt:       "confirmed_at",
	CanceledAt:        "canceled_at",
	CancelReasonCode:  "cancel_reason_code",
	IdempotencyKey:    "idempotency_key",
	RequestId:         "request_id",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewInventoryReservationDao 创建并返回该表的 DAO 实例。
func NewInventoryReservationDao(handlers ...gdb.ModelHandler) *InventoryReservationDao {
	return &InventoryReservationDao{
		group:    "default",
		table:    "inventory_reservation",
		columns:  inventoryReservationColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 使用的原始数据库管理对象。
func (dao *InventoryReservationDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 使用的表名。
func (dao *InventoryReservationDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 所有字段的名称集合。
func (dao *InventoryReservationDao) Columns() InventoryReservationColumns {
	return dao.columns
}

// Group 返回当前 DAO 使用的数据库配置组名。
func (dao *InventoryReservationDao) Group() string {
	return dao.group
}

// Ctx 为当前 DAO 创建并返回一个 Model，并自动绑定本次操作的上下文。
func (dao *InventoryReservationDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction 用提供的函数 f 包裹事务逻辑。
// 如果函数 f 返回非 nil 错误，则自动回滚事务并返回该错误。
// 如果函数 f 返回 nil，则自动提交事务并返回 nil。
//
// 注意：不要在函数 f 中显式提交或回滚事务，
// 该函数会自动处理相关操作。
func (dao *InventoryReservationDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
