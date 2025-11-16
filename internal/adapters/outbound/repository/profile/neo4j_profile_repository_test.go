package profile_test

import (
	"context"
	"testing"

	profileRepo "devconnectrelations/internal/adapters/outbound/repository/profile"
	"devconnectrelations/internal/domain/profile"
	"devconnectrelations/internal/tests"

	"github.com/stretchr/testify/require"
)

func TestCreateAndGetProfile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	driver, cleanup := tests.SetupNeo4j(t)

	defer cleanup()

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

	driver, cleanup := tests.SetupNeo4j(t)

	defer cleanup()

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
