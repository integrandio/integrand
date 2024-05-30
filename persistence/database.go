package persistence

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type NullFloat64 struct{ sql.NullFloat64 }

// NS extends sql null string
type NullString struct{ sql.NullString }

// MarshalJSON for NullFloat64
func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON for NullFloat64
func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nf.Float64)
	nf.Valid = (err == nil)
	return err
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}

type Datastore struct {
	db *sql.DB
	*sync.RWMutex
}

func setupConnection(isDevMode bool) (*Datastore, error) {
	var db_file string

	if isDevMode {
		db_file = ":memory:"
	} else {
		db_file = os.Getenv("DB_FILE_LOCATION")
	}
	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		log.Println(db_file)
		return nil, err
	}
	// Intialize our tables
	query, err := os.ReadFile("data/scripts/def.sql")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(string(query)); err != nil {
		return nil, err
	}

	integrandDB := &Datastore{
		db:      db,
		RWMutex: &sync.RWMutex{},
	}

	// Insert our root user into the db...
	// TODO: Clean this up
	plainPassword := os.Getenv("ROOT_PASSWORD")
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := User{
		Email:    os.Getenv("ROOT_EMAIL"),
		Password: string(password),
		AuthType: EMAIL,
	}
	_, err = integrandDB.CreateEmailUser(user)
	if err != nil {
		return nil, err
	}

	return integrandDB, nil
}
