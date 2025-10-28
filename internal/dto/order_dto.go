package dto

import (
	"time"
)

// CreateOrderInput - входные данные для создания заказа
type CreateOrderInput struct {
	OrderUID          string        `json:"order_uid"`
	TrackNumber       string        `json:"track_number"`
	Entry             string        `json:"entry"`
	Delivery          DeliveryInput `json:"delivery"`
	Payment           PaymentInput  `json:"payment"`
	Items             []ItemInput   `json:"items"`
	Locale            string        `json:"locale"`
	InternalSignature string        `json:"internal_signature"`
	CustomerID        string        `json:"customer_id"`
	DeliveryService   string        `json:"delivery_service"`
	Shardkey          string        `json:"shardkey"`
	SmID              int           `json:"sm_id"`
	DateCreated       time.Time     `json:"date_created"`
	OofShard          string        `json:"oof_shard"`
}

// DeliveryInput - входные данные доставки
type DeliveryInput struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

// PaymentInput - входные данные платежа
type PaymentInput struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

// ItemInput - входные данные товара
type ItemInput struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

// GetOrderInput - входные данные для получения заказа
type GetOrderInput struct {
	OrderUID string `json:"order_uid"`
}

// OrderOutput - выходные данные заказа
type OrderOutput struct {
	OrderUID          string         `json:"order_uid"`
	TrackNumber       string         `json:"track_number"`
	Entry             string         `json:"entry"`
	Delivery          DeliveryOutput `json:"delivery"`
	Payment           PaymentOutput  `json:"payment"`
	Items             []ItemOutput   `json:"items"`
	Locale            string         `json:"locale"`
	InternalSignature string         `json:"internal_signature"`
	CustomerID        string         `json:"customer_id"`
	DeliveryService   string         `json:"delivery_service"`
	Shardkey          string         `json:"shardkey"`
	SmID              int            `json:"sm_id"`
	DateCreated       time.Time      `json:"date_created"`
	OofShard          string         `json:"oof_shard"`
}

// DeliveryOutput - выходные данные доставки
type DeliveryOutput struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

// PaymentOutput - выходные данные платежа
type PaymentOutput struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

// ItemOutput - выходные данные товара
type ItemOutput struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
