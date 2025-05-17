package integration_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	integration_test_arrange "userservice/test/integration_test_common/arrange"
	integration_test_assert "userservice/test/integration_test_common/assert"

	database "userservice/internal/db"
	"userservice/internal/features/search_user"
	"userservice/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var db *database.Database
var loggerOutput bytes.Buffer
var controller *search_user.SearchUserController
var ginContext *gin.Context
var apiResponse *httptest.ResponseRecorder

func setUp() {
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)

	// Real infrastructure and services
	db = integration_test_arrange.CreateTestDatabase(ginContext)
	log.Logger = log.Output(&loggerOutput)
	controller = search_user.NewSearchUserController(search_user.NewSearchUserService(search_user.SearchUserRepository(*db)))
}

func tearDown() {
	db.Clean()
}

func TestGetUserProfileSnippets_WhenDatabaseReturnsSuccess(t *testing.T) {
	setUp()
	defer tearDown()
	populateUserPorfilesDb(t)
	query := "bo"
	limit := 5
	ginContext.Request, _ = http.NewRequest("GET", "/postLikes", nil)
	u := url.Values{}
	u.Add("query", query)
	u.Add("limit", strconv.Itoa(limit))
	ginContext.Request.URL.RawQuery = u.Encode()
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
					"username":  "aliceusername2",
					"name":  "bob2"
			},
			{
					"username":  "bobusername2",
					"name":  "alice2"
			},
			{
					"username":  "usernamebob3",
					"name":  "alice3"
			},
			{
					"username":  "aliceAndBobusername3",
					"name":  "aliceAndBob3"
			}
			]
		}
	}`

	controller.SearchUser(ginContext)

	integration_test_assert.AssertSuccessResult(t, apiResponse, expectedBodyResponse)
}

func populateUserPorfilesDb(t *testing.T) {
	existingUsers := []*model.UserProfile{
		{
			Username: "baliceusername0",
			Name:     "balice0",
			Bio:      "bob",
		},
		{
			Username: "bobusername1",
			Name:     "alice1",
			Bio:      "aaaaa",
		},
		{
			Username: "aliceusername5",
			Name:     "balice0",
			Bio:      "aaaaa",
		},
		{
			Username: "aliceusername2",
			Name:     "bob2",
			Bio:      "aaaaa",
		},
		{
			Username: "bobusername2",
			Name:     "alice2",
			Bio:      "aaaaa",
		},
		{
			Username: "usernamebob3",
			Name:     "alice3",
			Bio:      "aaaaa",
		},
		{
			Username: "baliceusername1",
			Name:     "alice1",
			Bio:      "bo",
		},
		{
			Username: "alicebusername1",
			Name:     "alice1",
			Bio:      "aaaaa",
		},
		{
			Username: "aliceAndBobusername3",
			Name:     "aliceAndBob3",
			Bio:      "aaaaa",
		},
		{
			Username: "aliceAndbobusername4",
			Name:     "aliceAndBob4",
			Bio:      "aaaaa",
		},
	}

	for _, existingUser := range existingUsers {
		integration_test_arrange.AddUserProfileToDatabase(t, existingUser)
	}
}
