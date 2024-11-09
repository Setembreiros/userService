package update_userprofile_test

import (
	"bytes"
	"errors"
	"testing"
	"userservice/internal/bus"
	update_userprofile "userservice/internal/features/update_useprofile"
	mock_update_userprofile "userservice/internal/features/update_useprofile/mock"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var serviceLoggerOutput bytes.Buffer
var serviceRepository *mock_update_userprofile.MockRepository
var serviceBus *bus.EventBus
var updateUserProfileService *update_userprofile.UpdateUserProfileService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceRepository = mock_update_userprofile.NewMockRepository(ctrl)
	serviceBus = &bus.EventBus{}
	log.Logger = log.Output(&serviceLoggerOutput)
	updateUserProfileService = update_userprofile.NewUpdateUserProfileService(serviceRepository, serviceBus)
}

func TestUpdateUserProfileWithService(t *testing.T) {
	setUpService(t)
	data := &update_userprofile.UserProfile{
		Username: "username1",
		FullName: "user name",
		Bio:      "",
		Link:     "",
	}
	serviceRepository.EXPECT().UpdateUserProfile(data, serviceBus)

	updateUserProfileService.UpdateUserProfile(data)

	assert.Contains(t, serviceLoggerOutput.String(), "User username1 was updated")
}

func TestErrorOnCreateNewUserProfileWithService(t *testing.T) {
	setUpService(t)
	data := &update_userprofile.UserProfile{
		Username: "username1",
		FullName: "user name",
		Bio:      "",
		Link:     "",
	}
	serviceRepository.EXPECT().UpdateUserProfile(data, serviceBus).Return(errors.New("some error"))

	updateUserProfileService.UpdateUserProfile(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Error updating userprofile for username username1")
}

func TestUpdateUserProfileImageWithService(t *testing.T) {
	setUpService(t)
	data := &update_userprofile.UserProfileImage{
		Username: "username1",
	}
	serviceRepository.EXPECT().GetPresignedUrlForUploading(data).Return("asdasd", nil)

	result, err := updateUserProfileService.UpdateUserProfileImage(data)

	assert.Contains(t, serviceLoggerOutput.String(), "PresignedUrl for the UserProfileImage of username1 was created")
	assert.NotEmpty(t, result)
	assert.Nil(t, err)
}

func TestErrorOnUpdateUserProfileImageWithService(t *testing.T) {
	setUpService(t)
	data := &update_userprofile.UserProfileImage{
		Username: "username1",
	}

	serviceRepository.EXPECT().GetPresignedUrlForUploading(data).Return("", errors.New("some error"))

	result, err := updateUserProfileService.UpdateUserProfileImage(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Error generating Pre-Signed URL")
	assert.Empty(t, result)
	assert.NotNil(t, err)
}
