// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AgentKnowledgeChunkDao is the data access object for the table agent_knowledge_chunk.
type AgentKnowledgeChunkDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  AgentKnowledgeChunkColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// AgentKnowledgeChunkColumns defines and stores column names for the table agent_knowledge_chunk.
type AgentKnowledgeChunkColumns struct {
	Id               string //
	ChunkNo          string //
	KnowledgeDocNo   string //
	SourceTypeCode   string //
	SourceId         string //
	SourceVersion    string //
	ChunkIndex       string //
	ChunkText        string //
	ChunkTextPreview string //
	VectorDocumentId string //
	MetadataJson     string //
	IsDeleted        string //
	IndexStatusCode  string //
	CreatedAt        string //
	UpdatedAt        string //
	DeletedAt        string //
}

// agentKnowledgeChunkColumns holds the columns for the table agent_knowledge_chunk.
var agentKnowledgeChunkColumns = AgentKnowledgeChunkColumns{
	Id:               "id",
	ChunkNo:          "chunk_no",
	KnowledgeDocNo:   "knowledge_doc_no",
	SourceTypeCode:   "source_type_code",
	SourceId:         "source_id",
	SourceVersion:    "source_version",
	ChunkIndex:       "chunk_index",
	ChunkText:        "chunk_text",
	ChunkTextPreview: "chunk_text_preview",
	VectorDocumentId: "vector_document_id",
	MetadataJson:     "metadata_json",
	IsDeleted:        "is_deleted",
	IndexStatusCode:  "index_status_code",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewAgentKnowledgeChunkDao creates and returns a new DAO object for table data access.
func NewAgentKnowledgeChunkDao(handlers ...gdb.ModelHandler) *AgentKnowledgeChunkDao {
	return &AgentKnowledgeChunkDao{
		group:    "default",
		table:    "agent_knowledge_chunk",
		columns:  agentKnowledgeChunkColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AgentKnowledgeChunkDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AgentKnowledgeChunkDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AgentKnowledgeChunkDao) Columns() AgentKnowledgeChunkColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AgentKnowledgeChunkDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AgentKnowledgeChunkDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AgentKnowledgeChunkDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
