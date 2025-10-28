package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/stan.go"

	"RWB_L0/pkg/logger"
	pkgnats "RWB_L0/pkg/nats"
)

// Consumer - NATS потребитель
type Consumer struct {
	subscriber *pkgnats.Subscriber
	handler    *Handler
	log        logger.Logger
	sub        stan.Subscription
}

// NewConsumer - создание consumer
func NewConsumer(subscriber *pkgnats.Subscriber, handler *Handler, log logger.Logger) *Consumer {
	return &Consumer{
		subscriber: subscriber,
		handler:    handler,
		log:        log,
	}
}

// Start - запуск подписки
func (c *Consumer) Start(ctx context.Context, subject string, durableName string) error {
	c.log.Info("Starting NATS consumer for subject: %s", subject)

	// Подписываемся с handler
	sub, err := c.subscriber.Subscribe(subject, durableName, c.handler.HandleOrderCreate)
	if err != nil {
		return fmt.Errorf("failed to subscribe to NATS: %w", err)
	}

	c.sub = sub
	c.log.Info("NATS consumer started successfully")

	// Ждём сигнала остановки
	<-ctx.Done()

	return c.Stop()
}

// Stop - остановка consumer
func (c *Consumer) Stop() error {
	c.log.Info("Stopping NATS consumer...")

	if c.sub != nil {
		if err := c.sub.Unsubscribe(); err != nil {
			c.log.Error("Failed to unsubscribe: %v", err)
			return err
		}
	}

	c.log.Info("NATS consumer stopped successfully")
	return nil
}
