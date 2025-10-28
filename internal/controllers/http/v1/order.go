package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"RWB_L0/internal/domain"
	"RWB_L0/internal/usecase"

	"github.com/go-chi/chi/v5"
)

// OrderHandler обрабатывает HTTP запросы для работы с заказами
type OrderHandler struct {
	orderUseCase usecase.OrderUseCaseInterface
}

// NewOrderHandler создаёт новый экземпляр OrderHandler
func NewOrderHandler(orderUseCase usecase.OrderUseCaseInterface) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

// GetByUID обрабатывает GET /api/v1/orders/:uid
func (h *OrderHandler) GetByUID(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "uid")

	order, err := h.orderUseCase.GetByUID(r.Context(), orderUID)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyOrderUID) {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Error: "Order UID is required",
			})
			return
		}
		// Любая другая ошибка - not found
		writeJSON(w, http.StatusNotFound, ErrorResponse{
			Error: "Order not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, order)
}

// HealthCheck обрабатывает GET /api/v1/health
func (h *OrderHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	stats := h.orderUseCase.GetCacheStats()

	writeJSON(w, http.StatusOK, HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Cache:     stats,
	})
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthResponse представляет ответ health check
type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Cache     map[string]interface{} `json:"cache"`
}

// writeJSON отправляет JSON ответ
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
