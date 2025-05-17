package search_user

import (
	"userservice/internal/model"

	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=service.go -destination=test/mock/service.go

type Repository interface {
	SearchUserProfileSnippets(query, lastUsername string, limit int) ([]*model.UserProfileSnippet, string, error)
}

type SearchUserService struct {
	repository Repository
}

func NewSearchUserService(repository Repository) *SearchUserService {
	return &SearchUserService{
		repository: repository,
	}
}

func (s *SearchUserService) SearchUserProfileSnippets(query, lastUsername string, limit int) ([]*model.UserProfileSnippet, string, error) {
	users, lastUsername, err := s.repository.SearchUserProfileSnippets(query, lastUsername, limit)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Error getting userprofile snippets for query %s with lastusername %s and limit %d", query, lastUsername, limit)
		return users, lastUsername, err
	}

	return users, lastUsername, err
}
