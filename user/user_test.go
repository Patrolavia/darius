package user

import (
	"database/sql"
	"log"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	con, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Cannot open test db: %s", err)
	}
	db = con
	if err := initTable(db, "AUTOINCREMENT"); err != nil {
		log.Fatalf("Cannot initialize user table: %s", err)
	}
}

func TestNew(t *testing.T) {
	u, err := New("test new", "new@test.com", "https://ronmi.tw/logo128.png")
	if err != nil {
		t.Fatalf("Where creating new user: %s", err)
	}

	if u.ID != 1 {
		t.Errorf("Expected user id starts from 1, got %d.", u.ID)
	}

	u, err = New("test new", "new@test.com", "https://ronmi.tw/logo128.png")
	if err == nil {
		t.Errorf("Insert same data into db but nothing wrong? got new user id %d", u.ID)
	}
}

func TestLoad(t *testing.T) {
	expected, err := New("test load", "load@test.com", "https://ronmi.tw/logo128.png")
	if err != nil {
		t.Fatalf("While creating user to be load: %s", err)
	}

	actual, err := Load(expected.ID)
	if err != nil {
		t.Fatalf("While loading just created user: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected loading %#v, got %#v", expected, actual)
	}
}

func TestFind(t *testing.T) {
	expected, err := New("test find", "find@test.com", "https://ronmi.tw/logo128.png")
	if err != nil {
		t.Fatalf("While creating user to be found: %s", err)
	}

	actual, err := Find(expected.Email)
	if err != nil {
		t.Fatalf("While finding just created user: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect finding %#v, got %#v", expected, actual)
	}
}

func TestSave(t *testing.T) {
	expected, err := New("test save", "save@test.com", "https://ronmi.tw/logo128.png")
	if err != nil {
		t.Fatalf("While creating user to be saved: %s", err)
	}

	expected.Name = "save"
	expected.Email = "save2@test.com"
	expected.Image = "https://ronmi.tw/logo64.png"
	if err := expected.Save(); err != nil {
		t.Fatalf("While saving user %s: %s", expected.Email, err)
	}

	actual, err := Load(expected.ID)
	if err != nil {
		t.Fatalf("While loading saved user %d: %s", expected.ID, err)
	}

	fn := func(actual, expect, msg string) {
		if actual != expect {
			t.Errorf("%s not saved: %s", msg, actual)
		}
	}

	fn(actual.Name, expected.Name, "Name")
	fn(actual.Email, expected.Email, "Email")
	fn(actual.Image, expected.Image, "Image")
}

func TestList(t *testing.T) {
	expected, err := New("test list", "list@test.com", "https://ronmi.tw/logo128.png")
	if err != nil {
		t.Fatalf("While creating user to be listed: %s", err)
	}

	users, err := List()
	if err != nil {
		t.Fatalf("While listing users: %s", err)
	}

	var has bool

	for _, u := range users {
		if reflect.DeepEqual(u, expected) {
			has = true
			break
		}
	}

	if !has {
		t.Errorf("Cannot find just created users from list: %d items", len(users))
	}
}
