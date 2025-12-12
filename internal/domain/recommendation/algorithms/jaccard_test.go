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
		{ID: 1, Score: 0.8, Name: "Gustavo"},
		{ID: 2, Score: 0.6, Name: "Akira"},
	}, nil
}

func (m *MockCityRelationRepository) CreateCityRelation(ctx context.Context, cityRelation *city.CityRelation) (*city.CityRelation, error) {
	return nil, nil
}

func (m *MockCityRelationRepository) GetCityRelatedToProfileId(ctx context.Context, profileId int64) ([]city.CityRelation, error) {
	return nil, nil
}

func (m *MockCityRelationRepository) GetCityRelatedToProfileIds(ctx context.Context, profileIds []int64) ([]city.CityRelation, error) {
	return nil, nil
}

type MockRelationsRepository struct{}

func (m *MockRelationsRepository) JaccardIndexByProfileId(ctx context.Context, profileId int64) ([]recommendation.Recommendation, error) {
	if profileId == 998 {
		return nil, errors.New("Error on repository")
	}
	if profileId == 666 {
		return []recommendation.Recommendation{
			{ID: 775, Score: 1, Name: "Errro"},
		}, nil
	}
	return []recommendation.Recommendation{
		{ID: 1, Score: 0.7, Name: "Gustavo"},
		{ID: 2, Score: 0.5, Name: "Akira"},
	}, nil
}

func (m *MockRelationsRepository) CreateRelation(context context.Context, relation relation.Relation) (relation.Relation, error) {
	return relation, nil
}
func (m *MockRelationsRepository) GetAllRelationsByFromId(ctx context.Context, fromId int64, offset int64, limit int64) ([]relation.Relation, error) {
	return nil, nil
}
func (m *MockRelationsRepository) AcceptRelation(ctx context.Context, fromId int64, toId int64) error {
	return nil
}
func (m *MockRelationsRepository) GetAllRelationPendingByFromId(ctx context.Context, fromId int64) ([]relation.Relation, error) {
	return nil, nil
}

func (m *MockRelationsRepository) CountRelationsByFromId(ctx context.Context, fromId int64) (int64, error) {
	return 0, nil
}

type MockStackRelationRepository struct{}

func (m *MockStackRelationRepository) JaccardIndexByProfileId(ctx context.Context, profileId int64) ([]recommendation.Recommendation, error) {
	if profileId == 999 {
		return nil, errors.New("Error on repository")
	}

	if profileId == 777 {
		return []recommendation.Recommendation{
			{ID: 776, Score: 1, Name: "Errro"},
		}, nil
	}
	return []recommendation.Recommendation{
		{ID: 1, Score: 0.9, Name: "Gustavo"},
		{ID: 2, Score: 0.4, Name: "Akira"},
	}, nil
}

func (m *MockStackRelationRepository) GetStackRelationByProfileIds(ctx context.Context, profileIds []int64) ([]stack.StackRelation, error) {
	return nil, nil
}

func (m *MockStackRelationRepository) CreateStackRelation(ctx context.Context, stackRelation *stack.StackRelation) (*stack.StackRelation, error) {
	return nil, nil
}

func (m *MockStackRelationRepository) DeleteStackRelation(ctx context.Context, stackName string, profileID int64) error {
	return nil
}

func (m *MockStackRelationRepository) GetStackRelationByProfileId(ctx context.Context, profileId int64) ([]stack.StackRelation, error) {
	return nil, nil
}

var (
	mockCityRepo      = &MockCityRelationRepository{}
	mockRelationsRepo = &MockRelationsRepository{}
	mockStacksRepo    = &MockStackRelationRepository{}
)

func TestJaccardAlgorithm(t *testing.T) {
	JaccardAlgorithm := algorithms.NewJaccardAlgorithm(mockCityRepo, mockRelationsRepo, mockStacksRepo)
	weights := []float64{0.5, 0.3, 0.2}

	recs, err := JaccardAlgorithm.Run(context.Background(), weights, 123)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedScores := map[int64]float64{
		1: 0.5*0.8 + 0.3*0.9 + 0.2*0.7,
		2: 0.5*0.6 + 0.3*0.4 + 0.2*0.5,
	}

	if len(recs) != 2 {
		t.Fatalf("Expected 2 recommendations, got %d", len(recs))
	}

	if recs[0].ID != 1 || recs[0].Name != "Gustavo" {
		t.Errorf("Expected first recommendation to be Gustavo with ID 1")
	}

	for _, rec := range recs {
		expectedScore := expectedScores[rec.ID]
		if rec.Score != expectedScore {
			t.Errorf("For ID %d expected score %.2f got %.2f", rec.ID, expectedScore, rec.Score)
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

func TestJaccardAlgorithm_RepositoryErrorStackBatch(t *testing.T) {
	mockStackWithError := &MockStackRelationRepository{}

	JaccardAlgorithm := algorithms.NewJaccardAlgorithm(mockCityRepo, mockRelationsRepo, mockStackWithError)
	weights := []float64{0.5, 0.3, 0.2}

	_, err := JaccardAlgorithm.Run(context.Background(), weights, 777)

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	expected := "Error on repository stacks"
	if err.Error() != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestJaccardAlgorithm_RepositoryErrorCityBatch(t *testing.T) {
	mockCity := &MockCityRelationRepository{}

	JaccardAlgorithm := algorithms.NewJaccardAlgorithm(mockCity, mockRelationsRepo, mockStacksRepo)
	weights := []float64{0.5, 0.3, 0.2}

	_, err := JaccardAlgorithm.Run(context.Background(), weights, 666)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	expected := "Error on repository city batch"
	if err.Error() != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, err.Error())
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
