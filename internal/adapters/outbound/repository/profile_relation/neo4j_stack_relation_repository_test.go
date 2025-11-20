package relation_test

import (
	relation "devconnectrelations/internal/adapters/outbound/repository/profile_relation"
	"devconnectrelations/internal/tests"
	"testing"
)

func TestCreateStackRelation(t *testing.T) {
	t.Parallel()
	tests.SeedStackRelationships(t, driver)
}

func TestJaccardStackByProfileId(t *testing.T) {
	t.Parallel()
	repo := relation.NewNeo4jStackRelationRepository(driver)
	jaccardIndices, err := repo.JaccardIndexByProfileId(t.Context(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(jaccardIndices) == 0 {
		t.Fatalf("expected at least one jaccard index, got 0")
	}
}
