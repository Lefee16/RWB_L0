package usecase

import (
	"RWB_L0/internal/domain"
	"context"
)

// OrderRepository - интерфейс для работы с БД
type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, orderUID string) (*domain.Order, error)
	GetAll(ctx context.Context) ([]*domain.Order, error)
	Delete(ctx context.Context, orderUID string) error
	Count(ctx context.Context) (int, error)
}

// Cache - интерфейс для работы с кэшем
type Cache interface {
	Set(orderUID string, order *domain.Order) error
	Get(orderUID string) (*domain.Order, error)
	Delete(orderUID string) error
	LoadAll(orders []*domain.Order) error
	GetAll() ([]*domain.Order, error)
	Count() int
	Clear() error
}
