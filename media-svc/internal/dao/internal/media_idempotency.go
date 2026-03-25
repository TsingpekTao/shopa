// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MediaIdempotencyDao 是 media_idempotency 表的数据访问对象。
type MediaIdempotencyDao struct {
	table    string                  // table 表示 DAO 所在的表名。
	group    string                  // group 表示数据库配置分组。
	columns  MediaIdempotencyColumns // columns 包含全部列名，方便复用。
	handlers []gdb.ModelHandler      // handlers 用于自定义模型处理。
}

// MediaIdempotencyColumns 定义并存储 media_idempotency 表的列名。
type MediaIdempotencyColumns struct {
	Id             string //
	ActorUserId    string // metadata 中的用户 ID。
	ActionCode     string // 示例：InitUpload/BatchBindAssetsToBiz/...
	IdempotencyKey string // x-idempotency-key 请求头。
	RequestHash    string // 请求有效载荷哈希。
	SceneCode      string //
	BizType        string //
	BizNo          string //
	Status         string //
	ResponseCode   string //
	ResponseJson   string //
	ExpiredAt      string //
	CreatedAt      string //
	UpdatedAt      string //
}

// mediaIdempotencyColumns 保存 media_idempotency 表的列名映射。
var mediaIdempotencyColumns = MediaIdempotencyColumns{
	Id:             "id",
	ActorUserId:    "actor_user_id",
	ActionCode:     "action_code",
	IdempotencyKey: "idempotency_key",
	RequestHash:    "request_hash",
	SceneCode:      "scene_code",
	BizType:        "biz_type",
	BizNo:          "biz_no",
	Status:         "status",
	ResponseCode:   "response_code",
	ResponseJson:   "response_json",
	ExpiredAt:      "expired_at",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewMediaIdempotencyDao 创建并返回 media_idempotency 表的数据访问对象。
func NewMediaIdempotencyDao(handlers ...gdb.ModelHandler) *MediaIdempotencyDao {
	return &MediaIdempotencyDao{
		group:    "default",
		table:    "media_idempotency",
		columns:  mediaIdempotencyColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 的底层数据库管理对象。
func (dao *MediaIdempotencyDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 操作的表名。
func (dao *MediaIdempotencyDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *MediaIdempotencyDao) Columns() MediaIdempotencyColumns {
	return dao.columns
}

// Group 返回 DAO 使用的数据库配置分组名称。
func (dao *MediaIdempotencyDao) Group() string {
	return dao.group
}

// Ctx 创建并返回当前 DAO 的 Model，并自动设置操作上下文。
func (dao *MediaIdempotencyDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MediaIdempotencyDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
