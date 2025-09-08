package main

import "testing"

func TestAdd(t *testing.T) {
	var op Operator = Add{}
	result := op.Apply(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add.Apply(2, 3) = %d; want %d", result, expected)
	}
}

func TestMultiply(t *testing.T) {
	var op Operator = Multiply{}
	result := op.Apply(4, 5)
	expected := 20

	if result != expected {
		t.Errorf("Multiply.Apply(4, 5) = %d; want %d", result, expected)
	}
}
