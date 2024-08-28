package persistence

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"-"` // Password is often omitted in JSON responses
	CreatedAt    time.Time `json:"createdAt"`
	LastModified time.Time `json:"lastModified,omitempty"`
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

func (dstore *Datastore) GetUserByEmail(email string) (User, error) {
	selectQuery := "SELECT id, email, password FROM users WHERE email=?;"
	dstore.RWMutex.RLock()
	row := dstore.db.QueryRow(selectQuery, email)
	dstore.RWMutex.RUnlock()

	var u User
	err := row.Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (dstore *Datastore) GetUsers() ([]User, error) {
	var users []User
	selectQuery := "SELECT id, email, created_at, last_modified FROM users;"
	dstore.RWMutex.RLock()
	rows, err := dstore.db.Query(selectQuery)
	dstore.RWMutex.RUnlock()
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.LastModified)
		if err != nil {
			rows.Close()
			return users, err
		}
		users = append(users, user)
	}
	rows.Close()
	return users, nil
}
func (dstore *Datastore) CreateUser(u User) (int, error) {
	insertQuery := `
	INSERT INTO users(email, password) VALUES(?, ?)
	RETURNING id;
	`
	var id int
	dstore.RWMutex.Lock()
	err := dstore.db.QueryRow(insertQuery, u.Email, u.Password).Scan(&id)
	dstore.RWMutex.Unlock()
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (dstore *Datastore) DeleteUserByID(id int) (int, error) {
	deleteQuery := "DELETE FROM users WHERE id=?"
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

func (dstore *Datastore) UpdateEmailUser(id int, newPassword string) (User, error) {
	updateQuery := `
		UPDATE users
		SET password = ?, last_modified = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING id, email, created_at, last_modified ;
	`

	dstore.RWMutex.Lock()
	row := dstore.db.QueryRow(updateQuery, newPassword, id)
	dstore.RWMutex.Unlock()
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.LastModified)
	if err != nil {
		return user, err
	}
	return user, nil
}
