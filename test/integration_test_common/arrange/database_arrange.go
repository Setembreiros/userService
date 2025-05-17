package integration_test_arrange

import (
	"context"
	"testing"
	"userservice/cmd/provider"
	database "userservice/internal/db"
	newuser "userservice/internal/features/new_user"
	"userservice/internal/model"
)

func CreateTestDatabase(ctx context.Context) *database.Database {
	provider := provider.NewProvider("test", "postgres://postgres:artis@localhost:5432/artis?search_path=public&sslmode=disable")
	sqlDb, err := provider.ProvideDb()
	if err != nil {
		panic(err)
	}

	return sqlDb
}

func AddUserProfileToDatabase(t *testing.T, userPorfile *model.UserProfile) {
	provider := provider.NewProvider("test", "postgres://postgres:artis@localhost:5432/artis?search_path=public&sslmode=disable")
	sqlDb, err := provider.ProvideDb()
	if err != nil {
		panic(err)
	}

	repository := newuser.NewUserRepository(*sqlDb)

	user := &newuser.User{
		Username: userPorfile.Username,
		Email:    userPorfile.Username + "@email.com",
		UserProfile: &newuser.UserProfile{
			FullName: userPorfile.Name,
			Bio:      userPorfile.Bio,
		},
	}

	err = repository.AddNewUser(user)
	if err != nil {
		panic(err)
	}
}
