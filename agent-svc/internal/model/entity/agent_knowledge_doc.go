// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentKnowledgeDoc is the golang structure for table agent_knowledge_doc.
type AgentKnowledgeDoc struct {
	Id                 uint64      `json:"id"                 orm:"id"                   ` //
	KnowledgeDocNo     string      `json:"knowledgeDocNo"     orm:"knowledge_doc_no"     ` //
	Title              string      `json:"title"              orm:"title"                ` //
	KnowledgeScopeCode string      `json:"knowledgeScopeCode" orm:"knowledge_scope_code" ` //
	SourceTypeCode     string      `json:"sourceTypeCode"     orm:"source_type_code"     ` //
	SourceId           string      `json:"sourceId"           orm:"source_id"            ` //
	SourceVersion      uint64      `json:"sourceVersion"      orm:"source_version"       ` //
	ShopNo             string      `json:"shopNo"             orm:"shop_no"              ` //
	ContentTypeCode    string      `json:"contentTypeCode"    orm:"content_type_code"    ` //
	ContentText        string      `json:"contentText"        orm:"content_text"         ` //
	AssetIdsJson       string      `json:"assetIdsJson"       orm:"asset_ids_json"       ` //
	TagsJson           string      `json:"tagsJson"           orm:"tags_json"            ` //
	StatusCode         string      `json:"statusCode"         orm:"status_code"          ` //
	RejectReasonCode   string      `json:"rejectReasonCode"   orm:"reject_reason_code"   ` //
	RejectComment      string      `json:"rejectComment"      orm:"reject_comment"       ` //
	Version            uint64      `json:"version"            orm:"version"              ` //
	PublishedAt        *gtime.Time `json:"publishedAt"        orm:"published_at"         ` //
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           ` //
	UpdatedAt          *gtime.Time `json:"updatedAt"          orm:"updated_at"           ` //
	DeletedAt          *gtime.Time `json:"deletedAt"          orm:"deleted_at"           ` //
}
