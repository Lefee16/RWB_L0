package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"RWB_L0/internal/domain"
	"RWB_L0/internal/dto"
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

func TestOrderUseCase_GetByUID(t *testing.T) {
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
	output, err := uc.GetByUID(context.Background(), "test123")
	if err != nil {
		t.Fatalf("GetByUID() error = %v", err)
	}

	if output.OrderUID != "test123" {
		t.Errorf("Expected OrderUID = test123, got %s", output.OrderUID)
	}

	// Проверяем, что заказ попал в кэш
	if cache.Count() != 1 {
		t.Errorf("Expected 1 order in cache, got %d", cache.Count())
	}
}

func TestOrderUseCase_GetByUID_EmptyUID(t *testing.T) {
	repo := NewMockRepository()
	cache := NewMockCache()
	uc := NewOrderUseCase(repo, cache)

	// Получаем с пустым UID
	_, err := uc.GetByUID(context.Background(), "")

	if err == nil {
		t.Error("Expected error for empty UID, got nil")
	}

	if !errors.Is(err, domain.ErrEmptyOrderUID) {
		t.Errorf("Expected ErrEmptyOrderUID, got %v", err)
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

func TestOrderUseCase_GetByUID_FromCache(t *testing.T) {
	repo := NewMockRepository()
	cache := NewMockCache()
	uc := NewOrderUseCase(repo, cache)

	// Создаём заказ и сохраняем только в кэш
	order, _ := domain.NewOrder("cached123", "TRACK456", "WBIL")
	_ = cache.Set("cached123", order)

	// Получаем заказ (должен взять из кэша, не трогая БД)
	output, err := uc.GetByUID(context.Background(), "cached123")
	if err != nil {
		t.Fatalf("GetByUID() error = %v", err)
	}

	if output.OrderUID != "cached123" {
		t.Errorf("Expected OrderUID = cached123, got %s", output.OrderUID)
	}

	// Проверяем, что в БД заказа нет (брали из кэша)
	if len(repo.orders) != 0 {
		t.Error("Order should not be in repository")
	}
}

func TestOrderUseCase_GetAll(t *testing.T) {
	repo := NewMockRepository()
	cache := NewMockCache()
	uc := NewOrderUseCase(repo, cache)

	// Добавляем заказы
	order1, _ := domain.NewOrder("order1", "TRACK1", "WBIL")
	order2, _ := domain.NewOrder("order2", "TRACK2", "WBIL")
	order3, _ := domain.NewOrder("order3", "TRACK3", "WBIL")

	_ = repo.Save(context.Background(), order1)
	_ = repo.Save(context.Background(), order2)
	_ = repo.Save(context.Background(), order3)

	// Получаем все заказы
	orders, err := uc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	if len(orders) != 3 {
		t.Errorf("Expected 3 orders, got %d", len(orders))
	}
}

func TestOrderUseCase_Create(t *testing.T) {
	repo := NewMockRepository()
	cache := NewMockCache()
	uc := NewOrderUseCase(repo, cache)

	// Создаём входной DTO
	input := &dto.CreateOrderInput{
		OrderUID:    "new-order-123",
		TrackNumber: "TRACK-NEW",
		Entry:       "WBIL",
		Delivery: dto.DeliveryInput{
			Name:    "Test User",
			Phone:   "+79001234567",
			Zip:     "123456",
			City:    "Moscow",
			Address: "Test Street",
			Region:  "Moscow",
			Email:   "test@test.com",
		},
		Payment: dto.PaymentInput{
			Transaction:  "TX-NEW",
			Currency:     "USD",
			Provider:     "PayPal",
			Amount:       1000,
			Bank:         "Sberbank",
			DeliveryCost: 100,
			GoodsTotal:   900,
		},
		Items: []dto.ItemInput{
			{
				ChrtID:      12345,
				TrackNumber: "TRACK-NEW",
				Price:       500,
				Name:        "Test Item",
				Sale:        10,
				Size:        "M",
				TotalPrice:  450,
				NmID:        54321,
				Brand:       "Test Brand",
				Status:      202,
			},
		},
		Locale:          "en",
		CustomerID:      "customer-1",
		DeliveryService: "DHL",
		Shardkey:        "shard-1",
		SmID:            999,
	}

	// Создаём заказ
	err := uc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Проверяем, что заказ в репозитории
	if len(repo.orders) != 1 {
		t.Errorf("Expected 1 order in repo, got %d", len(repo.orders))
	}

	// Проверяем, что заказ в кэше
	if cache.Count() != 1 {
		t.Errorf("Expected 1 order in cache, got %d", cache.Count())
	}
}

func TestOrderUseCase_GetCacheStats(t *testing.T) {
	repo := NewMockRepository()
	cache := NewMockCache()
	uc := NewOrderUseCase(repo, cache)

	// Добавляем заказы в кэш
	order1, _ := domain.NewOrder("order1", "TRACK1", "WBIL")
	order2, _ := domain.NewOrder("order2", "TRACK2", "WBIL")

	_ = cache.Set("order1", order1)
	_ = cache.Set("order2", order2)

	// Получаем статистику
	stats := uc.GetCacheStats()

	if stats["cached_orders"] != 2 {
		t.Errorf("Expected cached_orders = 2, got %v", stats["cached_orders"])
	}
}
