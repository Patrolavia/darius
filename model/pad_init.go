// This file is part of Darius. See License.txt for license information.

package model

import (
	"database/sql"
	"fmt"
)

// No "savePadQuery" because it requires transaction.
// Lacking some of tag or coop related queries because they need bulk execute.
var (
	newPadQuery,
	newTagQuery,
	newCoopQuery,
	loadPadQuery,
	deletePadQuery,
	deleteTagQuery,
	deletePadTagQuery,
	deleteCoopQuery,
	deletePadCoopQuery,
	updatePadQuery,
	listPadQuery,
	listTagQuery,
	listCoopQuery,
	findTagQuery,
	findCoopQuery *sql.Stmt
)

func initPadTable(db *sql.DB, ai string) (err error) {
	fn := func(query string) {
		if err == nil {
			_, err = db.Exec(query)
		}
	}

	pad := `CREATE TABLE IF NOT EXISTS pads (
id INTEGER PRIMARY KEY %s,
uid INTEGER,
title VARCHAR(64),
content TEXT,
html TEXT,
version INTEGER,
CONSTRAINT pad_user FOREIGN KEY (uid) REFERENCES users (id) ON DELETE RESTRICT ON UPDATE RESTRICT)`

	tag := `CREATE TABLE IF NOT EXISTS tags (
name VARCHAR(16),
pid INTEGER,
CONSTRAINT tag_pk PRIMARY KEY (name, pid),
CONSTRAINT tag_pad FOREIGN KEY (pid) REFERENCES pads (id) ON DELETE RESTRICT ON UPDATE RESTRICT)`

	coop := `CREATE TABLE IF NOT EXISTS coops (
uid INTEGER,
pid INTEGER,
CONSTRAINT coop_pad FOREIGN KEY (pid) REFERENCES pads (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
CONSTRAINT coop_user FOREIGN KEY (uid) REFERENCES users (id) ON DELETE RESTRICT ON UPDATE RESTRICT,
CONSTRAINT coop_pk PRIMARY KEY (uid, pid))`

	fn(fmt.Sprintf(pad, ai))
	fn(tag)
	fn(coop)

	cn := func(q string) (ret *sql.Stmt) {
		if err != nil {
			return nil
		}

		ret, err = db.Prepare(q)
		return
	}

	newPadQuery = cn(`INSERT INTO pads (uid,title,content,html,version) VALUES (?,?,?,?,1)`)
	newTagQuery = cn(`INSERT INTO tags (name,pid) VALUES (?,?)`)
	newCoopQuery = cn(`INSERT INTO coops (uid,pid) VALUES (?,?)`)
	loadPadQuery = cn(`SELECT uid,title,content,html,version FROM pads WHERE id=?`)
	deletePadQuery = cn(`DELETE FROM pads WHERE id=?`)
	deleteTagQuery = cn(`DELETE FROM tags WHERE name=? AND pid=?`)
	deletePadTagQuery = cn(`DELETE FROM tags WHERE pid=?`)
	deleteCoopQuery = cn(`DELETE FROM coops WHERE uid=? AND pid=?`)
	deletePadCoopQuery = cn(`DELETE FROM coops WHERE pid=?`)
	updatePadQuery = cn(`UPDATE pads SET title=?,content=?,html=?,version=version+1 WHERE id=? AND version=?`)
	listPadQuery = cn(`SELECT id,uid,title FROM pads`)
	listTagQuery = cn(`SELECT name,pid FROM tags`)
	listCoopQuery = cn(`SELECT uid,pid FROM coops`)
	findTagQuery = cn(`SELECT name FROM tags WHERE pid=?`)
	findCoopQuery = cn(`SELECT uid FROM coops WHERE pid=?`)
	return
}
