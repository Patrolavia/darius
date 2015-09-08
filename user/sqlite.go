// +build sqlite3

package user

import "database/sql"

func Init(db *sql.DB) error {
	return initTable(db, "AUTOINCREMENT")
}
