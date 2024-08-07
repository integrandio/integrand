package persistence

import (
	"log"
	"time"
)

type AuthType string

const (
	EMAIL  AuthType = "email"
	GOOGLE AuthType = "google"
	GITHUB AuthType = "github"
)

type User struct {
	ID           int
	Email        string
	AuthType     AuthType
	Password     string
	SocialID     string
	CreatedAt    time.Time
	LastModified time.Time
}

func (dstore *Datastore) getAllUsers() ([]User, error) {
	var users []User
	selectQuery := "SELECT id, email, auth_type, created_at, last_modified FROM users;"
	dstore.RWMutex.RLock()
	rows, err := dstore.db.Query(selectQuery)
	dstore.RWMutex.RUnlock()
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.AuthType, &user.CreatedAt, &user.LastModified)
		if err != nil {
			rows.Close()
			return users, err
		}
		users = append(users, user)
	}
	rows.Close()
	return users, nil
}

func (dstore *Datastore) GetUserByID(id int) (User, error) {
	selectQuery := "SELECT id, email, password FROM users WHERE id=?;"
	dstore.RWMutex.RLock()
	row := dstore.db.QueryRow(selectQuery, id)
	dstore.RWMutex.RUnlock()
	var u User
	err := row.Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (dstore *Datastore) GetEmailUser(email string) (User, error) {
	selectQuery := "SELECT id, email, password FROM users WHERE email=?;"
	dstore.RWMutex.RLock()
	row := dstore.db.QueryRow(selectQuery, email, EMAIL)
	dstore.RWMutex.RUnlock()

	var u User
	err := row.Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (dstore *Datastore) CreateEmailUser(u User) (int, error) {
	insertQuery := "INSERT INTO users(id, email, password, auth_type) VALUES(NULL, ?, ?, ?);"
	dstore.RWMutex.Lock()
	res, err := dstore.db.Exec(insertQuery, u.Email, u.Password, EMAIL)
	dstore.RWMutex.Unlock()
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (dstore *Datastore) GetSocialUser(email string) (User, error) {
	selectQuery := "SELECT id, socialID FROM users WHERE email=?;"
	dstore.RWMutex.RLock()
	row := dstore.db.QueryRow(selectQuery, email)
	dstore.RWMutex.RUnlock()

	var u User

	err := row.Scan(&u.ID, &u.SocialID)
	if err != nil {
		return u, err
	}
	u.AuthType = GOOGLE
	return u, nil
}

func (dstore *Datastore) CreateSocialUser(u User) (int, error) {
	log.Println(u.Email)
	insertQuery := "INSERT INTO users(id, email, auth_type, socialID) VALUES(NULL, ?, ?, ?);"
	dstore.RWMutex.Lock()
	res, err := dstore.db.Exec(insertQuery, u.Email, GOOGLE, u.SocialID)
	dstore.RWMutex.Unlock()
	if err != nil {
		return 0, err
	}
	rowsCreated, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(rowsCreated), nil
}
