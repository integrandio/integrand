package services

import (
	"errors"
	"integrand/persistence"
	"log/slog"
	"math/rand"
	"time"
)

const SLEEP_TIME int = 1
const MULTIPLYER int = 2
const MAX_BACKOFF int = 10

func Workflower() error {
	// Wait 5 seconds so we don't run into any race conditions
	time.Sleep(5 * time.Second)

	sleep_time := SLEEP_TIME
	for {
		for i, workflow := range Workflows {
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
					slog.Warn("Why")
					return err
				}
			}
			workflow.Call(bytes, workflow.SinkURL)
			Workflows[i].Offset++
			sleep_time = 1
		}
	}
}

func GetWorkflows() ([]Workflow, error) {
	return Workflows, nil
}

func DeleteWorkflow(id uint32) error {
	for i, workflow := range Workflows {
		if workflow.Id == id {
			Workflows = append(Workflows[:i], Workflows[i+1:]...)
			return nil
		}
	}
	return errors.New("workflow not found")
}

func UpdateWorkflow(id uint32) (*Workflow, error) {
	for i, workflow := range Workflows {
		if workflow.Id == id {
			Workflows[i].Enabled = !Workflows[i].Enabled
			return &Workflows[i], nil
		}
	}
	return nil, errors.New("workflow not found")
}

func GetWorkflow(id uint32) (*Workflow, error) {
	for _, workflow := range Workflows {
		if workflow.Id == id {
			return &workflow, nil
		}
	}
	return nil, errors.New("workflow not found")
}

func CreateWorkflow(topicName string, functionName string, sinkURL string) (*Workflow, error) {
	_, ok := FUNC_MAP[functionName]
	if !ok {
		slog.Error("function not found")
		return nil, errors.New("workflow with this functionName: " + functionName + " cannot be created")
	}

	newWorkflow := Workflow{
		Id:           rand.Uint32(),
		TopicName:    topicName,
		Offset:       0,
		FunctionName: functionName,
		Enabled:      true,
		SinkURL:      sinkURL,
	}

	Workflows = append(Workflows, newWorkflow)
	return &newWorkflow, nil
}
