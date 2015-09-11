package controller

import (
	"log"
	"net/http"

	"github.com/Patrolavia/mdpadgo/common"
	"github.com/Patrolavia/mdpadgo/user"
)

type User struct {
	SF     common.SessionFactory
	Config common.Config
}

func (uc *User) Users(w http.ResponseWriter, r *http.Request) {
	res := new(Response)
	userList, err := user.List()
	if err != nil {
		log.Printf("/api/users: While loading user list from db: %s", err)
		res.Fail("Cannot load user list").Do(w)
		return
	}

	res.Ok(userList).Do(w)
}

func (uc *User) Me(w http.ResponseWriter, r *http.Request) {
	res := new(Response)
	u, err := Me(uc.SF.Get(r))
	if err != nil {
		res.Fail("Not logged in").Do(w)
		return
	}

	res.Ok(u).Do(w)
}
