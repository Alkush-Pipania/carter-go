package rabbitmq

import (
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Connection wraps amqp.Connection with reconnection logic
type Connection struct {
	conn    *amqp.Connection
	url     string
	mu      sync.RWMutex
	closed  bool
	logger  *zap.Logger
	onClose chan *amqp.Error
}

// ConnectionConfig holds connection configuration
type ConnectionConfig struct {
	URL            string
	ReconnectDelay time.Duration // Base delay between reconnection attempts
	MaxReconnects  int           // Max reconnection attempts (0 = infinite)
}

// DefaultConfig returns sensible defaults
func DefaultConfig(url string) ConnectionConfig {
	return ConnectionConfig{
		URL:            url,
		ReconnectDelay: 5 * time.Second,
		MaxReconnects:  0, // infinite retries
	}
}

// NewConnection creates a new RabbitMQ connection with reconnection support
func NewConnection(cfg ConnectionConfig, logger *zap.Logger) (*Connection, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	c := &Connection{
		conn:    conn,
		url:     cfg.URL,
		logger:  logger,
		onClose: make(chan *amqp.Error),
	}

	// Register close notifier
	conn.NotifyClose(c.onClose)

	// Start reconnection goroutine
	go c.handleReconnect(cfg)

	logger.Info("RabbitMQ connection established")
	return c, nil
}

// handleReconnect monitors connection and reconnects if dropped
func (c *Connection) handleReconnect(cfg ConnectionConfig) {
	for {
		select {
		case err := <-c.onClose:
			if err == nil {
				// Graceful close, exit
				return
			}

			c.logger.Warn("RabbitMQ connection lost, attempting reconnect",
				zap.Error(err))

			c.reconnect(cfg)
		}
	}
}

// reconnect attempts to re-establish connection with exponential backoff
func (c *Connection) reconnect(cfg ConnectionConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	attempts := 0
	delay := cfg.ReconnectDelay

	for {
		attempts++
		c.logger.Info("Attempting to reconnect to RabbitMQ",
			zap.Int("attempt", attempts))

		conn, err := amqp.Dial(cfg.URL)
		if err == nil {
			// Close the previous onClose channel to avoid leaking goroutines
			oldOnClose := c.onClose
			if oldOnClose != nil {
				close(oldOnClose)
			}
			c.conn = conn
			c.onClose = make(chan *amqp.Error)
			conn.NotifyClose(c.onClose)
			c.logger.Info("RabbitMQ reconnection successful")
			return
		}

		c.logger.Error("Reconnection failed",
			zap.Error(err),
			zap.Duration("retry_in", delay))

		if cfg.MaxReconnects > 0 && attempts >= cfg.MaxReconnects {
			c.logger.Error("Max reconnection attempts reached",
				zap.Int("attempts", attempts))
			c.closed = true
			return
		}

		time.Sleep(delay)
		// Exponential backoff capped at 60 seconds
		delay = time.Duration(float64(delay) * 1.5)
		if delay > 60*time.Second {
			delay = 60 * time.Second
		}
	}
}

// Channel creates a new channel from the connection
func (c *Connection) Channel() (*amqp.Channel, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return nil, fmt.Errorf("connection is nil")
	}

	return c.conn.Channel()
}

// Close closes the connection gracefully
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.closed = true
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsClosed returns whether connection is closed
func (c *Connection) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed || c.conn.IsClosed()
}
