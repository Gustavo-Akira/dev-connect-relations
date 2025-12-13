package relation

import (
	"context"
	domain "devconnectrelations/internal/domain/profile_relation/relation"
	"devconnectrelations/internal/domain/recommendation"
	"errors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jRelationRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jRelationRepository(driver neo4j.DriverWithContext) *Neo4jRelationRepository {
	return &Neo4jRelationRepository{driver: driver}
}

func (r *Neo4jRelationRepository) CreateRelation(ctx context.Context, relation domain.Relation) (domain.Relation, error) {
	params := map[string]any{
		"fromId":   relation.FromID,
		"targetId": relation.ToID,
		"type":     relation.Type,
		"status":   relation.Status,
	}
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (fromPerson:Profile {id: $fromId}), (toPerson:Profile {id:$targetId}) OPTIONAL MATCH (fromPerson)-[r:Relation]-(toPerson) DELETE r MERGE (fromPerson)-[newR:Relation{type:$type, status:$status}]->(toPerson) RETURN fromPerson", params, neo4j.EagerResultTransformer)
	if err != nil {
		return domain.Relation{}, err
	}
	println(result.Records)
	return relation, nil
}

func (r *Neo4jRelationRepository) GetAllRelationsByFromId(ctx context.Context, fromId int64, offset int64, limit int64) ([]domain.Relation, error) {
	params := map[string]any{
		"fromId": fromId,
	}
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (fromPerson:Profile {id: $fromId})-[r:Relation]-(toPerson:Profile) RETURN r, toPerson, fromPerson", params, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	relations := make([]domain.Relation, 0)
	print(result.Records)
	for _, record := range result.Records {
		relationNode, _ := record.Get("r")
		toPersonNode, _ := record.Get("toPerson")
		fromPersonNode, _ := record.Get("fromPerson")
		fromPersonProps := fromPersonNode.(neo4j.Node).Props
		relationProps := relationNode.(neo4j.Relationship).Props
		toPersonProps := toPersonNode.(neo4j.Node).Props
		relation := domain.Relation{
			FromID:          fromId,
			FromProfileName: fromPersonProps["name"].(string),
			ToID:            toPersonProps["id"].(int64),
			ToProfileName:   toPersonProps["name"].(string),
			Type:            domain.RelationType(relationProps["type"].(string)),
			Status:          domain.RelationStatus(relationProps["status"].(string)),
		}
		relations = append(relations, relation)
	}
	return relations, nil
}

func (r Neo4jRelationRepository) AcceptRelation(ctx context.Context, fromId int64, toId int64) error {
	params := map[string]any{
		"fromId": fromId,
		"toId":   toId,
	}

	result, err := neo4j.ExecuteQuery(ctx, r.driver,
		"MATCH (fromPerson:Profile {id: $fromId})-[r:Relation {status: 'PENDING'}]->(toPerson:Profile {id: $toId}) SET r.status = 'ACCEPTED' RETURN r",
		params, neo4j.EagerResultTransformer)
	if err != nil {
		return err
	}
	if len(result.Records) == 0 {
		return errors.New("no pending relation found in the correct direction")
	}
	return nil
}

func (r *Neo4jRelationRepository) GetAllRelationPendingByFromId(ctx context.Context, fromId int64) ([]domain.Relation, error) {
	params := map[string]any{
		"fromId": fromId,
	}
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (fromPerson:Profile)-[r:Relation {status: 'PENDING'}]->(toPerson:Profile{id: $fromId}) RETURN r, fromPerson, toPerson", params, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	relations := make([]domain.Relation, 0)
	print(result.Records)
	for _, record := range result.Records {
		relationNode, _ := record.Get("r")
		fromPersonNode, _ := record.Get("fromPerson")
		toPersonNode, _ := record.Get("toPerson")
		relationProps := relationNode.(neo4j.Relationship).Props
		fromPersonProps := fromPersonNode.(neo4j.Node).Props
		toPersonProps := toPersonNode.(neo4j.Node).Props
		relation := domain.Relation{
			FromID:          fromPersonProps["id"].(int64),
			FromProfileName: fromPersonProps["name"].(string),
			ToID:            fromId,
			ToProfileName:   toPersonProps["name"].(string),
			Type:            domain.RelationType(relationProps["type"].(string)),
			Status:          domain.RelationStatus(relationProps["status"].(string)),
		}
		relations = append(relations, relation)
	}
	return relations, nil
}

func (r *Neo4jRelationRepository) JaccardIndexByProfileId(ctx context.Context, profileID int64) ([]recommendation.Recommendation, error) {
	params := map[string]any{
		"id": profileID,
	}

	query := `MATCH (p1:Profile {id: $id})-[r:Relation {type:"FRIEND", status:"ACCEPTED"}]-(p2:Profile)
WITH p1, COLLECT(p2.id) AS friends_p1

MATCH (p3:Profile)
WHERE p3.id <> p1.id
AND NOT EXISTS {
  MATCH (p1)-[r_check:Relation]-(p3)
  WHERE r_check.type IN ["FRIEND", "BLOCKED"]
}

MATCH (p3)-[r:Relation {type:"FRIEND", status:"ACCEPTED"}]-(p4:Profile)
WITH p1, p3, friends_p1, COLLECT(p4.id) AS friends_p3

WITH 
    p3,
    friends_p1,
    friends_p3,
    [x IN friends_p1 WHERE x IN friends_p3] AS inter,
    (friends_p1 + [x IN friends_p3 WHERE NOT x IN friends_p1]) AS uni

WHERE SIZE(uni) > 0 

RETURN 
    p3.id AS recommended_profile,
    (SIZE(inter) * 1.0 / SIZE(uni)) AS jaccard_friend,
	p3.name AS profile_name
ORDER BY jaccard_friend DESC
LIMIT 20`

	result, err := neo4j.ExecuteQuery(ctx, r.driver, query, params, neo4j.EagerResultTransformer)

	if err != nil {
		return nil, err
	}

	records := result.Records
	jaccardIndices := make([]recommendation.Recommendation, 0, len(records))
	for _, record := range records {
		jaccardIndex := record.Values[1].(float64)
		recommendedProfileID := record.Values[0].(int64)
		name := record.Values[2].(string)
		jaccardInde := recommendation.Recommendation{
			ID:    recommendedProfileID,
			Score: jaccardIndex,
			Name:  name,
		}
		jaccardIndices = append(jaccardIndices, jaccardInde)
	}
	return jaccardIndices, nil
}

func (r *Neo4jRelationRepository) CountRelationsByFromId(ctx context.Context, fromId int64) (int64, error) {
	params := map[string]any{
		"fromId": fromId,
	}
	query := "MATCH (n:Profile {id:$fromId}) RETURN COUNT { (n)-[:Relation]-() } AS total"
	result, err := neo4j.ExecuteQuery(ctx, r.driver, query, params, neo4j.EagerResultTransformer)
	if err != nil {
		return 0, err
	}

	v, ok := result.Records[0].Get("total")
	if !ok || v == nil {
		return 0, nil
	}
	return v.(int64), nil
}
