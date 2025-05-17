package search_user_test

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"userservice/internal/features/search_user"
	mock_search_user "userservice/internal/features/search_user/test/mock"
	"userservice/internal/model"

	"github.com/go-playground/assert/v2"
)

var controllerService *mock_search_user.MockControllerService
var controller *search_user.SearchUserController

func setUpController(t *testing.T) {
	SetUp(t)
	controllerService = mock_search_user.NewMockControllerService(ctrl)
	controller = search_user.NewSearchUserController(controllerService)
}

func TestSearchUser_WhenSuccess(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/userprofile-snippets", nil)
	expectedQuery := "bob"
	expectedLastUsername := "bobusername0"
	expectedLimit := 3
	u := url.Values{}
	u.Add("query", expectedQuery)
	u.Add("lastUsername", expectedLastUsername)
	u.Add("limit", strconv.Itoa(expectedLimit))
	ginContext.Request.URL.RawQuery = u.Encode()
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
	controllerService.EXPECT().SearchUserProfileSnippets(expectedQuery, expectedLastUsername, expectedLimit).Return(expectedUsers, "aliceAndBobusername3", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"users":[	
			{
					"username":  "bobusername1",
					"name":  "alice1"
			},
			{
					"username":  "aliceusername1",
					"name":  "bob2"
			},
			{
					"username":  "aliceAndBobusername3",
					"name":  "aliceAndBob3"
			}
			],
			"lastUsername":"aliceAndBobusername3"
		}
	}`

	controller.SearchUser(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestSearchUser_WhenSuccessWithDefaultPaginationParameters(t *testing.T) {
	setUpController(t)
	ginContext.Request, _ = http.NewRequest("GET", "/userprofile-snippets", nil)
	expectedQuery := ""
	expectedLastUsername := ""
	expectedLimit := 5
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
	controllerService.EXPECT().SearchUserProfileSnippets(expectedQuery, expectedLastUsername, expectedLimit).Return(expectedUsers, "aliceAndBobusername3", nil)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"users":[	
			{
					"username":  "bobusername1",
					"name":  "alice1"
			},
			{
					"username":  "aliceusername1",
					"name":  "bob2"
			},
			{
					"username":  "aliceAndBobusername3",
					"name":  "aliceAndBob3"
			}
			],
			"lastUsername":"aliceAndBobusername3"
		}
	}`

	controller.SearchUser(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}
