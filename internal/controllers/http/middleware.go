package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return middleware.Logger(next)
}

func (m *Middleware) Recoverer(next http.Handler) http.Handler {
	return middleware.Recoverer(next)
}

func (m *Middleware) Timeout(next http.Handler) http.Handler {
	return middleware.Timeout(30 * time.Second)(next)
}
