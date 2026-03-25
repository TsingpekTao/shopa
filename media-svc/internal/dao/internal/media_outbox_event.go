// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaOutboxEventDao 是 media_outbox_event 表的数据访问对象。
type MediaOutboxEventDao struct {
	table    string                  // table 表示 DAO 所在的表名。
	group    string                  // group 表示数据库配置分组。
	columns  MediaOutboxEventColumns // columns 包含全部列名，方便复用。
	handlers []gdb.ModelHandler      // handlers 用于自定义模型处理。
}

// MediaOutboxEventColumns 定义并存储 media_outbox_event 表的列名。
type MediaOutboxEventColumns struct {
	Id            string //
	EventId       string // 全局唯一事件 ID。
	EventType     string // 事件类型，例如 AssetProcessingCompleted/AssetProcessingFailed/...
	AggregateType string // 聚合类型：ASSET/BINDING/SCENE。
	AggregateKey  string // 聚合键：asset_id/biz_key/scene_code。
	RequestId     string //
	PayloadJson   string //
	Status        string //
	AvailableAt   string //
	SentAt        string //
	FailCount     string //
	LastError     string //
	CreatedAt     string //
	UpdatedAt     string //
}

// mediaOutboxEventColumns 保存 media_outbox_event 表的列名映射。
var mediaOutboxEventColumns = MediaOutboxEventColumns{
	Id:            "id",
	EventId:       "event_id",
	EventType:     "event_type",
	AggregateType: "aggregate_type",
	AggregateKey:  "aggregate_key",
	RequestId:     "request_id",
	PayloadJson:   "payload_json",
	Status:        "status",
	AvailableAt:   "available_at",
	SentAt:        "sent_at",
	FailCount:     "fail_count",
	LastError:     "last_error",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewMediaOutboxEventDao 创建并返回 media_outbox_event 表的数据访问对象。
func NewMediaOutboxEventDao(handlers ...gdb.ModelHandler) *MediaOutboxEventDao {
	return &MediaOutboxEventDao{
		group:    "default",
		table:    "media_outbox_event",
		columns:  mediaOutboxEventColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 的底层数据库管理对象。
func (dao *MediaOutboxEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 操作的表名。
func (dao *MediaOutboxEventDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *MediaOutboxEventDao) Columns() MediaOutboxEventColumns {
	return dao.columns
}

// Group 返回 DAO 使用的数据库配置分组名称。
func (dao *MediaOutboxEventDao) Group() string {
	return dao.group
}

// Ctx 创建并返回当前 DAO 的 Model，并自动设置操作上下文。
func (dao *MediaOutboxEventDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MediaOutboxEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
