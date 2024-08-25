package persistence

import (
	"time"
)

type ApiKey struct {
	Id        int       `json:"id,omitempty"`
	Key       string    `json:"key,omitempty"` // Update JSON tag to match the expected field name
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UserId    int       `json:"userId,omitempty"`
}

func (dstore *Datastore) GetApiKey(key string) (ApiKey, error) {
	selectQuery := "SELECT id, key, created_at, user_id FROM api_keys WHERE key=?"
	dstore.RWMutex.RLock()
	row := dstore.db.QueryRow(selectQuery, key)
	dstore.RWMutex.RUnlock()

	var api_key ApiKey

	err := row.Scan(&api_key.Id, &api_key.Key, &api_key.CreatedAt, &api_key.UserId)
	if err != nil {
		return api_key, err
	}
	return api_key, nil
}

// adds a new API key to the store
func (dstore *Datastore) InsertAPIKey(key string, userID int) (int, error) {
	insertQuery := "INSERT INTO api_keys(id, key, user_id) VALUES(NULL, ?, ?);"
	dstore.RWMutex.Lock()
	res, err := dstore.db.Exec(insertQuery, key, userID)
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

// removes an API key from the store
func (dstore *Datastore) DeleteAPIKey(key string, userID int) (int, error) {
	insertQuery := "DELETE FROM api_keys WHERE key=? AND user_id=?"
	dstore.RWMutex.Lock()
	res, err := dstore.db.Exec(insertQuery, key, userID)
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

func (dstore *Datastore) GetAPIKeys() ([]ApiKey, error) {
	selectQuery := "SELECT id, key, created_at FROM api_keys;"
	dstore.RWMutex.RLock()
	rows, err := dstore.db.Query(selectQuery)
	dstore.RWMutex.RUnlock()
	if err != nil {
		return nil, err
	}
	var apiKeys []ApiKey
	for rows.Next() {
		var apiKey ApiKey
		err := rows.Scan(&apiKey.Id, &apiKey.Key, &apiKey.CreatedAt)
		if err != nil {
			rows.Close()
			return nil, err
		}
		apiKeys = append(apiKeys, apiKey)
	}
	rows.Close()
	return apiKeys, nil
}
