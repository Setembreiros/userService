package get_userprofile_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	database "userservice/internal/db"
	get_userprofile "userservice/internal/features/get_userprofile"
	objectstorage "userservice/internal/objectStorage"
	mock_objectstorage "userservice/internal/objectStorage/mock"

	"go.uber.org/mock/gomock"
)

var osClient *mock_objectstorage.MockObjectStorageClient

func TestGetUserProfileRepository_GetUserProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ctrl := gomock.NewController(t)
	dataRepository := &database.Database{Client: db}
	osClient = mock_objectstorage.NewMockObjectStorageClient(ctrl)
	repository := get_userprofile.NewGetUserProfileRepository(dataRepository, objectstorage.NewObjectStorage(osClient))

	username := "testuser"
	expectedUserProfile := &get_userprofile.UserProfile{
		Username: "testuser",
		FullName: "Test User",
		Bio:      "Test Bio",
		Link:     "https://test.com",
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"username", "full_name", "bio", "link"}).
			AddRow(expectedUserProfile.Username, expectedUserProfile.FullName, expectedUserProfile.Bio, expectedUserProfile.Link)

		mock.ExpectQuery("SELECT (.+) FROM userservice.users (.+) WHERE u.username = \\$1").WithArgs(username).WillReturnRows(rows)

		userProfile, err := repository.GetUserProfile(username)

		assert.NoError(t, err)
		assert.Equal(t, expectedUserProfile, userProfile)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM userservice.users (.+) WHERE u.username = \\$1").WithArgs(username).WillReturnError(sql.ErrNoRows)

		userProfile, err := repository.GetUserProfile(username)

		assert.Error(t, err)
		assert.Nil(t, userProfile)
		assert.IsType(t, &database.NotFoundError{}, err)
	})
}

func TestGetUserProfileRepository_GetPresignedUrlForDownloading(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	osClient = mock_objectstorage.NewMockObjectStorageClient(ctrl)
	objectRepository := objectstorage.NewObjectStorage(osClient)
	dataRepository := &database.Database{} // No database interaction in this test
	repository := get_userprofile.NewGetUserProfileRepository(dataRepository, objectRepository)

	username := "testuser"
	expectedKey := username + "/IMAGEPROFILE/" + username
	expectedUrl := "https://example.com/presigned-url"

	t.Run("Success", func(t *testing.T) {
		osClient.EXPECT().GetPreSignedUrlForPuttingObject(expectedKey).Return(expectedUrl, nil)

		url, err := repository.GetPresignedUrlForDownloading(username)

		assert.NoError(t, err)
		assert.Equal(t, expectedUrl, url)
	})

	t.Run("Error", func(t *testing.T) {
		osClient.EXPECT().GetPreSignedUrlForPuttingObject(expectedKey).Return("", assert.AnError)

		url, err := repository.GetPresignedUrlForDownloading(username)

		assert.Error(t, err)
		assert.Empty(t, url)
	})
}
