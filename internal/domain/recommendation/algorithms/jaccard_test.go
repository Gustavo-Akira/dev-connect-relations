package algorithms_test

import (
	"context"
	"devconnectrelations/internal/domain/profile_relation/city"
	"devconnectrelations/internal/domain/profile_relation/relation"
	"devconnectrelations/internal/domain/profile_relation/stack"
	"devconnectrelations/internal/domain/recommendation"
	"devconnectrelations/internal/domain/recommendation/algorithms"
	"errors"
	"testing"
)

type MockCityRelationRepository struct{}

func (m *MockCityRelationRepository) JaccardIndexByProfileId(ctx context.Context, profileId int64) ([]recommendation.Recommendation, error) {
	if profileId == 997 {
		return nil, errors.New("Error on repository")
	}
	return []recommendation.Recommendation{
		{ID: 1, Score: 0.8},
		{ID: 2, Score: 0.6},
	}, nil
}

func (m *MockCityRelationRepository) CreateCityRelation(ctx context.Context, cityRelation *city.CityRelation) (*city.CityRelation, error) {
	return nil, nil
}

type MockRelationsRepository struct{}

func (m *MockRelationsRepository) JaccardIndexByProfileId(ctx context.Context, profileId int64) ([]recommendation.Recommendation, error) {
	if profileId == 998 {
		return nil, errors.New("Error on repository")
	}
	return []recommendation.Recommendation{
		{ID: 1, Score: 0.7},
		{ID: 2, Score: 0.5},
	}, nil
}

func (m *MockRelationsRepository) CreateRelation(context context.Context, relation relation.Relation) (relation.Relation, error) {
	return relation, nil
}
func (m *MockRelationsRepository) GetAllRelationsByFromId(ctx context.Context, fromId int64) ([]relation.Relation, error) {
	return nil, nil
}
func (m *MockRelationsRepository) AcceptRelation(ctx context.Context, fromId int64, toId int64) error {
	return nil
}
func (m *MockRelationsRepository) GetAllRelationPendingByFromId(ctx context.Context, fromId int64) ([]relation.Relation, error) {
	return nil, nil
}

type MockStackRelationRepository struct{}

func (m *MockStackRelationRepository) JaccardIndexByProfileId(ctx context.Context, profileId int64) ([]recommendation.Recommendation, error) {
	if profileId == 999 {
		return nil, errors.New("Error on repository")
	}
	return []recommendation.Recommendation{
		{ID: 1, Score: 0.9},
		{ID: 2, Score: 0.4},
	}, nil
}

func (m *MockStackRelationRepository) CreateStackRelation(ctx context.Context, stackRelation *stack.StackRelation) (*stack.StackRelation, error) {
	return nil, nil
}

func (m *MockStackRelationRepository) DeleteStackRelation(ctx context.Context, stackName string, profileID int64) error {
	return nil
}

var (
	mockCityRepo      = &MockCityRelationRepository{}
	mockRelationsRepo = &MockRelationsRepository{}
	mockStacksRepo    = &MockStackRelationRepository{}
)

func TestJaccardAlgorithm(t *testing.T) {
	JaccardAlgorithm := algorithms.NewJaccardAlgorithm(mockCityRepo, mockRelationsRepo, mockStacksRepo)
	weights := []float64{0.5, 0.3, 0.2}
	recommendations, err := JaccardAlgorithm.Run(context.Background(), weights, 123)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedScores := map[int64]float64{
		1: 0.5*0.8 + 0.3*0.9 + 0.2*0.7,
		2: 0.5*0.6 + 0.3*0.4 + 0.2*0.5,
	}
	for _, rec := range recommendations {
		expectedScore, exists := expectedScores[rec.ID]
		if !exists {
			t.Errorf("Unexpected recommendation ID %d", rec.ID)
		}
		if rec.Score != expectedScore {
			t.Errorf("For ID %d, expected score %f, got %f", rec.ID, expectedScore, rec.Score)
		}
	}
}

func TestJaccardAlgorithm_RepositoryError(t *testing.T) {
	JaccardAlgorithm := algorithms.NewJaccardAlgorithm(mockCityRepo, mockRelationsRepo, mockStacksRepo)
	weights := []float64{0.5, 0.3, 0.2}
	_, err := JaccardAlgorithm.Run(context.Background(), weights, 999)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	expectedErrorMessage := "Error on repository"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestJaccardAlgorithm_RepositoryErrorCity(t *testing.T) {
	JaccardAlgorithm := algorithms.NewJaccardAlgorithm(mockCityRepo, mockRelationsRepo, mockStacksRepo)
	weights := []float64{0.5, 0.3, 0.2}
	_, err := JaccardAlgorithm.Run(context.Background(), weights, 997)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	expectedErrorMessage := "Error on repository"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestJaccardAlgorithm_RepositoryErrorRelation(t *testing.T) {
	weights := []float64{0.5, 0.3, 0.2}
	JaccardAlgorithm := algorithms.NewJaccardAlgorithm(mockCityRepo, mockRelationsRepo, mockStacksRepo)
	_, err := JaccardAlgorithm.Run(context.Background(), weights, 998)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	expectedErrorMessage := "Error on repository"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}
