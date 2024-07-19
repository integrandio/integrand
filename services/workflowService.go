package services

import (
	"integrand/persistence"
	"log/slog"
	"time"
)

const SLEEP_TIME int = 1
const MULTIPLYER int = 2
const MAX_BACKOFF int = 10

func Workflower() error {
	workflow := Workflow{
		TopicName:    "test",
		Offset:       0,
		FunctionName: "ld_ld_sync",
	}

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
