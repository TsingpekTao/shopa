package sms

import (
	"context"
	"strings"

	"github.com/TsingpekTao/shopa/iam-svc/internal/consts"
	"github.com/gogf/gf/v2/frame/g"
)

// Config 聚合短信通道配置。
type Config struct {
	Provider string
	Aliyun   AliyunConfig
}

// AliyunConfig 描述阿里云短信配置。
type AliyunConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	SignName        string
	Endpoint        string

	TemplateCodeRegister      string
	TemplateCodeLogin         string
	TemplateCodeMFA           string
	TemplateCodeResetPassword string
}

// LoadFromConfig 按统一键读取短信配置。
func LoadFromConfig(ctx context.Context) Config {
	return Config{
		Provider: strings.ToLower(strings.TrimSpace(g.Cfg().MustGet(ctx, "sms.provider", consts.SmsProviderMock).String())),
		Aliyun: AliyunConfig{
			AccessKeyID:               strings.TrimSpace(g.Cfg().MustGet(ctx, "sms.aliyun.accessKeyId", "").String()),
			AccessKeySecret:           strings.TrimSpace(g.Cfg().MustGet(ctx, "sms.aliyun.accessKeySecret", "").String()),
			SignName:                  strings.TrimSpace(g.Cfg().MustGet(ctx, "sms.aliyun.signName", "").String()),
			Endpoint:                  strings.TrimSpace(g.Cfg().MustGet(ctx, "sms.aliyun.endpoint", "dysmsapi.aliyuncs.com").String()),
			TemplateCodeRegister:      strings.TrimSpace(g.Cfg().MustGet(ctx, "sms.aliyun.templateCodeRegister", "").String()),
			TemplateCodeLogin:         strings.TrimSpace(g.Cfg().MustGet(ctx, "sms.aliyun.templateCodeLogin", "").String()),
			TemplateCodeMFA:           strings.TrimSpace(g.Cfg().MustGet(ctx, "sms.aliyun.templateCodeMfa", "").String()),
			TemplateCodeResetPassword: strings.TrimSpace(g.Cfg().MustGet(ctx, "sms.aliyun.templateCodeResetPassword", "").String()),
		},
	}
}
