package services

import (
	"errors"
	"integrand/persistence"
	"integrand/utils"
)

func GetEventStreams() ([]persistence.TopicDetails, error) {
	return persistence.BROKER.GetTopics(), nil
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
		LastestOffset:  0,
		RetentionBytes: 1000,
	}
	return topicDetails, nil
}

func DeleteEventStream(topicName string) error {
	//Check if topic is being used...
	stickyConnections, err := GetStickyConnections()
	if err != nil {
		return err
	}
	topicIsBeingUsed := false
	for _, stickyConnection := range stickyConnections {
		if stickyConnection.TopicName == topicName {
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
