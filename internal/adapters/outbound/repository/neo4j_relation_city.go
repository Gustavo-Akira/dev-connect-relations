package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jRelationCityRepository struct {
	driver *neo4j.DriverWithContext
}

func NewNeo4jRelationCityRepository(driver *neo4j.DriverWithContext) *Neo4jRelationCityRepository {
	return &Neo4jRelationCityRepository{
		driver: driver,
	}
}

func (r *Neo4jRelationCityRepository) CreateCityRelation(ctx context.Context, city *entities.CityRelation) (*entities.CityRelation, error) {
	_, err := neo4j.ExecuteQuery(ctx, *r.driver, `MATCH (c:City {full_name: $cityFullName}),(p:Profile {id: $profileID})
	MERGE (p)-[:LIVES_IN]->(c)`, map[string]any{
		"cityFullName": city.CityFullName,
		"profileID":    city.ProfileID,
	}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	return city, nil
}
