package user

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// User struct
type User struct {
	FullName string `bson:"name"`
	Age      uint8  `bson:"page_count"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
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
		err = errors.New("Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return

}
