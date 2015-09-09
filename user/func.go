package user

import (
	"database/sql"
	"fmt"
)

var (
	insertQuery *sql.Stmt
	deleteQuery *sql.Stmt
	updateQuery *sql.Stmt
	emailQuery  *sql.Stmt
	idQuery     *sql.Stmt
	listQuery   *sql.Stmt
)

func initTable(db *sql.DB, ai string) (err error) {
	query := `CREATE TABLE IF NOT EXISTS users (
id INTEGER PRIMARY KEY %s,
name VARCHAR(32),
email VARCHAR(128) UNIQUE,
image TEXT)`
	query = fmt.Sprintf(query, ai)
	if _, err = db.Exec(query); err != nil {
		return
	}

	insertQuery, err = db.Prepare(`INSERT INTO users (name,email,image) VALUES (?,?,?)`)
	if err != nil {
		return
	}

	deleteQuery, err = db.Prepare(`DELETE FROM users WHERE id=?`)
	if err != nil {
		return
	}

	updateQuery, err = db.Prepare(`UPDATE users SET name=?,email=?,image=? WHERE id=?`)
	if err != nil {
		return
	}

	emailQuery, err = db.Prepare(`SELECT id,name,image FROM users WHERE email=?`)
	if err != nil {
		return
	}

	idQuery, err = db.Prepare(`SELECT name,email,image FROM users WHERE id=?`)
	if err != nil {
		return
	}

	listQuery, err = db.Prepare(`SELECT id,name,email,image FROM users`)
	return
}

// InitMysql initializes table using mysql syntax.
func InitMysql(db *sql.DB) error {
	return initTable(db, "AUTO_INCREMENT")
}

// InitSqlite3 initializes table using sqlite syntax.
func InitSqlite3(db *sql.DB) error {
	return initTable(db, "AUTOINCREMENT")
}

// List all users.
func List() (users []*User, err error) {
	rows, err := listQuery.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var (
			id    int
			name  string
			email string
			image string
		)
		rows.Scan(&id, &name, &email, &image)
		u := &User{id, name, email, image}
		users = append(users, u)
	}
	err = rows.Err()
	return
}

// New inserts a new record into db.
func New(name, email, image string) (u *User, err error) {
	res, err := insertQuery.Exec(name, email, image)
	if err != nil {
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		return
	}

	return &User{int(id), name, email, image}, nil
}

// Find a user from db by email.
func Find(email string) (u *User, err error) {
	row := emailQuery.QueryRow(email)
	var (
		id    int
		name  string
		image string
	)

	if err = row.Scan(&id, &name, &image); err != nil {
		return
	}

	return &User{id, name, email, image}, nil
}

// Load user record from db by id.
func Load(id int) (u *User, err error) {
	row := idQuery.QueryRow(id)
	var (
		name  string
		email string
		image string
	)

	if err = row.Scan(&name, &email, &image); err != nil {
		return
	}

	return &User{id, name, email, image}, nil
}

// Save user record into db.
func (u *User) Save() (err error) {
	res, err := updateQuery.Exec(u.Name, u.Email, u.Image, u.ID)
	if err != nil {
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rows != 1 {
		err = fmt.Errorf("When updating user#%d: affected %d row(s)", u.ID, rows)
	}
	return
}

// Delete user record from db.
func (u *User) Delete() (err error) {
	if u.ID == 0 {
		err = fmt.Errorf("User %s has been deleted, you can't delete it again.", u.Email)
		return
	}

	_, err = deleteQuery.Exec(u.ID)
	return
}
