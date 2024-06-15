package services

import (
	"encoding/json"
	"integrand/persistence"
	"integrand/utils"
	"log/slog"
	"time"
)

func MessageToSink(topicName string, value interface{}) error {
	resBytes, err := json.Marshal(value)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = persistence.BROKER.ProduceMessage(topicName, resBytes)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func GetStickyConnections() ([]persistence.StickyConnection, error) {
	return persistence.DATASTORE.GetAllStickyConnections()
}

func GetStickyConnection(stickyConnectionID string) (persistence.StickyConnection, error) {
	return persistence.DATASTORE.GetStickeyConnection(stickyConnectionID)
}

func CreateStickyConnection(id string, topicName string) (persistence.StickyConnection, error) {
	if id == "" {
		id = utils.RandomString(5)
	}
	if topicName == "" {
		// TODO: Should this be random?
		topicName = utils.RandomString(5)
	}

	connectionKey := utils.RandomString(8)

	sticky_connection := persistence.StickyConnection{
		RouteID:          id,
		ConnectionApiKey: connectionKey,
		LastModified:     time.Now(),
	}

	_, err := persistence.BROKER.GetTopic(topicName)
	if err != nil {
		// TODO: Create the topic if it does not exist
		_, err := persistence.BROKER.CreateTopic(topicName)
		if err != nil {
			return sticky_connection, err
		}
	}
	sticky_connection.TopicName = topicName
	// TODO: we're missing the timestamp that the db assigns. At some point lets fix this.
	_, err = persistence.DATASTORE.InsertStickyConnection(sticky_connection)
	if err != nil {
		return sticky_connection, err
	}
	return sticky_connection, nil
}

func RemoveStickyConnection(stickyConnectionID string) (int, error) {
	return persistence.DATASTORE.DeleteStickyConnection(stickyConnectionID)
}
