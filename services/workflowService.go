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
		currentWorkflows := append([]Workflow(nil), WORKFLOWS...)
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
				// This error is returned when we're given an offset thats ahead of the commitlog
				slog.Debug(err.Error())
				time.Sleep(time.Duration(sleep_time) * time.Second)
				continue
			} else if err.Error() == "offset does not exist" {
				// This error is returned when we look for an offset and it does not exist becuase it can't be avaliable in the commitlog
				slog.Warn(err.Error())
				time.Sleep(time.Duration(sleep_time) * time.Second)
				return // Exit the function, to be re-checked in the next cycle
			} else {
				slog.Error(err.Error())
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
	return WORKFLOWS, nil
}

func DeleteWorkflow(id int) error {
	workflowMu.Lock()
	defer workflowMu.Unlock()
	for i, workflow := range WORKFLOWS {
		if workflow.Id == id {
			WORKFLOWS = append(WORKFLOWS[:i], WORKFLOWS[i+1:]...)
			return nil
		}
	}
	return errors.New("workflow not found")
}

func UpdateWorkflow(id int) (*Workflow, error) {
	workflowMu.Lock()
	defer workflowMu.Unlock()
	for i, workflow := range WORKFLOWS {
		if workflow.Id == id {
			WORKFLOWS[i].Enabled = !WORKFLOWS[i].Enabled
			return &WORKFLOWS[i], nil
		}
	}
	return nil, errors.New("workflow not found")
}

func GetWorkflow(id int) (*Workflow, error) {
	workflowMu.Lock()
	defer workflowMu.Unlock()
	for _, workflow := range WORKFLOWS {
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
		Id:           rangeIn(0, 100),
		TopicName:    topicName,
		Offset:       topic.OldestOffset,
		FunctionName: functionName,
		Enabled:      true,
		SinkURL:      sinkURL,
	}

	workflowMu.Lock()
	WORKFLOWS = append(WORKFLOWS, newWorkflow)
	workflowMu.Unlock()
	return &newWorkflow, nil
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}
