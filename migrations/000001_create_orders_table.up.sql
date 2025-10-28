-- Создание нормализованной схемы БД (3NF)

-- 1. Таблица заказов (основная информация)
CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(50) NOT NULL,
    locale VARCHAR(10),
    internal_signature TEXT,
    customer_id VARCHAR(255),
    delivery_service VARCHAR(100),
    shardkey VARCHAR(10),
    sm_id INTEGER,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    oof_shard VARCHAR(10)
    );

-- 2. Таблица доставок (связь 1:1 с orders)
CREATE TABLE IF NOT EXISTS deliveries (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    zip VARCHAR(20) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    region VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    CONSTRAINT fk_delivery_order
    FOREIGN KEY (order_uid)
    REFERENCES orders(order_uid)
    ON DELETE CASCADE
    );

-- 3. Таблица платежей (связь 1:1 с orders)
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) UNIQUE NOT NULL,
    transaction VARCHAR(255) NOT NULL,
    request_id VARCHAR(255),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(100) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER DEFAULT 0,
    CONSTRAINT fk_payment_order
    FOREIGN KEY (order_uid)
    REFERENCES orders(order_uid)
    ON DELETE CASCADE
    );

-- 4. Таблица товаров (связь 1:N с orders)
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL,
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INTEGER DEFAULT 0,
    size VARCHAR(50),
    total_price INTEGER NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR(255),
    status INTEGER NOT NULL,
    CONSTRAINT fk_item_order
    FOREIGN KEY (order_uid)
    REFERENCES orders(order_uid)
    ON DELETE CASCADE
    );

-- 5. Индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_orders_track_number ON orders(track_number);
CREATE INDEX IF NOT EXISTS idx_orders_date_created ON orders(date_created DESC);
CREATE INDEX IF NOT EXISTS idx_deliveries_order_uid ON deliveries(order_uid);
CREATE INDEX IF NOT EXISTS idx_payments_order_uid ON payments(order_uid);
CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);
CREATE INDEX IF NOT EXISTS idx_items_chrt_id ON items(chrt_id);