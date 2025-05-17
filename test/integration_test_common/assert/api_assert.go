package integration_test_assert

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertSuccessResult(t *testing.T, apiResponse *httptest.ResponseRecorder, expectedBodyResponse string) {
	assert.Equal(t, 200, apiResponse.Code)
	assert.Equal(t, removeSpace(expectedBodyResponse), removeSpace(apiResponse.Body.String()))
}

func removeSpace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", ""), "\n", "")
}
