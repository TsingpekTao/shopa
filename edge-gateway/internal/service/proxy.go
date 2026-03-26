// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// 说明：如需手动维护此接口文件，可删除这些注释。
// ================================================================================

package service

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	// IProxy 定义网关反向代理门面。
	// handled=true 表示请求已被代理层消费，不应再继续后续 handler。
	IProxy interface {
		HandleProxyRequest(ctx context.Context, r *ghttp.Request) (bool, error)
	}
)

var (
	// localProxy 保存已注册的代理实现。
	localProxy IProxy
)

// Proxy 返回代理实现；未注册时 panic，避免无感知漏配。
func Proxy() IProxy {
	if localProxy == nil {
		panic("implement not found for interface IProxy, forgot register?")
	}
	return localProxy
}

// RegisterProxy 注册代理实现。
func RegisterProxy(i IProxy) {
	localProxy = i
}
