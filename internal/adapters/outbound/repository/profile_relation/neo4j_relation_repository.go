package relation

import (
	"context"
	domain "devconnectrelations/internal/domain/profile_relation/relation"
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
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (fromPerson:Profile {id: $fromId}), (toPerson:Profile {id:$targetId}) MERGE (fromPerson)-[r:Relation{type:$type, status:$status}]->(toPerson) RETURN fromPerson", params, neo4j.EagerResultTransformer)
	if err != nil {
		return domain.Relation{}, err
	}
	println(result.Records)
	return relation, nil
}

func (r *Neo4jRelationRepository) GetAllRelationsByFromId(ctx context.Context, fromId int64) ([]domain.Relation, error) {
	params := map[string]any{
		"fromId": fromId,
	}
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (fromPerson:Profile {id: $fromId})-[r:Relation]-(toPerson:Profile) RETURN r, toPerson", params, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	relations := make([]domain.Relation, 0)
	print(result.Records)
	for _, record := range result.Records {
		relationNode, _ := record.Get("r")
		toPersonNode, _ := record.Get("toPerson")
		relationProps := relationNode.(neo4j.Relationship).Props
		toPersonProps := toPersonNode.(neo4j.Node).Props
		relation := domain.Relation{
			FromID: fromId,
			ToID:   toPersonProps["id"].(int64),
			Type:   domain.RelationType(relationProps["type"].(string)),
			Status: domain.RelationStatus(relationProps["status"].(string)),
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
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (fromPerson:Profile)-[r:Relation {status: 'PENDING'}]->(toPerson:Profile{id: $fromId}) RETURN r, fromPerson", params, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	relations := make([]domain.Relation, 0)
	print(result.Records)
	for _, record := range result.Records {
		relationNode, _ := record.Get("r")
		fromPersonNode, _ := record.Get("fromPerson")
		relationProps := relationNode.(neo4j.Relationship).Props
		fromPersonProps := fromPersonNode.(neo4j.Node).Props
		relation := domain.Relation{
			FromID: fromPersonProps["id"].(int64),
			ToID:   fromId,
			Type:   domain.RelationType(relationProps["type"].(string)),
			Status: domain.RelationStatus(relationProps["status"].(string)),
		}
		relations = append(relations, relation)
	}
	return relations, nil
}
