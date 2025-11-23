package recommendation_test

import (
	"devconnectrelations/internal/domain/recommendation"
	"testing"
)

func TestCreateRecommencDation(t *testing.T) {
	rec := recommendation.CreateRecommendation(1, 0.85, "Gustavo")
	if rec.ID != 1 {
		t.Errorf("Expected ID to be 'rec1', got '%d'", rec.ID)
	}
	if rec.Name != "Gustavo" {
		t.Errorf("Expected ID to be 'Gustavo', got '%s'", rec.Name)
	}
}
