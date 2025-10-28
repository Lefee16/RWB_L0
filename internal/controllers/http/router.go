package http

import (
	"RWB_L0/internal/controllers/http/v1"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Router - HTTP роутер
type Router struct {
	mux *chi.Mux
}

// NewRouter - создание роутера
func NewRouter(orderHandler *v1.OrderHandler, webHandler *v1.WebHandler, mw *Middleware) *Router {
	r := chi.NewRouter()

	// Глобальные middleware
	r.Use(mw.Logger)
	r.Use(mw.Recoverer)
	r.Use(mw.Timeout)

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/orders/{order_uid}", orderHandler.GetOrder)
		r.Get("/health", orderHandler.HealthCheck)
	})

	// Web интерфейс
	r.Get("/", webHandler.IndexPage)
	r.Get("/orders/{order_uid}", webHandler.OrderPage)

	// Статика
	fileServer := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return &Router{mux: r}
}

// Handler - получить http.Handler
func (rt *Router) Handler() http.Handler {
	return rt.mux
}
