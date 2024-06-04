package update_userprofile

import (
	database "userservice/internal/db"

	"github.com/rs/zerolog/log"
)

type UpdateUserProfileRepository database.Database

func (r UpdateUserProfileRepository) UpdateUserProfile(data *UserProfile) error {
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

	return nil
}
