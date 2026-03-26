// =================================================================================
// 说明：该文件由 GoFrame CLI 工具生成，可按需修改。
// =================================================================================

package dao

import (
	"github.com/TsingpekTao/shopa/edge-gateway/internal/dao/internal"
)

// edgeProxyRouteDao 是 edge_proxy_route 表的数据访问对象。
// 可在此处添加自定义方法以扩展功能。
type edgeProxyRouteDao struct {
	*internal.EdgeProxyRouteDao
}

var (
	// EdgeProxyRoute 是 edge_proxy_route 表全局可用的操作对象。
	EdgeProxyRoute = edgeProxyRouteDao{internal.NewEdgeProxyRouteDao()}
)

// 在此添加自定义方法或业务逻辑。
