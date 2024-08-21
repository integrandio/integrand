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
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	AuthType     AuthType  `json:"authType"`
	Password     string    `json:"-"`                  // Password is often omitted in JSON responses
	SocialID     string    `json:"socialId,omitempty"` // Omits if empty
	CreatedAt    time.Time `json:"createdAt"`
	LastModified time.Time `json:"lastModified,omitempty"`
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

func (dstore *Datastore) GetEmailUsers() ([]User, error) {
	var users []User
	selectQuery := "SELECT id, email, auth_type, created_at, last_modified FROM users WHERE auth_type=?;"
	dstore.RWMutex.RLock()
	rows, err := dstore.db.Query(selectQuery, EMAIL)
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
func (dstore *Datastore) CreateEmailUser(u User) (int, error) {
	insertQuery := `
	INSERT INTO users(email, password, auth_type) VALUES(?, ?, ?)
	RETURNING id;
	`
	var id int
	dstore.RWMutex.Lock()
	err := dstore.db.QueryRow(insertQuery, u.Email, u.Password, EMAIL).Scan(&id)
	dstore.RWMutex.Unlock()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (dstore *Datastore) DeleteEmailUser(id int) (int, error) {
	deleteQuery := "DELETE FROM users WHERE id=? and auth_type=?"
	dstore.RWMutex.Lock()
	res, err := dstore.db.Exec(deleteQuery, id, EMAIL)
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

func (dstore *Datastore) UpdateEmailUser(id int, password string) (User, error) {
	updateQuery := `
		UPDATE users
		SET password = ?, last_modified = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING id, email, auth_type, created_at, last_modified ;
	`

	dstore.RWMutex.Lock()
	row := dstore.db.QueryRow(updateQuery, id, password)
	dstore.RWMutex.Unlock()
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.AuthType, &user.CreatedAt, &user.LastModified)
	if err != nil {
		return user, err
	}
	return user, nil
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
