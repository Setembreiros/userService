package newuser_test

import (
	"bytes"
	"errors"
	"testing"

	newuser "userservice/internal/features/new_user"
	mock_newuser "userservice/internal/features/new_user/mock"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var serviceLoggerOutput bytes.Buffer
var serviceRepository *mock_newuser.MockRepository
var newUserService *newuser.NewUserRegisteredService

func setUpService(t *testing.T) {
	ctrl := gomock.NewController(t)
	serviceRepository = mock_newuser.NewMockRepository(ctrl)
	log.Logger = log.Output(&serviceLoggerOutput)
	newUserService = newuser.NewService(serviceRepository)
}

func TestCreateNewUserWithService(t *testing.T) {
	setUpService(t)
	data := &newuser.User{
		Username: "username1",
		Email:    "username@email.com",
		UserType: "UA",
		Region:   "Vigo",
		UserProfile: &newuser.UserProfile{
			FullName: "user name",
			Bio:      "",
			Link:     "",
		},
	}
	serviceRepository.EXPECT().AddNewUser(data)

	newUserService.CreateNewUser(data)

	assert.Contains(t, serviceLoggerOutput.String(), "User username1 was added")
}

func TestErrorOnCreateNewUserWithService(t *testing.T) {
	setUpService(t)
	data := &newuser.User{
		Username: "username1",
		Email:    "username@email.com",
		UserType: "UA",
		Region:   "Vigo",
		UserProfile: &newuser.UserProfile{
			FullName: "user name",
			Bio:      "",
			Link:     "",
		},
	}
	serviceRepository.EXPECT().AddNewUser(data).Return(errors.New("some error"))

	newUserService.CreateNewUser(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Error adding user")
}
