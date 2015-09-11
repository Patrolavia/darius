package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/Patrolavia/mdpadgo/common"
	"github.com/Patrolavia/mdpadgo/pad"
)

type Pad struct {
	DB     *sql.DB
	SF     common.SessionFactory
	Config common.Config
}

func (pc *Pad) View(w http.ResponseWriter, r *http.Request) {
	res := new(Response)
	args := PathArg(r.URL.Path, "/api/pad/")
	if len(args) != 1 {
		log.Printf("Invalid path for /pad/pad: %s, %#v", r.URL.Path, args)
		res.Fail(r.URL.Path + " is invalid").Do(w)
		return
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		res.Fail(args[0] + " is not integer").Do(w)
		return
	}

	p, err := pad.Load(pid)
	if err != nil {
		log.Printf("Cannot load pad#%d from db: %s", pid, err)
		res.Failf("Cannot load pad#%d from database.", pid).Do(w)
		return
	}

	p.Sort()
	res.Ok(p).Do(w)
}

func (pc *Pad) List(w http.ResponseWriter, r *http.Request) {
	res := new(Response)
	pads, err := pad.List()
	if err != nil {
		log.Printf("Failed to load pad list: %s", err)
		res.Fail("Database error").Do(w)
		return
	}
	for _, v := range pads {
		v.Sort()
	}
	res.Ok(pads).Do(w)
}
