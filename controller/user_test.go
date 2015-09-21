package controller

import (
	"testing"

	"github.com/Patrolavia/jsonapi"
	"github.com/Patrolavia/mdpadgo/model"
)

func TestMe(t *testing.T) {
	u, err := model.NewUser("test me", "me@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		t.Fatalf("Cannot create user for testing me: %s", err)
	}

	resp, err := jsonapi.HandlerTest(uc().Me).Get("/api/user", "")
	if err != nil {
		t.Fatalf("While getting response of /api/user (not logged in): %s", err)
	}

	if !testResult(resp.Body, false) {
		t.Errorf("Got login data without logging in: %#v", resp.Body.String())
	}

	sess := session(uc().SF, t)
	sess.Login(u)
	sess.Save()
	if err := sess.Err(); err != nil {
		t.Fatalf("Failed to save session: %s", err)
	}
	cookie := sess.Cookie()

	resp, err = jsonapi.HandlerTest(uc().Me).Get("/api/user", cookie)
	if err != nil {
		t.Fatalf("While getting response of /api/user (not logged in): %s", err)
	}
	if !testResult(resp.Body, true) {
		t.Fatal("Cannot get login data")
	}

	if !testData(u.ID, resp.Body, "id") {
		t.Fatalf("User id mismatch, expected %d: %#v", u.ID, resp.Body.String())
	}
}

func TestUserList(t *testing.T) {
	u, err := model.NewUser("test userlist", "userlist@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		t.Fatalf("Cannot create user for testing userlist: %s", err)
	}

	resp, err := jsonapi.HandlerTest(uc().Users).Get("/api/users", "")
	if err != nil {
		t.Fatalf("While getting response of /api/users: %s", err)
	}

	if !testArrayHas(u.ID, resp.Body) {
		t.Errorf("Cannot find user#%d in response of /api/users: %s", u.ID, resp.Body.String())
	}
}
