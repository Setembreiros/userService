package get_userprofile_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"userservice/internal/bus"
	mock_bus "userservice/internal/bus/mock"
	database "userservice/internal/db"
	get_userprofile "userservice/internal/features/get_userprofile"
	mock_get_userprofile "userservice/internal/features/get_userprofile/mock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var controllerExternalBus *mock_bus.MockExternalBus
var controllerBus *bus.EventBus

func TestGetUserProfileController_GetUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_get_userprofile.NewMockRepository(ctrl)
	controllerExternalBus = mock_bus.NewMockExternalBus(ctrl)
	controllerBus = bus.NewEventBus(controllerExternalBus)
	controller := get_userprofile.NewGetUserProfileController(mockRepository, controllerBus)

	gin.SetMode(gin.TestMode)

	username := "testuser"
	expectedUserProfile := &get_userprofile.UserProfile{
		Username: "testuser",
		FullName: "Test User",
		Bio:      "Test Bio",
		Link:     "https://test.com",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepository.EXPECT().GetUserProfile(username).Return(expectedUserProfile, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "username", Value: username}}

		controller.GetUserProfile(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		expectedError := database.NewNotFoundError("userservice.users", username)
		mockRepository.EXPECT().GetUserProfile(username).Return(nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "username", Value: username}}

		controller.GetUserProfile(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		expectedError := errors.New("repository error")
		mockRepository.EXPECT().GetUserProfile(username).Return(nil, expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "username", Value: username}}

		controller.GetUserProfile(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
func TestGetUserProfileController_GetUserProfileImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mock_get_userprofile.NewMockRepository(ctrl)
	controllerExternalBus = mock_bus.NewMockExternalBus(ctrl)
	controllerBus = bus.NewEventBus(controllerExternalBus)
	controller := get_userprofile.NewGetUserProfileController(mockRepository, controllerBus)

	gin.SetMode(gin.TestMode)

	username := "testuser"
	expectedUserProfileImage := "https://test.com/image.jpg"

	t.Run("Success", func(t *testing.T) {
		mockRepository.EXPECT().GetPresignedUrlForDownloading(username).Return(expectedUserProfileImage, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "username", Value: username}}

		controller.GetUserProfileImage(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		expectedError := database.NewNotFoundError("userservice.users", username)
		mockRepository.EXPECT().GetPresignedUrlForDownloading(username).Return("", expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "username", Value: username}}

		controller.GetUserProfileImage(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		expectedError := errors.New("repository error")
		mockRepository.EXPECT().GetPresignedUrlForDownloading(username).Return("", expectedError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "username", Value: username}}

		controller.GetUserProfileImage(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
