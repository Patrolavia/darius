// This file is part of Darius. See License.txt for license information.

package controller

import (
	"fmt"
	"log"
	"testing"

	"github.com/Patrolavia/darius/model"
	"github.com/Patrolavia/jsonapi"
)

var peAclUser, peAclCoop, peAclGuest *model.User
var peAclPad *model.PadContent
var peAclURI string

func init() {
	var err error
	peAclUser, err = model.NewUser("test pe user", "pe_user@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		log.Fatalf("Cannot create user for testing pad edit acl: %s", err)
	}

	peAclCoop, err = model.NewUser("test pe coop", "pe_coop@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		log.Fatalf("Cannot create coop for testing pad edit acl: %s", err)
	}

	peAclGuest, err = model.NewUser("test pe guest", "pe_guest@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		log.Fatalf("Cannot create guest for testing pad edit acl: %s", err)
	}

	peAclPad, err = model.NewPad(db, peAclUser.ID, "pe acl", "test pe acl", nil, []int{peAclCoop.ID})
	if err != nil {
		log.Fatalf("Cannot create pad for testing pad edit act: %s", err)
	}

	peAclURI = fmt.Sprintf("/api/edit/%d", peAclPad.ID)
}

func TestPadEditAclNoLogin(t *testing.T) {
	param := map[string]interface{}{
		"title":      peAclPad.Title,
		"content":    peAclPad.Title,
		"tags":       peAclPad.Tags,
		"cooperator": peAclPad.Cooperators,
	}
	resp, err := jsonapi.HandlerTest(pc().Edit).PostJSON(peAclURI, "", param)
	if err != nil {
		t.Fatalf("Error occured when testing edit ACL without login: %s", err)
	}

	if !testResult(resp.Body, false) {
		t.Errorf("Edit success while no user logged in!: %s", resp.Body.String())
	}

	if !testData(1, resp.Body, "code") {
		t.Errorf("Expect error code 1 when edit without login, got response %s", resp.Body.String())
	}
}

func TestPadEditAclGuest(t *testing.T) {
	sess := session(pc().SF, t)
	sess.Login(peAclGuest)
	sess.Save()
	if err := sess.Err(); err != nil {
		t.Fatalf("While loging in as guest: %s", err)
	}
	param := map[string]interface{}{
		"title":      peAclPad.Title,
		"content":    peAclPad.Title,
		"tags":       peAclPad.Tags,
		"cooperator": peAclPad.Cooperators,
	}
	resp, err := jsonapi.HandlerTest(pc().Edit).PostJSON(peAclURI, sess.Cookie(), param)
	if err != nil {
		t.Fatalf("Error occured when testing edit ACL with guest: %s", err)
	}

	if !testResult(resp.Body, false) {
		t.Errorf("Edit success with guest: %s", resp.Body.String())
	}

	if !testData(3, resp.Body, "code") {
		t.Errorf("Expect error code 3 when edit with guest, got response %s", resp.Body.String())
	}
}

func TestPadEditAclCoop(t *testing.T) {
	sess := session(pc().SF, t)
	sess.Login(peAclCoop)
	sess.Save()
	if err := sess.Err(); err != nil {
		t.Fatalf("While loging in as coop: %s", err)
	}
	param := map[string]interface{}{
		"title":      peAclPad.Title,
		"content":    peAclPad.Title,
		"tags":       peAclPad.Tags,
		"cooperator": []int{peAclCoop.ID + 1},
		"version":    peAclPad.Version,
	}
	resp, err := jsonapi.HandlerTest(pc().Edit).PostJSON(peAclURI, sess.Cookie(), param)
	if err != nil {
		t.Fatalf("Error occured when testing edit ACL with coop: %s", err)
	}

	if !testResult(resp.Body, false) {
		t.Errorf("Edit cooperators success with coop: %s", resp.Body.String())
	}

	if !testData(3, resp.Body, "code") {
		t.Errorf("Expect error code 3 when edit cooperators with coop, got response %s", resp.Body.String())
	}
}
