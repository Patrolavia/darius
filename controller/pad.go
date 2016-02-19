// This file is part of Darius. See License.txt for license information.

package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	"github.com/Patrolavia/darius/common"
	"github.com/Patrolavia/darius/model"
	"github.com/Patrolavia/jsonapi"
)

// Pad represents a mdpad document, a.k.a. pad
type Pad struct {
	DB     *sql.DB
	SF     common.SessionFactory
	Config common.Config
}

// View handles pad render request
func (pc *Pad) View(w *json.Encoder, r *json.Decoder, h *jsonapi.HTTP) {
	res := new(Response)
	path := h.Request.URL.Path
	args := PathArg(path, "/api/pad/")
	if len(args) != 1 {
		log.Printf("Invalid path for /pad/pad: %s, %#v", path, args)
		res.Fail(path + " is invalid").Do(w)
		return
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		res.Fail(args[0] + " is not integer").Do(w)
		return
	}

	p, err := model.LoadPad(pid)
	if err != nil {
		log.Printf("Cannot load pad#%d from db: %s", pid, err)
		res.Failf("Cannot load pad#%d from database.", pid).Do(w)
		return
	}

	p.Sort()
	res.Ok(p).Do(w)
}

// List handles request to list all pads
func (pc *Pad) List(w *json.Encoder, r *json.Decoder, h *jsonapi.HTTP) {
	res := new(Response)
	pads, err := model.ListPad()
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
