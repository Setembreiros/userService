package get_userprofile

import (
	database "userservice/internal/db"

	"github.com/rs/zerolog/log"
)

type GetUserProfileRepository struct {
	dataRepository *database.Database
}

func NewGetUserProfileRepository(dataRepository *database.Database) *GetUserProfileRepository {
	return &GetUserProfileRepository{
		dataRepository: dataRepository,
	}
}

func (r *GetUserProfileRepository) GetUserProfile(username string) (*UserProfile, error) {
	var userProfile UserProfile
	err := r.dataRepository.Client.QueryRow(`
		SELECT u.username, up.full_name, up.bio, up.link
		FROM userservice.users u
		INNER JOIN userservice.user_profiles up ON u.id = up.user_id
		WHERE u.username = $1
	`, username).Scan(&userProfile.Username, &userProfile.FullName, &userProfile.Bio, &userProfile.Link)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting userprofile for username %s", username)
		return nil, database.NewNotFoundError("userservice.users", username)
	}

	return &userProfile, nil
}
