package auth

import (
	"errors"
	"log"
	"social-network/pkg/schema"
	tokenService "social-network/pkg/token"

	"github.com/gorilla/sessions"
)

func CheckSession(session *sessions.Session) (*schema.User, error) {
	var user schema.User
	tokenString, ok := session.Values["token"]
	if !ok {
		return &user, errors.New("invalid_token")
	}

	token, err := tokenService.Decode(tokenString.(string))
	if err != nil {
		log.Println("Token decode error", err)
		return &user, errors.New("invalid_token")
	}

	return token.User, nil
}
