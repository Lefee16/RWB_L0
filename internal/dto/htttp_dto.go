package dto

// ErrorResponse - стандартный ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// SuccessResponse - стандартный успешный ответ
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// GetOrderResponse - ответ с данными заказа
type GetOrderResponse struct {
	Order *OrderOutput `json:"order"`
}

// CreateOrderResponse - ответ при создании заказа
type CreateOrderResponse struct {
	Success  bool   `json:"success"`
	OrderUID string `json:"order_uid"`
	Message  string `json:"message"`
}

// HealthCheckResponse - ответ health check
type HealthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}
