// This file is part of Darius. See License.txt for license information.

package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/Patrolavia/darius/common"
	"github.com/Patrolavia/darius/model"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db  *sql.DB
	sf  common.SessionFactory
	cfg common.Config
)

func uc() *User {
	return &User{sf, cfg}
}

func pc() *Pad {
	return &Pad{db, sf, cfg}
}

func init() {
	var err error
	if db, err = sql.Open("sqlite3", ":memory:"); err != nil {
		log.Fatalf("Cannot open db connection: %s", err)
	}

	if err = model.InitSqlite3(db); err != nil {
		log.Fatalf("Cannot init table scheme: %s", err)
	}

	redis := os.Getenv("REDIS")
	if redis == "" {
		log.Fatalf("Please set environment variable REDIS.")
	}

	cfg = map[string]string{
		"RedisAddr":  redis,
		"SessSecret": "1234567890",
		"SessName":   "mdpadtest",
	}

	_sf, err := common.BuildSession(cfg)
	if err != nil {
		log.Fatalf("Cannot init session factory: %s", err)
	}
	sf = _sf
}

func testResult(body *bytes.Buffer, ok bool) bool {
	var res Response
	d := json.NewDecoder(strings.NewReader(body.String()))
	if err := d.Decode(&res); err != nil {
		return false
	}
	return res.Result == ok
}

func testArrayHas(id int, body *bytes.Buffer) (eq bool) {
	var res Response
	d := json.NewDecoder(strings.NewReader(body.String()))
	if err := d.Decode(&res); err != nil {
		return false
	}
	data, ok := res.Data.([]interface{})
	if !ok {
		return false
	}

	for _, v := range data {
		m, ok := v.(map[string]interface{})
		if !ok {
			break
		}
		if i, ok := m["id"]; ok && eqInt(id, i) {
			return true
		}
	}
	return false
}

func testHasData(body *bytes.Buffer, key string) (data interface{}, ok bool) {
	var res Response
	d := json.NewDecoder(strings.NewReader(body.String()))
	if err := d.Decode(&res); err != nil {
		return
	}

	m, ok := res.Data.(map[string]interface{})
	if !ok {
		return
	}

	data, ok = m[key]
	return
}

func testData(expect interface{}, body *bytes.Buffer, key string) (eq bool) {
	var res Response
	d := json.NewDecoder(strings.NewReader(body.String()))
	if err := d.Decode(&res); err != nil {
		return false
	}

	m, ok := res.Data.(map[string]interface{})
	if !ok {
		return
	}

	actual, ok := m[key]

	switch v := expect.(type) {
	case int:
		return eqInt(v, actual)
	case string:
		return eqStr(v, actual)
	}

	return eqRef(expect, actual)
}

func eqInt(expect int, actual interface{}) bool {
	a, ok := actual.(float64)
	return ok && float64(expect) == a
}

func eqStr(expect string, actual interface{}) bool {
	a, ok := actual.(string)
	return ok && expect == a
}

func eqRef(expect interface{}, actual interface{}) bool {
	return reflect.DeepEqual(expect, actual)
}

type sessTest struct {
	req *http.Request
	rec *httptest.ResponseRecorder
	common.Session
}

func session(sf common.SessionFactory, t *testing.T) *sessTest {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil
	}
	rec := httptest.NewRecorder()
	sess := sf.Get(req)
	if err != nil {
		t.Fatalf("Cannot initialize session")
	}
	return &sessTest{req, rec, sess}
}

func (s *sessTest) Save() {
	s.Session.Save(s.req, s.rec)
}

func (s *sessTest) Cookie() string {
	return s.rec.HeaderMap.Get("Set-Cookie")
}

func (s *sessTest) Login(u *model.User) {
	s.Session.Set("uid", fmt.Sprintf("%d", u.ID))
}
