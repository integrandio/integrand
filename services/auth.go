package services

import (
	"errors"
	"integrand/persistence"
	"integrand/persistence/apikeys"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func AuthorizeToken(headerValue string) error {
	splitToken := strings.Split(headerValue, "Bearer")
	if len(splitToken) != 2 {
		log.Println(splitToken)
		// Error: Bearer token not in proper format
		return errors.New("malformed token")
	}
	authToken := strings.TrimSpace(splitToken[1])
	if apikeys.IsAPIKeyValid(authToken) {
		return nil
	} else {
		return errors.New("invalid token")
	}
}

func EmailAuthenticate(Email string, password string) (persistence.User, error) {
	user, err := persistence.DATASTORE.GetEmailUser(Email)
	if err != nil {
		log.Println(err)
		return persistence.User{}, err
	}
	if checkPasswordHash(password, user.Password) {
		return user, nil
	} else {
		return persistence.User{}, errors.New("password not valid")
	}
}

func CreateNewEmailUser(email string, plainPassword string) (persistence.User, error) {
	var user persistence.User
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	user = persistence.User{
		Email:        email,
		Password:     string(password),
		AuthType:     "email",
		CreatedAt:    time.Now(),
		LastModified: time.Now(),
	}
	id, err := persistence.DATASTORE.CreateEmailUser(user)
	if err != nil {
		return user, nil
	}
	user.ID = id
	return user, nil
}

func GetSession(w http.ResponseWriter, r *http.Request) persistence.SessionDB {
	session, err := persistence.SESSION_MANAGER.SessionStart(w, r)
	if err != nil {
		log.Fatal(err)
	}
	return session
}

func CreateAPIKey(key string) {
	apikeys.AddAPIKey(key)
}

func DeleteAPIKey(key string) {
	apikeys.DeleteAPIKey(key)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
