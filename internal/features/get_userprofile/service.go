package get_userprofile

import (
	"userservice/internal/bus"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	GetUserProfile(username string) (*UserProfile, error)
}

type GetUserProfileService struct {
	repository Repository
	bus        *bus.EventBus
}

type UserProfile struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
	Link     string `json:"link"`
}

func NewGetUserProfileService(repository Repository, bus *bus.EventBus) *GetUserProfileService {
	return &GetUserProfileService{
		repository: repository,
		bus:        bus,
	}
}

func (s *GetUserProfileService) GetUserProfile(username string) (*UserProfile, error) {
	userProfile, err := s.repository.GetUserProfile(username)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting userprofile for username %s", username)
		return nil, err
	}

	log.Info().Msgf("User %s was found", username)
	return userProfile, nil
}
