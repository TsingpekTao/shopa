// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaConsumerEventDedupDao 是 media_consumer_event_dedup 表的数据访问对象。
type MediaConsumerEventDedupDao struct {
	table    string                         // table 表示 DAO 所在的表名。
	group    string                         // group 表示数据库配置分组。
	columns  MediaConsumerEventDedupColumns // columns 包含全部列名，方便复用。
	handlers []gdb.ModelHandler             // handlers 用于自定义模型处理。
}

// MediaConsumerEventDedupColumns 定义并存储 media_consumer_event_dedup 表的列名。
type MediaConsumerEventDedupColumns struct {
	Id           string //
	ConsumerName string //
	EventId      string //
	PayloadHash  string //
	FirstSeenAt  string //
}

// mediaConsumerEventDedupColumns 保存 media_consumer_event_dedup 表的列名映射。
var mediaConsumerEventDedupColumns = MediaConsumerEventDedupColumns{
	Id:           "id",
	ConsumerName: "consumer_name",
	EventId:      "event_id",
	PayloadHash:  "payload_hash",
	FirstSeenAt:  "first_seen_at",
}

// NewMediaConsumerEventDedupDao 创建并返回 media_consumer_event_dedup 表的数据访问对象。
func NewMediaConsumerEventDedupDao(handlers ...gdb.ModelHandler) *MediaConsumerEventDedupDao {
	return &MediaConsumerEventDedupDao{
		group:    "default",
		table:    "media_consumer_event_dedup",
		columns:  mediaConsumerEventDedupColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 的底层数据库管理对象。
func (dao *MediaConsumerEventDedupDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 操作的表名。
func (dao *MediaConsumerEventDedupDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *MediaConsumerEventDedupDao) Columns() MediaConsumerEventDedupColumns {
	return dao.columns
}

// Group 返回 DAO 使用的数据库配置分组名称。
func (dao *MediaConsumerEventDedupDao) Group() string {
	return dao.group
}

// Ctx 创建并返回当前 DAO 的 Model，并自动设置操作上下文。
func (dao *MediaConsumerEventDedupDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction 用参数函数 f 包裹事务逻辑。
// 如果 f 返回非 nil 错误则回滚并返回该错误。
// 如果 f 返回 nil 则提交并返回 nil。
//
// 注意：f 内请勿显式提交或回滚，事务由此方法自动管理。
func (dao *MediaConsumerEventDedupDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
