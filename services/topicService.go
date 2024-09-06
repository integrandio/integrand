package services

import (
	"errors"
	"integrand/persistence"
	"integrand/utils"
)

func GetEventStreams(userId int) ([]persistence.TopicDetails, error) {
	err := isUserAuthorized(userId, persistence.READ_TOPIC)
	if err != nil {
		return nil, err
	}
	return persistence.BROKER.GetTopics(), nil
}

func GetEventStream(topicName string, userId int) (persistence.TopicDetails, error) {
	err := isUserAuthorized(userId, persistence.READ_TOPIC)
	if err != nil {
		return persistence.TopicDetails{}, err
	}
	return persistence.BROKER.GetTopic(topicName)
}

func CreateEventStream(topicName string, userId int) (persistence.TopicDetails, error) {
	var topicDetails persistence.TopicDetails
	err := isUserAuthorized(userId, persistence.WRITE_TOPIC)
	if err != nil {
		return topicDetails, err
	}
	if topicName == "" {
		// TODO: Should this be random?
		topicName = utils.RandomString(5)
	}
	topic, err := persistence.BROKER.CreateTopic(topicName)
	if err != nil {
		return topicDetails, err
	}
	// Eventually replace with actual data...
	topicDetails = persistence.TopicDetails{
		TopicName:      topic.TopicName,
		OldestOffset:   0,
		NextOffset:     0,
		RetentionBytes: 1000,
	}
	return topicDetails, nil
}

func DeleteEventStream(topicName string, userId int) error {
	err := isUserAuthorized(userId, persistence.WRITE_TOPIC)
	if err != nil {
		return err
	}
	// Check if the topic is being used by any endpoint
	endpoints, err := persistence.DATASTORE.GetAllEndpoints()
	if err != nil {
		return err
	}
	for _, endpoint := range endpoints {
		if endpoint.TopicName == topicName {
			return errors.New("topic is being used by an endpoint")
		}
	}
	// Check if the topic is being used by any workflow
	workflows, err := persistence.DATASTORE.GetWorkflows()
	if err != nil {
		return err
	}
	for _, workflow := range workflows {
		if workflow.TopicName == topicName {
			return errors.New("topic is being used by a workflow")
		}
	}
	// Delete the topic if it is not being used
	return persistence.BROKER.DeleteTopic(topicName)
}

func GetEvent(topicName string, offset int) ([]byte, error) {
	return persistence.BROKER.ConsumeMessage(topicName, offset)
}
