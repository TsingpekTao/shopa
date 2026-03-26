package sms

import (
	"context"

	"github.com/TsingpekTao/shopa/iam-svc/internal/consts"
	"github.com/gogf/gf/v2/errors/gerror"
)

// NewSenderFromConfig 按配置创建短信发送器。
func NewSenderFromConfig(ctx context.Context) (Sender, error) {
	cfg := LoadFromConfig(ctx)
	switch cfg.Provider {
	case "", consts.SmsProviderMock:
		return NewMockSender(), nil
	case consts.SmsProviderAliyun:
		return NewAliyunSender(cfg.Aliyun)
	default:
		return nil, gerror.Newf("unsupported sms provider: %s", cfg.Provider)
	}
}
