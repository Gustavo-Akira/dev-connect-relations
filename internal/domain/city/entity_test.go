package city

import "testing"

func TestNewCity(t *testing.T) {
	city := NewCity(" New York ", " USA ", " NY ")
	if city.Name != "new york" {
		t.Errorf("expected 'new york', got '%s'", city.Name)
	}
	if city.Country != "usa" {
		t.Errorf("expected 'usa', got '%s'", city.Country)
	}
	if city.State != "ny" {
		t.Errorf("expected 'ny', got '%s'", city.State)
	}

	fullName := city.GetFullName()
	expectedFullName := "new york, ny, usa"
	if fullName != expectedFullName {
		t.Errorf("expected '%s', got '%s'", expectedFullName, fullName)
	}
}

func TestNormalizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  Hello World  ", "hello world"},
		{"TESTING", "testing"},
		{"  Mixed Case String  ", "mixed case string"},
		{"   ", ""},
		{"NoSpaces", "nospaces"},
	}
	for _, test := range tests {
		result := normalizeString(test.input)
		if result != test.expected {
			t.Errorf("normalizeString(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
