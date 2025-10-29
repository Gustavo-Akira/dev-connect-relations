package entities

import "testing"

func TestShouldReturnErrorWhen(t *testing.T) {
	var name string = ""
	stack, err := NewStack(name)
	if err == nil || stack != nil {
		t.Error("Should throw error when stack name is invalid")
	}
}
