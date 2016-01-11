package controller

import (
	"fmt"
	"log"
	"testing"

	"github.com/Patrolavia/darius/model"
	"github.com/Patrolavia/jsonapi"
)

var padTestUser *model.User

func init() {
	var err error
	padTestUser, err = model.NewUser("test pad", "pad@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		log.Fatalf("Cannot create user for testing pad: %s", err)
	}
}

func TestPad(t *testing.T) {
	p, err := model.NewPad(db, padTestUser.ID, "pad", "test pad", []string{}, []int{})
	if err != nil {
		t.Fatalf("Cannot create pad for later testing: %s", err)
	}

	resp, err := jsonapi.HandlerTest(pc().View).Get(fmt.Sprintf("/api/pad/%d", p.ID), "")
	if err != nil {
		t.Fatalf("While getting response of /api/pad/%d: %s", p.ID, err)
	}

	if !testResult(resp.Body, true) {
		t.Fatalf("/api/pad/%d returns an error: %s", p.ID, resp.Body.String())
	}

	if !testData(p.ID, resp.Body, "id") {
		t.Fatalf("/api/pad/%d returns wrong data: %s", p.ID, resp.Body.String())
	}
}

func TestPads(t *testing.T) {
	p, err := model.NewPad(db, padTestUser.ID, "list pad", "test list pad", []string{}, []int{})
	if err != nil {
		t.Fatalf("Cannot create pad for testing list: %s", err)
	}

	resp, err := jsonapi.HandlerTest(pc().List).Get("/api/pads", "")
	if err != nil {
		t.Fatalf("While getting response of /api/pads: %s", err)
	}

	if !testResult(resp.Body, true) {
		t.Fatalf("/api/pads resurns an error: %s", resp.Body.String())
	}
	if !testArrayHas(p.ID, resp.Body) {
		t.Fatalf("Cannot find pad#%d in response of /api/pads: %s", p.ID, resp.Body.String())
	}
}
