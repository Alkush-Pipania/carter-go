package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	ch       *amqp091.Channel
	exchange string
}

type ProducerConfig struct {
	Exchange     string
	ExchangeType string
	Durable      bool
}

func NewProducer(conn *Connection, cfg ProducerConfig) (*Producer, error) {
	ch, err := conn.Conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		cfg.Exchange,
		cfg.ExchangeType,
		cfg.Durable,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// create queue
	q, err := ch.QueueDeclare(
		"source.processor.queue", // Name
		true,                     // Durable
		false,                    // Delete when unused
		false,                    // Exclusive
		false,                    // No-wait
		nil,                      // Arguments
	)
	if err != nil {
		return nil, err
	}

	// bind it
	err = ch.QueueBind(
		q.Name,
		"source.process", // Routing Key
		cfg.Exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Producer{
		ch:       ch,
		exchange: cfg.Exchange,
	}, nil
}

func (p *Producer) Publish(
	ctx context.Context,
	routingKey string,
	payload any,
) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(
		ctx,
		p.exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *Producer) Close() error {
	return p.ch.Close()
}
