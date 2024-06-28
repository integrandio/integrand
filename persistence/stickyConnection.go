package persistence

import (
	"time"
)

// Dat sticky sticky
type StickyConnection struct {
	RouteID      string    `json:"id,omitempty"`
	Security_key string    `json:"securityKey,omitempty"`
	TopicName    string    `json:"topicName,omitempty"`
	LastModified time.Time `json:"lastModified,omitempty"`
	UserId       int       `json:"userId,omitempty"`
}

func (dstore *Datastore) GetStickeyConnectionBySecurityKey(id string, security_key string) (StickyConnection, error) {
	selectQuery := "SELECT id, security_key, topic_name, last_modified, user_id FROM stickey_connections WHERE id=? and security_key=?;"
	dstore.RWMutex.RLock()
	row := dstore.db.QueryRow(selectQuery, id, security_key)
	dstore.RWMutex.RUnlock()

	var stickey_connection StickyConnection

	err := row.Scan(&stickey_connection.RouteID, &stickey_connection.Security_key, &stickey_connection.TopicName, &stickey_connection.LastModified, &stickey_connection.UserId)
	if err != nil {
		return stickey_connection, err
	}
	return stickey_connection, nil
}

func (dstore *Datastore) GetStickeyConnectionByUser(id string, userId int) (StickyConnection, error) {
	selectQuery := "SELECT id, security_key, topic_name, last_modified, user_id FROM stickey_connections WHERE id=? and user_id=?;"
	dstore.RWMutex.RLock()
	row := dstore.db.QueryRow(selectQuery, id, userId)
	dstore.RWMutex.RUnlock()

	var stickey_connection StickyConnection

	err := row.Scan(&stickey_connection.RouteID, &stickey_connection.Security_key, &stickey_connection.TopicName, &stickey_connection.LastModified, &stickey_connection.UserId)
	if err != nil {
		return stickey_connection, err
	}
	return stickey_connection, nil
}

func (dstore *Datastore) GetAllStickyConnections(userId int) ([]StickyConnection, error) {
	stickey_connections := []StickyConnection{}
	selectQuery := "SELECT id, security_key, topic_name, last_modified FROM stickey_connections WHERE user_id=?;"
	dstore.RWMutex.RLock()
	rows, err := dstore.db.Query(selectQuery, userId)
	dstore.RWMutex.RUnlock()
	if err != nil {
		return stickey_connections, err
	}
	defer rows.Close()

	for rows.Next() {
		var stickey_connection StickyConnection
		err := rows.Scan(&stickey_connection.RouteID, &stickey_connection.Security_key, &stickey_connection.TopicName, &stickey_connection.LastModified)
		if err != nil {
			return stickey_connections, err
		}
		stickey_connections = append(stickey_connections, stickey_connection)
	}
	return stickey_connections, nil
}

func (dstore *Datastore) InsertStickyConnection(sticky_connection StickyConnection) (int, error) {
	insertQuery := "INSERT INTO stickey_connections(id, security_key, topic_name, user_id) VALUES(?, ?, ?, ?);"
	dstore.RWMutex.Lock()
	res, err := dstore.db.Exec(insertQuery, sticky_connection.RouteID, sticky_connection.Security_key, sticky_connection.TopicName, sticky_connection.UserId)
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

func (dstore *Datastore) DeleteStickyConnection(stickey_connection_id string, userId int) (int, error) {
	insertQuery := "DELETE FROM stickey_connections WHERE id=? AND user_id=?"
	dstore.RWMutex.Lock()
	res, err := dstore.db.Exec(insertQuery, stickey_connection_id, userId)
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
