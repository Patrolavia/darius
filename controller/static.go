package controller

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/Patrolavia/mdpadgo/common"
)

// Static serves static file.
//
// The index.html is treated as go html template, provided root url as only argument, aka. {{.}}
//
// Any path which cannot match to a file will fall back to index.html
type Static struct {
	Config common.Config
	cache  *template.Template
}

func (s Static) File(w http.ResponseWriter, r *http.Request) {
	if s.Config["FrontEnd"] == "" {
		http.Error(w, "Internal server error", 500)
		return
	}
	fn := s.Config["FrontEnd"] + r.URL.Path
	if info, err := os.Stat(fn); err != nil || r.URL.Path == "/index.html" || info.IsDir() {
		s.index(w)
		return
	}

	http.ServeFile(w, r, fn)
}

func (s Static) index(w http.ResponseWriter) {
	if s.cache == nil {
		fn := s.Config["FrontEnd"] + "/index.html"
		tmpl, err := template.ParseFiles(fn)
		if err != nil {
			log.Printf("Error parsing index.html: %s", err)
			http.Error(w, "Failed to load template", 500)
			return
		}
		s.cache = tmpl
	}
	if err := s.cache.Execute(w, s.Config["SiteRoot"]); err != nil {
		log.Printf("Failed to execute template: %s", err)
		http.Error(w, "Failed to execute template", 500)
	}
}
