package newuser_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	newuser "userservice/internal/features/new_user"
	mock_newuser "userservice/internal/features/new_user/mock"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var handlerLoggerOutput bytes.Buffer
var handlerRepository *mock_newuser.MockRepository
var handler *newuser.UserWasRegisteredEventHandler

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	handlerRepository = mock_newuser.NewMockRepository(ctrl)
	log.Logger = log.Output(&handlerLoggerOutput)
	handler = newuser.NewUserWasRegisteredEventHandler(handlerRepository)
}

func TestHandleUserWasRegisteredEventHandler(t *testing.T) {
	setUpHandler(t)
	data := &newuser.UserWasRegisteredEvent{
		Username: "username1",
		Email:    "user@email.com",
		UserType: "UA",
		Region:   "Vigo",
		FullName: "user lastname",
	}
	event, _ := json.Marshal(data)
	expectedUser := &newuser.User{
		Username: "username1",
		Email:    "user@email.com",
		UserType: "UA",
		Region:   "Vigo",
		UserProfile: &newuser.UserProfile{
			FullName: "user lastname",
			Bio:      "",
			Link:     "",
		},
	}
	handlerRepository.EXPECT().AddNewUser(expectedUser)

	handler.Handle(event)
}

func TestInvalidDataInUserWasRegisteredEventHandler(t *testing.T) {
	setUpHandler(t)
	invalidData := "invalid data"
	event, _ := json.Marshal(invalidData)

	handler.Handle(event)

	fmt.Println(handlerLoggerOutput.String())
	assert.Contains(t, handlerLoggerOutput.String(), "Invalid event data")
}
