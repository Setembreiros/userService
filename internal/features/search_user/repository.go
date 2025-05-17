package search_user

import (
	"fmt"
	database "userservice/internal/db"
	"userservice/internal/model"

	"github.com/rs/zerolog/log"
)

type SearchUserRepository database.Database

func (r SearchUserRepository) SearchUserProfileSnippets(query string, limit int) ([]*model.UserProfileSnippet, error) {
	baseQuery := `
		SELECT u.username, up.full_name
		FROM userservice.users u
		JOIN userservice.user_profiles up ON u.id = up.user_id
		WHERE (LOWER(u.username) LIKE LOWER($1) 
		   OR LOWER(up.full_name) LIKE LOWER($1))
		LIMIT $2
	`
	args := []interface{}{fmt.Sprintf("%%%s%%", query), limit}

	rows, err := r.Client.Query(baseQuery, args...)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error searching user profiles")
		return nil, err
	}
	defer rows.Close()

	var users []*model.UserProfileSnippet

	count := 0
	for rows.Next() {
		count++
		if count > limit {
			break
		}

		user := &model.UserProfileSnippet{}
		if err := rows.Scan(&user.Username, &user.Name); err != nil {
			log.Error().Stack().Err(err).Msg("Error scanning user profile snippet")
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Error().Stack().Err(err).Msg("Error iterating over results")
		return nil, err
	}

	return users, nil
}
