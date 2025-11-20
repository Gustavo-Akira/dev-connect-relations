package tests

import (
	"context"
	"devconnectrelations/internal/adapters/outbound/repository/city"
	"devconnectrelations/internal/adapters/outbound/repository/profile"
	relation "devconnectrelations/internal/adapters/outbound/repository/profile_relation"
	"devconnectrelations/internal/adapters/outbound/repository/stack"
	cityDomain "devconnectrelations/internal/domain/city"
	domain "devconnectrelations/internal/domain/profile"
	cityRelationDomain "devconnectrelations/internal/domain/profile_relation/city"
	stackRelationDomain "devconnectrelations/internal/domain/profile_relation/stack"
	stackDomain "devconnectrelations/internal/domain/stack"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/require"
	neo4jContainer "github.com/testcontainers/testcontainers-go/modules/neo4j"
)

func SetupNeo4j(t *testing.T) (neo4j.DriverWithContext, func()) {
	ctx := context.Background()

	container, err := neo4jContainer.Run(ctx,
		"neo4j:5.20",
		neo4jContainer.WithAdminPassword("Kabuto123"),
		neo4jContainer.WithLabsPlugin("apoc"),
	)
	require.NoError(t, err)

	uri, err := container.BoltUrl(ctx)
	require.NoError(t, err)

	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth("neo4j", "Kabuto123", ""))
	require.NoError(t, err)

	cleanup := func() {
		driver.Close(ctx)
		container.Terminate(ctx)
	}

	return driver, cleanup
}

func SeedProfiles(t *testing.T, driver neo4j.DriverWithContext) {
	ctx := context.Background()
	profile_repo := profile.NewNeo4jProfileRepository(driver)
	//Create a array and iterate it to create multiple profiles
	profiles := []domain.Profile{
		{ConnectId: 1, Name: "John Doe"},
		{ConnectId: 2, Name: "Jane Smith"},
		{ConnectId: 3, Name: "Alice Johnson"},
		{ConnectId: 4, Name: "Bob Brown"},
		{ConnectId: 5, Name: "Charlie Davis"},
	}

	for _, p := range profiles {
		_, err := profile_repo.CreateProfile(ctx, &p)
		require.NoError(t, err)
	}
}

func SeedCityRelationships(t *testing.T, driver neo4j.DriverWithContext) {
	ctx := context.Background()
	// Create cities and relationships with profile nodes based on profiles created on seed profiles
	city_repository := city.NewNeo4jCityRepository(driver)
	cities := []cityDomain.City{
		{Name: "New York", Country: "USA", State: "NY"},
		{Name: "Los Angeles", Country: "USA", State: "CA"},
		{Name: "Chicago", Country: "USA", State: "IL"},
	}

	for _, c := range cities {
		_, err := city_repository.CreateCity(ctx, c)
		require.NoError(t, err)
	}

	//Create cityRelationships
	city_relationship_repo := relation.NewNeo4jRelationCityRepository(&driver)
	profileIds := []int64{1, 2, 3, 4, 5}
	cityFullNames := []string{"New York,NY,USA", "Los Angeles,CA,USA", "Chicago,IL,USA"}

	for i, profileId := range profileIds {
		cityFullName := cityFullNames[i%len(cityFullNames)]
		cityRelation := cityRelationDomain.NewCityRelation(cityFullName, profileId)
		_, err := city_relationship_repo.CreateCityRelation(ctx, cityRelation)
		require.NoError(t, err)
	}
}

func SeedStackRelationships(t *testing.T, driver neo4j.DriverWithContext) {
	ctx := context.Background()
	stack_relation_repository := relation.NewNeo4jStackRelationRepository(driver)
	stack_repo := stack.NewNeo4jStackRepository(driver)
	stacks := []string{"go", "python", "javascript"}
	for _, stack_name := range stacks {
		create_stack, err := stackDomain.NewStack(stack_name)
		_, err = stack_repo.CreateStack(ctx, create_stack)
		require.NoError(t, err)
	}
	profileIds := []int64{1, 2, 3, 4, 5}
	// Now i want this to create similarity between profiles based on stacks
	for i, profileId := range profileIds {
		stackName := stacks[i%len(stacks)]
		stackRelation, _ := stackRelationDomain.NewStackRelation(stackName, profileId)
		_, err := stack_relation_repository.CreateStackRelation(ctx, stackRelation)
		require.NoError(t, err)
	}
}
