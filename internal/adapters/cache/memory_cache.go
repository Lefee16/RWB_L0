package cache

import (
	"fmt"
	"sync"

	"RWB_L0/internal/domain"
)

// MemoryCache - простой in-memory кэш с sync.RWMutex
type MemoryCache struct {
	mu     sync.RWMutex
	orders map[string]*domain.Order
}

// NewMemoryCache - создание нового кэша
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		orders: make(map[string]*domain.Order),
	}
}

// Set - добавить заказ в кэш
func (c *MemoryCache) Set(orderUID string, order *domain.Order) error {
	if orderUID == "" {
		return fmt.Errorf("order_uid cannot be empty")
	}
	if order == nil {
		return fmt.Errorf("order cannot be nil")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.orders[orderUID] = order
	return nil
}

// Get - получить заказ из кэша
func (c *MemoryCache) Get(orderUID string) (*domain.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	order, exists := c.orders[orderUID]
	if !exists {
		return nil, fmt.Errorf("order not found in cache")
	}

	return order, nil
}

// Delete - удалить заказ из кэша
func (c *MemoryCache) Delete(orderUID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.orders, orderUID)
	return nil
}

// LoadAll - загрузить все заказы в кэш (для восстановления из БД)
func (c *MemoryCache) LoadAll(orders []*domain.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, order := range orders {
		c.orders[order.OrderUID] = order
	}

	return nil
}

// GetAll - получить все заказы из кэша
func (c *MemoryCache) GetAll() ([]*domain.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	orders := make([]*domain.Order, 0, len(c.orders))
	for _, order := range c.orders {
		orders = append(orders, order)
	}

	return orders, nil
}

// Count - получить количество заказов в кэше
func (c *MemoryCache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.orders)
}

// Clear - очистить весь кэш
func (c *MemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.orders = make(map[string]*domain.Order)
	return nil
}
