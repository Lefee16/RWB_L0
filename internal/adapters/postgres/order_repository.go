package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"RWB_L0/internal/domain"
)

// OrderRepository - репозиторий для работы с заказами
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository - создание репозитория заказов
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Save - сохранить заказ в БД
func (r *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (
			order_uid, track_number, entry,
			delivery_name, delivery_phone, delivery_zip, delivery_city, 
			delivery_address, delivery_region, delivery_email,
			payment_transaction, payment_request_id, payment_currency,
			payment_provider, payment_amount, payment_payment_dt,
			payment_bank, payment_delivery_cost, payment_goods_total, payment_custom_fee,
			items, locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29
		)
		ON CONFLICT (order_uid) DO UPDATE SET
			track_number = EXCLUDED.track_number,
			entry = EXCLUDED.entry,
			delivery_name = EXCLUDED.delivery_name,
			delivery_phone = EXCLUDED.delivery_phone,
			items = EXCLUDED.items
	`

	// Сериализуем Items в JSONB
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		// Delivery
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
		// Payment
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
		// Other
		itemsJSON,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)

	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	return nil
}

// GetByID - получить заказ по ID
func (r *OrderRepository) GetByID(ctx context.Context, orderUID string) (*domain.Order, error) {
	query := `
		SELECT 
			order_uid, track_number, entry,
			delivery_name, delivery_phone, delivery_zip, delivery_city,
			delivery_address, delivery_region, delivery_email,
			payment_transaction, payment_request_id, payment_currency,
			payment_provider, payment_amount, payment_payment_dt,
			payment_bank, payment_delivery_cost, payment_goods_total, payment_custom_fee,
			items, locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
		WHERE order_uid = $1
	`

	order, err := r.scanOrder(r.db.QueryRowContext(ctx, query, orderUID))
	if err != nil {
		// Используем errors.Is() вместо ==
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEmptyOrderUID
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

// GetAll - получить все заказы
func (r *OrderRepository) GetAll(ctx context.Context) ([]*domain.Order, error) {
	query := `
		SELECT 
			order_uid, track_number, entry,
			delivery_name, delivery_phone, delivery_zip, delivery_city,
			delivery_address, delivery_region, delivery_email,
			payment_transaction, payment_request_id, payment_currency,
			payment_provider, payment_amount, payment_payment_dt,
			payment_bank, payment_delivery_cost, payment_goods_total, payment_custom_fee,
			items, locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
		ORDER BY date_created DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var orders []*domain.Order

	for rows.Next() {
		order, err := r.scanOrder(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return orders, nil
}

// Delete - удалить заказ (для тестов)
func (r *OrderRepository) Delete(ctx context.Context, orderUID string) error {
	query := `DELETE FROM orders WHERE order_uid = $1`

	_, err := r.db.ExecContext(ctx, query, orderUID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
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

// scanOrder - вспомогательная функция для сканирования строки в Order (убирает дублирование)
func (r *OrderRepository) scanOrder(scanner scanner) (*domain.Order, error) {
	var (
		order     domain.Order
		itemsJSON []byte
	)

	err := scanner.Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		// Delivery
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
		// Payment
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
		// Other
		&itemsJSON,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)

	if err != nil {
		return nil, err
	}

	// Десериализуем Items из JSONB
	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return &order, nil
}

// Scanner - интерфейс для унификации sql.Row и sql.Rows
type scanner interface {
	Scan(dest ...interface{}) error
}
