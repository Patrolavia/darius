package controller

import (
	"fmt"
	"log"
	"testing"

	"github.com/Patrolavia/darius/model"
	"github.com/Patrolavia/jsonapi"
)

var pdAclUser, pdAclCoop, pdAclGuest *model.User
var pdAclPad *model.PadContent
var pdAclUri string

func init() {
	var err error
	pdAclUser, err = model.NewUser("test pd user", "pd_user@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		log.Fatalf("Cannot create user for testing pad del acl: %s", err)
	}

	pdAclCoop, err = model.NewUser("test pd coop", "pd_coop@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		log.Fatalf("Cannot create coop for testing pad del acl: %s", err)
	}

	pdAclGuest, err = model.NewUser("test pd guest", "pd_guest@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		log.Fatalf("Cannot create guest for testing pad del acl: %s", err)
	}

	pdAclPad, err = model.NewPad(db, pdAclUser.ID, "pd acl", "test pd acl", nil, []int{pdAclCoop.ID})
	if err != nil {
		log.Fatalf("Cannot create pad for testing pad del act: %s", err)
	}

	pdAclUri = fmt.Sprintf("/api/delete/%d", pdAclPad.ID)
}

func TestPadDelAclNoLogin(t *testing.T) {
	resp, err := jsonapi.HandlerTest(pc().Delete).Get(pdAclUri, "")
	if err != nil {
		t.Fatalf("Error occured while testing pad delete ACL without login: %s", err)
	}

	if !testResult(resp.Body, false) {
		t.Errorf("Delete success without login: %s", resp.Body.String())
	}

	code, ok := testHasData(resp.Body, "code")
	if !ok {
		t.Errorf("Error code not found in deletion without login: %s", resp.Body.String())
	}
	if !eqInt(1, code) {
		t.Errorf("Expected error code is 1 when delete without login, got %d", code)
	}
}

func TestPadDelAclGuest(t *testing.T) {
	sess := session(pc().SF, t)
	sess.Login(pdAclGuest)
	sess.Save()
	if err := sess.Err(); err != nil {
		t.Fatalf("Failed to login with guest when test delete ACL: %s", err)
	}
	resp, err := jsonapi.HandlerTest(pc().Delete).Get(pdAclUri, sess.Cookie())
	if err != nil {
		t.Fatalf("Error occured while testing pad delete ACL with guest: %s", err)
	}

	if !testResult(resp.Body, false) {
		t.Errorf("Delete success with guest: %s", resp.Body.String())
	}

	code, ok := testHasData(resp.Body, "code")
	if !ok {
		t.Errorf("Error code not found in deletion with guest: %s", resp.Body.String())
	}
	if !eqInt(3, code) {
		t.Errorf("Expected error code is 3 when delete with guest, got %d", code)
	}
}

func TestPadDelAclCoop(t *testing.T) {
	sess := session(pc().SF, t)
	sess.Login(pdAclCoop)
	sess.Save()
	if err := sess.Err(); err != nil {
		t.Fatalf("Failed to login with coop when test delete ACL: %s", err)
	}
	resp, err := jsonapi.HandlerTest(pc().Delete).Get(pdAclUri, sess.Cookie())
	if err != nil {
		t.Fatalf("Error occured while testing pad delete ACL with coop: %s", err)
	}

	if !testResult(resp.Body, false) {
		t.Errorf("Delete success with coop: %s", resp.Body.String())
	}

	code, ok := testHasData(resp.Body, "code")
	if !ok {
		t.Errorf("Error code not found in deletion with coop: %s", resp.Body.String())
	}
	if !eqInt(3, code) {
		t.Errorf("Expected error code is 3 when delete with coop, got %d", code)
	}
}
