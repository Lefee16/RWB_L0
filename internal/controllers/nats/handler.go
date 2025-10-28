package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/stan.go"

	"RWB_L0/internal/dto"
	"RWB_L0/internal/usecase"
	"RWB_L0/pkg/logger"
)

// Handler - обработчик NATS сообщений
type Handler struct {
	orderUseCase *usecase.OrderUseCase
	log          logger.Logger
}

// NewHandler - создание handler
func NewHandler(orderUseCase *usecase.OrderUseCase, log logger.Logger) *Handler {
	return &Handler{
		orderUseCase: orderUseCase,
		log:          log,
	}
}

// HandleOrderCreate - обработка создания заказа
func (h *Handler) HandleOrderCreate(msg *stan.Msg) error {
	h.log.Debug("Received order creation message: sequence=%d", msg.Sequence)

	// Десериализуем JSON
	var input dto.CreateOrderInput
	if err := json.Unmarshal(msg.Data, &input); err != nil {
		h.log.Error("Failed to unmarshal order: %v", err)
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Валидация order_uid
	if input.OrderUID == "" {
		h.log.Error("Received order with empty order_uid")
		return fmt.Errorf("order_uid is required")
	}

	h.log.Info("Processing order: %s", input.OrderUID)

	// Создаём заказ через Use Case
	if err := h.orderUseCase.CreateOrder(context.Background(), &input); err != nil {
		h.log.Error("Failed to create order %s: %v", input.OrderUID, err)
		return fmt.Errorf("failed to create order: %w", err)
	}

	h.log.Info("Order %s successfully created and cached", input.OrderUID)
	return nil
}
