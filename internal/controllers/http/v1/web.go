package v1

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"

	"RWB_L0/internal/usecase"
)

type WebHandler struct {
	orderUseCase *usecase.OrderUseCase
	templates    *template.Template
}

func NewWebHandler(orderUseCase *usecase.OrderUseCase) *WebHandler {
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))

	return &WebHandler{
		orderUseCase: orderUseCase,
		templates:    tmpl,
	}
}

func (h *WebHandler) IndexPage(w http.ResponseWriter, r *http.Request) {
	if err := h.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *WebHandler) OrderPage(w http.ResponseWriter, r *http.Request) {
	orderUID := chi.URLParam(r, "order_uid")

	order, err := h.orderUseCase.GetOrderByID(r.Context(), orderUID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "order.html", order); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
