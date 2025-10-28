package domain

import "errors"

var (
	ErrEmptyDeliveryName = errors.New("delivery name cannot be empty")

	ErrEmptyDeliveryPhone = errors.New("delivery phone cannot be empty")

	ErrEmptyPaymentTransaction = errors.New("payment transaction cannot be empty")

	ErrInvalidPaymentAmount = errors.New("payment amount must be greater than 0")

	ErrEmptyItemName = errors.New("item name cannot be empty")

	ErrInvalidItemPrice = errors.New("item price must be greater than or equal to 0")

	ErrEmptyOrderUID = errors.New("order UID cannot be empty")

	ErrEmptyTrackNumber = errors.New("track number cannot be empty")

	ErrOrderNotFound = errors.New("order not found")
)
