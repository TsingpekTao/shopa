// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AgentKnowledgeDocDao is the data access object for the table agent_knowledge_doc.
type AgentKnowledgeDocDao struct {
	table    string                   // table is the underlying table name of the DAO.
	group    string                   // group is the database configuration group name of the current DAO.
	columns  AgentKnowledgeDocColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler       // handlers for customized model modification.
}

// AgentKnowledgeDocColumns defines and stores column names for the table agent_knowledge_doc.
type AgentKnowledgeDocColumns struct {
	Id                 string //
	KnowledgeDocNo     string //
	Title              string //
	KnowledgeScopeCode string //
	SourceTypeCode     string //
	SourceId           string //
	SourceVersion      string //
	ShopNo             string //
	ContentTypeCode    string //
	ContentText        string //
	AssetIdsJson       string //
	TagsJson           string //
	StatusCode         string //
	RejectReasonCode   string //
	RejectComment      string //
	Version            string //
	PublishedAt        string //
	CreatedAt          string //
	UpdatedAt          string //
	DeletedAt          string //
}

// agentKnowledgeDocColumns holds the columns for the table agent_knowledge_doc.
var agentKnowledgeDocColumns = AgentKnowledgeDocColumns{
	Id:                 "id",
	KnowledgeDocNo:     "knowledge_doc_no",
	Title:              "title",
	KnowledgeScopeCode: "knowledge_scope_code",
	SourceTypeCode:     "source_type_code",
	SourceId:           "source_id",
	SourceVersion:      "source_version",
	ShopNo:             "shop_no",
	ContentTypeCode:    "content_type_code",
	ContentText:        "content_text",
	AssetIdsJson:       "asset_ids_json",
	TagsJson:           "tags_json",
	StatusCode:         "status_code",
	RejectReasonCode:   "reject_reason_code",
	RejectComment:      "reject_comment",
	Version:            "version",
	PublishedAt:        "published_at",
	CreatedAt:          "created_at",
	UpdatedAt:          "updated_at",
	DeletedAt:          "deleted_at",
}

// NewAgentKnowledgeDocDao creates and returns a new DAO object for table data access.
func NewAgentKnowledgeDocDao(handlers ...gdb.ModelHandler) *AgentKnowledgeDocDao {
	return &AgentKnowledgeDocDao{
		group:    "default",
		table:    "agent_knowledge_doc",
		columns:  agentKnowledgeDocColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AgentKnowledgeDocDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AgentKnowledgeDocDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AgentKnowledgeDocDao) Columns() AgentKnowledgeDocColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AgentKnowledgeDocDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AgentKnowledgeDocDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AgentKnowledgeDocDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
