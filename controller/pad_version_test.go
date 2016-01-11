package controller

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/Patrolavia/darius/model"
	"github.com/Patrolavia/jsonapi"
)

var pvTestUser *model.User

func init() {
	var err error
	pvTestUser, err = model.NewUser("test pv", "pv@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		log.Fatalf("Cannot create user to test pad version control: %s", err)
	}
}

func TestPadVersionMismatch(t *testing.T) {
	pad, err := model.NewPad(db, pvTestUser.ID, "pm version", "test pm version", nil, nil)
	if err != nil {
		t.Fatalf("Cannot create pad for version testing: %s", err)
	}
	sess := session(pc().SF, t)
	sess.Login(pvTestUser)
	sess.Save()
	if err := sess.Err(); err != nil {
		t.Fatalf("Failed to login for testing version: %s", err)
	}
	cookie := sess.Cookie()
	uri := fmt.Sprintf("/api/edit/%d", pad.ID)
	param := map[string]interface{}{
		"title":   "pm version modified",
		"content": "content",
		"version": pad.Version + 1,
	}

	fn := func() {
		resp, err := jsonapi.HandlerTest(pc().Edit).PostJSON(uri, cookie, param)
		if err != nil {
			t.Fatalf("Cannot run /api/edit for version test: %s", err)
		}

		if !testResult(resp.Body, false) {
			t.Errorf("Version (%d) mismatch but pad is updates: %s", param["version"].(int), resp.Body.String())
		}
		code, ok := testHasData(resp.Body, "code")
		if !ok {
			t.Errorf("Not returning code wher version mismatch: %s", resp.Body.String())
		}
		var c int
		switch v := code.(type) {
		case float32, float64:
			c = int(reflect.ValueOf(v).Float())
		case int64, int32, int16, int8, int:
			c = int(reflect.ValueOf(v).Int())
		default:
			t.Fatalf("Returned code is not integer: %#v", code)
		}
		if c != 5 {
			t.Errorf("Error code is not 5: %d", c)
		}
	}

	fn()
	param["version"] = pad.Version - 1
	fn()
}
