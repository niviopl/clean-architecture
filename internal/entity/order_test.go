package entity

import "testing"

func TestNewOrder(t *testing.T) {
	order, err := NewOrder("1", 100, 5)
	if err != nil {
		t.Fatal(err)
	}
	if order.FinalPrice != 105 {
		t.Errorf("expected final price 105, got %f", order.FinalPrice)
	}
}

func TestNewOrder_InvalidID(t *testing.T) {
	_, err := NewOrder("", 100, 5)
	if err != ErrIDIsRequired {
		t.Errorf("expected ErrIDIsRequired, got %v", err)
	}
}

func TestNewOrder_InvalidPrice(t *testing.T) {
	_, err := NewOrder("1", 0, 5)
	if err != ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice, got %v", err)
	}
	_, err = NewOrder("1", -1, 5)
	if err != ErrInvalidPrice {
		t.Errorf("expected ErrInvalidPrice, got %v", err)
	}
}

func TestNewOrder_InvalidTax(t *testing.T) {
	_, err := NewOrder("1", 100, 0)
	if err != ErrInvalidTax {
		t.Errorf("expected ErrInvalidTax, got %v", err)
	}
	_, err = NewOrder("1", 100, -1)
	if err != ErrInvalidTax {
		t.Errorf("expected ErrInvalidTax, got %v", err)
	}
}

func TestCalculateFinalPrice(t *testing.T) {
	order := Order{ID: "1", Price: 10.5, Tax: 2.3}
	order.CalculateFinalPrice()
	if order.FinalPrice != 12.8 {
		t.Errorf("expected final price 12.8, got %f", order.FinalPrice)
	}
}
