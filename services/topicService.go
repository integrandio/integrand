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
	//Check if topic is being used...
	endpoints, err := GetEndpoints(userId)
	if err != nil {
		return err
	}
	topicIsBeingUsed := false
	for _, endpoint := range endpoints {
		if endpoint.TopicName == topicName {
			topicIsBeingUsed = true
			break
		}
	}

	if topicIsBeingUsed {
		return errors.New("topic is being used")
	}

	return persistence.BROKER.DeleteTopic(topicName)
}

func GetEvent(topicName string, offset int) ([]byte, error) {
	return persistence.BROKER.ConsumeMessage(topicName, offset)
}
