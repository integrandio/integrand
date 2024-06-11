package persistence

import (
	"integrand/utils"
	"log"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var DATASTORE *Datastore
var SESSION_MANAGER *SessionManager
var BROKER *Broker

// Data structure to store API keys
var API_KEYS struct {
	sync.RWMutex
	keys []string
}

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
		AuthType: EMAIL,
	}
	_, err = DATASTORE.CreateEmailUser(user)
	if err != nil {
		// if the user already exists, we don't want this program to crash
		slog.Error(err.Error())
	}

	API_KEYS.keys = make([]string, 0)
	// Generate and log an initial API key
	initialAPIKey := utils.RandomString(20)
	err = AddAPIKey(initialAPIKey)
	if err != nil {
		log.Fatal("Error generating initial API key: ", err)
	}
	log.Printf("Initial API Key: %s\n", initialAPIKey)
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
