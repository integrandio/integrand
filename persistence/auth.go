package persistence

func (dstore *Datastore) GetSecurablesByUserId(user_id int) ([]string, error) {
	securables := []string{}

	selectQuery := `SELECT securables.name FROM securables WHERE securables.id IN 
		(SELECT role_to_securable.securable_id FROM role_to_securable INNER JOIN user_to_role ON
			role_to_securable.role_id = user_to_role.role_id WHERE user_to_role.user_id = ?);`
	dstore.RWMutex.RLock()
	rows, err := dstore.db.Query(selectQuery, user_id)
	dstore.RWMutex.RUnlock()

	if err != nil {
		return securables, err
	}

	for rows.Next() {
		var securable string
		err := rows.Scan(&securable)
		if err != nil {
			return securables, err
		}
		securables = append(securables, securable)
	}
	return securables, err
}

type AuthRole string

const (
	SUPER_USER      AuthRole = "super_admin"
	INTEGRAND_ADMIN AuthRole = "integrand_admin"
	INTEGRAND_USER  AuthRole = "integrand_user"
)

type Securable string

const (
	READ_TOPIC     Securable = "read_topic"
	WRITE_TOPIC    Securable = "write_topic"
	READ_ENDPOINT  Securable = "read_endpoint"
	WRITE_ENDPOINT Securable = "write_endpoint"
	READ_WORKFLOW  Securable = "read_workflow"
	WRITE_WORKFLOW Securable = "write_workflow"
	READ_USER      Securable = "read_user"
	WRITE_USER     Securable = "write_user"
	READ_API_KEY   Securable = "read_api_key"
	WRITE_API_KEY  Securable = "write_api_key"
)

func (dstore *Datastore) createUserRole(userId int, role AuthRole) error {
	insert_query := `INSERT INTO user_to_role(user_id, role_id) SELECT users.id, roles.id FROM users INNER JOIN roles ON users.id = ? AND roles.name = ?;`

	dstore.RWMutex.Lock()
	_, err := dstore.db.Exec(insert_query, userId, role)
	dstore.RWMutex.Unlock()

	if err != nil {
		return err
	}
	return nil
}
