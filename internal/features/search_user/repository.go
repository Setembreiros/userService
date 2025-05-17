package search_user

import (
	"fmt"
	database "userservice/internal/db"
	"userservice/internal/model"

	"github.com/rs/zerolog/log"
)

type SearchUserRepository database.Database

func (r SearchUserRepository) SearchUserProfileSnippets(query string, lastUsername string, limit int) ([]*model.UserProfileSnippet, string, error) {
	baseQuery := `
		SELECT u.username, up.full_name
		FROM userservice.users u
		JOIN userservice.user_profiles up ON u.id = up.user_id
		WHERE (LOWER(u.username) LIKE LOWER($1) 
		   OR LOWER(up.full_name) LIKE LOWER($1) 
	`

	args := []interface{}{fmt.Sprintf("%%%s%%", query)}
	if lastUsername != "" {
		baseQuery += " AND u.username > $2"
		args = append(args, lastUsername)
	}

	baseQuery += " ORDER BY u.username ASC LIMIT $" + fmt.Sprintf("%d", len(args)+1)
	args = append(args, limit+1) // Solicitar un elemento extra para determinar se hai mais resultados

	rows, err := r.Client.Query(baseQuery, args...)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error searching user profiles")
		return nil, "", err
	}
	defer rows.Close()

	var users []*model.UserProfileSnippet
	var nextCursor string

	count := 0
	for rows.Next() {
		count++
		if count > limit {
			var username string
			if err := rows.Scan(&username, nil); err != nil {
				log.Error().Stack().Err(err).Msg("Error scanning next cursor")
				return nil, "", err
			}
			nextCursor = username
			break
		}

		user := &model.UserProfileSnippet{}
		if err := rows.Scan(&user.Username, &user.Name); err != nil {
			log.Error().Stack().Err(err).Msg("Error scanning user profile snippet")
			return nil, "", err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Error().Stack().Err(err).Msg("Error iterating over results")
		return nil, "", err
	}

	return users, nextCursor, nil
}
