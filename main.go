package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Patrolavia/mdpadgo/common"
	"github.com/Patrolavia/mdpadgo/controller"
	"github.com/Patrolavia/mdpadgo/pad"
	"github.com/Patrolavia/mdpadgo/user"
	"github.com/gorilla/context"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2/google"
)

var jsonFile string

func init() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s path/to/config.json", os.Args[0])
	}
	jsonFile = os.Args[1]
}

func main() {
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("Cannot read data from %s: %s", jsonFile, err)
	}

	cfg, err := common.JsonConfig(data)
	if err != nil {
		log.Fatalf("Cannot read configuration: %s", err)
	}

	db, err := cfg.DB()
	if err != nil {
		log.Fatalf("Cannot open db connection: %s", err)
	}

	if err := initDB(db, cfg["DBType"]); err != nil {
		log.Fatalf("Cannot initialize database: %s", err)
	}

	sf, err := common.BuildSession(cfg)
	if err != nil {
		log.Fatalf("Cannot prepare session store: %s", err)
	}

	if data, err = ioutil.ReadFile(cfg["GoogleKeyFile"]); err != nil {
		log.Fatalf("Cannot read data from %s: %s", cfg["GoogleKeyFile"], err)
	}

	googleConfig, err := google.ConfigFromJSON(data, "https://www.googleapis.com/auth/plus.login", "email")
	if err != nil {
		log.Fatalf("Cannot prepare google login info: %s", err)
	}

	ac := &controller.Auth{googleConfig, sf, cfg}
	http.HandleFunc("/auth/google", ac.Google)
	http.HandleFunc("/auth/google/oauth2callback", ac.GoogleCallback)
	http.HandleFunc("/api/paths", ac.Paths)
	http.HandleFunc("/api/logout", ac.Logout)

	uc := &controller.User{sf, cfg}
	http.HandleFunc("/api/user", uc.Me)
	http.HandleFunc("/api/users", uc.Users)

	pc := &controller.Pad{db, sf, cfg}
	http.HandleFunc("/api/create", pc.Create)
	http.HandleFunc("/api/pad/", pc.View)
	http.HandleFunc("/api/pads", pc.List)
	http.HandleFunc("/api/delete/", pc.Delete)
	http.HandleFunc("/api/edit/", pc.Edit)

	if err := http.ListenAndServe(cfg["Listen"], context.ClearHandler(http.DefaultServeMux)); err != nil {
		log.Fatalf("Cannot start http server at %s: %s", cfg["Listen"], err)
	}
}

func initDB(db *sql.DB, t string) (err error) {
	u := user.InitSqlite3
	p := pad.InitSqlite3

	if t == "mysql" {
		u = user.InitMysql
		p = user.InitMysql
	}

	err = u(db)
	if err == nil {
		err = p(db)
	}
	return
}