package events

import "testing"

func TestIsValidTitle(t *testing.T) {
	result := IsValidTitle("Привет, мир 0.9")
	if !result {
		t.Error("Ожидали положительный вариант валидации, но получен отрицательный")
	}
}
