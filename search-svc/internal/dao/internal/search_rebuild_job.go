// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SearchRebuildJobDao is the data access object for the table search_rebuild_job.
type SearchRebuildJobDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  SearchRebuildJobColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// SearchRebuildJobColumns defines and stores column names for the table search_rebuild_job.
type SearchRebuildJobColumns struct {
	Id           string //
	JobNo        string //
	ReasonCode   string //
	Status       string //
	TotalCount   string //
	SuccessCount string //
	FailCount    string //
	StartedAt    string //
	FinishedAt   string //
	CreatedAt    string //
	UpdatedAt    string //
}

// searchRebuildJobColumns holds the columns for the table search_rebuild_job.
var searchRebuildJobColumns = SearchRebuildJobColumns{
	Id:           "id",
	JobNo:        "job_no",
	ReasonCode:   "reason_code",
	Status:       "status",
	TotalCount:   "total_count",
	SuccessCount: "success_count",
	FailCount:    "fail_count",
	StartedAt:    "started_at",
	FinishedAt:   "finished_at",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewSearchRebuildJobDao creates and returns a new DAO object for table data access.
func NewSearchRebuildJobDao(handlers ...gdb.ModelHandler) *SearchRebuildJobDao {
	return &SearchRebuildJobDao{
		group:    "default",
		table:    "search_rebuild_job",
		columns:  searchRebuildJobColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SearchRebuildJobDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SearchRebuildJobDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SearchRebuildJobDao) Columns() SearchRebuildJobColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SearchRebuildJobDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SearchRebuildJobDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SearchRebuildJobDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
