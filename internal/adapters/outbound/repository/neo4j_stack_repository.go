package repository

import (
	"context"
	"devconnectrelations/internal/domain/entities"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jStackRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jStackRepository(driver neo4j.DriverWithContext) *Neo4jStackRepository {
	return &Neo4jStackRepository{driver: driver}
}

func (r *Neo4jStackRepository) CreateStack(ctx context.Context, stack *entities.Stack) (entities.Stack, error) {
	params := map[string]any{
		"name": stack.Name,
	}
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "CREATE (s:Stack {name: $name}) RETURN s.name AS name", params, neo4j.EagerResultTransformer)
	if err != nil {
		return entities.Stack{}, err
	}
	println(result.Records)
	return *stack, nil
}

func (r *Neo4jStackRepository) GetStackByName(ctx context.Context, name string) (entities.Stack, error) {
	params := map[string]any{
		"name": name,
	}
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (s:Stack {name: $name}) RETURN s.name AS name", params, neo4j.EagerResultTransformer)
	if err != nil {
		return entities.Stack{}, err
	}
	records := result.Records
	println(records)
	if len(records) == 0 {
		return entities.Stack{}, nil
	}

	stackName, _ := records[0].Get("name")
	return entities.Stack{Name: stackName.(string)}, nil
}

func (r *Neo4jStackRepository) DeleteStack(ctx context.Context, name string) error {
	params := map[string]any{
		"name": name,
	}
	_, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (s:Stack {name: $name}) DETACH DELETE s", params, neo4j.EagerResultTransformer)
	return err
}
