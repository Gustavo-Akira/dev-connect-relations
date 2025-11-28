package relation

import (
	"context"
	"devconnectrelations/internal/domain/profile_relation/city"
	"devconnectrelations/internal/domain/recommendation"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jRelationCityRepository struct {
	driver *neo4j.DriverWithContext
}

func NewNeo4jRelationCityRepository(driver *neo4j.DriverWithContext) *Neo4jRelationCityRepository {
	return &Neo4jRelationCityRepository{
		driver: driver,
	}
}

func (r *Neo4jRelationCityRepository) CreateCityRelation(ctx context.Context, city *city.CityRelation) (*city.CityRelation, error) {
	_, err := neo4j.ExecuteQuery(ctx, *r.driver, `MATCH (c:City {full_name: $cityFullName}),(p:Profile {id: $profileID})
	MERGE (p)-[:LIVES_IN]->(c)`, map[string]any{
		"cityFullName": city.CityFullName,
		"profileID":    city.ProfileID,
	}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}

	return city, nil
}

func (r *Neo4jRelationCityRepository) GetCityRelatedToProfileId(ctx context.Context, profileId int64) ([]city.CityRelation, error) {
	result, err := neo4j.ExecuteQuery(ctx, *r.driver, `MATCH (p:Profile {id: $id})-[:LIVES_IN]->(c:City) RETURN c.full_name AS cityFullName`, map[string]any{"id": profileId}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	records := result.Records
	results := make([]city.CityRelation, 0)
	for _, record := range records {
		city := city.CityRelation{
			CityFullName: record.Values[0].(string),
			ProfileID:    profileId,
		}
		results = append(results, city)
	}
	return results, nil
}

func (r *Neo4jRelationCityRepository) GetCityRelatedToProfileIds(ctx context.Context, profileIds []int64) ([]city.CityRelation, error) {
	result, err := neo4j.ExecuteQuery(ctx, *r.driver, `MATCH (p:Profile)-[:LIVES_IN]->(c:City)
WHERE p.id IN $ids
RETURN p.id AS profileId, c.full_name AS city`, map[string]any{"ids": profileIds}, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	records := result.Records
	results := make([]city.CityRelation, 0)
	for _, record := range records {
		relation := city.CityRelation{
			CityFullName: record.Values[1].(string),
			ProfileID:    record.Values[0].(int64),
		}
		results = append(results, relation)
	}
	return results, nil
}

func (r *Neo4jRelationCityRepository) JaccardIndexByProfileId(ctx context.Context, profileID int64) ([]recommendation.Recommendation, error) {
	params := map[string]any{
		"id": profileID,
	}
	query := `MATCH (p1:Profile {id: $id})-[:LIVES_IN]->(c:City)
WITH p1, COLLECT(c.full_name) AS cities_p1

MATCH (p2:Profile)-[:LIVES_IN]->(c2:City)
WHERE NOT EXISTS {
  MATCH (p1)-[r:Relation]-(p2)
  WHERE r.type IN ["FRIEND", "BLOCKED"]
} AND p1<> p2
WITH p1, p2, cities_p1, COLLECT(c2.full_name) AS cities_p2

WITH 
    p1,
    p2,
    cities_p1,
    cities_p2,
    [x IN cities_p1 WHERE x IN cities_p2] AS inter,
    (cities_p1 + [x IN cities_p2 WHERE NOT x IN cities_p1]) AS uni

// evita divisão por zero
WHERE SIZE(uni) > 0

RETURN 
    p2.id AS recommended_profile,
    (SIZE(inter) * 1.0 / SIZE(uni)) AS jaccard_city,
	p2.name AS profile_name
ORDER BY jaccard_city DESC
LIMIT 20`
	result, err := neo4j.ExecuteQuery(ctx, *r.driver, query, params, neo4j.EagerResultTransformer)
	if err != nil {
		return nil, err
	}
	records := result.Records
	jaccardIndices := make([]recommendation.Recommendation, 0, len(records))
	for _, record := range records {
		jaccardIndex := record.Values[1].(float64)
		profileID := record.Values[0].(int64)
		name := record.Values[2].(string)
		jaccardInde := recommendation.Recommendation{
			ID:    profileID,
			Score: jaccardIndex,
			Name:  name,
		}
		jaccardIndices = append(jaccardIndices, jaccardInde)
	}
	return jaccardIndices, nil
}
