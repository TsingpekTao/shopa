package mq

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// Publisher 定义 outbox 事件发布能力。
type Publisher interface {
	Enabled() bool
	PublishUserRegistered(ctx context.Context, eventID string, payload []byte) error
	Close() error
}

type noopPublisher struct{}

// NewNoopPublisher 创建空实现发布器。
func NewNoopPublisher() Publisher {
	return &noopPublisher{}
}

// Enabled 返回当前发布器是否启用。
func (p *noopPublisher) Enabled() bool {
	return false
}

// PublishUserRegistered 发布用户注册事件。
func (p *noopPublisher) PublishUserRegistered(ctx context.Context, eventID string, payload []byte) error {
	return nil
}

// Close 关闭发布器并释放资源。
func (p *noopPublisher) Close() error {
	return nil
}

type rabbitConf struct {
	Enabled                  bool
	URL                      string
	Exchange                 string
	RoutingKeyUserRegistered string
}

// loadRabbitConf 加载 RabbitMQ 发布配置。
func loadRabbitConf(ctx context.Context) rabbitConf {
	return rabbitConf{
		Enabled:                  g.Cfg().MustGet(ctx, "mq.rabbitmq.enabled", false).Bool(),
		URL:                      strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.url", "").String()),
		Exchange:                 strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.exchange", "shopa.user.events").String()),
		RoutingKeyUserRegistered: strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.routingKeyUserRegistered", "iam.user.registered").String()),
	}
}

// NewPublisherFromConfig 按配置创建发布器。
func NewPublisherFromConfig(ctx context.Context) (Publisher, error) {
	conf := loadRabbitConf(ctx)
	if !conf.Enabled {
		return NewNoopPublisher(), nil
	}
	return newRabbitPublisher(conf)
}
