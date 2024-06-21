package services

import (
	"encoding/json"
	"integrand/persistence"
	"integrand/utils"
	"log/slog"
	"time"
)

func MessageToSink(topicName string, userId int, value interface{}) error {
	resBytes, err := json.Marshal(value)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = persistence.BROKER.ProduceMessage(topicName, userId, resBytes)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func GetStickyConnections(userId int) ([]persistence.StickyConnection, error) {
	return persistence.DATASTORE.GetAllStickyConnections(userId)
}

func GetStickyConnection(userId int, stickyConnectionID string) (persistence.StickyConnection, error) {
	return persistence.DATASTORE.GetStickeyConnectionByUser(stickyConnectionID, userId)
}

func GetStickyConnectionBySecurityKey(stickyConnectionID string, security_key string) (persistence.StickyConnection, error) {
	return persistence.DATASTORE.GetStickeyConnectionBySecurityKey(stickyConnectionID, security_key)
}

func CreateStickyConnection(userId int, stickyConnectionID string, topicName string) (persistence.StickyConnection, error) {
	if stickyConnectionID == "" {
		stickyConnectionID = utils.RandomString(5)
	}
	if topicName == "" {
		// TODO: Should this be random?
		topicName = utils.RandomString(5)
	}

	security_key := utils.RandomString(8)

	sticky_connection := persistence.StickyConnection{
		RouteID:      stickyConnectionID,
		Security_key: security_key,
		LastModified: time.Now(),
		UserId:       userId,
	}

	_, err := persistence.BROKER.GetTopic(topicName, userId)
	if err != nil {
		// TODO: Create the topic if it does not exist
		_, err := persistence.BROKER.CreateTopic(topicName, userId)
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

func RemoveStickyConnection(userId int, stickyConnectionID string) (int, error) {
	return persistence.DATASTORE.DeleteStickyConnection(stickyConnectionID, userId)
}
