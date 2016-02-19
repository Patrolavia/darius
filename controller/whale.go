package controller

import (
	"encoding/json"
	"os/exec"

	"github.com/Patrolavia/darius/common"
	"github.com/Patrolavia/darius/model"
	"github.com/Patrolavia/jsonapi"
)

// Whale is super whale power
type Whale struct {
	SF common.SessionFactory
}

// Whale handles 404 page
func (c *Whale) Whale(w *json.Encoder, r *json.Decoder, h *jsonapi.HTTP) {
	res := new(Response)
	p := &model.PadContent{
		Pad: &model.Pad{
			Title:       "Not found",
			Tags:        []string{},
			Cooperators: []int{},
		},
		Content: `# What?
Even a wise whale like Darius cannot understand your request.

He doesn't have secret whale power.`,
	}

	cmd := exec.Command("fortune")
	if msg, err := cmd.Output(); err == nil {
		p.Title = "Wise words from Darius the whale"
		p.Content = "```\n" + string(msg) + "\n```"
	}

	p.Render()
	res.Ok(p).Do(w)
}
