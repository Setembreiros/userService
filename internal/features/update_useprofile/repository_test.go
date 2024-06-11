package update_userprofile_test

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"

	"userservice/internal/bus"
	mock_bus "userservice/internal/bus/mock"
	database "userservice/internal/db"
	update_userprofile "userservice/internal/features/update_useprofile"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var sqlMock sqlmock.Sqlmock
var updateUserProfileRepository update_userprofile.UpdateUserProfileRepository
var repositoryBus *bus.EventBus
var repositoryExternalBus *mock_bus.MockExternalBus
var repositoryLoggerOutput bytes.Buffer

func setUp(t *testing.T) {
	var db *sql.DB
	db, sqlMock, _ = sqlmock.New()
	database := &database.Database{
		Client: db,
	}
	ctrl := gomock.NewController(t)
	repositoryExternalBus = mock_bus.NewMockExternalBus(ctrl)
	repositoryBus = bus.NewEventBus(repositoryExternalBus)
	log.Logger = log.Output(&serviceLoggerOutput)
	updateUserProfileRepository = update_userprofile.UpdateUserProfileRepository(*database)
}

func TestUpdateUserProfileInRepository(t *testing.T) {
	setUp(t)
	data := &update_userprofile.UserProfile{
		Username: "username1",
		FullName: "user name",
		Bio:      "O mellor usuario do mundo",
		Link:     "www.exemplo.com",
	}
	sqlMock.ExpectBegin()
	expectedRows := sqlmock.NewRows([]string{"user_id"})
	expectedRowData := []driver.Value{"1"}
	expectedRows.AddRow(expectedRowData...)
	sqlMock.ExpectExec("UPDATE userservice.user_profiles").WithArgs("username1", "user name", "O mellor usuario do mundo", "www.exemplo.com").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectCommit()
	expectedUserProfileUpdatedEvent := &update_userprofile.UserProfileUpdatedEvent{
		Username: data.Username,
		Bio:      data.Bio,
		Link:     data.Link,
		FullName: data.FullName,
	}
	expectedEvent, _ := createEvent("UserProfileUpdatedEvent", expectedUserProfileUpdatedEvent)
	repositoryExternalBus.EXPECT().Publish(expectedEvent)

	err := updateUserProfileRepository.UpdateUserProfile(data, repositoryBus)

	assert.Nil(t, err)
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestErrorOnUpdateUserProfileInRepositoryWhenUpdatingUserProfilesTable(t *testing.T) {
	setUp(t)
	data := &update_userprofile.UserProfile{
		Username: "username1",
		FullName: "user name",
		Bio:      "O mellor usuario do mundo",
		Link:     "www.exemplo.com",
	}
	sqlMock.ExpectBegin()
	expectedRows := sqlmock.NewRows([]string{"user_id"})
	expectedRowData := []driver.Value{"1"}
	expectedRows.AddRow(expectedRowData...)
	sqlMock.ExpectExec("UPDATE userservice.user_profiles").WithArgs("username1", "user name", "O mellor usuario do mundo", "www.exemplo.com").WillReturnError(errors.New("failed"))
	sqlMock.ExpectRollback()

	err := updateUserProfileRepository.UpdateUserProfile(data, repositoryBus)

	assert.Contains(t, serviceLoggerOutput.String(), "Update user profile failed")
	assert.NotNil(t, err)
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestNotFoundErrorOnUpdateUserProfileInRepositoryWhenUpdatingUserProfilesTable(t *testing.T) {
	setUp(t)
	data := &update_userprofile.UserProfile{
		Username: "username1",
		FullName: "user name",
		Bio:      "O mellor usuario do mundo",
		Link:     "www.exemplo.com",
	}
	sqlMock.ExpectBegin()
	expectedRows := sqlmock.NewRows([]string{"user_id"})
	expectedRowData := []driver.Value{"1"}
	expectedRows.AddRow(expectedRowData...)
	sqlMock.ExpectExec("UPDATE userservice.user_profiles").WithArgs("username1", "user name", "O mellor usuario do mundo", "www.exemplo.com").WillReturnResult(sqlmock.NewResult(0, 0))
	sqlMock.ExpectRollback()
	var expectedNotFoundError *database.NotFoundError

	err := updateUserProfileRepository.UpdateUserProfile(data, repositoryBus)

	assert.Contains(t, serviceLoggerOutput.String(), "Update user profile failed")
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &expectedNotFoundError)
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestErrorOnUpdateUserProfileInRepositoryWhenPublishingEvent(t *testing.T) {
	setUp(t)
	data := &update_userprofile.UserProfile{
		Username: "username1",
		FullName: "user name",
		Bio:      "O mellor usuario do mundo",
		Link:     "www.exemplo.com",
	}
	sqlMock.ExpectBegin()
	expectedRows := sqlmock.NewRows([]string{"user_id"})
	expectedRowData := []driver.Value{"1"}
	expectedRows.AddRow(expectedRowData...)
	sqlMock.ExpectExec("UPDATE userservice.user_profiles").WithArgs("username1", "user name", "O mellor usuario do mundo", "www.exemplo.com").WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectRollback()
	expectedUserProfileUpdatedEvent := &update_userprofile.UserProfileUpdatedEvent{
		Username: data.Username,
		Bio:      data.Bio,
		Link:     data.Link,
		FullName: data.FullName,
	}
	expectedEvent, _ := createEvent("UserProfileUpdatedEvent", expectedUserProfileUpdatedEvent)
	repositoryExternalBus.EXPECT().Publish(expectedEvent).Return(errors.New("failed"))

	err := updateUserProfileRepository.UpdateUserProfile(data, repositoryBus)

	assert.Contains(t, serviceLoggerOutput.String(), "Publishing UserProfileUpdatedEvent failed")
	assert.NotNil(t, err)
	if err := sqlMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func createEvent(eventName string, eventData any) (*bus.Event, error) {
	dataEvent, err := serialize(eventData)
	if err != nil {
		return nil, err
	}

	return &bus.Event{
		Type: eventName,
		Data: dataEvent,
	}, nil
}

func serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}
