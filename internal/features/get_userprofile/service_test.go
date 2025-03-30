package get_userprofile_test

import (
	"bytes"
	"errors"
	"testing"
	"userservice/internal/bus"
	mock_bus "userservice/internal/bus/mock"
	get_userprofile "userservice/internal/features/get_userprofile"
	mock_get_userprofile "userservice/internal/features/get_userprofile/mock"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var serviceExternalBus *mock_bus.MockExternalBus
var serviceBus *bus.EventBus

var serviceLoggerOutput bytes.Buffer
var mockRepository *mock_get_userprofile.MockRepository
var service *get_userprofile.GetUserProfileService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceLoggerOutput.Truncate(0)
	log.Logger = log.Output(&serviceLoggerOutput)

	mockRepository = mock_get_userprofile.NewMockRepository(ctrl)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus = bus.NewEventBus(controllerExternalBus)
	service = get_userprofile.NewGetUserProfileService(mockRepository, serviceBus)
}

func TestGetUserProfileService_GetUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepository := mock_get_userprofile.NewMockRepository(ctrl)
	serviceExternalBus = mock_bus.NewMockExternalBus(ctrl)
	serviceBus = bus.NewEventBus(controllerExternalBus)
	service := get_userprofile.NewGetUserProfileService(mockRepository, serviceBus)

	username := "testuser"
	expectedUserProfile := &get_userprofile.UserProfile{
		Username: "testuser",
		FullName: "Test User",
		Bio:      "Test Bio",
		Link:     "https://test.com",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepository.EXPECT().GetUserProfile(username).Return(expectedUserProfile, nil)

		userProfile, err := service.GetUserProfile(username)

		assert.NoError(t, err)
		assert.Equal(t, expectedUserProfile, userProfile)
	})

	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("repository error")
		mockRepository.EXPECT().GetUserProfile(username).Return(nil, expectedError)

		userProfile, err := service.GetUserProfile(username)

		assert.Error(t, err)
		assert.Nil(t, userProfile)
		assert.Equal(t, expectedError, err)
	})
}
