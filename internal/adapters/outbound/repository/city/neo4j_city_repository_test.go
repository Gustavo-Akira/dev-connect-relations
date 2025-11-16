package city

import (
	"context"
	domain "devconnectrelations/internal/domain/city"
	"devconnectrelations/internal/tests"
	"testing"
)

func TestSaveNewCityEntityWithValidInputs(t *testing.T) {

	driver, cleanup := tests.SetupNeo4j(t)
	defer cleanup()
	repo := NewNeo4jCityRepository(driver)

	city := domain.City{
		Name:    "San Francisco",
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

	driver, cleanup := tests.SetupNeo4j(t)
	defer cleanup()
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

	driver, cleanup := tests.SetupNeo4j(t)
	defer cleanup()
	repo := NewNeo4jCityRepository(driver)
	_, err := repo.GetCityByFullName(context.Background(), "NonExistingCity, XX, YY")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

}
