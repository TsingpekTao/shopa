// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AgentKnowledgeDoc is the golang structure of table agent_knowledge_doc for DAO operations like Where/Data.
type AgentKnowledgeDoc struct {
	g.Meta             `orm:"table:agent_knowledge_doc, do:true"`
	Id                 any         //
	KnowledgeDocNo     any         //
	Title              any         //
	KnowledgeScopeCode any         //
	SourceTypeCode     any         //
	SourceId           any         //
	SourceVersion      any         //
	ShopNo             any         //
	ContentTypeCode    any         //
	ContentText        any         //
	AssetIdsJson       any         //
	TagsJson           any         //
	StatusCode         any         //
	RejectReasonCode   any         //
	RejectComment      any         //
	Version            any         //
	PublishedAt        *gtime.Time //
	CreatedAt          *gtime.Time //
	UpdatedAt          *gtime.Time //
	DeletedAt          *gtime.Time //
}
