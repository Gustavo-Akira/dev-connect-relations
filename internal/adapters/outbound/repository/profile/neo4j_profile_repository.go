package profile

import (
	"context"
	domain "devconnectrelations/internal/domain/profile"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jProfileRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jProfileRepository(driver neo4j.DriverWithContext) *Neo4jProfileRepository {
	return &Neo4jProfileRepository{driver: driver}
}

func (r *Neo4jProfileRepository) CreateProfile(ctx context.Context, profile *domain.Profile) (domain.Profile, error) {
	params := map[string]any{
		"id":   profile.ConnectId,
		"name": profile.Name,
	}
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "CREATE (p:Profile {id: $id, name: $name}) RETURN p.id AS id, p.name AS name", params, neo4j.EagerResultTransformer)
	if err != nil {
		return domain.Profile{}, err
	}
	println(result.Records)
	return *profile, nil
}

func (r *Neo4jProfileRepository) DeleteProfile(ctx context.Context, id int64) error {
	params := map[string]any{
		"id": id,
	}
	_, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (p:Profile {id: $id}) DELETE p;", params, neo4j.EagerResultTransformer)
	return err
}

func (r *Neo4jProfileRepository) GetProfileByID(ctx context.Context, profileId int64) (*domain.Profile, error) {
	params := map[string]any{
		"id": profileId,
	}
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (p:Profile {id: $id}) RETURN p.id AS id, p.name AS name", params, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	records := result.Records
	if len(records) == 0 {
		return nil, nil
	}
	record := records[0]
	id, _ := record.Get("id")
	name, _ := record.Get("name")
	return domain.NewProfile(id.(int64), name.(string))
}
