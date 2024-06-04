package update_userprofile

import "github.com/rs/zerolog/log"

type Repository interface {
	UpdateUserProfile(data *UserProfile) error
}

type UpdateUserProfileService struct {
	repository Repository
}

type UserProfile struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
	Link     string `json:"link"`
}

func NewUpdateUserProfileService(repository Repository) *UpdateUserProfileService {
	return &UpdateUserProfileService{
		repository: repository,
	}
}

func (s *UpdateUserProfileService) UpdateUserProfile(userPorfile *UserProfile) error {
	err := s.repository.UpdateUserProfile(userPorfile)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error updating userprofile for username %s", userPorfile.Username)
		return err
	}

	return nil
}
