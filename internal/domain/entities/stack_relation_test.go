package entities

import "testing"

func TestNewStackRelation(t *testing.T) {
	stackName := "Go"
	var userID int64 = 12
	stackRelation, stackError := NewStackRelation(stackName, userID)
	if stackError != nil {
		t.Errorf("unexpected error: %v", stackError)
	}
	if stackRelation.StackName != stackName {
		t.Errorf("expected StackName %s, got %s", stackName, stackRelation.StackName)
	}
	if stackRelation.ProfileID != userID {
		t.Errorf("expected UserID %d, got %d", userID, stackRelation.ProfileID)
	}
}

func TestNewStackRelation_InvalidParameters(t *testing.T) {
	_, stackError := NewStackRelation("", 0)
	if stackError == nil {
		t.Error("expected error for invalid parameters, got nil")
	}
	_, stackError = NewStackRelation("Python", 0)
	if stackError == nil {
		t.Error("expected error for invalid userID, got nil")
	}
	_, stackError = NewStackRelation("", 5)
	if stackError == nil {
		t.Error("expected error for invalid stackName, got nil")
	}
}
