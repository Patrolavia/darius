package controller

import (
	"testing"

	"github.com/Patrolavia/jsonapi"
	"github.com/Patrolavia/darius/model"
)

func TestMe(t *testing.T) {
	u, err := model.NewUser("test me", "me@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		t.Fatalf("Cannot create user for testing me: %s", err)
	}

	resp, err := jsonapi.HandlerTest(uc().Me).Get("/api/me", "")
	if err != nil {
		t.Fatalf("While getting response of /api/me (not logged in): %s", err)
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

	resp, err = jsonapi.HandlerTest(uc().Me).Get("/api/me", cookie)
	if err != nil {
		t.Fatalf("While getting response of /api/me (not logged in): %s", err)
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

func TestUser(t *testing.T) {
	u, err := model.NewUser("test user", "user@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		t.Fatalf("Cannot create user for testing user: %s", err)
	}
	param := map[string]interface{}{
		"userid": []int{u.ID},
	}

	resp, err := jsonapi.HandlerTest(uc().User).PostJSON("/api/user", "", param)
	if err != nil {
		t.Fatalf("While getting response of /api/user: %s", err)
	}

	if !testResult(resp.Body, true) {
		t.Errorf("Error occurs when fetching user info: %s", resp.Body.String())
	}

	if !testArrayHas(u.ID, resp.Body) {
		t.Errorf("Did not fetch what we want: %s", resp.Body.String())
	}
}
