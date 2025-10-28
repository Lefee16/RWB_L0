package v1

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"RWB_L0/internal/dto"
	"RWB_L0/internal/usecase"
)

type OrderHandler struct {
	orderUseCase *usecase.OrderUseCase
}

func NewOrderHandler(orderUseCase *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{orderUseCase: orderUseCase}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "order_uid")

	if orderUID == "" {
		respondWithError(w, http.StatusBadRequest, "order_uid is required")
		return
	}

	order, err := h.orderUseCase.GetOrderByID(r.Context(), orderUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	respondWithJSON(w, http.StatusOK, dto.GetOrderResponse{Order: order})
}

func (h *OrderHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	stats := h.orderUseCase.GetCacheStats()

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"cache":     stats,
	})
}
