package v1

import (
	"errors"
	"html/template"
	"net/http"

	"RWB_L0/internal/domain"
	"RWB_L0/internal/usecase"

	"github.com/go-chi/chi/v5"
)

// WebHandler обрабатывает HTTP запросы для веб-интерфейса
type WebHandler struct {
	orderUseCase usecase.OrderUseCaseInterface
	templates    *template.Template
}

// NewWebHandler создаёт новый экземпляр WebHandler
func NewWebHandler(orderUseCase usecase.OrderUseCaseInterface) *WebHandler {
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))

	return &WebHandler{
		orderUseCase: orderUseCase,
		templates:    tmpl,
	}
}

// IndexPage обрабатывает GET / - главная страница
func (h *WebHandler) IndexPage(w http.ResponseWriter, r *http.Request) {
	if err := h.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// OrderPage обрабатывает GET /orders/:order_uid - страница заказа
func (h *WebHandler) OrderPage(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "order_uid")

	// Используем GetByUID вместо GetOrderByID
	order, err := h.orderUseCase.GetByUID(r.Context(), orderUID)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyOrderUID) {
			http.Error(w, "Order UID is required", http.StatusBadRequest)
			return
		}
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "order.html", order); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
