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
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (fromPerson:Profile {id: $fromId}), (toPerson:Profile {id:$targetId}) MERGE (fromPerson)-[r:Relation{type:$type, status:$status}]->(toPerson) RETURN fromPerson", params, neo4j.EagerResultTransformer)
	if err != nil {
		return entities.Relation{}, err
	}
	println(result.Records)
	return relation, nil
}

func (r *Neo4jRelationRepository) GetAllRelationsByFromId(ctx context.Context, fromId int64) ([]entities.Relation, error) {
	params := map[string]any{
		"fromId": fromId,
	}
	result, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (fromPerson:Profile {id: $fromId})-[r:Relation]-(toPerson:Profile) RETURN r, toPerson", params, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	relations := make([]entities.Relation, 0)
	print(result.Records)
	for _, record := range result.Records {
		relationNode, _ := record.Get("r")
		toPersonNode, _ := record.Get("toPerson")
		relationProps := relationNode.(neo4j.Relationship).Props
		toPersonProps := toPersonNode.(neo4j.Node).Props
		relation := entities.Relation{
			FromID: fromId,
			ToID:   toPersonProps["id"].(int64),
			Type:   entities.RelationType(relationProps["type"].(string)),
			Status: entities.RelationStatus(relationProps["status"].(string)),
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

	_, err := neo4j.ExecuteQuery(ctx, r.driver, "MATCH (fromPerson:Profile {id: $fromId})-[r:Relation]->(toPerson:Profile {id:$toId}) SET r.status = 'ACCEPTED' RETURN r", params, neo4j.EagerResultTransformer)
	if err != nil {
		return err
	}
	return nil
}
