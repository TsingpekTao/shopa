// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ReviewSpuSummaryDao is the data access object for the table review_spu_summary.
type ReviewSpuSummaryDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  ReviewSpuSummaryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// ReviewSpuSummaryColumns defines and stores column names for the table review_spu_summary.
type ReviewSpuSummaryColumns struct {
	Id               string //
	SpuNo            string //
	TotalReviews     string //
	Score1Count      string //
	Score2Count      string //
	Score3Count      string //
	Score4Count      string //
	Score5Count      string //
	AvgScoreX100     string //
	PositiveRateX100 string //
	CreatedAt        string //
	UpdatedAt        string //
	DeletedAt        string //
}

// reviewSpuSummaryColumns holds the columns for the table review_spu_summary.
var reviewSpuSummaryColumns = ReviewSpuSummaryColumns{
	Id:               "id",
	SpuNo:            "spu_no",
	TotalReviews:     "total_reviews",
	Score1Count:      "score_1_count",
	Score2Count:      "score_2_count",
	Score3Count:      "score_3_count",
	Score4Count:      "score_4_count",
	Score5Count:      "score_5_count",
	AvgScoreX100:     "avg_score_x100",
	PositiveRateX100: "positive_rate_x100",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewReviewSpuSummaryDao creates and returns a new DAO object for table data access.
func NewReviewSpuSummaryDao(handlers ...gdb.ModelHandler) *ReviewSpuSummaryDao {
	return &ReviewSpuSummaryDao{
		group:    "default",
		table:    "review_spu_summary",
		columns:  reviewSpuSummaryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ReviewSpuSummaryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ReviewSpuSummaryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ReviewSpuSummaryDao) Columns() ReviewSpuSummaryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ReviewSpuSummaryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ReviewSpuSummaryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ReviewSpuSummaryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
