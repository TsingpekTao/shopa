// =================================================================================
// 说明：该文件由 GoFrame CLI 工具生成，可按需修改。
// =================================================================================

package dao

import (
	"github.com/TsingpekTao/shopa/edge-gateway/internal/dao/internal"
)

// edgeRequestAuditDao 是 edge_request_audit 表的数据访问对象。
// 可在此处添加自定义方法以扩展功能。
type edgeRequestAuditDao struct {
	*internal.EdgeRequestAuditDao
}

var (
	// EdgeRequestAudit 是 edge_request_audit 表全局可用的操作对象。
	EdgeRequestAudit = edgeRequestAuditDao{internal.NewEdgeRequestAuditDao()}
)

// 在此添加自定义方法或业务逻辑。
