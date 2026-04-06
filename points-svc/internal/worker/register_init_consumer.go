package worker

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/TsingpekTao/shopa/points-svc/internal/service"
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

func StartRegisterInitConsumer(ctx context.Context) {
	registerInitOnce.Do(func() {
		go runRegisterInitConsumer()
	})
}

func runRegisterInitConsumer() {
	ctx := context.Background()
	conf := loadRegisterInitConsumerConf(ctx)
	if !conf.Enabled {
		g.Log().Info(ctx, "[points-svc] register-init consumer disabled")
		return
	}

	for {
		if err := consumeLoop(ctx, conf); err != nil {
			g.Log().Errorf(ctx, "[points-svc] register-init consumer failed: %+v", err)
			time.Sleep(5 * time.Second)
			continue
		}
	}
}

func consumeLoop(ctx context.Context, conf registerInitConsumerConf) error {
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

	for msg := range deliveries {
		if handleErr := handleRegisterMessage(ctx, msg.Body); handleErr != nil {
			g.Log().Errorf(ctx, "[points-svc] handle register-init message failed: %+v", handleErr)
			_ = msg.Nack(false, true)
			continue
		}
		_ = msg.Ack(false)
	}
	return nil
}

func handleRegisterMessage(ctx context.Context, body []byte) error {
	var payload struct {
		UserID uint64 `json:"user_id"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return err
	}
	return service.Points().InitAccountFromRegisterEvent(ctx, payload.UserID)
}

func loadRegisterInitConsumerConf(ctx context.Context) registerInitConsumerConf {
	return registerInitConsumerConf{
		Enabled:    g.Cfg().MustGet(ctx, "mq.rabbitmq.enabled", false).Bool(),
		URL:        strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.url", "").String()),
		Exchange:   strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.exchange", "shopa.user.events").String()),
		Queue:      strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.queueRegisterInit", "points.register-init.q").String()),
		RoutingKey: strings.TrimSpace(g.Cfg().MustGet(ctx, "mq.rabbitmq.routingKeyUserRegistered", "iam.user.registered").String()),
	}
}
