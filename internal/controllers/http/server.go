package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Server - HTTP сервер
type Server struct {
	server *http.Server
}

// NewServer - создание HTTP сервера
func NewServer(host string, port int, router *Router) *Server {
	return &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", host, port),
			Handler:      router.Handler(),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Start - запуск сервера
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown - graceful shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// GetAddress - получить адрес сервера
func (s *Server) GetAddress() string {
	return s.server.Addr
}
