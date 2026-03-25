// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserProfileDao 是 user_profile 表的数据访问对象。
type UserProfileDao struct {
	table    string             // table 是 DAO 所在的底层表名。
	group    string             // group 是当前 DAO 使用的数据库配置组名。
	columns  UserProfileColumns // columns 缓存了表中所有列名，便于复用。
	handlers []gdb.ModelHandler // handlers 用于对模型的自定义修改。
}

// UserProfileColumns 定义 user_profile 表的列名。
type UserProfileColumns struct {
	UserId             string // iam-svc 提供的用户 ID
	DisplayName        string // 展示名
	AvatarAssetId      string // media-svc 的头像资源 ID
	AvatarUrl          string // 头像 URL 缓存
	Gender             string // 0=未指定,1=男,2=女,3=其他
	Birthday           string // 出生日期
	Locale             string // 区域设置
	Timezone           string // 时区
	MarketingOptIn     string // 营销订阅
	DefaultAddressId   string // 兼容字段（已弃用）
	ProfileVersion     string // 画像聚合的乐观锁版本
	AddressBookVersion string // 地址簿的 CAS 版本
	DisplayNameSource  string // 1=SYSTEM_INIT,2=USER_SET
	LastInitEventAt    string // 最近 register-init 事件时间
	LastInitEventId    string // 最近 register-init 事件 ID
	Ext                string // 扩展字段
	CreatedAt          string //
	UpdatedAt          string //
}

// userProfileColumns 保存该表的列名。
var userProfileColumns = UserProfileColumns{
	UserId:             "user_id",
	DisplayName:        "display_name",
	AvatarAssetId:      "avatar_asset_id",
	AvatarUrl:          "avatar_url",
	Gender:             "gender",
	Birthday:           "birthday",
	Locale:             "locale",
	Timezone:           "timezone",
	MarketingOptIn:     "marketing_opt_in",
	DefaultAddressId:   "default_address_id",
	ProfileVersion:     "profile_version",
	AddressBookVersion: "address_book_version",
	DisplayNameSource:  "display_name_source",
	LastInitEventAt:    "last_init_event_at",
	LastInitEventId:    "last_init_event_id",
	Ext:                "ext",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
}

// NewUserProfileDao 创建并返回 user_profile 表的 DAO 实例。
func NewUserProfileDao(handlers ...gdb.ModelHandler) *UserProfileDao {
	return &UserProfileDao{
		group:    "default",
		table:    "user_profile",
		columns:  userProfileColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 使用的底层数据库对象。
func (dao *UserProfileDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 所使用的表名。
func (dao *UserProfileDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *UserProfileDao) Columns() UserProfileColumns {
	return dao.columns
}

// Group 返回数据库配置组名。
func (dao *UserProfileDao) Group() string {
	return dao.group
}

// Ctx 为当前 DAO 创建并返回一个上下文已设置的 ORM 模型。
func (dao *UserProfileDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserProfileDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
