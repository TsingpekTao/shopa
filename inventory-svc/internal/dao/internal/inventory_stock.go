// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// InventoryStockDao 是 inventory_stock 表的数据访问对象。
type InventoryStockDao struct {
	table    string                // table 是 DAO 的底层表名。
	group    string                // group 是当前 DAO 使用的数据库配置组名。
	columns  InventoryStockColumns // columns 持有该表所有字段名，便于直接调用。
	handlers []gdb.ModelHandler    // handlers 用于自定义模型处理器。
}

// InventoryStockColumns 定义并存储 inventory_stock 表的字段名。
type InventoryStockColumns struct {
	Id           string //
	SkuNo        string //
	SpuNo        string //
	ShopNo       string //
	TotalQty     string //
	LockedQty    string //
	AvailableQty string //
	StockVersion string // Monotonic version for projection events
	StockStatus  string //
	IsHot        string //
	RowVersion   string // Internal optimistic lock version
	LastEventId  string //
	CreatedAt    string //
	UpdatedAt    string //
}

// inventoryStockColumns 保持表 inventory_stock 的字段信息。
var inventoryStockColumns = InventoryStockColumns{
	Id:           "id",
	SkuNo:        "sku_no",
	SpuNo:        "spu_no",
	ShopNo:       "shop_no",
	TotalQty:     "total_qty",
	LockedQty:    "locked_qty",
	AvailableQty: "available_qty",
	StockVersion: "stock_version",
	StockStatus:  "stock_status",
	IsHot:        "is_hot",
	RowVersion:   "row_version",
	LastEventId:  "last_event_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewInventoryStockDao 创建并返回该表的 DAO 实例。
func NewInventoryStockDao(handlers ...gdb.ModelHandler) *InventoryStockDao {
	return &InventoryStockDao{
		group:    "default",
		table:    "inventory_stock",
		columns:  inventoryStockColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 使用的原始数据库管理对象。
func (dao *InventoryStockDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 使用的表名。
func (dao *InventoryStockDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 所有字段的名称集合。
func (dao *InventoryStockDao) Columns() InventoryStockColumns {
	return dao.columns
}

// Group 返回当前 DAO 使用的数据库配置组名。
func (dao *InventoryStockDao) Group() string {
	return dao.group
}

// Ctx 为当前 DAO 创建并返回一个 Model，并自动绑定本次操作的上下文。
func (dao *InventoryStockDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *InventoryStockDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
