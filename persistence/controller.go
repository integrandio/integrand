package persistence

import (
	"errors"
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
var apiKeys struct {
	sync.RWMutex
	keys []string
}

func init() {
	apiKeys.keys = make([]string, 0)
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

	// Generate and log an initial API key
	initialAPIKey := utils.RandomString(32)
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
		return err
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

// adds a new API key to the store
func AddAPIKey(key string) error {
	apiKeys.Lock()
	defer apiKeys.Unlock()
	for _, k := range apiKeys.keys {
		if k == key {
			return errors.New("API key already exists")
		}
	}
	apiKeys.keys = append(apiKeys.keys, key)
	return nil
}

// removes an API key from the store
func DeleteAPIKey(key string) error {
	apiKeys.Lock()
	defer apiKeys.Unlock()
	for i, k := range apiKeys.keys {
		if k == key {
			apiKeys.keys = append(apiKeys.keys[:i], apiKeys.keys[i+1:]...)
			return nil
		}
	}
	return errors.New("API key not found")
}

// checks if an API key is valid
func IsAPIKeyValid(key string) bool {
	apiKeys.RLock()
	defer apiKeys.RUnlock()
	for _, k := range apiKeys.keys {
		if k == key {
			return true
		}
	}
	return false
}
