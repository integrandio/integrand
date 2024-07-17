package workflows

import "log"

func ExecuteWorkflow(bytes []byte) error {
	log.Println("Executing")
	log.Println(bytes)
	return nil
}
