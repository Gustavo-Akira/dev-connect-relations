package tests

import (
	"context"
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
