package get_userprofile_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	database "userservice/internal/db"
	get_userprofile "userservice/internal/features/get_userprofile"
)

func TestGetUserProfileRepository_GetUserProfile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dataRepository := &database.Database{Client: db}
	repository := get_userprofile.NewGetUserProfileRepository(dataRepository)

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
