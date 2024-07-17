package services

import (
	"integrand/persistence"
	"integrand/workflows"
)

func Workflower() error {
	topicName := "myTopic"
	offset := 0
	bytes, err := persistence.BROKER.ConsumeMessage(topicName, offset)
	if err != nil {
		return err
	}
	workflows.ExecuteWorkflow(bytes)
	return nil
}
