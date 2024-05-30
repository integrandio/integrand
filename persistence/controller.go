package persistence

import (
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var DATASTORE *Datastore
var SESSION_MANAGER *SessionManager
var BROKER *Broker

func Initialize() {
	devModeString := os.Getenv("DEV_MODE")
	devMode, err := strconv.ParseBool(devModeString)
	if err != nil {
		slog.Error("Invalid value for DEV_MODE variable")
		devMode = false
	}
	DATASTORE, err = setupConnection(devMode)
	if err != nil {
		log.Fatal("Unable to setup database connection: ", err)
	}

	SESSION_MANAGER = NewSessionManager("integrand_session", 3600)

	err = initialize_broker()
	if err != nil {
		log.Fatal(err)
	}
}

func initialize_broker() error {
	var err error
	BROKER, err = NewBroker("data/commitlog")
	if err != nil {
		log.Fatal(err)
	}
	all_sticky_connections, err := DATASTORE.GetAllStickyConnections()
	if err != nil {
		return err
	}
	for _, stickyConnection := range all_sticky_connections {
		_, err := BROKER.GetTopic(stickyConnection.TopicName)
		if err != nil {
			_, err = BROKER.CreateTopic(stickyConnection.TopicName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func AuthorizeToken(headerValue string) error {
	splitToken := strings.Split(headerValue, "Bearer")
	if len(splitToken) != 2 {
		log.Println(splitToken)
		// Error: Bearer token not in proper format
		return errors.New("malformed token")
	}
	authToken := strings.TrimSpace(splitToken[1])
	// TODO: Replace auth token with database logic
	if authToken == "11111" {
		return nil
	} else {
		return errors.New("invalid token")
	}
}

func EmailAuthenticate(Email string, password string) (User, error) {
	user, err := DATASTORE.GetEmailUser(Email)
	if err != nil {
		log.Println(err)
		return User{}, err
	}
	if checkPasswordHash(password, user.Password) {
		return user, nil
	} else {
		return User{}, errors.New("password not valid")
	}
}

func GetSession(w http.ResponseWriter, r *http.Request) SessionDB {
	session, err := SESSION_MANAGER.SessionStart(w, r)
	if err != nil {
		log.Fatal(err)
	}
	return session
}
