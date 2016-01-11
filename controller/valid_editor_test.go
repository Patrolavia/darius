package controller

import (
	"testing"

	"github.com/Patrolavia/darius/model"
	"github.com/Patrolavia/jsonapi"
)

func TestValidEditor(t *testing.T) {
	cfg := map[string]string{
		"RedisAddr":    ":6379",
		"SessSecret":   "1234567890",
		"SessName":     "mdpadtest",
		"ValidEditors": "valid@patrolavia.com",
	}
	pc := &Pad{db, sf, cfg}

	u, err := model.NewUser("invalid editor", "invalid@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		t.Fatalf("Cannot create invalid user to test valid editor: %s", err)
	}
	vu, err := model.NewUser("valid editor", "valid@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		t.Fatalf("Cannot create valid user to test valid editor: %s", err)
	}

	sess := session(sf, t)
	sess.Login(u)
	sess.Save()
	if err := sess.Err(); err != nil {
		t.Fatalf("Cannot login as invalid editor: %s", err)
	}

	pad := map[string]interface{}{
		"title":   "invalid edit",
		"content": "test invalid edit",
	}

	resp, err := jsonapi.HandlerTest(pc.Create).PostJSON("/api/create", sess.Cookie(), pad)
	if err != nil {
		t.Fatalf("Failed to get response of creating pad with invalid editor: %s", err)
	}

	if !testResult(resp.Body, false) {
		t.Errorf("api call should fail when invalid editor tries to create pad: %s", resp.Body.String())
	}

	if !testData(3, resp.Body, "code") {
		t.Errorf("api call should return errcode 1 when invalid editor tries to create pad: %s", resp.Body.String())
	}

	sess.Login(vu)
	sess.Save()
	if err := sess.Err(); err != nil {
		t.Fatalf("Cannot login as valid editor: %s", err)
	}

	resp, err = jsonapi.HandlerTest(pc.Create).PostJSON("/api/create", sess.Cookie(), pad)
	if err != nil {
		t.Fatalf("Failed to get response of creating pad with invalid editor: %s", err)
	}

	if !testResult(resp.Body, true) {
		t.Errorf("api call should be ok when valid editor tries to create pad: %s", resp.Body.String())
	}

	if !testData(0, resp.Body, "code") {
		t.Errorf("api call should return errcode 0 when valid editor tries to create pad: %s", resp.Body.String())
	}
}
