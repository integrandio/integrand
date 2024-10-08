package services

import (
	"errors"
	"integrand/persistence"
	"integrand/utils"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func AuthenticateCookie(w http.ResponseWriter, r *http.Request) (int, error) {
	sess := GetSession(w, r)
	userInterface, err := sess.Get("userID")
	if err != nil {
		log.Fatalln(err)
	}
	if userInterface == nil {
		return 0, errors.New("cookie not valid")
	}
	userID, ok := userInterface.(float64)
	if !ok {
		log.Fatal("Unable to cast user value to float")
	}
	return int(userID), nil
}

func isUserAuthorized(userId int, requiredSecurable persistence.Securable) error {
	securables, err := persistence.DATASTORE.GetSecurablesByUserId(userId)
	if err != nil {
		return err
	}
	contained := utils.Contains(securables, string(requiredSecurable))
	if !contained {
		return errors.New("user is not authorized")
	}
	return nil
}

func ListAPIKeys() ([]persistence.ApiKey, error) {
	return persistence.DATASTORE.GetAPIKeys()
}

func AuthenticateToken(headerValue string) (persistence.ApiKey, error) {
	splitToken := strings.Split(headerValue, "Bearer")
	if len(splitToken) != 2 {
		log.Println(splitToken)
		// Error: Bearer token not in proper format
		return persistence.ApiKey{}, errors.New("malformed token")
	}
	authToken := strings.TrimSpace(splitToken[1])
	apiKey, err := persistence.DATASTORE.GetApiKey(authToken)
	if err != nil {
		return persistence.ApiKey{}, err
	}
	return apiKey, nil
}

func EmailAuthenticate(Email string, password string) (persistence.User, error) {
	user, err := persistence.DATASTORE.GetUserByEmail(Email)
	if err != nil {
		slog.Error(err.Error())
		return persistence.User{}, err
	}
	if checkPasswordHash(password, user.Password) {
		return user, nil
	} else {
		return persistence.User{}, errors.New("password not valid")
	}
}

func GetSession(w http.ResponseWriter, r *http.Request) persistence.SessionDB {
	session, err := persistence.SESSION_MANAGER.SessionStart(w, r)
	if err != nil {
		log.Fatal(err)
	}
	return session
}

func CreateAPIKey(userId int) (string, error) {
	for {
		key := utils.RandomString(20)
		_, err := persistence.DATASTORE.InsertAPIKey(key, userId)
		if err != nil {
			if err.Error() == "API key already exists" {
				continue
			}
			return "", err
		}
		return key, nil
	}
}

func DeleteAPIKey(key string, userId int) (int, error) {
	return persistence.DATASTORE.DeleteAPIKey(key, userId)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
