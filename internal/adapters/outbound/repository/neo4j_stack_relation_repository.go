package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4JStackRelationRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jStackRelationRepository(driver neo4j.DriverWithContext) *Neo4JStackRelationRepository {
	return &Neo4JStackRelationRepository{driver: driver}
}

func (r *Neo4JStackRelationRepository) CreateStackRelation(ctx context.Context, stackRelation *entities.StackRelation) (*entities.StackRelation, error) {
	params := map[string]any{
		"stackName": stackRelation.StackName,
		"profileID": stackRelation.ProfileID,
	}
	_, err := neo4j.ExecuteQuery(ctx, r.driver,
		`MATCH (s:Stack {name: $stackName}), (p:Profile {id: $profileID})
		 Merge (p)-[:USES]->(s)`, params, neo4j.EagerResultTransformer)

	if err != nil {
		return nil, err
	}
	return stackRelation, nil
}
