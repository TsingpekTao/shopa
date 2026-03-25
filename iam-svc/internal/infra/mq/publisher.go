package mq

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// Publisher 鐎规矮绠?outbox 閸欐垵绔烽崳銊﹀复閸欙絻鈧
type Publisher interface {
	Enabled() bool
	PublishUserRegistered(ctx context.Context, eventID string, payload []byte) error
	Close() error
}

type noopPublisher struct{}

// NewNoopPublisher 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func NewNoopPublisher() Publisher {
	return &noopPublisher{}
}

// Enabled 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (p *noopPublisher) Enabled() bool {
	return false
}

// PublishUserRegistered 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (p *noopPublisher) PublishUserRegistered(ctx context.Context, eventID string, payload []byte) error {
	return nil
}

// Close 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (p *noopPublisher) Close() error {
	return nil
}

type rabbitConf struct {
	Enabled                  bool
	URL                      string
	Exchange                 string
	RoutingKeyUserRegistered string
}

// loadRabbitConf 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func loadRabbitConf(ctx context.Context) rabbitConf {
	return rabbitConf{
		Enabled:                  g.Cfg().MustGet(ctx, "mq.rabbitmq.enabled", false).Bool(),
		URL:                      strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.url", "").String()),
		Exchange:                 strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.exchange", "shopa.user.events").String()),
		RoutingKeyUserRegistered: strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.routingKeyUserRegistered", "iam.user.registered").String()),
	}
}

// NewPublisherFromConfig 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func NewPublisherFromConfig(ctx context.Context) (Publisher, error) {
	conf := loadRabbitConf(ctx)
	if !conf.Enabled {
		return NewNoopPublisher(), nil
	}
	return newRabbitPublisher(conf)
}
