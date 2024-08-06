package persistence

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SessionManager struct {
	cookieName  string     //private cookiename
	lock        sync.Mutex // protects session
	maxlifetime int64
}

func NewSessionManager(cookieName string, maxlifetime int64) *SessionManager {
	return &SessionManager{
		cookieName:  cookieName,
		maxlifetime: maxlifetime,
	}
}

func (manager *SessionManager) sessionId() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		log.Println(err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (manager *SessionManager) SessionStart(w http.ResponseWriter, r *http.Request) (SessionDB, error) {
	var session SessionDB
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid, err := manager.sessionId()
		if err != nil {
			return session, err
		}
		v := make(map[string]interface{})
		session, err = DATASTORE.CreateSession(sid, v)
		if err != nil {
			return session, err
		}
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, err = DATASTORE.GetSession(sid)
		if err != nil {
			if err == sql.ErrNoRows {
				// If there is no rows returned, we need to setup a new session and create a new cookie
				sid, err := manager.sessionId()
				if err != nil {
					return session, err
				}
				v := make(map[string]interface{})
				session, err = DATASTORE.CreateSession(sid, v)
				if err != nil {
					return session, err
				}
				cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
				http.SetCookie(w, &cookie)
			}
			return session, nil
		}
	}
	return session, nil
}

// Destroy sessionid
func (manager *SessionManager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		DATASTORE.DeleteSession(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (manager *SessionManager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	DATASTORE.GcSessions(manager.maxlifetime)
	time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })
}

type SessionDB struct {
	ID           string
	value        map[string]interface{}
	TimeAccessed time.Time
}

func (sessionDB *SessionDB) Set(key string, value interface{}) error {
	sessionDB.value[key] = value
	err := DATASTORE.UpdateSession(sessionDB.ID, sessionDB.value)
	return err
}

func (sessionDB *SessionDB) Get(key string) (interface{}, error) {
	// update session to update time....
	err := DATASTORE.UpdateSession(sessionDB.ID, nil)
	if err != nil {
		return nil, err
	}
	if v, ok := sessionDB.value[key]; ok {
		return v, nil
	}
	return nil, nil
}

func (sessionDB *SessionDB) Delete(key string) error {
	delete(sessionDB.value, key)
	err := DATASTORE.UpdateSession(sessionDB.ID, sessionDB.value)
	return err
}

func (st *SessionDB) SessionID() string {
	return st.ID
}

// Get session by id
func (dstore *Datastore) GetSession(sessionID string) (SessionDB, error) {
	selectQuery := "SELECT id, value, time_accessed FROM sessions WHERE id=?;"
	dstore.RWMutex.RLock()
	row := dstore.db.QueryRow(selectQuery, sessionID)
	dstore.RWMutex.RUnlock()

	var sdb SessionDB
	var valueString string
	err := row.Scan(&sdb.ID, &valueString, &sdb.TimeAccessed)
	if err != nil {
		log.Println(err)
		return sdb, err
	}

	// Unmarshal string to our type
	err = json.Unmarshal([]byte(valueString), &sdb.value)
	if err != nil {
		log.Println(err)
		return sdb, err
	}

	return sdb, nil
}

func (dstore *Datastore) CreateSession(sessionID string, value map[string]interface{}) (SessionDB, error) {
	var sdb SessionDB
	valueJson, err := json.Marshal(value)
	if err != nil {
		return sdb, err
	}
	insertQuery := "INSERT INTO sessions(id, value) VALUES (?, json(?)) RETURNING time_accessed;"
	dstore.RWMutex.Lock()
	row := dstore.db.QueryRow(insertQuery, sessionID, string(valueJson))
	dstore.RWMutex.Unlock()
	err = row.Scan(&sdb.TimeAccessed)
	if err != nil {
		return sdb, err
	}
	sdb.ID = sessionID
	sdb.value = value
	return sdb, err
}

func (dstore *Datastore) UpdateSession(sessionID string, value map[string]interface{}) error {
	beginingUpdateQuery := "UPDATE sessions SET "
	endUpdateQuery := "time_accessed=? WHERE ID=?;"
	var interf []any
	if value != nil {
		beginingUpdateQuery = beginingUpdateQuery + "value=?, "
		valueJSON, err := json.Marshal(value)
		if err != nil {
			return err
		}
		interf = append(interf, string(valueJSON))
	}
	updateQuery := beginingUpdateQuery + endUpdateQuery
	interf = append(interf, time.Now())
	interf = append(interf, sessionID)
	dstore.RWMutex.Lock()
	_, err := dstore.db.Exec(updateQuery, interf...)
	dstore.RWMutex.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (dstore *Datastore) DeleteSession(sessionID string) error {
	deleteQuery := "DELETE FROM sessions ID=?;"
	dstore.RWMutex.Lock()
	_, err := dstore.db.Exec(deleteQuery, sessionID)
	dstore.RWMutex.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (dstore *Datastore) GcSessions(maxlifetime int64) error {
	selectQuery := "SELECT id, time_accessed FROM sessions;"
	deleteQuery := "DELETE FROM sessions ID=?;"
	dstore.RWMutex.Lock()
	rows, err := dstore.db.Query(selectQuery)
	dstore.RWMutex.Unlock()
	if err != nil {
		return err
	}
	for rows.Next() {
		var sdb SessionDB
		err = rows.Scan(&sdb.ID, &sdb.TimeAccessed)
		if err != nil {
			rows.Close()
			return err
		}
		if (sdb.TimeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			_, err := dstore.db.Exec(deleteQuery, sdb.ID)
			if err != nil {
				return err
			}
		}
	}
	rows.Close()
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
