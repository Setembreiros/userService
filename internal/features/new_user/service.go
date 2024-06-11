package newuser

import "github.com/rs/zerolog/log"

type NewUserRegisteredService struct {
	repository Repository
}

type Repository interface {
	AddNewUser(data *User) error
}

type User struct {
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	UserType    string       `json:"user_type"`
	Region      string       `json:"region"`
	UserProfile *UserProfile `json:"user_profile"`
}

type UserProfile struct {
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
	Link     string `json:"link"`
}

func NewService(repository Repository) *NewUserRegisteredService {
	return &NewUserRegisteredService{
		repository: repository,
	}
}

func (s *NewUserRegisteredService) CreateNewUser(data *User) {
	err := s.repository.AddNewUser(data)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error adding user")
		return
	}

	log.Info().Msgf("User %s was added", data.Username)
}
