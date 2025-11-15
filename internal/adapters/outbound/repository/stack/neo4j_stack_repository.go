package stack

import (
	"context"
	domain "devconnectrelations/internal/domain/stack"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jStackRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jStackRepository(driver neo4j.DriverWithContext) *Neo4jStackRepository {
	return &Neo4jStackRepository{driver: driver}
}

func (r *Neo4jStackRepository) CreateStack(ctx context.Context, stack *domain.Stack) (domain.Stack, error) {
	params := map[string]any{"name": stack.Name}
	// use MERGE to avoid duplicates
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MERGE (s:Stack {name: $name}) RETURN s.name AS name", params, neo4j.EagerResultTransformer)
	if err != nil {
		return domain.Stack{}, err
	}
	fmt.Println("CreateStack result records:", result.Records)
	return *stack, nil
}

func (r *Neo4jStackRepository) GetStackByName(ctx context.Context, name string) (domain.Stack, error) {
	params := map[string]any{
		"name": name,
	}
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (s:Stack {name: $name}) RETURN s.name AS name", params, neo4j.EagerResultTransformer)
	if err != nil {
		return domain.Stack{}, err
	}
	records := result.Records
	println(records)
	if len(records) == 0 {
		return domain.Stack{}, nil
	}

	stackName, _ := records[0].Get("name")
	return domain.Stack{Name: stackName.(string)}, nil
}

func (r *Neo4jStackRepository) DeleteStack(ctx context.Context, name string) error {
	params := map[string]any{
		"name": name,
	}
	_, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (s:Stack {name: $name}) DETACH DELETE s", params, neo4j.EagerResultTransformer)
	return err
}
