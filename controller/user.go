package controller

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/Patrolavia/jsonapi"
	"github.com/Patrolavia/mdpadgo/common"
	"github.com/Patrolavia/mdpadgo/model"
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
	args := PathArg(h.Request.URL.Path, "/api/user/")
	if len(args) != 1 {
		res.Fail("No such user").Do(w)
		return
	}

	uid, err := strconv.Atoi(args[0])
	if err != nil {
		res.Fail("No such user").Do(w)
		return
	}

	u, err := model.LoadUser(uid)
	if err != nil {
		log.Printf("Error loading user from db: %s", err)
		res.Fail("Error loading user from db").Do(w)
		return
	}
	res.Ok(u).Do(w)
}
