package relation_test

import (
	relation "devconnectrelations/internal/adapters/outbound/repository/profile_relation"
	"devconnectrelations/internal/tests"
	"os"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	driver neo4j.DriverWithContext
)

func TestMain(m *testing.M) {
	test := &testing.T{}
	d, cleanup := tests.SetupNeo4j(test)
	driver = d
	tests.SeedProfiles(test, driver)
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func TestCreateCityRelation(t *testing.T) {
	t.Parallel()

	tests.SeedCityRelationships(t, driver)
}

func TestJaccardIndexByProfileId(t *testing.T) {
	t.Parallel()
	tests.SeedProfiles(t, driver)
	tests.SeedCityRelationships(t, driver)

	repo := relation.NewNeo4jRelationCityRepository(&driver)
	jaccardIndices, err := repo.JaccardIndexByProfileId(t.Context(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(jaccardIndices) == 0 {
		t.Fatalf("expected at least one jaccard index, got 0")
	}

}
