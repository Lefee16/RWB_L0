package usecase

import (
	"context"
	"fmt"

	"RWB_L0/internal/domain"
	"RWB_L0/internal/dto"
)

// Проверка на этапе компиляции, что OrderUseCase реализует OrderUseCaseInterface
var _ OrderUseCaseInterface = (*OrderUseCase)(nil)

// OrderUseCase реализует бизнес-логику работы с заказами
type OrderUseCase struct {
	repo  OrderRepository
	cache Cache
}

// NewOrderUseCase создаёт новый экземпляр OrderUseCase
func NewOrderUseCase(repo OrderRepository, cache Cache) *OrderUseCase {
	return &OrderUseCase{
		repo:  repo,
		cache: cache,
	}
}

// Create создаёт новый заказ
func (uc *OrderUseCase) Create(ctx context.Context, input *dto.CreateOrderInput) error {
	// Конвертируем DTO в доменную модель
	order, err := input.ToDomain()
	if err != nil {
		return fmt.Errorf("failed to convert DTO: %w", err)
	}

	// Валидация доменной модели
	if err := order.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Сохраняем в БД
	if err := uc.repo.Save(ctx, order); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// Сохраняем в кэш (игнорируем ошибку кэша)
	if err := uc.cache.Set(order.OrderUID, order); err != nil {
		_ = err // Логируем, но не возвращаем ошибку
	}

	return nil
}

// GetByUID получает заказ по UID
func (uc *OrderUseCase) GetByUID(ctx context.Context, orderUID string) (*dto.OrderOutput, error) {
	// Проверяем пустой UID
	if orderUID == "" {
		return nil, domain.ErrEmptyOrderUID
	}

	// Пытаемся получить из кэша
	cached, err := uc.cache.Get(orderUID)
	if err == nil && cached != nil {
		return dto.FromDomain(cached), nil
	}

	// Если не в кэше - идём в БД
	order, err := uc.repo.GetByID(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Кэшируем для следующего раза
	_ = uc.cache.Set(orderUID, order)

	return dto.FromDomain(order), nil
}

// GetAll получает все заказы
func (uc *OrderUseCase) GetAll(ctx context.Context) ([]*dto.OrderOutput, error) {
	orders, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders: %w", err)
	}

	// Конвертируем в DTO
	result := make([]*dto.OrderOutput, 0, len(orders))
	for _, order := range orders {
		result = append(result, dto.FromDomain(order))
	}

	return result, nil
}

// RestoreCache восстанавливает кэш из БД при старте приложения
func (uc *OrderUseCase) RestoreCache(ctx context.Context) error {
	orders, err := uc.repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to restore cache: %w", err)
	}

	// Используем LoadAll из Cache интерфейса
	if err := uc.cache.LoadAll(orders); err != nil {
		return fmt.Errorf("failed to load cache: %w", err)
	}

	return nil
}

// GetCacheStats возвращает статистику кэша
func (uc *OrderUseCase) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"cached_orders": uc.cache.Count(),
	}
}
