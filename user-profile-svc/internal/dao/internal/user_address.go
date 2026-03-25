// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserAddressDao 是 user_address 表的数据访问对象。
type UserAddressDao struct {
	table    string             // table 是 DAO 所在的底层表名。
	group    string             // group 是当前 DAO 使用的数据库配置组名。
	columns  UserAddressColumns // columns 缓存了表中所有列名，便于复用。
	handlers []gdb.ModelHandler // handlers 用于对模型的自定义修改。
}

// UserAddressColumns 定义 user_address 表的列名。
type UserAddressColumns struct {
	AddressId             string // 地址 ID
	UserId                string // 所属用户 ID
	Status                string // 1=活跃,2=已删除,3=已替换
	AddressVersion        string // 地址行的乐观锁版本
	ReplacedFromAddressId string // 替换来源地址链
	Label                 string // 地址标签
	ReceiverName          string // 收件人姓名
	ReceiverPhone         string // 收件人电话
	CountryCode           string // ISO 国家码
	ProvinceCode          string // 省份编码
	ProvinceName          string // 省份名称
	CityCode              string // 城市编码
	CityName              string // 城市名称
	DistrictCode          string // 区县编码
	DistrictName          string // 区县名称
	Street                string // 街道/乡镇
	Detail                string // 详细地址
	PostalCode            string // 邮政编码
	IsDefault             string // 1=默认,NULL=非默认
	DefaultSlot           string // 默认唯一槽位
	Latitude              string // 纬度
	Longitude             string // 经度
	Ext                   string // 扩展字段
	CreatedAt             string //
	UpdatedAt             string //
}

// userAddressColumns 保存该表的列名。
var userAddressColumns = UserAddressColumns{
	AddressId:             "address_id",
	UserId:                "user_id",
	Status:                "status",
	AddressVersion:        "address_version",
	ReplacedFromAddressId: "replaced_from_address_id",
	Label:                 "label",
	ReceiverName:          "receiver_name",
	ReceiverPhone:         "receiver_phone",
	CountryCode:           "country_code",
	ProvinceCode:          "province_code",
	ProvinceName:          "province_name",
	CityCode:              "city_code",
	CityName:              "city_name",
	DistrictCode:          "district_code",
	DistrictName:          "district_name",
	Street:                "street",
	Detail:                "detail",
	PostalCode:            "postal_code",
	IsDefault:             "is_default",
	DefaultSlot:           "default_slot",
	Latitude:              "latitude",
	Longitude:             "longitude",
	Ext:                   "ext",
	CreatedAt:             "created_at",
	UpdatedAt:             "updated_at",
}

// NewUserAddressDao 创建并返回 user_address 表的 DAO 实例。
func NewUserAddressDao(handlers ...gdb.ModelHandler) *UserAddressDao {
	return &UserAddressDao{
		group:    "default",
		table:    "user_address",
		columns:  userAddressColumns,
		handlers: handlers,
	}
}

// DB 返回当前 DAO 使用的底层数据库对象。
func (dao *UserAddressDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table 返回当前 DAO 所使用的表名。
func (dao *UserAddressDao) Table() string {
	return dao.table
}

// Columns 返回当前 DAO 的所有列名。
func (dao *UserAddressDao) Columns() UserAddressColumns {
	return dao.columns
}

// Group 返回数据库配置组名。
func (dao *UserAddressDao) Group() string {
	return dao.group
}

// Ctx 为当前 DAO 创建并返回一个上下文已设置的 ORM 模型。
func (dao *UserAddressDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserAddressDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
