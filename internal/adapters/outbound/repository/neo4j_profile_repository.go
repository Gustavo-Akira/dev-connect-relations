package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jProfileRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jProfileRepository(driver neo4j.DriverWithContext) *Neo4jProfileRepository {
	return &Neo4jProfileRepository{driver: driver}
}

func (r *Neo4jProfileRepository) CreateProfile(ctx context.Context, profile *entities.Profile) (entities.Profile, error) {
	params := map[string]any{
		"id":   profile.ConnectId,
		"name": profile.Name,
	}
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "CREATE (p:Profile {id: $id, name: $name}) RETURN p.id AS id, p.name AS name", params, neo4j.EagerResultTransformer)
	if err != nil {
		return entities.Profile{}, err
	}
	println(result.Records)
	return *profile, nil
}
