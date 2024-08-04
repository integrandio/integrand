package persistence

import (
	"errors"
	"reflect"
	"time"
)

type Workflow struct {
	Id           int       `json:"id"`
	TopicName    string    `json:"topicName"`
	Offset       int       `json:"offset"`
	FunctionName string    `json:"functionName"`
	Enabled      bool      `json:"enabled"`
	SinkURL      string    `json:"sinkURL"`
	LastModified time.Time `json:"lastModified,omitempty"`
}

type funcMap map[string]interface{}

var FUNC_MAP = funcMap{}

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

func (dstore *Datastore) GetWorkflow(id int) (Workflow, error) {
	selectQuery := "SELECT id, topic_name, offset, function_name, enabled, sink_url, last_modified FROM workflows WHERE id=?;"
	dstore.RWMutex.RLock()
	row := dstore.db.QueryRow(selectQuery, id)
	dstore.RWMutex.RUnlock()

	var workflow Workflow

	err := row.Scan(&workflow.Id, &workflow.TopicName, &workflow.Offset, &workflow.FunctionName, &workflow.Enabled, &workflow.SinkURL, &workflow.LastModified)
	if err != nil {
		return workflow, err
	}
	return workflow, nil
}

func (dstore *Datastore) GetWorkflows() ([]Workflow, error) {
	workflows := []Workflow{}
	selectQuery := "SELECT id, topic_name, offset, function_name, enabled, sink_url, last_modified FROM workflows;"
	dstore.RWMutex.RLock()
	rows, err := dstore.db.Query(selectQuery)
	if err != nil {
		return workflows, err
	}

	defer rows.Close()
	for rows.Next() {
		var workflow Workflow
		err := rows.Scan(&workflow.Id, &workflow.TopicName, &workflow.Offset, &workflow.FunctionName, &workflow.Enabled, &workflow.SinkURL, &workflow.LastModified)
		if err != nil {
			return workflows, err
		}
		workflows = append(workflows, workflow)
	}
	dstore.RWMutex.RUnlock()
	err = rows.Err()
	if err != nil {
		return workflows, err
	}
	return workflows, nil
}

func (dstore *Datastore) GetEnabledWorkflows() ([]Workflow, error) {
	workflows := []Workflow{}
	selectQuery := "SELECT id, topic_name, offset, function_name, enabled, sink_url, last_modified FROM workflows WHERE enabled=true;"
	dstore.RWMutex.RLock()
	rows, err := dstore.db.Query(selectQuery)

	if err != nil {
		return workflows, err
	}
	defer rows.Close()
	for rows.Next() {
		var workflow Workflow
		err := rows.Scan(&workflow.Id, &workflow.TopicName, &workflow.Offset, &workflow.FunctionName, &workflow.Enabled, &workflow.SinkURL, &workflow.LastModified)
		if err != nil {
			return workflows, err
		}
		workflows = append(workflows, workflow)
	}
	dstore.RWMutex.RUnlock()
	return workflows, nil
}

func (dstore *Datastore) InsertWorkflow(workflow Workflow) (time.Time, error) {
	insertQuery := `
		INSERT INTO workflows(id, topic_name, offset, function_name, sink_url, last_modified)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		RETURNING id, last_modified;
	`

	dstore.RWMutex.Lock()
	row := dstore.db.QueryRow(insertQuery, workflow.Id, workflow.TopicName, workflow.Offset, workflow.FunctionName, workflow.SinkURL)
	dstore.RWMutex.Unlock()

	var insertedID int
	var lastModified time.Time

	err := row.Scan(&insertedID, &lastModified)
	if err != nil {
		return time.Time{}, err
	}

	return lastModified, nil
}

func (dstore *Datastore) DeleteWorkflow(id int) (int, error) {
	deleteQuery := "DELETE FROM workflows WHERE id=?"
	dstore.RWMutex.Lock()
	res, err := dstore.db.Exec(deleteQuery, id)
	dstore.RWMutex.Unlock()
	if err != nil {
		return 0, err
	}
	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsDeleted), nil
}

func (dstore *Datastore) UpdateWorkflow(id int) (Workflow, error) {
	updateQuery := `
		UPDATE workflows
		SET enabled = NOT enabled, last_modified = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING id, topic_name, offset, function_name, enabled, sink_url, last_modified;
	`

	dstore.RWMutex.Lock()
	row := dstore.db.QueryRow(updateQuery, id)
	dstore.RWMutex.Unlock()
	var workflow Workflow
	err := row.Scan(&workflow.Id, &workflow.TopicName, &workflow.Offset, &workflow.FunctionName, &workflow.Enabled, &workflow.SinkURL, &workflow.LastModified)
	if err != nil {
		return workflow, err
	}
	return workflow, err
}
