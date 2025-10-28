package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"RWB_L0/internal/domain"
)

// OrderRepository - репозиторий для работы с заказами (3NF)
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository - создание репозитория заказов
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Save - сохранить заказ в БД (нормализованная структура - 4 таблицы)
func (r *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	// Начинаем транзакцию
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
		}
	}(tx)

	// 1. Сохраняем основную информацию о заказе
	queryOrder := `
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO UPDATE SET
			track_number = EXCLUDED.track_number,
			entry = EXCLUDED.entry,
			locale = EXCLUDED.locale
	`
	_, err = tx.ExecContext(ctx, queryOrder,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	// 2. Сохраняем информацию о доставке
	queryDelivery := `
		INSERT INTO deliveries (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (order_uid) DO UPDATE SET
			name = EXCLUDED.name,
			phone = EXCLUDED.phone,
			zip = EXCLUDED.zip,
			city = EXCLUDED.city,
			address = EXCLUDED.address,
			region = EXCLUDED.region,
			email = EXCLUDED.email
	`
	_, err = tx.ExecContext(ctx, queryDelivery,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to save delivery: %w", err)
	}

	// 3. Сохраняем информацию о платеже
	queryPayment := `
		INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO UPDATE SET
			transaction = EXCLUDED.transaction,
			amount = EXCLUDED.amount,
			payment_dt = EXCLUDED.payment_dt
	`
	_, err = tx.ExecContext(ctx, queryPayment,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("failed to save payment: %w", err)
	}

	// 4. Удаляем старые товары (для обновления)
	_, err = tx.ExecContext(ctx, "DELETE FROM items WHERE order_uid = $1", order.OrderUID)
	if err != nil {
		return fmt.Errorf("failed to delete old items: %w", err)
	}

	// 5. Сохраняем товары
	queryItem := `
		INSERT INTO items (
			order_uid, chrt_id, track_number, price, name, sale,
			size, total_price, nm_id, brand, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, queryItem,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price,
			item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to save item: %w", err)
		}
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByID - получить заказ по ID (JOIN 4 таблиц)
func (r *OrderRepository) GetByID(ctx context.Context, orderUID string) (*domain.Order, error) {
	// Получаем основную информацию о заказе
	queryOrder := `
		SELECT 
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
		WHERE order_uid = $1
	`

	var order domain.Order
	err := r.db.QueryRowContext(ctx, queryOrder, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Получаем информацию о доставке
	queryDelivery := `
		SELECT name, phone, zip, city, address, region, email
		FROM deliveries
		WHERE order_uid = $1
	`
	err = r.db.QueryRowContext(ctx, queryDelivery, orderUID).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}

	// Получаем информацию о платеже
	queryPayment := `
		SELECT transaction, request_id, currency, provider, amount, payment_dt,
		       bank, delivery_cost, goods_total, custom_fee
		FROM payments
		WHERE order_uid = $1
	`
	err = r.db.QueryRowContext(ctx, queryPayment, orderUID).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Получаем товары
	queryItems := `
		SELECT chrt_id, track_number, price, name, sale, size,
		       total_price, nm_id, brand, status
		FROM items
		WHERE order_uid = $1
	`
	rows, err := r.db.QueryContext(ctx, queryItems, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	order.Items = make([]domain.Item, 0)
	for rows.Next() {
		var item domain.Item
		err = rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID,
			&item.Brand, &item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("items rows error: %w", err)
	}

	return &order, nil
}

// GetAll - получить все заказы
func (r *OrderRepository) GetAll(ctx context.Context) ([]*domain.Order, error) {
	query := `
		SELECT order_uid, track_number, entry, locale, internal_signature,
		       customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
		ORDER BY date_created DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var orders []*domain.Order

	for rows.Next() {
		var order domain.Order
		err = rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		// Получаем связанные данные для каждого заказа
		fullOrder, err := r.GetByID(ctx, order.OrderUID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, fullOrder)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return orders, nil
}

// Delete - удалить заказ (каскадное удаление из всех связанных таблиц)
func (r *OrderRepository) Delete(ctx context.Context, orderUID string) error {
	query := `DELETE FROM orders WHERE order_uid = $1`

	result, err := r.db.ExecContext(ctx, query, orderUID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}

// Count - получить количество заказов
func (r *OrderRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM orders`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count orders: %w", err)
	}

	return count, nil
}
