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

func TestGetStackRelationByProfileId(t *testing.T) {
	t.Parallel()

	tests.SeedProfiles(t, driver)
	tests.SeedStackRelationships(t, driver)

	repo := relation.NewNeo4jStackRelationRepository(driver)

	relations, err := repo.GetStackRelationByProfileId(t.Context(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(relations) == 0 {
		t.Fatalf("expected at least one stack relation, got 0")
	}

	for _, r := range relations {
		if r.ProfileID != 1 {
			t.Errorf("expected ProfileID to be 1, got %d", r.ProfileID)
		}
		if r.StackName == "" {
			t.Errorf("expected StackName not to be empty")
		}
	}
}

func TestGetStackRelationByProfileIds(t *testing.T) {
	t.Parallel()

	tests.SeedProfiles(t, driver)
	tests.SeedStackRelationships(t, driver)

	repo := relation.NewNeo4jStackRelationRepository(driver)

	relations, err := repo.GetStackRelationByProfileIds(t.Context(), []int64{1, 2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(relations) == 0 {
		t.Fatalf("expected at least one stack relation, got 0")
	}

	for _, r := range relations {
		if r.StackName == "" {
			t.Errorf("expected StackName not to be empty")
		}
	}
}
