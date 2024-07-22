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
	err = persistence.BROKER.ProduceMessage(topicName, resBytes)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func GetEndpoints(userId int) ([]persistence.Endpoint, error) {
	return persistence.DATASTORE.GetAllEndpoints(userId)
}

func GetEndpoint(EndpointID string, userId int) (persistence.Endpoint, error) {
	return persistence.DATASTORE.GetEndpointByUser(EndpointID, userId)
}

func GetEndpointBySecurityKey(EndpointID string, security_key string) (persistence.Endpoint, error) {
	return persistence.DATASTORE.GetEndpointBySecurityKey(EndpointID, security_key)
}

func CreateEndpoint(userId int, EndpointID string, topicName string) (persistence.Endpoint, error) {
	if EndpointID == "" {
		EndpointID = utils.RandomString(5)
	}
	if topicName == "" {
		// TODO: Should this be random?
		topicName = utils.RandomString(5)
	}

	security_key := utils.RandomString(8)

	sticky_connection := persistence.Endpoint{
		RouteID:      EndpointID,
		Security_key: security_key,
		LastModified: time.Now(),
		UserId:       userId,
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
	_, err = persistence.DATASTORE.InsertEndpoint(sticky_connection)
	if err != nil {
		return sticky_connection, err
	}
	return sticky_connection, nil
}

func RemoveEndpoint(userId int, EndpointID string) (int, error) {
	return persistence.DATASTORE.DeleteEndpoint(EndpointID, userId)
}
