// This file is part of Darius. See License.txt for license information.

package controller

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/Patrolavia/darius/common"
	"github.com/Patrolavia/darius/model"
	"golang.org/x/oauth2"
	plus "google.golang.org/api/plus/v1"
)

// Auth is controller handling auth requests
type Auth struct {
	GoogleConfig *oauth2.Config
	SF           common.SessionFactory
	Config       common.Config
}

func (ac *Auth) googleConfig(r *http.Request) *oauth2.Config {
	ret := *ac.GoogleConfig // create a copy
	ret.RedirectURL = ac.Config.URL("/auth/google/oauth2callback")
	return &ret
}

// Google handles google oauth login request
func (ac *Auth) Google(w http.ResponseWriter, r *http.Request) {
	sess := ac.SF.Get(r)
	stat := fmt.Sprint(rand.Int())
	sess.Set("login_token", stat)
	sess.Save(r, w)
	if err := sess.Err(); err != nil {
		common.Fatalf(w, err, "Cannot save session")
		return
	}
	url := ac.googleConfig(r).AuthCodeURL(stat)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback is callback for google oauth
func (ac *Auth) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	sess := ac.SF.Get(r)
	stat := sess.Get("login_token")
	sess.Unset("login_token")
	if err := sess.Err(); err != nil {
		common.Fatalf(w, err, "Cannot read token from session")
		return
	}
	conf := ac.googleConfig(r)
	_ = r.ParseForm()
	code := r.Form.Get("code")
	state := r.Form.Get("state")
	if state != stat {
		http.Error(w, "Token mismatch!", http.StatusBadRequest)
		return
	}
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		common.Errorf(err, "Error occurs when exchanging token")
		http.Redirect(w, r, ac.Config.URL("/"), http.StatusTemporaryRedirect)
		return
	}
	client := conf.Client(oauth2.NoContext, tok)

	p, _ := plus.New(client)
	me, err := p.People.Get("me").Do()
	if err != nil {
		common.Errorf(err, "Error occurs when getting user info")
		http.Redirect(w, r, ac.Config.URL("/"), http.StatusTemporaryRedirect)
		return
	}

	var email string
	for _, e := range me.Emails {
		if email == "" {
			email = e.Value
		}
		if e.Type == "account" {
			email = e.Value
			break
		}
	}
	u, err := model.FindUser(email)
	if err != nil {
		common.Errorf(err, "Faile to save user, trying create one for %s.", email)
		u, err = model.NewUser(me.DisplayName, email, me.Image.Url)
		if err != nil {
			common.Fatalf(w, err, "Failed to login, please try again later.")
			return
		}
	}

	sess.Set("uid", fmt.Sprint(u.ID))
	sess.Save(r, w)
	if err := sess.Err(); err != nil {
		common.Fatal(w, err, "Cannot save user id in session")
		return
	}

	http.Redirect(w, r, ac.Config.URL("/"), http.StatusTemporaryRedirect)
}

// Paths returns available login method to user
func (ac *Auth) Paths(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h["Content-Type"] = []string{"application/json"}
	w.WriteHeader(200)
	path := `["google"]`
	w.Write([]byte(path))
}

// Logout current user
func (ac *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	sess := ac.SF.Get(r)
	sess.Unset("uid")
	sess.Save(r, w)
	if err := sess.Err(); err != nil {
		log.Printf("While user logout: %s", err)
	}
	http.Redirect(w, r, ac.Config.URL("/"), http.StatusTemporaryRedirect)
}
