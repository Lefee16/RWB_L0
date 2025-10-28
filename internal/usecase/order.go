package usecase

import (
	"context"
	"fmt"

	"RWB_L0/internal/domain"
	"RWB_L0/internal/dto"
)

// OrderUseCase - основная бизнес-логика работы с заказами
type OrderUseCase struct {
	repo  OrderRepository
	cache Cache
}

// NewOrderUseCase - создание use case
func NewOrderUseCase(repo OrderRepository, cache Cache) *OrderUseCase {
	return &OrderUseCase{
		repo:  repo,
		cache: cache,
	}
}

// CreateOrder - создать новый заказ
func (uc *OrderUseCase) CreateOrder(ctx context.Context, input *dto.CreateOrderInput) error {
	// 1. Конвертируем DTO → Domain (с валидацией!)
	order, err := input.ToDomain()
	if err != nil {
		return fmt.Errorf("invalid order data: %w", err)
	}

	// 2. Дополнительная валидация через Domain метод
	if err := order.Validate(); err != nil {
		return fmt.Errorf("order validation failed: %w", err)
	}

	// 3. Сохраняем в БД
	if err := uc.repo.Save(ctx, order); err != nil {
		return fmt.Errorf("failed to save order to database: %w", err)
	}

	// 4. Добавляем в кэш
	if err := uc.cache.Set(order.OrderUID, order); err != nil {
		// Логируем ошибку, но не падаем - БД важнее кэша
		return fmt.Errorf("failed to cache order: %w", err)
	}

	return nil
}

// GetOrderByID - получить заказ по ID
func (uc *OrderUseCase) GetOrderByID(ctx context.Context, orderUID string) (*dto.OrderOutput, error) {
	if orderUID == "" {
		return nil, domain.ErrEmptyOrderUID
	}

	// 1. Пытаемся получить из кэша (быстро!)
	order, err := uc.cache.Get(orderUID)
	if err == nil {
		// Нашли в кэше - конвертируем в DTO и возвращаем
		return dto.FromDomain(order), nil
	}

	// 2. Не нашли в кэше - идём в БД
	order, err = uc.repo.GetByID(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// 3. Добавляем в кэш для следующего раза
	_ = uc.cache.Set(order.OrderUID, order)

	// 4. Конвертируем Domain → DTO
	return dto.FromDomain(order), nil
}

// GetAllOrders - получить все заказы
func (uc *OrderUseCase) GetAllOrders(ctx context.Context) ([]*dto.OrderOutput, error) {
	// 1. Пытаемся получить из кэша
	orders, err := uc.cache.GetAll()
	if err == nil && len(orders) > 0 {
		// Есть в кэше - конвертируем и возвращаем
		return uc.convertOrdersToDTO(orders), nil
	}

	// 2. Нет в кэше - получаем из БД
	orders, err = uc.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	// 3. Конвертируем в DTO
	return uc.convertOrdersToDTO(orders), nil
}

// RestoreCache - восстановить кэш из БД (при старте сервиса)
func (uc *OrderUseCase) RestoreCache(ctx context.Context) error {
	// 1. Получаем все заказы из БД
	orders, err := uc.repo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to load orders from database: %w", err)
	}

	// 2. Загружаем в кэш
	if err := uc.cache.LoadAll(orders); err != nil {
		return fmt.Errorf("failed to load orders into cache: %w", err)
	}

	return nil
}

// GetCacheStats - получить статистику кэша
func (uc *OrderUseCase) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"cached_orders": uc.cache.Count(),
	}
}

// convertOrdersToDTO - конвертировать массив Domain → DTO
func (uc *OrderUseCase) convertOrdersToDTO(orders []*domain.Order) []*dto.OrderOutput {
	result := make([]*dto.OrderOutput, len(orders))
	for i, order := range orders {
		result[i] = dto.FromDomain(order)
	}
	return result
}
