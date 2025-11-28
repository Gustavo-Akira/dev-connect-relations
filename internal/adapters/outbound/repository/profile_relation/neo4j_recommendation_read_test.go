package relation_test

import (
	"context"
	outboundRelation "devconnectrelations/internal/adapters/outbound/repository/profile_relation"
	cityRelation "devconnectrelations/internal/domain/profile_relation/city"
	relation "devconnectrelations/internal/domain/profile_relation/stack"
	"devconnectrelations/internal/domain/recommendation"
	"errors"
	"testing"
)

type MockCityRelationRepository struct {
	result []cityRelation.CityRelation
	err    error
}

func (m *MockCityRelationRepository) GetCityRelatedToProfileIds(ctx context.Context, ids []int64) ([]cityRelation.CityRelation, error) {
	return m.result, m.err
}
func (m *MockCityRelationRepository) JaccardIndexByProfileId(ctx context.Context, profileId int64) ([]recommendation.Recommendation, error) {
	return nil, nil
}

func (m *MockCityRelationRepository) CreateCityRelation(ctx context.Context, cityRelation *cityRelation.CityRelation) (*cityRelation.CityRelation, error) {
	return nil, nil
}

func (m *MockCityRelationRepository) GetCityRelatedToProfileId(ctx context.Context, profileId int64) ([]cityRelation.CityRelation, error) {
	return nil, nil
}

type MockStackRelationRepository struct {
	result []relation.StackRelation
	err    error
}

func (m *MockStackRelationRepository) GetStackRelationByProfileIds(ctx context.Context, ids []int64) ([]relation.StackRelation, error) {
	return m.result, m.err
}

func (m *MockStackRelationRepository) CreateStackRelation(ctx context.Context, stackRelation *relation.StackRelation) (*relation.StackRelation, error) {
	return nil, nil
}

func (m *MockStackRelationRepository) DeleteStackRelation(ctx context.Context, stackName string, profileID int64) error {
	return nil
}

func (m *MockStackRelationRepository) GetStackRelationByProfileId(ctx context.Context, profileId int64) ([]relation.StackRelation, error) {
	return nil, nil
}

func (m *MockStackRelationRepository) JaccardIndexByProfileId(ctx context.Context, profileId int64) ([]recommendation.Recommendation, error) {
	return nil, nil
}

func TestEnrichRecommendations_Success(t *testing.T) {
	ctx := context.Background()

	cityMock := &MockCityRelationRepository{
		result: []cityRelation.CityRelation{
			{ProfileID: 1, CityFullName: "Tokyo"},
			{ProfileID: 2, CityFullName: "Osaka"},
		},
	}

	stackMock := &MockStackRelationRepository{
		result: []relation.StackRelation{
			{ProfileID: 1, StackName: "Go"},
			{ProfileID: 1, StackName: "React"},
			{ProfileID: 2, StackName: "Java"},
		},
	}

	repo := outboundRelation.CreateNeo4jRecommendationRepository(cityMock, stackMock)
	agg := []recommendation.AggregatedScore{
		{ID: 1, Score: 0.9, Name: "Gustavo"},
		{ID: 2, Score: 0.7, Name: "Akira"},
	}

	result, err := repo.EnrichRecommendations(ctx, agg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}

	if result[0].City != "Tokyo" || result[1].City != "Osaka" {
		t.Errorf("city mapping failed")
	}

	if len(result[0].Stacks) != 2 {
		t.Errorf("expected stacks for ID 1")
	}
}

func TestEnrichRecommendations_ErrorCities(t *testing.T) {
	ctx := context.Background()

	cityMock := &MockCityRelationRepository{
		err: errors.New("city error"),
	}

	stackMock := &MockStackRelationRepository{}

	repo := outboundRelation.CreateNeo4jRecommendationRepository(cityMock, stackMock)

	agg := []recommendation.AggregatedScore{{ID: 1}}

	_, err := repo.EnrichRecommendations(ctx, agg)

	if err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func TestEnrichRecommendations_ErrorStacks(t *testing.T) {
	ctx := context.Background()

	cityMock := &MockCityRelationRepository{
		result: []cityRelation.CityRelation{{ProfileID: 1, CityFullName: "Tokyo"}},
	}

	stackMock := &MockStackRelationRepository{
		err: errors.New("stack error"),
	}

	repo := outboundRelation.CreateNeo4jRecommendationRepository(cityMock, stackMock)

	agg := []recommendation.AggregatedScore{{ID: 1}}

	_, err := repo.EnrichRecommendations(ctx, agg)

	if err == nil {
		t.Fatalf("expected stack error but got nil")
	}
}

func TestEnrichRecommendations_EmptyFields(t *testing.T) {
	ctx := context.Background()

	cityMock := &MockCityRelationRepository{
		result: []cityRelation.CityRelation{}, // nenhum perfil tem cidade
	}

	stackMock := &MockStackRelationRepository{
		result: []relation.StackRelation{}, // nenhum stack tbm
	}

	repo := outboundRelation.CreateNeo4jRecommendationRepository(cityMock, stackMock)

	agg := []recommendation.AggregatedScore{
		{ID: 10, Score: 0.3, Name: "NoData"},
	}

	result, err := repo.EnrichRecommendations(ctx, agg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result[0].City != "" {
		t.Errorf("expected empty city but got: %s", result[0].City)
	}

	if len(result[0].Stacks) != 0 {
		t.Errorf("expected no stacks")
	}
}
