package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

// TxFunc - функция, выполняемая в транзакции
type TxFunc func(*sql.Tx) error

// WithTransaction - выполнить функцию в транзакции с автоматическим commit/rollback
func WithTransaction(ctx context.Context, db *sql.DB, fn TxFunc) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Если функция вернёт ошибку - откатываем транзакцию
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Выполняем функцию
	err = fn(tx)
	if err != nil {
		return err
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
