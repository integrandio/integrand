package persistence

import (
	"log"
	"log/slog"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
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

	// Insert our root user into the db...
	// TODO: Clean this up, Should errors cause the
	plainPassword := os.Getenv("ROOT_PASSWORD")
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	user := User{
		Email:    os.Getenv("ROOT_EMAIL"),
		Password: string(password),
	}

	id, err := DATASTORE.CreateUser(user)
	if err != nil {
		// if the user already exists, we don't want this program to crash
		slog.Error(err.Error())
	}
	user.ID = id
	err = DATASTORE.createUserRole(user.ID, SUPER_USER)
	if err != nil {
		slog.Error(err.Error())
	}
	apiKey := os.Getenv("INITIAL_API_KEY")
	// Generate and log an initial API key
	_, err = DATASTORE.InsertAPIKey(apiKey, user.ID)
	if err != nil {
		// if the API Key already exists, we don't want this program to crash
		slog.Error(err.Error())
	}
	log.Printf("Initial API Key: %s\n", apiKey)

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
	endpoints, err := DATASTORE.GetAllEndpoints()
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	for _, endpoint := range endpoints {
		_, err := BROKER.GetTopic(endpoint.TopicName)
		if err != nil {
			_, err = BROKER.CreateTopic(endpoint.TopicName)
			if err != nil {
				slog.Error(err.Error())
				return err
			}
		}
	}

	return nil
}
