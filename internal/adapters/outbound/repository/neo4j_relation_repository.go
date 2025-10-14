package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jRelationRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jRelationRepository(driver neo4j.DriverWithContext) *Neo4jRelationRepository {
	return &Neo4jRelationRepository{driver: driver}
}

func (r *Neo4jRelationRepository) CreateRelation(ctx context.Context, relation entities.Relation) (entities.Relation, error) {
	params := map[string]any{
		"fromId":   relation.FromID,
		"targetId": relation.ToID,
		"type":     relation.Type,
		"status":   relation.Status,
	}
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (fromPerson:Person {id: $fromId}), (toPerson:Person {id:$targetId}) MERGE (fromPerson)-[r:Relation{type:$type, status:$status}]->(toPerson) RETURN fromPerson", params, neo4j.EagerResultTransformer)
	if err != nil {
		return entities.Relation{}, err
	}
	println(result.Records)
	return relation, nil
}
