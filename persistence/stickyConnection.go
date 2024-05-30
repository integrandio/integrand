package persistence

import (
	"time"
)

// Dat sticky sticky
type StickyConnection struct {
	RouteID          string    `json:"id,omitempty"`
	ConnectionApiKey string    `json:"connectionKey,omitempty"`
	TopicName        string    `json:"topicName,omitempty"`
	LastModified     time.Time `json:"lastModified,omitempty"`
}

func (dstore *Datastore) GetStickeyConnection(id string) (StickyConnection, error) {
	selectQuery := "SELECT id, connection_api_key, topic_name, last_modified FROM stickey_connections WHERE id=?;"
	dstore.RWMutex.RLock()
	row := dstore.db.QueryRow(selectQuery, id)
	dstore.RWMutex.RUnlock()

	var stickey_connection StickyConnection

	err := row.Scan(&stickey_connection.RouteID, &stickey_connection.ConnectionApiKey, &stickey_connection.TopicName, &stickey_connection.LastModified)
	if err != nil {
		return stickey_connection, err
	}
	return stickey_connection, nil
}

func (dstore *Datastore) GetAllStickyConnections() ([]StickyConnection, error) {
	stickey_connections := []StickyConnection{}
	selectQuery := "SELECT id, connection_api_key, topic_name, last_modified FROM stickey_connections;"
	dstore.RWMutex.RLock()
	rows, err := dstore.db.Query(selectQuery)
	dstore.RWMutex.RUnlock()
	if err != nil {
		return stickey_connections, err
	}
	defer rows.Close()

	for rows.Next() {
		var stickey_connection StickyConnection
		err := rows.Scan(&stickey_connection.RouteID, &stickey_connection.ConnectionApiKey, &stickey_connection.TopicName, &stickey_connection.LastModified)
		if err != nil {
			return stickey_connections, err
		}
		stickey_connections = append(stickey_connections, stickey_connection)
	}

	return stickey_connections, nil
}

func (dstore *Datastore) InsertStickyConnection(sticky_connection StickyConnection) (int, error) {
	insertQuery := "INSERT INTO stickey_connections(id, connection_api_key, topic_name) VALUES(?, ?, ?);"
	dstore.RWMutex.Lock()
	res, err := dstore.db.Exec(insertQuery, sticky_connection.RouteID, sticky_connection.ConnectionApiKey, sticky_connection.TopicName)
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

func (dstore *Datastore) DeleteStickyConnection(stickey_connection_id string) (int, error) {
	insertQuery := "DELETE FROM stickey_connections WHERE id=?"
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
