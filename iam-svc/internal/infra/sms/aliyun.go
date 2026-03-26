package sms

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/TsingpekTao/shopa/iam-svc/internal/consts"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gogf/gf/v2/errors/gerror"
)

type aliyunSender struct {
	client *dysmsapi.Client
	config AliyunConfig
}

// NewAliyunSender 创建阿里云短信发送器。
func NewAliyunSender(cfg AliyunConfig) (Sender, error) {
	if strings.TrimSpace(cfg.AccessKeyID) == "" || strings.TrimSpace(cfg.AccessKeySecret) == "" {
		return nil, gerror.New("sms.aliyun accessKeyId/accessKeySecret is required")
	}
	if strings.TrimSpace(cfg.SignName) == "" {
		return nil, gerror.New("sms.aliyun signName is required")
	}
	if strings.TrimSpace(cfg.Endpoint) == "" {
		cfg.Endpoint = "dysmsapi.aliyuncs.com"
	}
	client, err := dysmsapi.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(cfg.AccessKeyID),
		AccessKeySecret: tea.String(cfg.AccessKeySecret),
		Endpoint:        tea.String(cfg.Endpoint),
	})
	if err != nil {
		return nil, gerror.Wrap(err, "init aliyun sms client failed")
	}
	return &aliyunSender{
		client: client,
		config: cfg,
	}, nil
}

func (s *aliyunSender) Provider() string {
	return consts.SmsProviderAliyun
}

func (s *aliyunSender) SendCode(ctx context.Context, req *SendCodeRequest) (*SendCodeResult, error) {
	if req == nil {
		return nil, gerror.New("sms request is nil")
	}
	templateCode, err := s.templateCodeByScene(req.Scene)
	if err != nil {
		return nil, err
	}
	templateParamBytes, err := json.Marshal(map[string]string{
		"code": req.Code,
	})
	if err != nil {
		return nil, gerror.Wrap(err, "marshal sms template param failed")
	}
	resp, err := s.client.SendSms(&dysmsapi.SendSmsRequest{
		PhoneNumbers: tea.String(req.Phone),
		SignName:     tea.String(s.config.SignName),
		TemplateCode: tea.String(templateCode),
		TemplateParam: tea.String(
			string(templateParamBytes),
		),
	})
	if err != nil {
		return nil, gerror.Wrap(err, "aliyun send sms failed")
	}
	if resp == nil || resp.Body == nil {
		return nil, gerror.New("aliyun send sms empty response")
	}

	result := &SendCodeResult{
		Provider:   consts.SmsProviderAliyun,
		BizID:      tea.StringValue(resp.Body.BizId),
		RawCode:    tea.StringValue(resp.Body.Code),
		RawMessage: tea.StringValue(resp.Body.Message),
	}
	if !strings.EqualFold(result.RawCode, "OK") {
		return result, gerror.Newf("aliyun send sms failed: code=%s message=%s", result.RawCode, result.RawMessage)
	}
	return result, nil
}

func (s *aliyunSender) templateCodeByScene(scene string) (string, error) {
	var tpl string
	switch strings.ToLower(strings.TrimSpace(scene)) {
	case SceneRegister:
		tpl = s.config.TemplateCodeRegister
	case SceneLogin:
		tpl = s.config.TemplateCodeLogin
	case SceneMFA:
		tpl = s.config.TemplateCodeMFA
	case SceneResetPassword:
		tpl = s.config.TemplateCodeResetPassword
	default:
		return "", gerror.Newf("unsupported sms scene: %s", scene)
	}
	if strings.TrimSpace(tpl) == "" {
		return "", gerror.Newf("missing aliyun template code for scene: %s", scene)
	}
	return tpl, nil
}
