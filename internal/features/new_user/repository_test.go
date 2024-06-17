package newuser_test

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	database "userservice/internal/db"
	newuser "userservice/internal/features/new_user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var sqlMock sqlmock.Sqlmock
var newUserRepository newuser.NewUserRepository
var repositoryLoggerOutput bytes.Buffer

func setUp(t *testing.T) {
	var db *sql.DB
	db, sqlMock, _ = sqlmock.New()
	database := &database.Database{
		Client: db,
	}
	log.Logger = log.Output(&serviceLoggerOutput)
	newUserRepository = newuser.NewUserRepository(*database)
}

func TestAddNewUserInRepository(t *testing.T) {
	setUp(t)
	data := &newuser.User{
		Username: "username1",
		Email:    "username@email.com",
		UserType: "UA",
		Region:   "Vigo",
		UserProfile: &newuser.UserProfile{
			FullName: "user name",
			Bio:      "O mellor usuario do mundo",
			Link:     "www.exemplo.com",
		},
	}
	sqlMock.ExpectBegin()
	expectedRows := sqlmock.NewRows([]string{"user_id"})
	expectedRowData := []driver.Value{"1"}
	expectedRows.AddRow(expectedRowData...)
	sqlMock.ExpectQuery("INSERT INTO userservice.users").WithArgs("username1", "username@email.com", "Vigo", "UA").WillReturnRows(expectedRows)
	sqlMock.ExpectExec("INSERT INTO userservice.user_profiles").WithArgs(1, "user name", "O mellor usuario do mundo", "www.exemplo.com").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectCommit()

	err := newUserRepository.AddNewUser(data)

	assert.Nil(t, err)
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestErrorOnAddNewUserInRepositoryWhenInsertIntoUsersTable(t *testing.T) {
	setUp(t)
	data := &newuser.User{
		Username: "username1",
		Email:    "username@email.com",
		UserType: "UA",
		Region:   "Vigo",
		UserProfile: &newuser.UserProfile{
			FullName: "user name",
			Bio:      "O mellor usuario do mundo",
			Link:     "www.exemplo.com",
		},
	}
	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	err := newUserRepository.AddNewUser(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Insert user failed")
	assert.NotNil(t, err)
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestErrorOnAddNewUserInRepositoryWhenInsertIntoUserProfilesTable(t *testing.T) {
	setUp(t)
	data := &newuser.User{
		Username: "username1",
		Email:    "username@email.com",
		UserType: "UA",
		Region:   "Vigo",
		UserProfile: &newuser.UserProfile{
			FullName: "user name",
			Bio:      "O mellor usuario do mundo",
			Link:     "www.exemplo.com",
		},
	}
	sqlMock.ExpectBegin()
	expectedRows := sqlmock.NewRows([]string{"user_id"})
	expectedRowData := []driver.Value{"1"}
	expectedRows.AddRow(expectedRowData...)
	sqlMock.ExpectQuery("INSERT INTO userservice.users").WithArgs("username1", "username@email.com", "Vigo", "UA").WillReturnRows(expectedRows)
	sqlMock.ExpectExec("INSERT INTO userservice.user_profiles").WithArgs(1, "user name", "O mellor usuario do mundo", "www.exemplo.com").WillReturnError(errors.New("failed"))

	sqlMock.ExpectRollback()

	err := newUserRepository.AddNewUser(data)

	assert.Contains(t, serviceLoggerOutput.String(), "Insert userProfile failed")
	assert.NotNil(t, err)
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
