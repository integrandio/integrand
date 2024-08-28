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

func GetEndpoints(userID int) ([]persistence.Endpoint, error) {
	err := isUserAuthorized(userID, "read_endpoint")
	if err != nil {
		return nil, err
	}
	return persistence.DATASTORE.GetAllEndpoints()
}

func GetEndpoint(EndpointID string, userID int) (persistence.Endpoint, error) {
	err := isUserAuthorized(userID, "read_endpoint")
	if err != nil {
		return persistence.Endpoint{}, err
	}
	return persistence.DATASTORE.GetEndpoint(EndpointID)
}

func GetEndpointBySecurityKey(EndpointID string, security_key string) (persistence.Endpoint, error) {
	return persistence.DATASTORE.GetEndpointBySecurityKey(EndpointID, security_key)
}

func CreateEndpoint(userId int, EndpointID string, topicName string, userID int) (persistence.Endpoint, error) {
	err := isUserAuthorized(userID, "write_endpoint")
	if err != nil {
		return persistence.Endpoint{}, err
	}
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
	}

	_, err = persistence.BROKER.GetTopic(topicName)
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

func RemoveEndpoint(EndpointID string, userID int) (int, error) {
	err := isUserAuthorized(userID, "write_endpoint")
	if err != nil {
		return 0, err
	}
	return persistence.DATASTORE.DeleteEndpoint(EndpointID)
}
