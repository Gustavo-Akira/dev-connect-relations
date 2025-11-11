package service

import (
	"context"
	"devconnectrelations/internal/domain/entities"
	"errors"
	"testing"
)

type MockCityRepository struct{}

func (m *MockCityRepository) CreateCity(ctx context.Context, city entities.City) (entities.City, error) {
	if city.Name == "ERROR" {
		return entities.City{}, errors.New("An error occurred")
	}
	return city, nil
}

func (m *MockCityRepository) GetCityByFullName(ctx context.Context, fullName string) (*entities.City, error) {
	return &entities.City{
		Name:    "MockCity",
		State:   "MC",
		Country: "MockCountry",
	}, nil
}

func TestCityService_CreateCity(t *testing.T) {
	mockRepo := &MockCityRepository{}
	cityService := NewCityService(mockRepo)
	t.Run("should create city", func(t *testing.T) {
		city := entities.City{
			Name:    "San Francisco",
			State:   "CA",
			Country: "USA",
		}

		createdCity, err := cityService.CreateCity(context.Background(), city)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if createdCity.Name != city.Name || createdCity.State != city.State || createdCity.Country != city.Country {
			t.Fatalf("expected city %v, got %v", city, createdCity)
		}
	})
	t.Run("should return error when repository fails", func(t *testing.T) {
		city := entities.City{
			Name:    "ERROR",
			State:   "CA",
			Country: "USA",
		}
		_, err := cityService.CreateCity(context.Background(), city)
		if err == nil {
			t.Fatalf("expected error, got none")
		}
	})
}
