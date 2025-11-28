package recommendation_test

import (
	"context"
	"devconnectrelations/internal/domain/recommendation"
	"errors"
	"testing"
)

type MockRecommendationAlgorithm struct{}

func (m *MockRecommendationAlgorithm) Run(ctx context.Context, weights []float64, profileId int64) ([]recommendation.AggregatedScore, error) {
	if profileId == 0 {
		return []recommendation.AggregatedScore{}, errors.New("invalid profile ID")
	}
	return []recommendation.AggregatedScore{
		{ID: 2, Score: 0.8},
		{ID: 3, Score: 0.6},
	}, nil
}

func TestGetRecommendationByProfileId(t *testing.T) {
	mockAlgorithm := &MockRecommendationAlgorithm{}
	service := &recommendation.RecommendationService{
		RecommendationAlgorithm: mockAlgorithm,
	}
	profileID := int64(1)
	recommendations, err := service.GetRecommendationByProfileId(context.Background(), profileID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedCount := 2
	if len(recommendations) != expectedCount {
		t.Fatalf("Expected %d recommendations, got %d", expectedCount, len(recommendations))
	}
	expectedFirstID := int64(2)
	if recommendations[0].ID != expectedFirstID {
		t.Errorf("Expected first recommendation ID to be %d, got %d", expectedFirstID, recommendations[0].ID)
	}
	expectedSecondID := int64(3)
	if recommendations[1].ID != expectedSecondID {
		t.Errorf("Expected second recommendation ID to be %d, got %d", expectedSecondID, recommendations[1].ID)
	}
}

func TestGetRecommendationByProfileId_InvalidProfileID(t *testing.T) {
	mockAlgorithm := &MockRecommendationAlgorithm{}
	service := &recommendation.RecommendationService{
		RecommendationAlgorithm: mockAlgorithm,
	}
	profileID := int64(0)
	_, err := service.GetRecommendationByProfileId(context.Background(), profileID)
	if err == nil {
		t.Fatal("Expected error for invalid profile ID, got nil")
	}
	expectedErrorMessage := "invalid profile ID"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message to be '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}
