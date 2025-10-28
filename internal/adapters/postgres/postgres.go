package postgres

import (
	"database/sql"
	"fmt"
)

// Repository - базовая структура для работы с PostgreSQL
type Repository struct {
	db *sql.DB
}

// New - создание репозитория
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetDB - получить подключение к БД
func (r *Repository) GetDB() *sql.DB {
	return r.db
}

// Close - закрыть подключение
func (r *Repository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// Ping - проверка подключения
func (r *Repository) Ping() error {
	if err := r.db.Ping(); err != nil {
		return fmt.Errorf("postgres ping failed: %w", err)
	}
	return nil
}
