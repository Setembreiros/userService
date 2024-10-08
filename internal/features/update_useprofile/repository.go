package update_userprofile

import (
	"strings"
	"userservice/internal/bus"
	database "userservice/internal/db"
	objectstorage "userservice/internal/objectStorage"

	"github.com/rs/zerolog/log"
)

type UpdateUserProfileRepository struct {
	dataRepository   *database.Database
	objectRepository *objectstorage.ObjectStorage
}

type UserProfileUpdatedEvent struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Link     string `json:"link"`
	FullName string `json:"full_name"`
}

type UserProfileImageMetadata struct {
	UserProfileImageId string `json:"user_profile_image_id"`
	Username           string `json:"username"`
	FileType           string `json:"file_type"`
}

func NewUpdateUserProfileRepository(dataRepository *database.Database, objectRepository *objectstorage.ObjectStorage) *UpdateUserProfileRepository {
	return &UpdateUserProfileRepository{
		dataRepository:   dataRepository,
		objectRepository: objectRepository,
	}
}

func (r UpdateUserProfileRepository) UpdateUserProfile(data *UserProfile, bus *bus.EventBus) error {
	tx, err := r.dataRepository.Client.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	result, err := tx.Exec(`
		WITH user_cte AS (
			SELECT id
			FROM userservice.users
			WHERE username = $1
		)
		UPDATE userservice.user_profiles
		SET full_name = $2, 
			bio = $3, 
			link = $4
		WHERE user_id = (SELECT id FROM user_cte);
		`,
		data.Username, data.FullName, data.Bio, data.Link)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Update user profile failed")
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		err = database.NewNotFoundError("userservice.users", data.Username)
		log.Error().Stack().Err(err).Msg("Update user profile failed")
		return err
	}

	err = publishUserProfileUpdated(data, bus)
	if err != nil {
		return err
	}

	return nil
}

func publishUserProfileUpdated(data *UserProfile, bus *bus.EventBus) error {
	userProfileUpdatedEvent := &UserProfileUpdatedEvent{
		Username: data.Username,
		Bio:      data.Bio,
		Link:     data.Link,
		FullName: data.FullName,
	}

	err := bus.Publish("UserProfileUpdatedEvent", userProfileUpdatedEvent)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Publishing UserProfileUpdatedEvent failed")
		return err
	}

	return nil
}

func (r UpdateUserProfileRepository) SaveUserProfileImageMetaData(userProfileImage *UserProfileImage) error {
	data := &UserProfileImageMetadata{
		UserProfileImageId: generateUserProfileImageId(userProfileImage),
		Username:           userProfileImage.Username,
		FileType:           userProfileImage.FileType,
	}
	tx, err := r.dataRepository.Client.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	result, err := tx.Exec(`
		WITH user_cte AS (
			SELECT id
			FROM userservice.users
			WHERE username = $1
		)
		UPDATE userservice.user_profiles
		SET profile_image_file_type = $2
		WHERE user_id = (SELECT id FROM user_cte);
		`,
		data.Username, data.FileType)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Update user profile failed")
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		err = database.NewNotFoundError("userservice.users", data.Username)
		log.Error().Stack().Err(err).Msg("Update user profile failed")
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *UpdateUserProfileRepository) GetPresignedUrlForUploading(userProfileImage *UserProfileImage) (string, error) {
	key := userProfileImage.Username + "/" + userProfileImage.FileType + "/" + generateUserProfileImageId(userProfileImage) + "." + userProfileImage.FileType
	return r.objectRepository.Client.GetPreSignedUrlForPuttingObject(key)
}

func generateUserProfileImageId(userProfileImage *UserProfileImage) string {
	userProfileImageId := userProfileImage.Username + "-profile_image-" + userProfileImage.FileType
	return strings.ReplaceAll(strings.ReplaceAll(userProfileImageId, " ", "_"), "\t", "_")
}
