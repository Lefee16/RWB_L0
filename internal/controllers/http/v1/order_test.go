package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"RWB_L0/internal/domain"
	"RWB_L0/internal/dto"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderUseCase - мок для OrderUseCaseInterface
type MockOrderUseCase struct {
	mock.Mock
}

func (m *MockOrderUseCase) Create(ctx context.Context, input *dto.CreateOrderInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockOrderUseCase) GetByUID(ctx context.Context, orderUID string) (*dto.OrderOutput, error) {
	args := m.Called(ctx, orderUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.OrderOutput), args.Error(1)
}

func (m *MockOrderUseCase) GetAll(ctx context.Context) ([]*dto.OrderOutput, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.OrderOutput), args.Error(1)
}

func (m *MockOrderUseCase) RestoreCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockOrderUseCase) GetCacheStats() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

// TestOrderHandler_GetByUID_Success тестирует успешное получение заказа
func TestOrderHandler_GetByUID_Success(t *testing.T) {
	// Arrange
	mockUseCase := new(MockOrderUseCase)
	handler := NewOrderHandler(mockUseCase)

	expectedOrder := &dto.OrderOutput{
		OrderUID:    "test-uid-123",
		TrackNumber: "TRACK123",
		Entry:       "WBIL",
	}

	mockUseCase.On("GetByUID", mock.Anything, "test-uid-123").Return(expectedOrder, nil)

	// Создаём HTTP запрос
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/test-uid-123", nil)
	w := httptest.NewRecorder()

	// Добавляем URL параметр через chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("uid", "test-uid-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handler.GetByUID(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.OrderOutput
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "test-uid-123", response.OrderUID)
	assert.Equal(t, "TRACK123", response.TrackNumber)

	mockUseCase.AssertExpectations(t)
}

// TestOrderHandler_GetByUID_EmptyUID тестирует запрос с пустым UID
func TestOrderHandler_GetByUID_EmptyUID(t *testing.T) {
	// Arrange
	mockUseCase := new(MockOrderUseCase)
	handler := NewOrderHandler(mockUseCase)

	mockUseCase.On("GetByUID", mock.Anything, "").Return(nil, domain.ErrEmptyOrderUID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("uid", "")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handler.GetByUID(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "required")

	mockUseCase.AssertExpectations(t)
}

// TestOrderHandler_GetByUID_NotFound тестирует запрос несуществующего заказа
func TestOrderHandler_GetByUID_NotFound(t *testing.T) {
	// Arrange
	mockUseCase := new(MockOrderUseCase)
	handler := NewOrderHandler(mockUseCase)

	mockUseCase.On("GetByUID", mock.Anything, "nonexistent").Return(nil, errors.New("order not found"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/nonexistent", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("uid", "nonexistent")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handler.GetByUID(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "not found")

	mockUseCase.AssertExpectations(t)
}

// TestOrderHandler_HealthCheck тестирует health check endpoint
func TestOrderHandler_HealthCheck(t *testing.T) {
	// Arrange
	mockUseCase := new(MockOrderUseCase)
	handler := NewOrderHandler(mockUseCase)

	stats := map[string]interface{}{
		"cached_orders": 42,
	}
	mockUseCase.On("GetCacheStats").Return(stats)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	// Act
	handler.HealthCheck(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, float64(42), response.Cache["cached_orders"])
	assert.NotZero(t, response.Timestamp)

	mockUseCase.AssertExpectations(t)
}
