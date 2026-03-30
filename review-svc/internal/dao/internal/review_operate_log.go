// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReviewOperateLogDao is the data access object for the table review_operate_log.
type ReviewOperateLogDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  ReviewOperateLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// ReviewOperateLogColumns defines and stores column names for the table review_operate_log.
type ReviewOperateLogColumns struct {
	Id           string //
	ReviewNo     string //
	OperatorType string //
	OperatorId   string //
	ActionCode   string //
	ReasonCode   string //
	Remark       string //
	SnapshotJson string //
	CreatedAt    string //
}

// reviewOperateLogColumns holds the columns for the table review_operate_log.
var reviewOperateLogColumns = ReviewOperateLogColumns{
	Id:           "id",
	ReviewNo:     "review_no",
	OperatorType: "operator_type",
	OperatorId:   "operator_id",
	ActionCode:   "action_code",
	ReasonCode:   "reason_code",
	Remark:       "remark",
	SnapshotJson: "snapshot_json",
	CreatedAt:    "created_at",
}

// NewReviewOperateLogDao creates and returns a new DAO object for table data access.
func NewReviewOperateLogDao(handlers ...gdb.ModelHandler) *ReviewOperateLogDao {
	return &ReviewOperateLogDao{
		group:    "default",
		table:    "review_operate_log",
		columns:  reviewOperateLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ReviewOperateLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ReviewOperateLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ReviewOperateLogDao) Columns() ReviewOperateLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ReviewOperateLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ReviewOperateLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ReviewOperateLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
