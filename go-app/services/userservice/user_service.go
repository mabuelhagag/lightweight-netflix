package userservice

import (
	"go-app/domain/user"
	"go-app/repositories/userrepo"
)

// UserService interface
type UserService interface {
	CreateUser(user *user.User) (*user.User, error)
}

type userService struct {
	Repo userrepo.Repo
}

// NewUserService will instantiate User Service
func NewUserService(
	repo userrepo.Repo,
) UserService {

	return &userService{
		Repo: repo,
	}
}

func (us *userService) CreateUser(user *user.User) (*user.User, error) {
	return us.Repo.CreateUser(user)
}
