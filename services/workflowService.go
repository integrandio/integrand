package services

import (
	"errors"
	"integrand/persistence"
	"log/slog"
)

func GetWorkflow(userID int, workflowID int) (persistence.Workflow, error) {
	err := isUserAuthorized(userID, "read_workflow")
	if err != nil {
		return persistence.Workflow{}, err
	}
	return persistence.DATASTORE.GetWorkflow(workflowID)
}

func GetWorkflows(userID int) ([]persistence.Workflow, error) {
	err := isUserAuthorized(userID, "read_workflow")
	if err != nil {
		return nil, err
	}
	return persistence.DATASTORE.GetWorkflows()
}

func GetEnabledWorkflows(userID int) ([]persistence.Workflow, error) {
	err := isUserAuthorized(userID, "read_workflow")
	if err != nil {
		return nil, err
	}
	return persistence.DATASTORE.GetEnabledWorkflows()
}

func DeleteWorkflow(userID int, workflowID int) (int, error) {
	err := isUserAuthorized(userID, "write_workflow")
	if err != nil {
		return 0, err
	}
	return persistence.DATASTORE.DeleteWorkflow(workflowID)
}

func UpdateWorkflow(userID int, workflowID int) (persistence.Workflow, error) {
	err := isUserAuthorized(userID, "write_workflow")
	if err != nil {
		return persistence.Workflow{}, err
	}
	workflow, err := persistence.DATASTORE.UpdateWorkflow(workflowID)
	if err != nil {
		slog.Error("Failed to update workflow", "error", err)
		return persistence.Workflow{}, err
	}
	return workflow, nil
}

func CreateWorkflow(topicName string, functionName string, sinkURL string, userID int) (persistence.Workflow, error) {
	err := isUserAuthorized(userID, "write_workflow")
	if err != nil {
		return persistence.Workflow{}, err
	}

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
		TopicName:    topicName,
		Offset:       topic.OldestOffset,
		FunctionName: functionName,
		Enabled:      true,
		SinkURL:      sinkURL,
	}

	id, last_modified, err := persistence.DATASTORE.InsertWorkflow(newWorkflow)
	if err != nil {
		return newWorkflow, err
	}
	newWorkflow.Id = id
	newWorkflow.LastModified = last_modified
	return newWorkflow, nil
}

func GetAvaliableWorkflowFunctions(userID int) ([]string, error) {
	err := isUserAuthorized(userID, "read_workflow")
	if err != nil {
		return []string{}, err
	}
	keys := make([]string, len(persistence.FUNC_MAP))
	i := 0
	for f := range persistence.FUNC_MAP {
		keys[i] = f
		i++
	}
	return keys, nil
}
