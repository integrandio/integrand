package services

import (
	"errors"
	"log"
	"reflect"
)

type Workflow struct {
	TopicName    string
	Offset       int
	FunctionName string
	Enabled      bool
}

type funcMap map[string]interface{}

var FUNC_MAP = funcMap{}

func init() {
	// Register all of our functions
	FUNC_MAP = map[string]interface{}{
		"ld_ld_sync": ld_ld_sync,
	}
}

func (workflow Workflow) Call(params ...interface{}) (result interface{}, err error) {
	f := reflect.ValueOf(FUNC_MAP[workflow.FunctionName])
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is out of index")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	res := f.Call(in)
	result = res[0].Interface()
	return
}

func ld_ld_sync(bytes []byte) error {
	log.Println("Executing")
	log.Println(string(bytes))
	return nil
}
