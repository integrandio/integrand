package persistence

import (
	"time"
)

// Dat sticky sticky
type Endpoint struct {
	RouteID      string    `json:"id,omitempty"`
	Security_key string    `json:"securityKey,omitempty"`
	TopicName    string    `json:"topicName,omitempty"`
	LastModified time.Time `json:"lastModified,omitempty"`
}

func (dstore *Datastore) GetEndpointBySecurityKey(id string, security_key string) (Endpoint, error) {
	selectQuery := "SELECT id, security_key, topic_name, last_modified FROM endpoints WHERE id=? and security_key=?;"
	dstore.RWMutex.RLock()
	row := dstore.db.QueryRow(selectQuery, id, security_key)
	dstore.RWMutex.RUnlock()

	var stickey_connection Endpoint
	err := row.Scan(&stickey_connection.RouteID, &stickey_connection.Security_key, &stickey_connection.TopicName, &stickey_connection.LastModified)
	if err != nil {
		return stickey_connection, err
	}
	return stickey_connection, nil
}

func (dstore *Datastore) GetEndpoint(id string) (Endpoint, error) {
	selectQuery := "SELECT id, security_key, topic_name, last_modified FROM endpoints WHERE id=?;"
	dstore.RWMutex.RLock()
	row := dstore.db.QueryRow(selectQuery, id)
	dstore.RWMutex.RUnlock()

	var stickey_connection Endpoint

	err := row.Scan(&stickey_connection.RouteID, &stickey_connection.Security_key, &stickey_connection.TopicName, &stickey_connection.LastModified)
	if err != nil {
		return stickey_connection, err
	}
	return stickey_connection, nil
}

func (dstore *Datastore) GetAllEndpoints() ([]Endpoint, error) {
	endpoints := []Endpoint{}
	selectQuery := "SELECT id, security_key, topic_name, last_modified FROM endpoints;"
	dstore.RWMutex.RLock()
	rows, err := dstore.db.Query(selectQuery)
	dstore.RWMutex.RUnlock()
	if err != nil {
		return endpoints, err
	}
	for rows.Next() {
		var stickey_connection Endpoint
		err := rows.Scan(&stickey_connection.RouteID, &stickey_connection.Security_key, &stickey_connection.TopicName, &stickey_connection.LastModified)
		if err != nil {
			rows.Close()
			return endpoints, err
		}
		endpoints = append(endpoints, stickey_connection)
	}
	rows.Close()
	return endpoints, nil
}

func (dstore *Datastore) InsertEndpoint(sticky_connection Endpoint) (int, error) {
	insertQuery := "INSERT INTO endpoints(id, security_key, topic_name) VALUES(?, ?, ?);"
	dstore.RWMutex.Lock()
	res, err := dstore.db.Exec(insertQuery, sticky_connection.RouteID, sticky_connection.Security_key, sticky_connection.TopicName)
	dstore.RWMutex.Unlock()
	if err != nil {
		return 0, err
	}
	rowsCreated, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsCreated), nil
}

func (dstore *Datastore) DeleteEndpoint(stickey_connection_id string) (int, error) {
	insertQuery := "DELETE FROM endpoints WHERE id=?"
	dstore.RWMutex.Lock()
	res, err := dstore.db.Exec(insertQuery, stickey_connection_id)
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
