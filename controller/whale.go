package controller

import (
	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/Patrolavia/jsonapi"
	"github.com/Patrolavia/darius/common"
	"github.com/Patrolavia/darius/model"
)

type Whale struct {
	SF common.SessionFactory
}

func (c *Whale) Whale(w *json.Encoder, r *json.Decoder, h *jsonapi.HTTP) {
	u, err := Me(c.SF.Get(h.Request))
	if err != nil {
		http.Error(h.ResponseWriter, "Not found", 404)
		return
	}

	res := new(Response)
	p := &model.PadContent{
		Pad: &model.Pad{
			UID:         u.ID,
			Title:       "Whale not found",
			Tags:        []string{},
			Cooperators: []int{},
		},
		Content: "# We don't have secret whalepower.",
	}

	cmd := exec.Command("fortune")
	if msg, err := cmd.Output(); err == nil {
		p.Title = "Wise words from whale"
		p.Content = "```\n" + string(msg) + "\n```"
	}

	p.Render()
	res.Ok(p).Do(w)
}
