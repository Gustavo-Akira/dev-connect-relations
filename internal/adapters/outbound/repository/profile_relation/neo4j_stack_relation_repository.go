package relation

import (
	"context"
	"devconnectrelations/internal/domain/profile_relation/stack"
	"devconnectrelations/internal/domain/recommendation"
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
func (r *Neo4JStackRelationRepository) GetStackRelationByProfileId(ctx context.Context, profileId int64) ([]stack.StackRelation, error) {
	st, err := neo4j.ExecuteQuery(ctx, r.driver, `MATCH (p:Profile {id:$id}) -[r:USES]->(s:Stack) RETURN s.name`, map[string]any{"id": profileId}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	records := st.Records
	result := make([]stack.StackRelation, 0)
	for _, record := range records {
		name := record.Values[0].(string)
		relation := stack.StackRelation{
			StackName: name,
			ProfileID: profileId,
		}
		result = append(result, relation)
	}
	return result, nil
}

func (r *Neo4JStackRelationRepository) GetStackRelationByProfileIds(ctx context.Context, profileIds []int64) (map[int64][]string, error) {
	st, err := neo4j.ExecuteQuery(ctx, r.driver, `MATCH (p:Profile) -[r:USES]->(s:Stack) WHERE p.id IN $ids RETURN s.name, p.id`, map[string]any{"ids": profileIds}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	records := st.Records
	result := make(map[int64][]string, 0)
	for _, record := range records {
		name := record.Values[0].(string)
		relation := stack.StackRelation{
			StackName: name,
			ProfileID: record.Values[1].(int64),
		}
		result[relation.ProfileID] = append(result[relation.ProfileID], relation.StackName)
	}
	return result, nil
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

func (r *Neo4JStackRelationRepository) JaccardIndexByProfileId(ctx context.Context, profileID int64) ([]recommendation.Recommendation, error) {
	params := map[string]any{
		"id": profileID,
	}
	query := `MATCH (p1:Profile {id:$id}) -[:USES]-> (s1:Stack) WITH p1,collect(s1.name) AS stacks_p1 
MATCH (p2: Profile) -[:USES]->(s2:Stack)
WHERE NOT EXISTS {
  MATCH (p1)-[r:Relation]-(p2)
  WHERE r.type IN ["FRIEND", "BLOCKED"]
} AND p1<> p2
WITH p2,collect(s2.name) AS stacks_p2,p1,stacks_p1

WITH p1,p2,stacks_p1,stacks_p2,
[x IN stacks_p1 WHERE x IN stacks_p2] AS inter,
(stacks_p1 + [x IN stacks_p2 WHERE NOT x IN stacks_p1]) as uni
RETURN 
    p2.id AS recommended_profile,
    (SIZE(inter) * 1.0 / SIZE(uni)) AS jaccard_stack,
	p2.name AS profile_name
ORDER BY jaccard_stack DESC
LIMIT 20`
	result, err := neo4j.ExecuteQuery(ctx, r.driver, query, params, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	jaccardIndices := make([]recommendation.Recommendation, 0)
	records := result.Records
	for _, record := range records {
		jaccardIndex := record.Values[1].(float64)
		recommendedProfileID := record.Values[0].(int64)
		name := record.Values[2].(string)
		jaccardIndices = append(jaccardIndices, recommendation.Recommendation{
			ID:    recommendedProfileID,
			Score: jaccardIndex,
			Name:  name,
		})
	}
	return jaccardIndices, nil
}
