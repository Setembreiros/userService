package update_userprofile

import (
	"userservice/internal/bus"
	database "userservice/internal/db"

	"github.com/rs/zerolog/log"
)

type UpdateUserProfileRepository database.Database

type UserProfileUpdatedEvent struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Link     string `json:"link"`
	FullName string `json:"full_name"`
}

func (r UpdateUserProfileRepository) UpdateUserProfile(data *UserProfile, bus *bus.EventBus) error {
	tx, err := r.Client.Begin()
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
