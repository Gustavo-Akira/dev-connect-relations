package entities

import (
	"testing"
)

func TestShouldDeclareProfileWithAllValidProperties(t *testing.T) {
	name := "Akira"
	var id int64 = 1
	profile, err := NewProfile(id, name)
	if profile == nil || err != nil {
		t.Error("Should name and id be valid and create the profile")
	}
}

func TestShouldReturnErrorWithInvalidName(t *testing.T) {
	name := ""
	var id int64 = 1
	profile, err := NewProfile(id, name)
	if profile != nil || err.Error() != "name cannot be empty" {
		t.Error("Should return name is invalid")
	}
}

func TestShouldReturnErrorWithInvalidId(t *testing.T) {
	name := "akira"
	var id int64 = 0
	profile, err := NewProfile(id, name)
	if profile != nil || err.Error() != "id cannot be 0" {
		t.Error("Should return id is invalid")
	}
}
