package search_user_test

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
)

var ctrl *gomock.Controller
var loggerOutput bytes.Buffer
var apiResponse *httptest.ResponseRecorder
var ginContext *gin.Context

func SetUp(t *testing.T) {
	ctrl = gomock.NewController(t)
	log.Logger = log.Output(&loggerOutput)
	gin.SetMode(gin.TestMode)
	apiResponse = httptest.NewRecorder()
	ginContext, _ = gin.CreateTestContext(apiResponse)
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
