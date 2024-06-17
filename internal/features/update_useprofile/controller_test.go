package update_userprofile_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"userservice/internal/bus"
	database "userservice/internal/db"
	update_userprofile "userservice/internal/features/update_useprofile"
	mock_update_userprofile "userservice/internal/features/update_useprofile/mock"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var controllerLoggerOutput bytes.Buffer
var controllerRepository *mock_update_userprofile.MockRepository
var controllerBus *bus.EventBus
var controller *update_userprofile.PutUserProfileController
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func setUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	log.Logger = log.Output(&controllerLoggerOutput)
	controllerRepository = mock_update_userprofile.NewMockRepository(ctrl)
	controllerBus = &bus.EventBus{}
	controller = update_userprofile.NewPutUserProfileController(controllerRepository, controllerBus)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func TestPutUserProfile(t *testing.T) {
	setUpHandler(t)
	newUserProfile := &update_userprofile.UserProfile{
		Username: "username1",
		FullName: "user name",
		Bio:      "O mellor usuario do mundo",
		Link:     "www.exemplo.com",
	}
	data, _ := serializeData(newUserProfile)
	ginContext.Request = httptest.NewRequest(http.MethodPut, "/userprofile", bytes.NewBuffer(data))

	controllerRepository.EXPECT().UpdateUserProfile(newUserProfile, controllerBus)
	expectedBodyResponse := `{
		"error": false,
		"message": "200 OK",
		"content": {
			"username": "username1",
			"full_name": "user name",
			"bio": "O mellor usuario do mundo",
			"link": "www.exemplo.com"
		}
	}`

	controller.PutUserProfile(ginContext)

	assert.Equal(t, apiResponse.Code, 200)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestUserNotFoundOnPutUserProfile(t *testing.T) {
	setUpHandler(t)
	noExistingUserProfile := &update_userprofile.UserProfile{
		Username: "noUsername",
		FullName: "",
		Bio:      "",
		Link:     "",
	}
	data, _ := serializeData(noExistingUserProfile)
	ginContext.Request = httptest.NewRequest(http.MethodPut, "/userprofile", bytes.NewBuffer(data))
	expectedNotFoundError := &database.NotFoundError{}
	controllerRepository.EXPECT().UpdateUserProfile(noExistingUserProfile, controllerBus).Return(expectedNotFoundError)
	expectedBodyResponse := `{
		"error": true,
		"message": "User Profile not found for username ` + noExistingUserProfile.Username + `",
		"content":null
	}`

	controller.PutUserProfile(ginContext)

	assert.Equal(t, apiResponse.Code, 404)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func TestInternalServerOnGetUserProfile(t *testing.T) {
	setUpHandler(t)
	newUserProfile := &update_userprofile.UserProfile{
		Username: "username1",
		FullName: "",
		Bio:      "",
		Link:     "",
	}
	data, _ := serializeData(newUserProfile)
	ginContext.Request = httptest.NewRequest(http.MethodPut, "/userprofile", bytes.NewBuffer(data))
	ginContext.Request.Body = io.NopCloser(bytes.NewBuffer(data))
	expectedError := errors.New("some error")
	controllerRepository.EXPECT().UpdateUserProfile(newUserProfile, controllerBus).Return(expectedError)
	expectedBodyResponse := `{
		"error": true,
		"message": "` + expectedError.Error() + `",
		"content":null
	}`

	controller.PutUserProfile(ginContext)

	assert.Equal(t, apiResponse.Code, 500)
	assert.Equal(t, removeSpace(apiResponse.Body.String()), removeSpace(expectedBodyResponse))
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}

func serializeData(data any) ([]byte, error) {
	return json.Marshal(data)
}
