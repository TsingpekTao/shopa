// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentKnowledgeChunk is the golang structure for table agent_knowledge_chunk.
type AgentKnowledgeChunk struct {
	Id               uint64      `json:"id"               orm:"id"                 ` //
	ChunkNo          string      `json:"chunkNo"          orm:"chunk_no"           ` //
	KnowledgeDocNo   string      `json:"knowledgeDocNo"   orm:"knowledge_doc_no"   ` //
	SourceTypeCode   string      `json:"sourceTypeCode"   orm:"source_type_code"   ` //
	SourceId         string      `json:"sourceId"         orm:"source_id"          ` //
	SourceVersion    uint64      `json:"sourceVersion"    orm:"source_version"     ` //
	ChunkIndex       uint        `json:"chunkIndex"       orm:"chunk_index"        ` //
	ChunkText        string      `json:"chunkText"        orm:"chunk_text"         ` //
	ChunkTextPreview string      `json:"chunkTextPreview" orm:"chunk_text_preview" ` //
	VectorDocumentId string      `json:"vectorDocumentId" orm:"vector_document_id" ` //
	MetadataJson     string      `json:"metadataJson"     orm:"metadata_json"      ` //
	IsDeleted        int         `json:"isDeleted"        orm:"is_deleted"         ` //
	IndexStatusCode  string      `json:"indexStatusCode"  orm:"index_status_code"  ` //
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         ` //
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"         ` //
	DeletedAt        *gtime.Time `json:"deletedAt"        orm:"deleted_at"         ` //
}
