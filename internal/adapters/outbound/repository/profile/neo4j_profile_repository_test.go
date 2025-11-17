package profile_test

import (
	"context"
	"os"
	"testing"

	profileRepo "devconnectrelations/internal/adapters/outbound/repository/profile"
	"devconnectrelations/internal/domain/profile"
	"devconnectrelations/internal/tests"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/require"
)

var (
	driver neo4j.DriverWithContext
)

func TestMain(m *testing.M) {
	d, cleanup := tests.SetupNeo4j(&testing.T{})
	driver = d
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func TestCreateAndGetProfile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := profileRepo.NewNeo4jProfileRepository(driver)

	entity := &profile.Profile{
		ConnectId: 123,
		Name:      "Gustavo",
	}

	created, err := repo.CreateProfile(ctx, entity)
	require.NoError(t, err)
	require.Equal(t, entity.Name, created.Name)

	fetched, err := repo.GetProfileByID(ctx, entity.ConnectId)
	require.NoError(t, err)
	require.NotNil(t, fetched)

	require.Equal(t, entity.ConnectId, fetched.ConnectId)
	require.Equal(t, entity.Name, fetched.Name)
}

func TestDeleteProfile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := profileRepo.NewNeo4jProfileRepository(driver)

	entity := &profile.Profile{
		ConnectId: 999,
		Name:      "ToDelete",
	}

	_, err := repo.CreateProfile(ctx, entity)
	require.NoError(t, err)

	err = repo.DeleteProfile(ctx, entity.ConnectId)
	require.NoError(t, err)

	result, err := repo.GetProfileByID(ctx, entity.ConnectId)
	require.NoError(t, err)
	require.Nil(t, result)
}
