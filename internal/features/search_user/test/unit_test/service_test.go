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
	lastUsername := "bobusername0"
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
	expectedLastUsername := "aliceAndBobusername3"
	serviceRepository.EXPECT().SearchUserProfileSnippets(query, lastUsername, limit).Return(expectedUsers, expectedLastUsername, nil)

	users, lastUsername, err := userProfileService.SearchUserProfileSnippets(query, lastUsername, limit)
	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedUsers, users)
	assert.Equal(t, expectedLastUsername, lastUsername)
}

func TestErrorOnSearchUserProfileSnippetsWithService(t *testing.T) {
	setUpService(t)
	query := "bo"
	lastUsername := "bobusername0"
	limit := 5
	expectedUsers := []*model.UserProfileSnippet{}
	expectedLastUsername := ""
	serviceRepository.EXPECT().SearchUserProfileSnippets(query, lastUsername, limit).Return(expectedUsers, expectedLastUsername, errors.New("some error"))

	users, lastUsername, err := userProfileService.SearchUserProfileSnippets(query, lastUsername, limit)

	assert.Contains(t, loggerOutput.String(), fmt.Sprintf("Error getting userprofile snippets for query %s with lastusername %s and limit %d", query, lastUsername, limit))
	assert.NotNil(t, err)
	assert.ElementsMatch(t, expectedUsers, users)
	assert.Equal(t, expectedLastUsername, lastUsername)
}
