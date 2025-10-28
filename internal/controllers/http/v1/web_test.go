package v1

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"RWB_L0/internal/domain"
	"RWB_L0/internal/dto"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWebHandler_IndexPage(t *testing.T) {
	// Arrange
	mockUseCase := new(MockOrderUseCase)

	// Создаём простой mock template
	tmpl := template.Must(template.New("index.html").Parse("<html>Index Page</html>"))

	handler := &WebHandler{
		orderUseCase: mockUseCase,
		templates:    tmpl,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Act
	handler.IndexPage(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWebHandler_OrderPage_Success(t *testing.T) {
	// Arrange
	mockUseCase := new(MockOrderUseCase)

	// Создаём mock template
	tmpl := template.Must(template.New("order.html").Parse("<html>Order: {{.OrderUID}}</html>"))

	handler := &WebHandler{
		orderUseCase: mockUseCase,
		templates:    tmpl,
	}

	expectedOrder := &dto.OrderOutput{
		OrderUID:    "test-uid",
		TrackNumber: "TRACK123",
		Entry:       "WBIL",
	}

	mockUseCase.On("GetByUID", mock.Anything, "test-uid").Return(expectedOrder, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders/test-uid", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("order_uid", "test-uid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handler.OrderPage(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test-uid")

	mockUseCase.AssertExpectations(t)
}

func TestWebHandler_OrderPage_EmptyUID(t *testing.T) {
	// Arrange
	mockUseCase := new(MockOrderUseCase)
	tmpl := template.Must(template.New("order.html").Parse("<html>Order</html>"))

	handler := &WebHandler{
		orderUseCase: mockUseCase,
		templates:    tmpl,
	}

	mockUseCase.On("GetByUID", mock.Anything, "").Return(nil, domain.ErrEmptyOrderUID)

	req := httptest.NewRequest(http.MethodGet, "/orders/", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("order_uid", "")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handler.OrderPage(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockUseCase.AssertExpectations(t)
}

func TestWebHandler_OrderPage_NotFound(t *testing.T) {
	// Arrange
	mockUseCase := new(MockOrderUseCase)
	tmpl := template.Must(template.New("order.html").Parse("<html>Order</html>"))

	handler := &WebHandler{
		orderUseCase: mockUseCase,
		templates:    tmpl,
	}

	mockUseCase.On("GetByUID", mock.Anything, "nonexistent").Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/orders/nonexistent", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("order_uid", "nonexistent")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handler.OrderPage(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockUseCase.AssertExpectations(t)
}
