package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Producer publishes messages to RabbitMQ
type Producer struct {
	conn     *Connection
	channel  *amqp.Channel
	exchange string
	mu       sync.RWMutex
	logger   *zap.Logger
	confirms chan amqp.Confirmation
}

// ProducerConfig holds producer configuration
type ProducerConfig struct {
	Exchange     string
	ExchangeType string // "direct", "topic", "fanout"
	Durable      bool
}

// NewProducer creates a new producer
func NewProducer(conn *Connection, cfg ProducerConfig, logger *zap.Logger) (*Producer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		cfg.Exchange,     // name
		cfg.ExchangeType, // type
		cfg.Durable,      // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Enable publisher confirms for reliability
	err = ch.Confirm(false)
	if err != nil {
		return nil, fmt.Errorf("failed to enable confirms: %w", err)
	}

	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 100))

	p := &Producer{
		conn:     conn,
		channel:  ch,
		exchange: cfg.Exchange,
		logger:   logger,
		confirms: confirms,
	}

	logger.Info("RabbitMQ producer created",
		zap.String("exchange", cfg.Exchange))

	return p, nil
}

// Message represents a message to be published
type Message struct {
	RoutingKey    string
	CorrelationID string
	Body          interface{}
	Headers       map[string]interface{}
}

// Publish sends a message to the exchange
func (p *Producer) Publish(ctx context.Context, msg Message) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	body, err := json.Marshal(msg.Body)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	publishing := amqp.Publishing{
		ContentType:   "application/json",
		DeliveryMode:  amqp.Persistent, // Message survives broker restart
		Timestamp:     time.Now(),
		CorrelationId: msg.CorrelationID,
		Body:          body,
	}

	if msg.Headers != nil {
		publishing.Headers = amqp.Table(msg.Headers)
	}

	err = p.channel.PublishWithContext(
		ctx,
		p.exchange,     // exchange
		msg.RoutingKey, // routing key
		false,          // mandatory
		false,          // immediate
		publishing,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	// Wait for confirmation
	select {
	case confirm := <-p.confirms:
		if !confirm.Ack {
			return fmt.Errorf("message was nacked by broker")
		}
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return fmt.Errorf("confirmation timeout")
	}

	p.logger.Debug("Message published",
		zap.String("routing_key", msg.RoutingKey),
		zap.String("correlation_id", msg.CorrelationID))

	return nil
}

// Close closes the producer channel
func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.channel != nil {
		return p.channel.Close()
	}
	return nil
}
