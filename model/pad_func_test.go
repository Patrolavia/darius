package model

import (
	"database/sql"
	"log"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func (o *PadContent) equals(n *PadContent) (ret bool) {
	if o.ID != n.ID || o.UID != n.UID || o.Title != n.Title || o.Content != n.Content || o.HTML != n.HTML {
		return
	}

	if !reflect.DeepEqual(o.Tags, n.Tags) {
		return
	}

	return reflect.DeepEqual(o.Cooperators, n.Cooperators)
}

var (
	db   *sql.DB
	u    *User
	coop *User
)

func init() {
	con, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Error opening db connection: %s", err)
	}
	db = con
	if err := InitSqlite3(db); err != nil {
		log.Fatalf("Error preparing table and data: %s", err)
	}
	u, err = NewUser("Test user", "user@test.com", "https://ronmi.tw/logo128.png")
	if err != nil {
		log.Fatalf("Error preparing test user: %s", err)
	}
	coop, err = NewUser("Test coop", "coop@test.com", "https://ronmi.tw/logo128.png")
	if err != nil {
		log.Fatalf("Error preparing test cooperator: %s", err)
	}
}

func TestNewPad(t *testing.T) {
	pad, err := NewPad(db, u.ID, "title", "content", []string{"tag1", "tag2"}, []int{coop.ID})
	if err != nil {
		t.Fatalf("Error creating pad: %s", err)
	}

	i := func(actual, expect int, msg string) {
		if expect != actual {
			t.Errorf("%s: expected %d, got %d.", msg, expect, actual)
		}
	}
	s := func(actual, expect, msg string) {
		if expect != actual {
			t.Errorf("%s: expected %s, got %s.", msg, expect, actual)
		}
	}

	i(pad.UID, u.ID, "Pad owner")
	s(pad.Title, "title", "Pad title")
	s(pad.Content, "content", "Pad content")
	i(len(pad.Tags), 2, "Pag tag count")
	s(pad.Tags[0], "tag1", "Pag tag#1")
	s(pad.Tags[1], "tag2", "Pag tag#2")
	i(len(pad.Cooperators), 1, "Pad coop count")
	i(pad.Cooperators[0], coop.ID, "Pad coop #1")
}

func TestLoad(t *testing.T) {
	pad, err := NewPad(db, u.ID, "load", "content", []string{"tag1", "tag2"}, []int{coop.ID})
	if err != nil {
		t.Fatalf("Error creating pad for later loading: %s", err)
	}

	actual, err := LoadPad(pad.ID)
	if err != nil {
		t.Fatalf("Error loading just created pad#%d: %s", pad.ID, err)
	}
	if !reflect.DeepEqual(*actual, *pad) {
		t.Errorf("Loaded pad different with original pad: %#v", *actual)
	}
}

func TestSave(t *testing.T) {
	pad, err := NewPad(db, u.ID, "save", "content", []string{"tag1", "tag2"}, []int{coop.ID})
	if err != nil {
		t.Fatalf("Error creating pad for later saving: %s", err)
	}

	pad.Title = "save2"
	if err := pad.Save(db); err != nil {
		t.Fatalf("Cannot save modified title: %s", err)
	}

	actual, err := LoadPad(pad.ID)
	if err != nil {
		t.Fatalf("Error loading after modifing title: %s", err)
	}
	if !actual.equals(pad) {
		t.Errorf("Does not load right data after modifing title: %#v", *actual)
	}

	pad = actual
	pad.Tags = []string{"tag3"}
	if err := pad.Save(db); err != nil {
		t.Fatalf("Cannot save changed tag: %s", err)
	}
	actual, err = LoadPad(pad.ID)
	if err != nil {
		t.Fatalf("Error loading after changing tag: %s", err)
	}
	if !actual.equals(pad) {
		t.Errorf("Does not load right data after changing tag: %#v", actual.Pad)
	}

	pad = actual
	pad.Tags = []string{}
	if err := pad.Save(db); err != nil {
		t.Fatalf("Cannot save deleted tag: %s", err)
	}
	actual, err = LoadPad(pad.ID)
	if err != nil {
		t.Fatalf("Error loading after deleting all tags: %s", err)
	}
	if !actual.equals(pad) {
		t.Errorf("Does not load right data after deleting all tags: %#v", actual.Pad)
	}

	pad = actual
	pad.Tags = []string{"tag4"}
	if err := pad.Save(db); err != nil {
		t.Fatalf("Cannot save inserted tag: %s", err)
	}
	actual, err = LoadPad(pad.ID)
	if err != nil {
		t.Fatalf("Error loading after inserting tag: %s", err)
	}
	if !actual.equals(pad) {
		t.Errorf("Does not load right data after inserting tag: %#v", actual.Pad)
	}
}
