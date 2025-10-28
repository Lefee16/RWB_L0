package usecase

import (
	"context"

	"RWB_L0/internal/domain"
	"RWB_L0/internal/dto"
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

// OrderUseCaseInterface определяет контракт для бизнес-логики заказов
type OrderUseCaseInterface interface {
	// Create создаёт новый заказ
	Create(ctx context.Context, input *dto.CreateOrderInput) error

	// GetByUID получает заказ по UID
	GetByUID(ctx context.Context, orderUID string) (*dto.OrderOutput, error)

	// GetAll получает все заказы
	GetAll(ctx context.Context) ([]*dto.OrderOutput, error)

	// RestoreCache восстанавливает кэш из БД при старте приложения
	RestoreCache(ctx context.Context) error

	// GetCacheStats возвращает статистику кэша
	GetCacheStats() map[string]interface{}
}
