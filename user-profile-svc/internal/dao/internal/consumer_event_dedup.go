// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ConsumerEventDedupDao 是 consumer_event_dedup 表的数据访问对象。
type ConsumerEventDedupDao struct {
	table    string                    // table 是 DAO 所在的底层表名。
	group    string                    // group 是当前 DAO 使用的数据库配置组名。
	columns  ConsumerEventDedupColumns // columns 缓存了表中所有列名，便于复用。
	handlers []gdb.ModelHandler        // handlers 用于对模型的自定义修改。
}

// ConsumerEventDedupColumns 定义 consumer_event_dedup 表的列名。
type ConsumerEventDedupColumns struct {
	Id           string // 主键
	ConsumerName string // 消费者唯一标识
	EventId      string // 幂等事件 ID
	EventType    string // 事件类型
	ProcessedAt  string // 处理时间
	CreatedAt    string // 创建时间
}

// consumerEventDedupColumns 存储该表的列名。
var consumerEventDedupColumns = ConsumerEventDedupColumns{
	Id:           "id",
	ConsumerName: "consumer_name",
	EventId:      "event_id",
	EventType:    "event_type",
	ProcessedAt:  "processed_at",
	CreatedAt:    "created_at",
}

// NewConsumerEventDedupDao 创建并返回 consumer_event_dedup 表的 DAO 实例。
func NewConsumerEventDedupDao(handlers ...gdb.ModelHandler) *ConsumerEventDedupDao {
	return &ConsumerEventDedupDao{
		group:    "default",
		table:    "consumer_event_dedup",
		columns:  consumerEventDedupColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 使用的底层数据库对象。
func (dao *ConsumerEventDedupDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 所使用的表名。
func (dao *ConsumerEventDedupDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *ConsumerEventDedupDao) Columns() ConsumerEventDedupColumns {
	return dao.columns
}

// Group 返回数据库配置组名。
func (dao *ConsumerEventDedupDao) Group() string {
	return dao.group
}

// Ctx 为当前 DAO 创建并返回一个上下文已设置的 ORM 模型。
func (dao *ConsumerEventDedupDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction 用于包裹事务逻辑，执行传入的 f。
// 若 f 返回非空错误，事务会回滚并原样返回该错误。
// 若 f 返回 nil，则事务提交并返回 nil。
//
// 注意：f 内无需手动提交或回滚，函数会自动处理事务。
func (dao *ConsumerEventDedupDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
