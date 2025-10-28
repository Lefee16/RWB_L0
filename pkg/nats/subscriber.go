package nats

import (
	"fmt"

	"github.com/nats-io/stan.go"
)

// MessageHandler - функция-обработчик сообщений
type MessageHandler func(msg *stan.Msg) error

// Subscriber - подписка на каналы NATS
type Subscriber struct {
	client *Client
}

// NewSubscriber - создание subscriber
func NewSubscriber(client *Client) *Subscriber {
	return &Subscriber{client: client}
}

// Subscribe - подписка на канал с обработчиком
func (s *Subscriber) Subscribe(subject string, durableName string, handler MessageHandler) (stan.Subscription, error) {
	sub, err := s.client.conn.Subscribe(
		subject,
		func(msg *stan.Msg) {
			// Вызываем пользовательский обработчик
			if err := handler(msg); err != nil {
				// Логируем ошибку, но не подтверждаем сообщение
				// NATS попробует доставить повторно
				return
			}
			// Подтверждаем успешную обработку
			_ = msg.Ack()
		},
		stan.DurableName(durableName), // Durable subscription
		stan.SetManualAckMode(),       // Ручное подтверждение
		stan.DeliverAllAvailable(),    // Получить все непрочитанные сообщения
		stan.MaxInflight(1),           // Обрабатывать по одному
	)

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return sub, nil
}

// SubscribeFromLatest - подписка только на новые сообщения
func (s *Subscriber) SubscribeFromLatest(subject string, durableName string, handler MessageHandler) (stan.Subscription, error) {
	sub, err := s.client.conn.Subscribe(
		subject,
		func(msg *stan.Msg) {
			if err := handler(msg); err != nil {
				return
			}
			_ = msg.Ack()
		},
		stan.DurableName(durableName),
		stan.SetManualAckMode(),
		stan.StartWithLastReceived(), // Только новые сообщения
	)

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return sub, nil
}
