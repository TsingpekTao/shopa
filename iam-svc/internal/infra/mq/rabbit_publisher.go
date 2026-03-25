package mq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitPublisher struct {
	conf rabbitConf

	mu   sync.Mutex
	conn *amqp.Connection
	ch   *amqp.Channel
}

// newRabbitPublisher 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func newRabbitPublisher(conf rabbitConf) (Publisher, error) {
	p := &rabbitPublisher{conf: conf}
	if err := p.ensureConnected(); err != nil {
		return nil, err
	}
	return p, nil
}

// Enabled 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (p *rabbitPublisher) Enabled() bool {
	return true
}

// PublishUserRegistered 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (p *rabbitPublisher) PublishUserRegistered(ctx context.Context, eventID string, payload []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.ensureConnected(); err != nil {
		return err
	}

	err := p.ch.PublishWithContext(ctx, p.conf.Exchange, p.conf.RoutingKeyUserRegistered, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         payload,
		DeliveryMode: amqp.Persistent,
		MessageId:    eventID,
		Type:         "UserRegisteredV1",
		Timestamp:    time.Now().UTC(),
	})
	if err != nil {
		_ = p.closeLocked()
		return err
	}
	return nil
}

// Close 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (p *rabbitPublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.closeLocked()
}

// ensureConnected 纭繚鍓嶇疆鐘舵€佸瓨鍦紝涓嶆弧瓒冲垯鑷姩琛ラ綈銆
func (p *rabbitPublisher) ensureConnected() error {
	if p.conn != nil && !p.conn.IsClosed() && p.ch != nil {
		return nil
	}

	if p.conf.URL == "" {
		return fmt.Errorf("rabbitmq url is empty")
	}
	conn, err := amqp.Dial(p.conf.URL)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}
	if err = ch.ExchangeDeclare(p.conf.Exchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}

	p.conn = conn
	p.ch = ch
	return nil
}

// closeLocked 瀹炵幇璇ュ嚱鏁板搴旂殑鏍稿績涓氬姟閫昏緫銆
func (p *rabbitPublisher) closeLocked() error {
	var firstErr error
	if p.ch != nil {
		if err := p.ch.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		p.ch = nil
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		p.conn = nil
	}
	return firstErr
}
