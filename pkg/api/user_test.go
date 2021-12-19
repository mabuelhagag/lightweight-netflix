package api_test

import (
	"errors"
	"lightweight-netflix/pkg/api"
	"reflect"
	"testing"
)

type mockUserRepo struct{}

func (m mockUserRepo) CreateUser(request api.NewUserRequest) error {
	if request.FullName == "test user already created" {
		return errors.New("repository - user already exists in database")
	}

	return nil
}

func (m mockUserRepo) UpdateActivityLevel(request api.UpdateActivityLevelRequest) error {
	panic("implement me")
}

func TestCreateNewUser(t *testing.T) {
	mockRepo := mockUserRepo{}
	mockUserService := api.NewUserService(&mockRepo)

	tests := []struct {
		name    string
		request api.NewUserRequest
		want    error
	}{
		{name: "should create a new user successfully", request: api.NewUserRequest{
			FullName: "test user",
			Age:      20,
			Email:    "test_user@gmail.com",
		}, want: nil},
		{name: "should return an error because of missing email", request: api.NewUserRequest{
			FullName: "test user",
			Age:      20,
			Email:    "",
		}, want: errors.New("user service - email required")},
		{name: "should return an error because of missing name", request: api.NewUserRequest{
			FullName: "",
			Age:      20,
			Email:    "test_user@gmail.com",
		}, want: errors.New("user service - name required")},
		{name: "should return error from database because user already exists", request: api.NewUserRequest{
			FullName: "test user already created",
			Age:      20,
			Email:    "test_user@gmail.com",
		}, want: errors.New("repository - user already exists in database")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := mockUserService.New(test.request)

			if !reflect.DeepEqual(err, test.want) {
				t.Errorf("test: %v failed. got: %v, wanted: %v", test.name, err, test.want)
			}
		})
	}
}
