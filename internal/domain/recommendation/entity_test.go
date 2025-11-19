package recommendation

import "testing"

func TestCreateRecommencDation(t *testing.T) {
	rec := CreateRecommendation(1, 0.85)
	if rec.ID != 1 {
		t.Errorf("Expected ID to be 'rec1', got '%d'", rec.ID)
	}
}
