-- Создание таблицы orders
CREATE TABLE IF NOT EXISTS orders (
                                      order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(50) NOT NULL,
    delivery_name VARCHAR(255) NOT NULL,
    delivery_phone VARCHAR(50) NOT NULL,
    delivery_zip VARCHAR(20) NOT NULL,
    delivery_city VARCHAR(255) NOT NULL,
    delivery_address VARCHAR(500) NOT NULL,
    delivery_region VARCHAR(255) NOT NULL,
    delivery_email VARCHAR(255) NOT NULL,
    payment_transaction VARCHAR(255) NOT NULL,
    payment_request_id VARCHAR(255) NOT NULL,
    payment_currency VARCHAR(10) NOT NULL,
    payment_provider VARCHAR(100) NOT NULL,
    payment_amount INT NOT NULL,
    payment_payment_dt BIGINT NOT NULL,
    payment_bank VARCHAR(100) NOT NULL,
    payment_delivery_cost INT NOT NULL,
    payment_goods_total INT NOT NULL,
    payment_custom_fee INT NOT NULL,
    items JSONB NOT NULL,
    locale VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(500),
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(100) NOT NULL,
    shardkey VARCHAR(10) NOT NULL,
    sm_id BIGINT NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Создание индекса для быстрого поиска по order_uid
CREATE INDEX IF NOT EXISTS idx_orders_order_uid ON orders(order_uid);

-- Создание индекса для поиска по track_number
CREATE INDEX IF NOT EXISTS idx_orders_track_number ON orders(track_number);
