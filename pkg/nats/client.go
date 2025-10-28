package nats

import (
	"fmt"

	"github.com/nats-io/stan.go"
)

// Client - NATS Streaming клиент
type Client struct {
	conn stan.Conn
}

// Config - конфигурация NATS
type Config struct {
	URL       string
	ClusterID string
	ClientID  string
}

// New - создание NATS клиента
func New(cfg *Config) (*Client, error) {
	conn, err := stan.Connect(
		cfg.ClusterID,
		cfg.ClientID,
		stan.NatsURL(cfg.URL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &Client{conn: conn}, nil
}

// Close - закрытие соединения
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetConnection - получить низкоуровневое соединение
func (c *Client) GetConnection() stan.Conn {
	return c.conn
}
