package update_userprofile

import (
	"userservice/internal/bus"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	UpdateUserProfile(data *UserProfile, bus *bus.EventBus) error
}

type UpdateUserProfileService struct {
	repository Repository
	bus        *bus.EventBus
}

type UserProfile struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
	Link     string `json:"link"`
}

func NewUpdateUserProfileService(repository Repository, bus *bus.EventBus) *UpdateUserProfileService {
	return &UpdateUserProfileService{
		repository: repository,
		bus:        bus,
	}
}

func (s *UpdateUserProfileService) UpdateUserProfile(userPorfile *UserProfile) error {
	err := s.repository.UpdateUserProfile(userPorfile, s.bus)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error updating userprofile for username %s", userPorfile.Username)
		return err
	}

	log.Info().Msgf("User %s was updated", userPorfile.Username)
	return nil
}
