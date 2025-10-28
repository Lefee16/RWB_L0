package domain

import "errors"

var (
	// ErrEmptyOrderUID - пустой ID заказа
	ErrEmptyOrderUID = errors.New("order_uid cannot be empty")

	// ErrEmptyTrackNumber - пустой трек-номер
	ErrEmptyTrackNumber = errors.New("track_number cannot be empty")

	// ErrInvalidTotal - некорректная сумма заказа
	ErrInvalidTotal = errors.New("total amount must be greater than 0")

	// ErrEmptyDeliveryName - пустое имя получателя
	ErrEmptyDeliveryName = errors.New("delivery name cannot be empty")

	// ErrEmptyDeliveryPhone - пустой телефон
	ErrEmptyDeliveryPhone = errors.New("delivery phone cannot be empty")

	// ErrEmptyPaymentTransaction - пустая транзакция
	ErrEmptyPaymentTransaction = errors.New("payment transaction cannot be empty")

	// ErrInvalidPaymentAmount - некорректная сумма платежа
	ErrInvalidPaymentAmount = errors.New("payment amount must be greater than 0")

	// ErrEmptyItemName - пустое название товара
	ErrEmptyItemName = errors.New("item name cannot be empty")

	// ErrInvalidItemPrice - некорректная цена товара
	ErrInvalidItemPrice = errors.New("item price must be greater than or equal to 0")
)
