package service

import (
	"context"
	"devconnectrelations/internal/domain/entities"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProfileRepository struct {
	mock.Mock
}

func (m *MockProfileRepository) CreateProfile(ctx context.Context, profile *entities.Profile) (entities.Profile, error) {
	args := m.Called(ctx, profile)

	return args.Get(0).(entities.Profile), args.Error(1)
}

func (m *MockProfileRepository) DeleteProfile(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func TestProfileService_CreateProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("should create profile successfully", func(t *testing.T) {
		mockRepo := new(MockProfileRepository)
		profileService := CreateNewProfileService(mockRepo)

		inputProfile := &entities.Profile{
			ConnectId: 311,
			Name:      "Software Developer",
		}

		expectedProfile := entities.Profile{
			ConnectId: 311,
			Name:      "Software Developer",
		}

		mockRepo.On("CreateProfile", ctx, inputProfile).Return(expectedProfile, nil).Once()

		createdProfile, err := profileService.CreateProfile(ctx, inputProfile)

		assert.NoError(t, err)
		assert.Equal(t, expectedProfile, createdProfile)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := new(MockProfileRepository)
		profileService := CreateNewProfileService(mockRepo)

		inputProfile := &entities.Profile{
			ConnectId: 311,
			Name:      "Software Developer",
		}

		expectedError := errors.New("database error")

		mockRepo.On("CreateProfile", ctx, inputProfile).Return(entities.Profile{}, expectedError).Once()

		createdProfile, err := profileService.CreateProfile(ctx, inputProfile)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, createdProfile)

		mockRepo.AssertExpectations(t)
	})
}

func TestProfileServic_DeleteProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete with success when profile is ok", func(t *testing.T) {
		mockRepo := new(MockProfileRepository)
		service := CreateNewProfileService(mockRepo)
		mockRepo.On("DeleteProfile", ctx, int32(1)).Return(nil)
		err := service.DeleteProfile(ctx, int32(1))
		assert.Empty(t, err)
	})

	t.Run("should return error when profile return error", func(t *testing.T) {
		mockRepo := new(MockProfileRepository)
		service := CreateNewProfileService(mockRepo)
		mockRepo.On("DeleteProfile", ctx, int32(1)).Return(errors.New("Some error"))
		err := service.DeleteProfile(ctx, int32(1))
		assert.Equal(t, err.Error(), "Some error")
	})

	t.Run("should return error when id is 0", func(t *testing.T) {
		mockRepo := new(MockProfileRepository)
		service := CreateNewProfileService(mockRepo)
		err := service.DeleteProfile(ctx, int32(0))
		assert.Equal(t, err.Error(), "id to delete cannot be 0")
	})
}
