package users

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User struct
type User struct {
	mgm.DefaultModel `bson:",inline"`
	FullName         string `bson:"name"`
	Age              uint8  `bson:"age"`
	Email            string `bson:"email"`
	Password         string `bson:"password"`
}

func (model *User) Saving() error {

	// Check if the user exists; TODO: create email index in users collection
	user := &User{}
	err := mgm.Coll(model).First(bson.M{"email": model.Email}, user)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("unable to check for user existence")
		}

	}
	if user.Email != "" {
		return errors.New("email already used by another user")
	}

	// Hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(model.Password), 14)
	if err != nil {
		return err
	}
	model.Password = string(bytes)

	// Call the DefaultModel Creating hook
	if err := model.DefaultModel.Creating(); err != nil {
		return err
	}
	return nil
}

// UserInput represents createUser body format
type UserInput struct {
	FullName             string `json:"full_name" mod:"trim" binding:"required"`
	Age                  uint8  `json:"age" binding:"required,gt=0,lt=100"`
	Email                string `json:"email" mod:"trim,lcase" binding:"required"`
	Password             string `json:"password" binding:"required,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

// LoginInfoInput represents
type LoginInfoInput struct {
	Email    string `json:"email" mod:"trim" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInfoOutput struct {
	Token string `json:"token"`
}

// JwtWrapper wraps the signing key and the issuer
type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JwtClaim adds email as a claim to the token
type JwtClaim struct {
	Email string
	jwt.StandardClaims
}

// AppJwtWrapper TODO: add secret key to config
var AppJwtWrapper = JwtWrapper{
	SecretKey:       "verysecretkey",
	Issuer:          "AuthService",
	ExpirationHours: 24,
}

// GenerateToken generates a jwt token
func (j *JwtWrapper) GenerateToken(email string) (signedToken string, err error) {
	claims := &JwtClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}

	return
}

// ValidateToken validates the jwt token
func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return

}
