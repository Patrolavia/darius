package controller

import (
	"encoding/json"
	"log"

	"github.com/Patrolavia/jsonapi"
	"github.com/Patrolavia/darius/common"
	"github.com/Patrolavia/darius/model"
)

type User struct {
	SF     common.SessionFactory
	Config common.Config
}

func (uc *User) Users(w *json.Encoder, r *json.Decoder, h *jsonapi.HTTP) {
	res := new(Response)
	userList, err := model.ListUser()
	if err != nil {
		log.Printf("/api/users: While loading user list from db: %s", err)
		res.Fail("Cannot load user list").Do(w)
		return
	}

	res.Ok(userList).Do(w)
}

func (uc *User) Me(w *json.Encoder, r *json.Decoder, h *jsonapi.HTTP) {
	res := new(Response)
	u, err := Me(uc.SF.Get(h.Request))
	if err != nil {
		res.Fail("Not logged in").Do(w)
		return
	}

	res.Ok(u).Do(w)
}

func (uc *User) User(w *json.Encoder, r *json.Decoder, h *jsonapi.HTTP) {
	res := new(Response)
	var args map[string]interface{}
	if err := r.Decode(&args); err != nil {
		res.Fail("Arguments not in JSON format.")
		return
	}

	ids, ok := args["userid"]
	if !ok {
		res.Fail("No user id passed.").Do(w)
		return
	}

	uids, ok := ids.([]interface{})
	if !ok || len(uids) < 1 {
		res.Fail("No user id passed.").Do(w)
		return
	}

	ret := make([]*model.User, 0, len(uids))
	for _, uid := range uids {
		u, err := model.LoadUser(int(uid.(float64))) // json numbers converts to float64 in go
		if err != nil {
			log.Printf("Error loading user from db: %s", err)
			res.Fail("Error loading user from db").Do(w)
			return
		}
		ret = append(ret, u)
	}
	res.Ok(ret).Do(w)
}
