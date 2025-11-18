package recommendation

import "testing"

func TestCreateRecommencDation(t *testing.T) {
	rec := CreateRecommendation("rec1", 0.85)
	if rec.ID != "rec1" {
		t.Errorf("Expected ID to be 'rec1', got '%s'", rec.ID)
	}
}
