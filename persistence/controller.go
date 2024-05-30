package persistence

import (
	"log"
	"log/slog"
	"os"
	"strconv"
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
