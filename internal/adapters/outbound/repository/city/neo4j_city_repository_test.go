package city

import (
	"context"
	domain "devconnectrelations/internal/domain/city"
	"devconnectrelations/internal/tests"
	"os"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	driver neo4j.DriverWithContext
)

func TestMain(m *testing.M) {
	d, cleanup := tests.SetupNeo4j(&testing.T{})
	driver = d
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func TestSaveNewCityEntityWithValidInputs(t *testing.T) {
	t.Parallel()
	repo := NewNeo4jCityRepository(driver)

	city := domain.City{
		Name:    "San Francisco 1",
		State:   "CA",
		Country: "USA",
	}

	createdCity, err := repo.CreateCity(context.Background(), city)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if createdCity.Name != city.Name || createdCity.State != city.State || createdCity.Country != city.Country {
		t.Fatalf("expected city %v, got %v", city, createdCity)
	}
}

func TestGetCityByFullNameWithExistingCity(t *testing.T) {
	t.Parallel()
	repo := NewNeo4jCityRepository(driver)

	city := domain.City{
		Name:    "Los Angeles",
		State:   "CA",
		Country: "USA",
	}

	_, err := repo.CreateCity(context.Background(), city)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedCity, err := repo.GetCityByFullName(context.Background(), city.GetFullName())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrievedCity.Name != city.Name || retrievedCity.State != city.State || retrievedCity.Country != city.Country {
		t.Fatalf("expected city %v, got %v", city, retrievedCity)
	}

}

func TestGetCityByFullNameWithNonExistingCity(t *testing.T) {
	t.Parallel()
	repo := NewNeo4jCityRepository(driver)
	_, err := repo.GetCityByFullName(context.Background(), "NonExistingCity, XX, YY")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

}
