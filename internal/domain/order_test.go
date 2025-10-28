package domain

import (
	"errors"
	"testing"
)

func TestNewOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderUID    string
		trackNumber string
		entry       string
		wantErr     error
	}{
		{
			name:        "Valid order",
			orderUID:    "test123",
			trackNumber: "TRACK123",
			entry:       "WBIL",
			wantErr:     nil,
		},
		{
			name:        "Empty order_uid",
			orderUID:    "",
			trackNumber: "TRACK123",
			entry:       "WBIL",
			wantErr:     ErrEmptyOrderUID,
		},
		{
			name:        "Empty track_number",
			orderUID:    "test123",
			trackNumber: "",
			entry:       "WBIL",
			wantErr:     ErrEmptyTrackNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(tt.orderUID, tt.trackNumber, tt.entry)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if order.OrderUID != tt.orderUID {
					t.Errorf("OrderUID = %v, want %v", order.OrderUID, tt.orderUID)
				}
				if order.TrackNumber != tt.trackNumber {
					t.Errorf("TrackNumber = %v, want %v", order.TrackNumber, tt.trackNumber)
				}
			}
		})
	}
}

func TestOrder_AddItem(t *testing.T) {
	order, _ := NewOrder("test123", "TRACK123", "WBIL")

	// Валидный товар
	validItem := Item{
		ChrtID: 123,
		Name:   "Test Item",
		Price:  1000,
	}

	err := order.AddItem(validItem)
	if err != nil {
		t.Errorf("AddItem() error = %v, want nil", err)
	}

	if order.GetItemsCount() != 1 {
		t.Errorf("GetItemsCount() = %v, want 1", order.GetItemsCount())
	}

	// Невалидный товар
	invalidItem := Item{
		ChrtID: 456,
		Name:   "", // Пустое имя
		Price:  1000,
	}

	err = order.AddItem(invalidItem)
	if !errors.Is(err, ErrEmptyItemName) {
		t.Errorf("AddItem() error = %v, want %v", err, ErrEmptyItemName)
	}
}

func TestDelivery_Validate(t *testing.T) {
	tests := []struct {
		name     string
		delivery Delivery
		wantErr  error
	}{
		{
			name: "Valid delivery",
			delivery: Delivery{
				Name:  "Test User",
				Phone: "+79001234567",
			},
			wantErr: nil,
		},
		{
			name: "Empty name",
			delivery: Delivery{
				Name:  "",
				Phone: "+79001234567",
			},
			wantErr: ErrEmptyDeliveryName,
		},
		{
			name: "Empty phone",
			delivery: Delivery{
				Name:  "Test User",
				Phone: "",
			},
			wantErr: ErrEmptyDeliveryPhone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.delivery.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPayment_Validate(t *testing.T) {
	tests := []struct {
		name    string
		payment Payment
		wantErr error
	}{
		{
			name: "Valid payment",
			payment: Payment{
				Transaction: "test123",
				Amount:      1000,
			},
			wantErr: nil,
		},
		{
			name: "Empty transaction",
			payment: Payment{
				Transaction: "",
				Amount:      1000,
			},
			wantErr: ErrEmptyPaymentTransaction,
		},
		{
			name: "Invalid amount",
			payment: Payment{
				Transaction: "test123",
				Amount:      0,
			},
			wantErr: ErrInvalidPaymentAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payment.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
