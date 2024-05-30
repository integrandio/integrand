package persistence

import (
	"encoding/gob"
	"errors"
	"log"
	"log/slog"
	"os"

	"integrand/commitlog"
	"integrand/utils"
)

type Broker struct {
	BaseDirectory string
	Topics        []Topic
}

type Topic struct {
	TopicName      string
	TopicDirectory string
	commitlog      *commitlog.Commitlog
}

func NewBroker(directory string) (*Broker, error) {
	broker := Broker{
		BaseDirectory: directory,
	}
	topics, err := broker.loadTopicSnapshot()
	if err != nil {
		// Check the error here, only want to proceed if we hit certain errors
		slog.Warn(err.Error())
		topics = []Topic{}
	}

	broker.Topics = topics
	return &broker, nil
}

func (broker *Broker) ProduceMessage(topicName string, message []byte) error {
	for _, topic := range broker.Topics {
		if topic.TopicName == topicName {
			err := topic.commitlog.Append(message)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("topic not found")
}

func (broker *Broker) ConsumeMessage(topicName string, offset int) ([]byte, error) {
	for _, topic := range broker.Topics {
		if topic.TopicName == topicName {
			msg, err := topic.commitlog.Read(offset)
			if err != nil {
				return []byte{}, err
			}
			return msg, nil
		}
	}
	return []byte{}, errors.New("topic not found")
}

func (broker *Broker) CreateTopic(topicName string) (Topic, error) {
	var topic Topic
	//Check if the topic name exists
	for _, topic := range broker.Topics {
		if topic.TopicName == topicName {
			return topic, errors.New("topic already exists")
		}
	}
	//create our struct, add it to the array
	topicDirectory := broker.BaseDirectory + "/" + utils.RandomString(5)
	cl, err := commitlog.New(topicDirectory)
	if err != nil {
		slog.Error("Unable to initilize commitlog")
		return topic, err
	}
	topic = Topic{
		TopicName:      topicName,
		TopicDirectory: topicDirectory,
		commitlog:      cl,
	}
	broker.Topics = append(broker.Topics, topic)
	err = broker.takeTopicSnapshot() //Should this be in a go routine?
	if err != nil {
		return topic, err
	}
	return topic, nil
}

func (broker *Broker) DeleteTopic(topicName string) error {
	var foundTopicIndex int
	topicFound := false
	//Check if the topic name exists
	for i, topic := range broker.Topics {
		if topic.TopicName == topicName {
			topicFound = true
			foundTopicIndex = i
			break
		}
	}
	if topicFound {
		broker.Topics[foundTopicIndex] = broker.Topics[len(broker.Topics)-1]
		broker.Topics = broker.Topics[:len(broker.Topics)-1]
		log.Println(broker.Topics)
	} else {
		return errors.New("topic does not exist, unable to delete")
	}
	//We need to also delete the commitlog...
	err := broker.takeTopicSnapshot() //Should this be in a go routine?
	if err != nil {
		return err
	}
	return nil
}

func (broker *Broker) takeTopicSnapshot() error {
	snapshotPath := broker.BaseDirectory + "/topics.gob"

	file, err := os.OpenFile(snapshotPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(broker.Topics)
	if err != nil {
		return err
	}
	return nil
}

func (broker *Broker) loadTopicSnapshot() ([]Topic, error) {
	var snapTopics []Topic
	snapshotPath := broker.BaseDirectory + "/topics.gob"
	file, err := os.Open(snapshotPath)
	if err != nil {
		slog.Error(err.Error())
		return snapTopics, err
	}
	decoder := gob.NewDecoder(file)

	err = decoder.Decode(&snapTopics)
	if err != nil {
		slog.Error(err.Error())
		return snapTopics, err
	}

	for i, snapTopic := range snapTopics {
		cl, err := commitlog.New(snapTopic.TopicDirectory)
		if err != nil {
			slog.Error(err.Error())
			return snapTopics, err
		}
		snapTopics[i].commitlog = cl
	}
	return snapTopics, nil
}

type TopicDetails struct {
	TopicName      string `json:"topicName"`
	OldestOffset   int    `json:"oldestOffset"`
	LastestOffset  int    `json:"latestOffset"`
	RetentionBytes int    `json:"retentionBytes"`
}

func (broker *Broker) GetTopics() []TopicDetails {
	topicDetails := []TopicDetails{}
	for i := range broker.Topics {
		topicDetails = append(topicDetails, broker.getTopicDetails(i))
	}
	return topicDetails
}

func (broker *Broker) GetTopic(topicName string) (TopicDetails, error) {
	var topicDetails TopicDetails
	//Check if the topic name exists
	for i, topic := range broker.Topics {
		if topic.TopicName == topicName {
			return broker.getTopicDetails(i), nil
		}
	}
	return topicDetails, errors.New("topic does not exist")
}

func (broker *Broker) getTopicDetails(topicIndex int) TopicDetails {
	topic := broker.Topics[topicIndex]
	details := topic.commitlog.GetCommitlogDetails()
	return TopicDetails{
		TopicName:      topic.TopicName,
		OldestOffset:   details.OldestOffset,
		LastestOffset:  details.LastestOffset,
		RetentionBytes: details.RetentionBytes,
	}
}
