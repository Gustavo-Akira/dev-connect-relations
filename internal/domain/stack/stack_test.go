package stack

import "testing"

func TestShouldReturnErrorWhen(t *testing.T) {
	var name string = ""
	stack, err := NewStack(name)
	if err == nil || stack != nil {
		t.Error("Should throw error when stack name is invalid")
	}
}

func TestShouldCreateStackAndLowercaseName(t *testing.T) {
	var name string = "GOLANG"
	stack, err := NewStack(name)
	if err != nil {
		t.Error("Should not throw error when stack name is valid")
	}
	if stack.Name != "golang" {
		t.Errorf("Expected stack name to be 'golang', got '%s'", stack.Name)
	}
}
