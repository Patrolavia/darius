package model

import (
	"database/sql"
	"sort"

	bf "github.com/Ronmi/blackfriday"
)

func html(title, content string) string {
	buf := []byte(content)
	html_opt := 0 |
		bf.HTML_USE_SMARTYPANTS |
		bf.HTML_SMARTYPANTS_FRACTIONS |
		bf.HTML_SMARTYPANTS_LATEX_DASHES |
		bf.HTML_FOOTNOTE_RETURN_LINKS
	render := bf.HtmlRenderer(html_opt, title, "")
	res := bf.MarkdownOptions(buf, render, bf.Options{
		Extensions: 0 |
			bf.EXTENSION_NO_INTRA_EMPHASIS |
			bf.EXTENSION_TABLES |
			bf.EXTENSION_FENCED_CODE |
			bf.EXTENSION_AUTOLINK |
			bf.EXTENSION_STRIKETHROUGH |
			bf.EXTENSION_SPACE_HEADERS |
			bf.EXTENSION_HEADER_IDS |
			bf.EXTENSION_BACKSLASH_LINE_BREAK |
			bf.EXTENSION_DEFINITION_LISTS |
			bf.EXTENSION_HARD_LINE_BREAK |
			bf.EXTENSION_TITLEBLOCK |
			bf.EXTENSION_FOOTNOTES,
	})
	return string(res)
}

// New creates a new record in db, it uses transaction so you have to pass db connection to it.
func NewPad(db *sql.DB, uid int, title, content string, tags []string, coops []int) (pad *PadContent, err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	html := html(title, content)
	q := tx.Stmt(newPadQuery)
	res, err := q.Exec(uid, title, content, html)
	if err != nil {
		return
	}

	pid64, err := res.LastInsertId()
	if err != nil {
		return
	}
	pid := int(pid64)

	q = tx.Stmt(newTagQuery)
	for _, tag := range tags {
		if _, err = q.Exec(tag, pid); err != nil {
			return
		}
	}

	q = tx.Stmt(newCoopQuery)
	for _, coop := range coops {
		if _, err = q.Exec(coop, pid); err != nil {
			return
		}
	}

	if err = tx.Commit(); err != nil {
		return
	}

	pad = &PadContent{
		&Pad{pid, uid, title, make([]string, len(tags)), tags, make([]int, len(coops)), coops},
		content, html, 1,
	}
	for k, v := range pad.oldTags {
		pad.Tags[k] = v
	}
	for k, v := range pad.oldCoops {
		pad.Cooperators[k] = v
	}
	return
}

// Load record from db by id
func LoadPad(id int) (pad *PadContent, err error) {
	row := loadPadQuery.QueryRow(id)
	var (
		uid                  int
		title, content, html string
		version              int
		tags                 []string
		coops                []int
	)
	if err = row.Scan(&uid, &title, &content, &html, &version); err != nil {
		return
	}

	rows, err := findTagQuery.Query(id)
	if err != nil {
		return
	}

	for rows.Next() {
		var name string
		rows.Scan(&name)
		tags = append(tags, name)
	}
	if err = rows.Err(); err != nil {
		return
	}

	if rows, err = findCoopQuery.Query(id); err != nil {
		return
	}

	for rows.Next() {
		var coop int
		rows.Scan(&coop)
		coops = append(coops, coop)
	}
	if err = rows.Err(); err != nil {
		return
	}

	pad = &PadContent{
		&Pad{id, uid, title, make([]string, len(tags)), tags, make([]int, len(coops)), coops},
		content, html, version,
	}
	for k, v := range pad.oldTags {
		pad.Tags[k] = v
	}
	for k, v := range pad.oldCoops {
		pad.Cooperators[k] = v
	}
	return
}

// Save pad to db, it needs db connection to do transaction.
func (pad *PadContent) Save(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	q := tx.Stmt(updatePadQuery)
	res, err := q.Exec(pad.Title, pad.Content, html(pad.Title, pad.Content), pad.ID, pad.Version)
	if err != nil {
		return
	}

	row, err := res.RowsAffected()
	if err != nil {
		return
	}
	if row != 1 {
		return VersionError(pad.Version)
	}

	itag, dtag := pad.tagDiff()
	q = tx.Stmt(newTagQuery)
	for _, t := range itag {
		if _, err = q.Exec(t, pad.ID); err != nil {
			return
		}
	}
	q = tx.Stmt(deleteTagQuery)
	for _, t := range dtag {
		if _, err = q.Exec(t, pad.ID); err != nil {
			return
		}
	}

	icoop, dcoop := pad.coopDiff()
	q = tx.Stmt(newCoopQuery)
	for _, t := range icoop {
		if _, err = q.Exec(t, pad.ID); err != nil {
			return
		}
	}
	q = tx.Stmt(deleteCoopQuery)
	for _, t := range dcoop {
		if _, err = q.Exec(t, pad.ID); err != nil {
			return
		}
	}

	tx.Commit()
	pad.Version++
	return
}

func (pad *PadContent) Render() {
	pad.HTML = html(pad.Title, pad.Content)
}

func (pad *Pad) Delete(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	// delete tags
	delTag := tx.Stmt(deletePadTagQuery)
	if _, err = delTag.Exec(pad.ID); err != nil {
		return
	}

	// delete coops
	delCoop := tx.Stmt(deletePadCoopQuery)
	if _, err = delCoop.Exec(pad.ID); err != nil {
		return
	}

	// delete pad
	delPad := tx.Stmt(deletePadQuery)
	if _, err = delPad.Exec(pad.ID); err == nil {
		pad.ID = 0
		tx.Commit()
	}

	return
}

func (pad *Pad) CoopModified() bool {
	a, b := pad.coopDiff()
	return len(a) == 0 && len(b) == 0
}

func (pad *Pad) Sort() {
	sort.Sort(sort.IntSlice(pad.Cooperators))
	sort.Sort(sort.StringSlice(pad.Tags))
}
