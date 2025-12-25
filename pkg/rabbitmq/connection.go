package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	Conn *amqp091.Connection
}

type Config struct {
	URL string
}

func NewConnection(cfg Config) (*Connection, error) {
	conn, err := amqp091.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}
	return &Connection{Conn: conn}, nil
}

func (c *Connection) Close() error {
	return c.Conn.Close()
}
