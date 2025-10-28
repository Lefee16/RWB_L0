package nats

import (
	"encoding/json"
	"fmt"
)

// Publisher - публикация сообщений в NATS
type Publisher struct {
	client *Client
}

// NewPublisher - создание publisher

// Publish - отправка сообщения в канал
func (p *Publisher) Publish(subject string, data interface{}) error {
	// Сериализуем в JSON
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Публикуем в NATS
	if err := p.client.conn.Publish(subject, payload); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// PublishBytes - отправка байтов в канал
func (p *Publisher) PublishBytes(subject string, data []byte) error {
	if err := p.client.conn.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}
