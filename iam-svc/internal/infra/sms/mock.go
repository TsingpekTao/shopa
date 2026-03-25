package sms

import (
	"context"

	"github.com/TsingpekTao/shopa/iam-svc/internal/consts"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type mockSender struct{}

// NewMockSender 创建本地开发用 mock 通道。
func NewMockSender() Sender {
	return &mockSender{}
}

func (s *mockSender) Provider() string {
	return consts.SmsProviderMock
}

func (s *mockSender) SendCode(ctx context.Context, req *SendCodeRequest) (*SendCodeResult, error) {
	if req == nil {
		return nil, gerror.New("sms request is nil")
	}
	// mock 场景打印验证码，方便开发联调（生产环境请勿使用 mock）。
	g.Log().Infof(ctx, "[iam-sms][mock] scene=%s phone=%s code=%s request_id=%s", req.Scene, req.Phone, req.Code, req.RequestID)
	return &SendCodeResult{
		Provider:   consts.SmsProviderMock,
		BizID:      "",
		RawCode:    "OK",
		RawMessage: "mock accepted",
	}, nil
}
