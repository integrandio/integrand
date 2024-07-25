package services

import (
	"errors"
	"integrand/persistence"
	"log/slog"
	"time"
)

const SLEEP_TIME int = 1
const MULTIPLYER int = 2
const MAX_BACKOFF int = 10

func init() {
	// Creating first workflow through function
	_, err := CreateWorkflow("test", "ld_ld_sync")
	if err != nil {
		slog.Error("Failed to create the first workflow", "error", err)
		return
	}
}

func Workflower() error {
	// Wait 5 seconds so we don't run into any race conditions
	time.Sleep(5 * time.Second)

	workflow := Workflows[0]

	sleep_time := SLEEP_TIME
	for {
		bytes, err := persistence.BROKER.ConsumeMessage(workflow.TopicName, workflow.Offset)
		if err != nil {
			if err.Error() == "offset out of bounds" {
				slog.Warn(err.Error())
				time.Sleep(time.Duration(sleep_time) * time.Second)
				if sleep_time < MAX_BACKOFF {
					sleep_time = sleep_time * MULTIPLYER
				}
				continue
			} else {
				return err
			}
		}
		workflow.Call(bytes)
		workflow.Offset++
		sleep_time = 1
	}
}

func GetWorkflows() ([]Workflow, error) {
	return Workflows, nil
}

func DeleteWorkflow(topicName string) error {
	for i, workflow := range Workflows {
		if workflow.TopicName == topicName {
			Workflows = append(Workflows[:i], Workflows[i+1:]...)
			return nil
		}
	}
	return errors.New("workflow not found")
}

func UpdateWorkflow(topicName string) (*Workflow, error) {
	for i, workflow := range Workflows {
		if workflow.TopicName == topicName {
			Workflows[i].Enabled = !Workflows[i].Enabled
			return &Workflows[i], nil
		}
	}
	return nil, errors.New("workflow not found")
}

func GetWorkflow(topicName string) (*Workflow, error) {
	for _, workflow := range Workflows {
		if workflow.TopicName == topicName {
			return &workflow, nil
		}
	}
	return nil, errors.New("workflow not found")
}

func CreateWorkflow(topicName string, functionName string) (*Workflow, error) {
	// We also should check if function exists in our function map
	_, ok := FUNC_MAP[functionName]
	if !ok {
		slog.Error("function not found")
		return nil, errors.New("workflow with this functionName: " + functionName + " cannot be created")
	}

	for _, workflow := range Workflows {
		if workflow.TopicName == topicName {
			return nil, errors.New("workflow with this topicName already exists")
		}
	}

	newWorkflow := Workflow{
		TopicName:    topicName,
		Offset:       0,
		FunctionName: functionName,
		Enabled:      true,
	}

	Workflows = append(Workflows, newWorkflow)
	return &newWorkflow, nil
}
