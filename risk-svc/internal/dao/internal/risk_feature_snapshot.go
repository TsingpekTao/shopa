// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RiskFeatureSnapshotDao is the data access object for the table risk_feature_snapshot.
type RiskFeatureSnapshotDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  RiskFeatureSnapshotColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// RiskFeatureSnapshotColumns defines and stores column names for the table risk_feature_snapshot.
type RiskFeatureSnapshotColumns struct {
	Id           string //
	UserId       string //
	RiskScore    string //
	TagsJson     string //
	FeaturesJson string //
	UpdatedAt    string //
}

// riskFeatureSnapshotColumns holds the columns for the table risk_feature_snapshot.
var riskFeatureSnapshotColumns = RiskFeatureSnapshotColumns{
	Id:           "id",
	UserId:       "user_id",
	RiskScore:    "risk_score",
	TagsJson:     "tags_json",
	FeaturesJson: "features_json",
	UpdatedAt:    "updated_at",
}

// NewRiskFeatureSnapshotDao creates and returns a new DAO object for table data access.
func NewRiskFeatureSnapshotDao(handlers ...gdb.ModelHandler) *RiskFeatureSnapshotDao {
	return &RiskFeatureSnapshotDao{
		group:    "default",
		table:    "risk_feature_snapshot",
		columns:  riskFeatureSnapshotColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RiskFeatureSnapshotDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RiskFeatureSnapshotDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RiskFeatureSnapshotDao) Columns() RiskFeatureSnapshotColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RiskFeatureSnapshotDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RiskFeatureSnapshotDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *RiskFeatureSnapshotDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
