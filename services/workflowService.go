package services

import (
	"errors"
	"integrand/persistence"
	"log/slog"
	"math/rand"
)

func GetWorkflows() ([]persistence.Workflow, error) {
	return persistence.DATASTORE.GetWorkflows()
}

func GetEnabledWorkflows() ([]persistence.Workflow, error) {
	return persistence.DATASTORE.GetEnabledWorkflows()
}

func DeleteWorkflow(id int) (int, error) {
	return persistence.DATASTORE.DeleteWorkflow(id)
}

func UpdateWorkflow(id int) (persistence.Workflow, error) {
	workflow, err := persistence.DATASTORE.UpdateWorkflow(id)
	if err != nil {
		slog.Error("Failed to update workflow", "error", err)
		return persistence.Workflow{}, err
	}
	return workflow, nil
}
func GetWorkflow(id int) (persistence.Workflow, error) {
	return persistence.DATASTORE.GetWorkflow(id)
}

func CreateWorkflow(topicName string, functionName string, sinkURL string) (persistence.Workflow, error) {
	_, ok := persistence.FUNC_MAP[functionName]
	if !ok {
		slog.Error("function not found")
		return persistence.Workflow{}, errors.New("workflow with this functionName: " + functionName + " cannot be created")
	}
	// Get topic to use its offset for workflow creation
	topic, err := persistence.BROKER.GetTopic(topicName)
	if err != nil {
		slog.Error("topic with topicName " + topicName + " not found")
		return persistence.Workflow{}, errors.New("workflow with this functionName: " + functionName + " cannot be created")
	}

	newWorkflow := persistence.Workflow{
		Id:           rangeIn(0, 100),
		TopicName:    topicName,
		Offset:       topic.OldestOffset,
		FunctionName: functionName,
		Enabled:      true,
		SinkURL:      sinkURL,
	}

	last_modified, err := persistence.DATASTORE.InsertWorkflow(newWorkflow)
	if err != nil {
		return newWorkflow, err
	}
	newWorkflow.LastModified = last_modified
	return newWorkflow, nil
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}
