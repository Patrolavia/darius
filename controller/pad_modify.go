package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Patrolavia/mdpadgo/model"
)

func (pc *Pad) Create(w http.ResponseWriter, r *http.Request) {
	res := new(Response)
	u, err := Me(pc.SF.Get(r))
	if err != nil {
		res.Err(1, "Not logged in").Do(w)
		return
	}

	rawData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Cannot read post data from /api/create: %s", err)
		res.Err(2, "Data error").Do(w)
		return
	}

	var data CreateRequest
	if err := json.Unmarshal(rawData, &data); err != nil {
		log.Printf("Failed to parse postdata [%s]: %s", string(rawData), err)
		res.Err(2, "Data error").Do(w)
		return
	}

	p, err := model.NewPad(pc.DB, u.ID, data.Title, data.Content, data.Tags, data.Coops)
	if err != nil {
		log.Printf("Failed creating pad: %s", err)
		res.Err(2, "Database error").Do(w)
		return
	}

	resData := map[string]int{
		"code": 0,
		"id":   p.ID,
	}
	res.Ok(resData).Do(w)
}

func (pc *Pad) Delete(w http.ResponseWriter, r *http.Request) {
	res := new(Response)
	u, err := Me(pc.SF.Get(r))
	if err != nil {
		res.Err(1, "Not logged in").Do(w)
		return
	}

	args := PathArg(r.URL.Path, "/api/delete/")
	if len(args) != 1 {
		res.Err(2, "No such pad").Do(w)
		return
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		res.Err(2, "No such pad").Do(w)
		return
	}

	p, err := model.LoadPad(pid)
	if err != nil {
		res.Err(2, "No such pad").Do(w)
		return
	}

	if p.UID != u.ID {
		res.Err(3, "Not owner").Do(w)
		return
	}

	if err := p.Delete(); err != nil {
		log.Printf("Cannot delete pad#%d: %s", p.ID, err)
		res.Err(4, "Cannot delete from database").Do(w)
		return
	}

	res.Ok(nil).Do(w)
}

func (pc *Pad) Edit(w http.ResponseWriter, r *http.Request) {
	res := new(Response)
	u, err := Me(pc.SF.Get(r))
	if err != nil {
		res.Err(1, "Not logged in").Do(w)
		return
	}

	args := PathArg(r.URL.Path, "/api/edit/")
	if len(args) != 1 {
		res.Err(2, "No such pad").Do(w)
		return
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		res.Err(2, "No such pad").Do(w)
		return
	}

	p, err := model.LoadPad(pid)
	if err != nil {
		res.Err(2, "No such pad").Do(w)
		return
	}

	if p.UID != u.ID {
		if p.CoopModified() {
			res.Err(3, "Not owner").Do(w)
			return
		}
		var isCoop bool
		for _, c := range p.Cooperators {
			if c == u.ID {
				isCoop = true
				break
			}
		}
		if !isCoop {
			res.Err(3, "Not cooperator").Do(w)
			return

		}
	}

	rawData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res.Err(4, "Data error").Do(w)
		return
	}

	var data EditRequest
	if err := json.Unmarshal(rawData, &data); err != nil {
		res.Err(4, "Data error").Do(w)
		return
	}

	p.Title = data.Title
	p.Content = data.Content
	p.Tags = data.Tags
	p.Cooperators = data.Coops
	p.Version = data.Version
	err = p.Save(pc.DB)
	if err != nil {
		switch v := err.(type) {
		case model.VersionError:
			res.Err(5, "Version mismatch").Do(w)
		default:
			log.Printf("Error saving pad#%d: %s", p.ID, v)
			res.Err(4, "Unable to save pad into db").Do(w)
		}
		return
	}

	res.Ok(map[string]int{"code": 0}).Do(w)
}
