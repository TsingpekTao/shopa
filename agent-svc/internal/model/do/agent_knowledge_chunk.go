// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentKnowledgeChunk is the golang structure of table agent_knowledge_chunk for DAO operations like Where/Data.
type AgentKnowledgeChunk struct {
	g.Meta           `orm:"table:agent_knowledge_chunk, do:true"`
	Id               any         //
	ChunkNo          any         //
	KnowledgeDocNo   any         //
	SourceTypeCode   any         //
	SourceId         any         //
	SourceVersion    any         //
	ChunkIndex       any         //
	ChunkText        any         //
	ChunkTextPreview any         //
	VectorDocumentId any         //
	MetadataJson     any         //
	IsDeleted        any         //
	IndexStatusCode  any         //
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
	DeletedAt        *gtime.Time //
}
