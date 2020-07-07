package token

import (
	"errors"
	"social-network/pkg/schema"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	key = []byte("mySuperSecretKeyLol")
)

type CustomClaims struct {
	User *schema.User
	jwt.StandardClaims
}

type Authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *schema.User) (string, error)
}

type token struct {
}

func Decode(token string) (*CustomClaims, error) {
	tokenType, err := jwt.ParseWithClaims(string(token), &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if tokenType == nil {
		return nil, errors.New("invalid_token")
	}

	claims, ok := tokenType.Claims.(*CustomClaims)

	if ok && tokenType.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func Encode(user *schema.User, rememberMe bool) (string, error) {
	expireToken := time.Now().Add(time.Hour * 24).Unix()

	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "token.social.network",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}
