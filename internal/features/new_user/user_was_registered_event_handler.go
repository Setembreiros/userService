package newuser

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type UserWasRegisteredEvent struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	Region   string `json:"region"`
	FullName string `json:"full_name"`
}

type UserWasRegisteredEventService interface {
	CreateNewUser(data *User)
}

type UserWasRegisteredEventHandler struct {
	service UserWasRegisteredEventService
}

func NewUserWasRegisteredEventHandler(repository Repository) *UserWasRegisteredEventHandler {
	return &UserWasRegisteredEventHandler{
		service: NewService(repository),
	}
}

func (handler *UserWasRegisteredEventHandler) Handle(event []byte) {
	var userWasRegisteredEvent UserWasRegisteredEvent
	log.Info().Msg("Handling UserWasRegisteredEvent")

	err := Decode(event, &userWasRegisteredEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Invalid event data")
		return
	}

	data := mapData(userWasRegisteredEvent)
	handler.service.CreateNewUser(data)
}

func mapData(event UserWasRegisteredEvent) *User {
	return &User{
		Username: event.Username,
		Email:    event.Email,
		UserType: event.UserType,
		Region:   event.Region,
		UserProfile: &UserProfile{
			FullName: event.FullName,
			Bio:      "",
			Link:     "",
		},
	}
}

func Decode(datab []byte, data any) error {
	return json.Unmarshal(datab, &data)
}
