package update_userprofile

import (
	"userservice/internal/bus"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	UpdateUserProfile(data *UserProfile, bus *bus.EventBus) error
	GetPresignedUrlForUploading(userProfileImage *UserProfileImage) (string, error)
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

type UserProfileImage struct {
	Username string `json:"username"`
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

func (s *UpdateUserProfileService) UpdateUserProfileImage(userProfileImage *UserProfileImage) (string, error) {
	chError := make(chan error, 2)
	chResult := make(chan string, 1)

	go s.GeneratePreSignedUrl(userProfileImage, chResult, chError)

	numberOfTasks := 1
	for i := 0; i < numberOfTasks; i++ {
		err := <-chError
		if err != nil {
			return "", err
		}
	}

	result := <-chResult
	log.Info().Msgf("PresignedUrl for the UserProfileImage of %s was created", userProfileImage.Username)
	return result, nil
}

func (s *UpdateUserProfileService) GeneratePreSignedUrl(userProfileImage *UserProfileImage, chResult chan<- string, chError chan<- error) {
	presignedUrl, err := s.repository.GetPresignedUrlForUploading(userProfileImage)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error generating Pre-Signed URL")
		chError <- err
	}

	chError <- nil
	chResult <- presignedUrl
}
