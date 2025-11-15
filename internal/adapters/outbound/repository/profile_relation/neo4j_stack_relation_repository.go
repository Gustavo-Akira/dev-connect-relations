package relation

import (
	"context"
	"devconnectrelations/internal/domain/profile_relation/stack"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4JStackRelationRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jStackRelationRepository(driver neo4j.DriverWithContext) *Neo4JStackRelationRepository {
	return &Neo4JStackRelationRepository{driver: driver}
}

func (r *Neo4JStackRelationRepository) CreateStackRelation(ctx context.Context, stackRelation *stack.StackRelation) (*stack.StackRelation, error) {
	params := map[string]any{
		"stackName": stackRelation.StackName,
		"profileID": stackRelation.ProfileID,
	}
	// retornamos a relação criada para validar
	st, err := neo4j.ExecuteQuery(ctx, r.driver,
		`MATCH (s:Stack {name: $stackName}), (p:Profile {id: $profileID})
         MERGE (p)-[r:USES]->(s)
         RETURN r, p, s`, params, neo4j.EagerResultTransformer)
	fmt.Println("CreateStackRelation cypher params:", params)
	fmt.Println("CreateStackRelation result:", st)
	if err != nil {
		return nil, err
	}
	return stackRelation, nil
}

func (r *Neo4JStackRelationRepository) DeleteStackRelation(ctx context.Context, stackName string, profileID int64) error {
	params := map[string]any{
		"stackName": stackName,
		"profileID": profileID,
	}
	_, err := neo4j.ExecuteQuery(ctx, r.driver,
		`MATCH (p:Profile {id: $profileID})-[r:USES]->(s:Stack {name: $stackName})
		 DELETE r`, params, neo4j.EagerResultTransformer)
	return err
}
