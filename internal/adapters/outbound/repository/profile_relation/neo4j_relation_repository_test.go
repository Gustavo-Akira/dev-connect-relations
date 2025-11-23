package relation_test

import (
	relation_repo "devconnectrelations/internal/adapters/outbound/repository/profile_relation"
	"devconnectrelations/internal/tests"
	"testing"
)

func TestCreateRelation(t *testing.T) {
	tests.SeedRelations(t, driver)
}

func TestProfileJaccardIndexByProfileId(t *testing.T) {
	tests.SeedRelations(t, driver)
	repo := relation_repo.NewNeo4jRelationRepository(driver)
	reccomendations, err := repo.JaccardIndexByProfileId(t.Context(), 1)
	if err != nil {
		t.Errorf("Error getting jaccard index: %v", err)
	}

	if len(reccomendations) == 0 {
		t.Errorf("Expected at least one recommendation, got none")
	}
	t.Logf("Reccomendations: %v", reccomendations)
}
