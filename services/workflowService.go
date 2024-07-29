package services

import (
	"errors"
	"integrand/persistence"
	"log/slog"
	"math/rand"
	"sync"
	"time"
)

const SLEEP_TIME int = 1
const MULTIPLYER int = 2
const MAX_BACKOFF int = 10

var (
	workflowMu sync.Mutex
)

func Workflower() error {
	for {
		workflowMu.Lock()
		currentWorkflows := append([]Workflow(nil), Workflows...)
		workflowMu.Unlock()

		var wg sync.WaitGroup
		for _, workflow := range currentWorkflows {
			wg.Add(1)
			go processWorkflow(&wg, workflow)
		}
		wg.Wait()

	}
}

func processWorkflow(wg *sync.WaitGroup, workflow Workflow) {
	defer wg.Done()
	sleep_time := SLEEP_TIME
	for {
		if !workflow.Enabled {
			return
		}

		bytes, err := persistence.BROKER.ConsumeMessage(workflow.TopicName, workflow.Offset)
		if err != nil {
			if err.Error() == "offset out of bounds" {
				slog.Warn(err.Error())
				time.Sleep(time.Duration(sleep_time) * time.Second)
				return // Exit the function, to be re-checked in the next cycle
			} else if err.Error() == "offset does not exist" {
				// I think this means no message in given topic?
				slog.Warn(err.Error())
				time.Sleep(time.Duration(sleep_time) * time.Second)
				return // Exit the function, to be re-checked in the next cycle
			} else {
				slog.Warn(err.Error())
				return // Something's wrong
			}
		}
		workflow.Call(bytes, workflow.SinkURL)
		workflow.Offset++
		sleep_time = SLEEP_TIME
	}
}

func GetWorkflows() ([]Workflow, error) {
	workflowMu.Lock()
	defer workflowMu.Unlock()
	return Workflows, nil
}

func DeleteWorkflow(id int) error {
	workflowMu.Lock()
	defer workflowMu.Unlock()
	for i, workflow := range Workflows {
		if workflow.Id == id {
			Workflows = append(Workflows[:i], Workflows[i+1:]...)
			return nil
		}
	}
	return errors.New("workflow not found")
}

func UpdateWorkflow(id int) (*Workflow, error) {
	workflowMu.Lock()
	defer workflowMu.Unlock()
	for i, workflow := range Workflows {
		if workflow.Id == id {
			Workflows[i].Enabled = !Workflows[i].Enabled
			return &Workflows[i], nil
		}
	}
	return nil, errors.New("workflow not found")
}

func GetWorkflow(id int) (*Workflow, error) {
	workflowMu.Lock()
	defer workflowMu.Unlock()
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
	// Get topic to use its offset for workflow creation
	topic, err := persistence.BROKER.GetTopic(topicName)
	if err != nil {
		slog.Error("topic with topicName " + topicName + " not found")
		return nil, errors.New("workflow with this functionName: " + functionName + " cannot be created")
	}

	newWorkflow := Workflow{
		Id:           rand.Int(),
		TopicName:    topicName,
		Offset:       topic.OldestOffset,
		FunctionName: functionName,
		Enabled:      true,
		SinkURL:      sinkURL,
	}

	workflowMu.Lock()
	Workflows = append(Workflows, newWorkflow)
	workflowMu.Unlock()
	return &newWorkflow, nil
}
