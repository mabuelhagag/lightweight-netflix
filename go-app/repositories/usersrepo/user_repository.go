package usersrepo

import (
	"errors"
	"github.com/kamva/mgm/v3"
	"go-app/definitions/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Repo Interface
type Repo interface {
	CreateUser(user *users.User) (*users.User, error)
	CheckPassword(user *users.LoginInfoInput) error
}

type usersRepo struct {
	db *mongo.Client
}

// NewUsersRepo will instantiate User Repository
func NewUsersRepo(db *mongo.Client) Repo {
	return &usersRepo{
		db: db,
	}
}

func (b *usersRepo) CreateUser(user *users.User) (*users.User, error) {
	err := mgm.Coll(user).Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (b *usersRepo) CheckPassword(input *users.LoginInfoInput) (err error) {
	var userFound = &users.User{}
	err = mgm.Coll(userFound).First(bson.M{"email": input.Email}, userFound)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("user does not exist")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(input.Password))
	if err != nil {
		return errors.New("email and password combination is incorrect")
	}

	return nil
}
