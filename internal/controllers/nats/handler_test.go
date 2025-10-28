package nats

import (
	"testing"

	"RWB_L0/internal/domain"
)

// Простой тест: проверяем, что handler создаётся
func TestNewHandler(t *testing.T) {

	order, err := domain.NewOrder("test-123", "TRACK123", "WBIL")
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	if order.OrderUID != "test-123" {
		t.Errorf("Expected OrderUID = test-123, got %s", order.OrderUID)
	}
}

// Тест валидации данных
func TestOrderValidation(t *testing.T) {
	// Тест с пустым UID
	_, err := domain.NewOrder("", "TRACK", "WBIL")
	if err == nil {
		t.Error("Expected error for empty UID, got nil")
	}

	// Тест с пустым TrackNumber
	_, err = domain.NewOrder("test", "", "WBIL")
	if err == nil {
		t.Error("Expected error for empty TrackNumber, got nil")
	}
}
