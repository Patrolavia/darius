package controller

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"testing"

	"github.com/Patrolavia/jsonapi"
	"github.com/Patrolavia/darius/model"
)

var pmTestUser *model.User
var pmCoop *model.User

func init() {
	var err error
	pmTestUser, err = model.NewUser("test pm", "pm@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		log.Fatalf("Cannot create user for testing pad modification: %s", err)
	}

	pmCoop, err = model.NewUser("test pm coop", "pm_coop@patrolavia.com", "https://patrolavia.com/logo128.png")
	if err != nil {
		log.Fatalf("Cannot create cooperator for testing pad modification: %s", err)
	}
}

func realTestCreate(t *testing.T, tags []string, coops []int) {
	t.Logf("Test pad creation with tags %v and coops %v", tags, coops)
	randomString := rand.Int()
	randomTitle := fmt.Sprintf("pm create-%d", randomString)
	param := map[string]interface{}{
		"title":      randomTitle,
		"content":    fmt.Sprintf("test pm create-%d", randomString),
		"tags":       tags,
		"cooperator": coops,
	}
	sess := session(pc().SF, t)
	sess.Login(pmTestUser)
	sess.Save()
	if err := sess.Err(); err != nil {
		t.Fatalf("Failed to login for creating pad: %s", err)
	}
	cookie := sess.Cookie()

	resp, err := jsonapi.HandlerTest(pc().Create).PostJSON("/api/create", cookie, param)
	if err != nil {
		t.Fatalf("While getting response of pad creation: %s", err)
	}

	if !testResult(resp.Body, true) {
		t.Fatalf("Returns error when creating pad: %s", resp.Body.String())
	}

	data, ok := testHasData(resp.Body, "id")
	if !ok {
		t.Fatalf("No pad returned from creating pad: %s", resp.Body.String())
	}

	var pid int
	switch v := data.(type) {
	case float64, float32:
		pid = int(reflect.ValueOf(v).Float())
	case int64, int32, int16, int8, int:
		pid = int(reflect.ValueOf(v).Int())
	default:
		t.Fatalf("returned pid is not integer: %#v", data)
	}

	pad, err := model.LoadPad(pid)
	if err != nil {
		t.Fatalf("Failed to save created pad into db: %s", err)
	}
	if pad.Title != randomTitle {
		t.Errorf("Expected title [%s], got [%s]", randomTitle, pad.Title)
	}
	ts := tags
	if tags == nil {
		ts = []string{}
	}
	if !reflect.DeepEqual(pad.Tags, ts) {
		t.Errorf("Expected tags %v, got %v", tags, pad.Tags)
	}
	c := coops
	if coops == nil {
		c = []int{}
	}
	if !reflect.DeepEqual(pad.Cooperators, c) {
		t.Errorf("Expected coops %v, got %v", coops, pad.Cooperators)
	}
}

func TestCreatePad(t *testing.T) {
	type params struct {
		t []string
		c []int
	}

	arr := []params{
		{t: nil, c: nil},
		{t: []string{"tag1"}, c: nil},
		{t: nil, c: []int{pmCoop.ID}},
		{t: []string{}, c: []int{}},
		{t: []string{"tag1"}, c: []int{}},
		{t: []string{}, c: []int{pmCoop.ID}},
		{t: []string{"tag1"}, c: []int{pmCoop.ID}},
	}

	for _, v := range arr {
		realTestCreate(t, v.t, v.c)
	}
}

func TestDeletePad(t *testing.T) {
	sess := session(pc().SF, t)
	sess.Login(pmTestUser)
	sess.Save()
	if err := sess.Err(); err != nil {
		t.Fatalf("Failed to login for deleting pad: %s", err)
	}
	cookie := sess.Cookie()
	pad, err := model.NewPad(db, pmTestUser.ID, "pm del", "test pm del", nil, nil)
	if err != nil {
		t.Fatalf("Cannot create pad for later deletion: %s", err)
	}

	resp, err := jsonapi.HandlerTest(pc().Delete).Get(fmt.Sprintf("/api/delete/%d", pad.ID), cookie)
	if err != nil {
		t.Fatalf("Error occured when delete pad: %s", err)
	}

	if !testResult(resp.Body, true) {
		t.Errorf("Failed to delete pad#%d: returned %s", pad.ID, resp.Body.String())
	}
}

func TestEditPad(t *testing.T) {
	sess := session(pc().SF, t)
	sess.Login(pmTestUser)
	sess.Save()
	if err := sess.Err(); err != nil {
		t.Fatalf("Failed to login for modifing pad: %s", err)
	}
	cookie := sess.Cookie()
	pad, err := model.NewPad(db, pmTestUser.ID, "pm edit", "test pm edit", nil, nil)
	if err != nil {
		t.Fatalf("Cannot create pad for later modify: %s", err)
	}
	type params struct {
		t []string
		c []int
	}

	arr := []params{
		{t: []string{"tag1"}, c: nil},
		{t: nil, c: []int{pmCoop.ID}},
		{t: []string{}, c: []int{}},
		{t: []string{"tag1"}, c: []int{}},
		{t: []string{}, c: []int{pmCoop.ID}},
		{t: []string{"tag1"}, c: []int{pmCoop.ID}},
		{t: nil, c: nil},
		{t: nil, c: nil}, // this line tests if we did not modify anything.
	}
	for _, v := range arr {
		realTestEdit(t, pad.ID, v.t, v.c, cookie)
	}
}

func realTestEdit(t *testing.T, pid int, tags []string, coops []int, cookie string) {
	pad, err := model.LoadPad(pid)
	if err != nil {
		t.Fatalf("Cannot load pad to edit: %s", err)
	}
	t.Logf("Test pad modification with tags %#v and coops %#v", tags, coops)
	ts := tags
	if tags == nil {
		ts = []string{}
	}
	c := coops
	if coops == nil {
		c = []int{}
	}
	randomTitle := fmt.Sprintf("randomTitle-%d", rand.Int())
	uri := fmt.Sprintf("/api/edit/%d", pad.ID)
	param := map[string]interface{}{
		"title":      randomTitle,
		"content":    pad.Content,
		"tags":       tags,
		"cooperator": coops,
		"version":    pad.Version,
	}

	resp, err := jsonapi.HandlerTest(pc().Edit).PostJSON(uri, cookie, param)
	if err != nil {
		t.Fatalf("Cannot save pad modification to db: %s", err)
	}
	if !testResult(resp.Body, true) {
		t.Errorf("/apd/edit returns false: %s", resp.Body.String())
	}
	actual, err := model.LoadPad(pad.ID)
	if err != nil {
		t.Fatalf("Cannot load just updated pad#%d: %s", pad.ID, err)
	}
	if !reflect.DeepEqual(actual.Tags, ts) {
		t.Errorf("Expected tags %v, got %v", tags, pad.Tags)
	}
	if !reflect.DeepEqual(actual.Cooperators, c) {
		t.Errorf("Expected coops %v, got %v", coops, pad.Cooperators)
	}

}
