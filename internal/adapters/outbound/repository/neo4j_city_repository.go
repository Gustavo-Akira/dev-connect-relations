package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jCityRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jCityRepository(driver neo4j.DriverWithContext) *Neo4jCityRepository {
	return &Neo4jCityRepository{driver: driver}
}

func (r *Neo4jCityRepository) CreateCity(ctx context.Context, city entities.City) (entities.City, error) {
	params := map[string]any{
		"name":      city.Name,
		"state":     city.State,
		"country":   city.Country,
		"full_name": city.GetFullName(),
	}

	_, err := neo4j.ExecuteQuery(ctx, r.driver, "CREATE (c:City {name: $name, state: $state, country: $country, full_name: $full_name}) RETURN c", params, neo4j.EagerResultTransformer)
	if err != nil {
		return entities.City{}, err
	}
	return city, nil
}
