package cache

import (
	"errors"
	"sync"
	"time"

	"RWB_L0/internal/domain"
)

// cacheEntry представляет запись в кэше с временем истечения
type cacheEntry struct {
	order      *domain.Order
	expiresAt  time.Time
	lastAccess time.Time // Для LRU
}

// MemoryCache - in-memory кэш с TTL и лимитами
type MemoryCache struct {
	mu       sync.RWMutex
	data     map[string]*cacheEntry
	maxSize  int
	ttl      time.Duration
	stopChan chan struct{}
}

// NewMemoryCacheWithConfig создаёт кэш с кастомными настройками
func NewMemoryCacheWithConfig(maxSize int, ttl time.Duration) *MemoryCache {
	cache := &MemoryCache{
		data:     make(map[string]*cacheEntry),
		maxSize:  maxSize,
		ttl:      ttl,
		stopChan: make(chan struct{}),
	}

	// Запускаем фоновую очистку устаревших записей каждые 5 минут
	go cache.cleanupExpired()

	return cache
}

// Set добавляет заказ в кэш
func (c *MemoryCache) Set(orderUID string, order *domain.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Проверяем лимит размера
	if len(c.data) >= c.maxSize {
		// Если достигнут лимит, удаляем самую старую запись (LRU)
		c.evictOldest()
	}

	// Добавляем запись с TTL
	c.data[orderUID] = &cacheEntry{
		order:      order,
		expiresAt:  time.Now().Add(c.ttl),
		lastAccess: time.Now(),
	}

	return nil
}

// Get получает заказ из кэша
func (c *MemoryCache) Get(orderUID string) (*domain.Order, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.data[orderUID]
	if !exists {
		return nil, errors.New("order not found in cache")
	}

	// Проверяем, не истёк ли TTL
	if time.Now().After(entry.expiresAt) {
		// Удаляем устаревшую запись
		delete(c.data, orderUID)
		return nil, errors.New("order expired in cache")
	}

	// Обновляем время последнего доступа (для LRU)
	entry.lastAccess = time.Now()

	return entry.order, nil
}

// Delete удаляет заказ из кэша
func (c *MemoryCache) Delete(orderUID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.data[orderUID]; !exists {
		return errors.New("order not found in cache")
	}

	delete(c.data, orderUID)
	return nil
}

// LoadAll загружает все заказы в кэш
func (c *MemoryCache) LoadAll(orders []*domain.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, order := range orders {
		// Проверяем лимит при загрузке
		if len(c.data) >= c.maxSize {
			break
		}

		c.data[order.OrderUID] = &cacheEntry{
			order:      order,
			expiresAt:  time.Now().Add(c.ttl),
			lastAccess: time.Now(),
		}
	}

	return nil
}

// GetAll возвращает все заказы из кэша
func (c *MemoryCache) GetAll() ([]*domain.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	orders := make([]*domain.Order, 0, len(c.data))
	now := time.Now()

	for _, entry := range c.data {
		// Пропускаем устаревшие записи
		if now.After(entry.expiresAt) {
			continue
		}
		orders = append(orders, entry.order)
	}

	return orders, nil
}

// Count возвращает количество записей в кэше
func (c *MemoryCache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Считаем только не устаревшие записи
	count := 0
	now := time.Now()

	for _, entry := range c.data {
		if now.Before(entry.expiresAt) {
			count++
		}
	}

	return count
}

// Clear очищает весь кэш
func (c *MemoryCache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*cacheEntry)
	return nil
}

// evictOldest удаляет самую старую запись (LRU)
func (c *MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, entry := range c.data {
		if first || entry.lastAccess.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.lastAccess
			first = false
		}
	}

	if oldestKey != "" {
		delete(c.data, oldestKey)
	}
}

// cleanupExpired периодически удаляет устаревшие записи
func (c *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.stopChan:
			return
		}
	}
}

// removeExpired удаляет все устаревшие записи
func (c *MemoryCache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.data {
		if now.After(entry.expiresAt) {
			delete(c.data, key)
		}
	}
}

// countValid подсчитывает количество действительных (не устаревших) записей
func (c *MemoryCache) countValid() int {
	count := 0
	now := time.Now()

	for _, entry := range c.data {
		if now.Before(entry.expiresAt) {
			count++
		}
	}

	return count
}

// Close закрывает кэш и останавливает фоновую очистку
func (c *MemoryCache) Close() {
	close(c.stopChan)
}
