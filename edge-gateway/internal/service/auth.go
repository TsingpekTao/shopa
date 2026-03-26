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
	// IAuth 定义网关鉴权门面能力：
	// - ExtractAccessToken: 从 HTTP 请求提取 access token。
	// - VerifyAccessToken: 调 IAM 校验 token 并返回统一用户身份结果。
	IAuth interface {
		ExtractAccessToken(r *ghttp.Request) string
		VerifyAccessToken(ctx context.Context, accessToken string) (*VerifyResult, error)
	}
)

var (
	// localAuth 保存已注册的 IAuth 实现。
	localAuth IAuth
)

// Auth 返回鉴权服务实现；未注册时直接 panic，避免静默失败。
func Auth() IAuth {
	if localAuth == nil {
		panic("implement not found for interface IAuth, forgot register?")
	}
	return localAuth
}

// RegisterAuth 在 init 阶段注册鉴权实现。
func RegisterAuth(i IAuth) {
	localAuth = i
}
