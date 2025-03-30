package get_userprofile

import (
	"userservice/internal/bus"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=mock/service.go

type Repository interface {
	GetUserProfile(username string) (*UserProfile, error)
	GetPresignedUrlForDownloading(username string) (string, error)
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

func (s *GetUserProfileService) GeneratePreSignedUrl(username string, chResult chan<- string, chError chan<- error) {
	presignedUrl, err := s.repository.GetPresignedUrlForDownloading(username)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error generating Pre-Signed URL")
		chError <- err
	}

	chError <- nil
	chResult <- presignedUrl
}

func (s *GetUserProfileService) GetUserProfileImage(username string) (string, error) {
	chError := make(chan error, 2)
	chResult := make(chan string, 1)

	go s.GeneratePreSignedUrl(username, chResult, chError)

	numberOfTasks := 1
	for i := 0; i < numberOfTasks; i++ {
		err := <-chError
		if err != nil {
			return "", err
		}
	}

	result := <-chResult
	log.Info().Msgf("PresignedUrl for the UserProfileImage of %s was created", username)
	return result, nil
}
