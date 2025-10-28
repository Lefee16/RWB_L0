package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"RWB_L0/internal/domain"
)

// MockRepository - мок репозитория для тестов
type MockRepository struct {
	orders map[string]*domain.Order
	err    error
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		orders: make(map[string]*domain.Order),
	}
}

func (m *MockRepository) Save(_ context.Context, order *domain.Order) error {
	if m.err != nil {
		return m.err
	}
	m.orders[order.OrderUID] = order
	return nil
}

func (m *MockRepository) GetByID(_ context.Context, orderUID string) (*domain.Order, error) {
	if m.err != nil {
		return nil, m.err
	}
	order, exists := m.orders[orderUID]
	if !exists {
		return nil, domain.ErrEmptyOrderUID
	}
	return order, nil
}

func (m *MockRepository) GetAll(_ context.Context) ([]*domain.Order, error) {
	if m.err != nil {
		return nil, m.err
	}
	orders := make([]*domain.Order, 0, len(m.orders))
	for _, order := range m.orders {
		orders = append(orders, order)
	}
	return orders, nil
}

func (m *MockRepository) Delete(_ context.Context, orderUID string) error {
	delete(m.orders, orderUID)
	return nil
}

func (m *MockRepository) Count(_ context.Context) (int, error) {
	return len(m.orders), nil
}

// MockCache - мок кэша для тестов
type MockCache struct {
	orders map[string]*domain.Order
	err    error
}

func NewMockCache() *MockCache {
	return &MockCache{
		orders: make(map[string]*domain.Order),
	}
}

func (m *MockCache) Set(orderUID string, order *domain.Order) error {
	if m.err != nil {
		return m.err
	}
	m.orders[orderUID] = order
	return nil
}

func (m *MockCache) Get(orderUID string) (*domain.Order, error) {
	if m.err != nil {
		return nil, m.err
	}
	order, exists := m.orders[orderUID]
	if !exists {
		return nil, errors.New("not found")
	}
	return order, nil
}

func (m *MockCache) Delete(orderUID string) error {
	delete(m.orders, orderUID)
	return nil
}

func (m *MockCache) LoadAll(orders []*domain.Order) error {
	for _, order := range orders {
		m.orders[order.OrderUID] = order
	}
	return nil
}

func (m *MockCache) GetAll() ([]*domain.Order, error) {
	orders := make([]*domain.Order, 0, len(m.orders))
	for _, order := range m.orders {
		orders = append(orders, order)
	}
	return orders, nil
}

func (m *MockCache) Count() int {
	return len(m.orders)
}

func (m *MockCache) Clear() error {
	m.orders = make(map[string]*domain.Order)
	return nil
}

// Тесты

func TestOrderUseCase_GetOrderByID(t *testing.T) {
	repo := NewMockRepository()
	cache := NewMockCache()
	uc := NewOrderUseCase(repo, cache)

	// Создаём тестовый заказ
	order, _ := domain.NewOrder("test123", "TRACK123", "WBIL")
	order.DateCreated = time.Now()

	delivery, _ := domain.NewDelivery("John Doe", "+79001234567", "123456", "Moscow", "Red Square", "Moscow Region", "test@test.com")
	order.Delivery = *delivery

	payment, _ := domain.NewPayment("TX123", "REQ123", "USD", "PayPal", "Sberbank", 1000, time.Now().Unix(), 100, 900, 0)
	order.Payment = *payment

	// Сохраняем в репозиторий
	_ = repo.Save(context.Background(), order)

	// Получаем через use case
	output, err := uc.GetOrderByID(context.Background(), "test123")
	if err != nil {
		t.Fatalf("GetOrderByID() error = %v", err)
	}

	if output.OrderUID != "test123" {
		t.Errorf("Expected OrderUID = test123, got %s", output.OrderUID)
	}

	// Проверяем, что заказ попал в кэш
	if cache.Count() != 1 {
		t.Errorf("Expected 1 order in cache, got %d", cache.Count())
	}
}

func TestOrderUseCase_RestoreCache(t *testing.T) {
	repo := NewMockRepository()
	cache := NewMockCache()
	uc := NewOrderUseCase(repo, cache)

	// Добавляем заказы в репозиторий
	order1, _ := domain.NewOrder("order1", "TRACK1", "WBIL")
	order2, _ := domain.NewOrder("order2", "TRACK2", "WBIL")

	_ = repo.Save(context.Background(), order1)
	_ = repo.Save(context.Background(), order2)

	// Восстанавливаем кэш
	err := uc.RestoreCache(context.Background())
	if err != nil {
		t.Fatalf("RestoreCache() error = %v", err)
	}

	// Проверяем, что оба заказа в кэше
	if cache.Count() != 2 {
		t.Errorf("Expected 2 orders in cache, got %d", cache.Count())
	}
}

func TestOrderUseCase_GetOrderByID_FromCache(t *testing.T) {
	repo := NewMockRepository()
	cache := NewMockCache()
	uc := NewOrderUseCase(repo, cache)

	// Создаём заказ и сохраняем только в кэш
	order, _ := domain.NewOrder("cached123", "TRACK456", "WBIL")
	_ = cache.Set("cached123", order)

	// Получаем заказ (должен взять из кэша, не трогая БД)
	output, err := uc.GetOrderByID(context.Background(), "cached123")
	if err != nil {
		t.Fatalf("GetOrderByID() error = %v", err)
	}

	if output.OrderUID != "cached123" {
		t.Errorf("Expected OrderUID = cached123, got %s", output.OrderUID)
	}

	// Проверяем, что в БД заказа нет (брали из кэша)
	if len(repo.orders) != 0 {
		t.Error("Order should not be in repository")
	}
}
