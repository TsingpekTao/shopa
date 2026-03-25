package worker

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/TsingpekTao/shopa/user-profile-svc/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	amqp "github.com/rabbitmq/amqp091-go"
)

type registerInitConsumerConf struct {
	Enabled    bool
	URL        string
	Exchange   string
	Queue      string
	RoutingKey string
}

var registerInitOnce sync.Once

// StartRegisterInitConsumer 只启动一次 register-init 消费者。
func StartRegisterInitConsumer(ctx context.Context) {
	registerInitOnce.Do(func() {
		go runRegisterInitConsumer()
	})
}

// runRegisterInitConsumer 负责持续执行注册事件消费循环。
func runRegisterInitConsumer() {
	ctx := context.Background()
	conf := loadRegisterInitConsumerConf(ctx)
	if !conf.Enabled {
		g.Log().Info(ctx, "[user-profile-svc] register-init consumer disabled")
		return
	}

	for {
		if err := consumeRegisterInitLoop(ctx, conf); err != nil {
			g.Log().Errorf(ctx, "[user-profile-svc] register-init consumer loop failed: %+v", err)
			time.Sleep(5 * time.Second)
			continue
		}
	}
}

// consumeRegisterInitLoop 连接 RabbitMQ 并消费 register-init 消息。
func consumeRegisterInitLoop(ctx context.Context, conf registerInitConsumerConf) error {
	conn, err := amqp.Dial(conf.URL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err = ch.ExchangeDeclare(conf.Exchange, "topic", true, false, false, false, nil); err != nil {
		return err
	}
	queue, err := ch.QueueDeclare(conf.Queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	if err = ch.QueueBind(queue.Name, conf.RoutingKey, conf.Exchange, false, nil); err != nil {
		return err
	}

	deliveries, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	svc := service.UserProfile()
	for msg := range deliveries {
		if handleErr := handleRegisterInitMessage(ctx, svc, msg); handleErr != nil {
			g.Log().Errorf(ctx, "[user-profile-svc] handle register init message failed: %+v", handleErr)
			_ = msg.Nack(false, true)
			continue
		}
		_ = msg.Ack(false)
	}
	return nil
}

// handleRegisterInitMessage 解析消息并调用服务处理事件。
func handleRegisterInitMessage(ctx context.Context, svc service.IUserProfile, msg amqp.Delivery) error {
	var payload struct {
		EventID         string `json:"event_id"`
		EventVersion    string `json:"event_version"`
		UserID          uint64 `json:"user_id"`
		InitDisplayName string `json:"init_display_name"`
		OccurredAt      string `json:"occurred_at"`
	}
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		return err
	}

	occurredAt := time.Now().UTC()
	rawOccurredAt := strings.TrimSpace(payload.OccurredAt)
	if rawOccurredAt != "" {
		if t, err := time.Parse(time.RFC3339Nano, rawOccurredAt); err == nil {
			occurredAt = t.UTC()
		} else if t2, err2 := time.Parse(time.RFC3339, rawOccurredAt); err2 == nil {
			occurredAt = t2.UTC()
		}
	}

	return svc.ApplyRegisterInitEvent(ctx, service.RegisterInitEvent{
		EventID:         payload.EventID,
		EventVersion:    payload.EventVersion,
		UserID:          payload.UserID,
		InitDisplayName: payload.InitDisplayName,
		OccurredAt:      occurredAt,
	})
}

// loadRegisterInitConsumerConf 读取 register-init 消费者的配置。
func loadRegisterInitConsumerConf(ctx context.Context) registerInitConsumerConf {
	return registerInitConsumerConf{
		Enabled:    g.Cfg().MustGet(ctx, "mq.rabbitmq.enabled", false).Bool(),
		URL:        strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.url", "").String()),
		Exchange:   strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.exchange", "shopa.user.events").String()),
		Queue:      strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.queueRegisterInit", "user-profile.register-init.q").String()),
		RoutingKey: strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.routingKeyUserRegistered", "iam.user.registered").String()),
	}
}
