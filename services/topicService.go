package services

import (
	"errors"
	"integrand/persistence"
	"integrand/utils"
)

func GetEventStreams(userId int) ([]persistence.TopicDetails, error) {
	return persistence.BROKER.GetTopics(userId), nil
}

func GetEventStream(topicName string) (persistence.TopicDetails, error) {
	return persistence.BROKER.GetTopic(topicName)
}

func CreateEventStream(topicName string) (persistence.TopicDetails, error) {
	var topicDetails persistence.TopicDetails
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
	// Check if the topic is being used by any endpoint
	endpoints, err := GetEndpoints(userId)
	if err != nil {
		return err
	}
	for _, endpoint := range endpoints {
		if endpoint.TopicName == topicName {
			return errors.New("topic is being used by an endpoint")
		}
	}
	// Check if the topic is being used by any workflow
	workflows, err := GetWorkflows()
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
