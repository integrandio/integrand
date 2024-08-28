package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
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
		db_file = "file::memory:?cache=shared"
	} else {
		db_file = os.Getenv("DB_FILE_LOCATION")
	}
	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Intialize our tables
	query, err := os.ReadFile("data/scripts/def.sql")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(string(query)); err != nil {
		log.Println(err)
		return nil, err
	}

	//Intialize our auth table
	query, err = os.ReadFile("data/scripts/auth.sql")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if _, err := db.Exec(string(query)); err != nil {
		log.Println(err)
		return nil, err
	}

	// Intialize our tables
	query, err = os.ReadFile("data/scripts/auth_setup.sql")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(string(query)); err != nil {
		log.Println(err)
		return nil, err
	}

	// tables := []string{"roles", "securables", "role_to_securable", "user_to_role"}
	// for _, table := range tables {
	// 	// Execute the query
	// 	fmt.Println("Table ", table)
	// 	select_query := fmt.Sprintf("SELECT * FROM %s", table)
	// 	debug_table(db, select_query)
	// }.

	integrandDB := &Datastore{
		db:      db,
		RWMutex: &sync.RWMutex{},
	}

	return integrandDB, nil
}

func debug_table(db *sql.DB, select_query string) {
	rows, err := db.Query(select_query)
	if err != nil {
		panic(err.Error())
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		// Print data
		for i, value := range values {
			switch value := value.(type) {
			case nil:
				fmt.Println(columns[i], ": NULL")

			case []byte:
				fmt.Println(columns[i], ": ", string(value))

			default:
				fmt.Println(columns[i], ": ", value)
			}
		}
	}
	fmt.Println("-----------------------------------")
}
