package search_user_test

import (
	"errors"
	"fmt"
	"testing"

	"userservice/internal/features/search_user"
	mock_search_user "userservice/internal/features/search_user/test/mock"
	"userservice/internal/model"

	"github.com/stretchr/testify/assert"
)

var serviceRepository *mock_search_user.MockRepository
var userProfileService *search_user.SearchUserService

func setUpService(t *testing.T) {
	SetUp(t)
	serviceRepository = mock_search_user.NewMockRepository(ctrl)
	userProfileService = search_user.NewSearchUserService(serviceRepository)
}

func TestGetSearchUserProfileSnippetsWithService(t *testing.T) {
	setUpService(t)
	query := "bo"
	limit := 5
	expectedUsers := []*model.UserProfileSnippet{
		{
			Username: "bobusername1",
			Name:     "alice1",
		},
		{
			Username: "aliceusername1",
			Name:     "bob2",
		},
		{
			Username: "aliceAndBobusername3",
			Name:     "aliceAndBob3",
		},
	}
	serviceRepository.EXPECT().SearchUserProfileSnippets(query, limit).Return(expectedUsers, nil)

	users, err := userProfileService.SearchUserProfileSnippets(query, limit)
	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedUsers, users)
}

func TestErrorOnSearchUserProfileSnippetsWithService(t *testing.T) {
	setUpService(t)
	query := "bo"
	limit := 5
	expectedUsers := []*model.UserProfileSnippet{}
	serviceRepository.EXPECT().SearchUserProfileSnippets(query, limit).Return(expectedUsers, errors.New("some error"))

	users, err := userProfileService.SearchUserProfileSnippets(query, limit)

	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Error getting userprofile snippets for query %s with limit %d", query, limit))
	assert.NotNil(t, err)
	assert.ElementsMatch(t, expectedUsers, users)
}
