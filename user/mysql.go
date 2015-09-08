// +build mysql

package user

import "database/sql"

func Init(db *sql.DB) error {
	return initTable(db, "AUTO_INCREMENT")
}
