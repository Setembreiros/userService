package newuser

import (
	database "userservice/internal/db"

	"github.com/rs/zerolog/log"
)

type NewUserRepository database.Database

func (r NewUserRepository) AddNewUser(data *User) error {
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

	var userID int

	err = tx.QueryRow("INSERT INTO userservice.users (username, email, region, user_type) VALUES ($1, $2, $3, $4) RETURNING id",
		data.Username, data.Email, data.Region, data.UserType).Scan(&userID)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Insert user failed")
		return err
	}

	_, err = tx.Exec("INSERT INTO userservice.user_profiles (user_id, full_name, bio, link) VALUES ($1, $2, $3, $4)",
		userID, data.UserProfile.FullName, data.UserProfile.Bio, data.UserProfile.Link)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Insert userProfile failed")
		return err
	}

	return nil
}
