package search_user

import (
	"userservice/internal/model"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	SearchUserProfileSnippets(query string, limit int) ([]*model.UserProfileSnippet, error)
}

type SearchUserService struct {
	repository Repository
}

func NewSearchUserService(repository Repository) *SearchUserService {
	return &SearchUserService{
		repository: repository,
	}
}

func (s *SearchUserService) SearchUserProfileSnippets(query string, limit int) ([]*model.UserProfileSnippet, error) {
	users, err := s.repository.SearchUserProfileSnippets(query, limit)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting userprofile snippets for query %s with limit %d", query, limit)
		return users, err
	}

	return users, err
}
