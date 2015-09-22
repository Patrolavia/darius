package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Patrolavia/jsonapi"
	"github.com/Patrolavia/mdpadgo/common"
	"github.com/Patrolavia/mdpadgo/controller"
	"github.com/Patrolavia/mdpadgo/model"
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

	whale := &controller.Whale{sf}
	jsonapi.HandleFunc("/api/whale", whale.Whale)

	ac := &controller.Auth{googleConfig, sf, cfg}
	http.HandleFunc("/auth/google", ac.Google)
	http.HandleFunc("/auth/google/oauth2callback", ac.GoogleCallback)
	http.HandleFunc("/api/paths", ac.Paths)
	http.HandleFunc("/api/logout", ac.Logout)

	uc := &controller.User{sf, cfg}
	jsonapi.HandleFunc("/api/me", uc.Me)
	jsonapi.HandleFunc("/api/user/", uc.User)
	jsonapi.HandleFunc("/api/users", uc.Users)

	pc := &controller.Pad{db, sf, cfg}
	jsonapi.HandleFunc("/api/create", pc.Create)
	jsonapi.HandleFunc("/api/pad/", pc.View)
	jsonapi.HandleFunc("/api/pads", pc.List)
	jsonapi.HandleFunc("/api/delete/", pc.Delete)
	jsonapi.HandleFunc("/api/edit/", pc.Edit)

	sc := &controller.Static{Config: cfg}
	http.HandleFunc("/", sc.File)

	if err := http.ListenAndServe(cfg["Listen"], context.ClearHandler(http.DefaultServeMux)); err != nil {
		log.Fatalf("Cannot start http server at %s: %s", cfg["Listen"], err)
	}
}

func initDB(db *sql.DB, t string) (err error) {
	m := model.InitSqlite3

	if t == "mysql" {
		m = model.InitMysql
	}

	return m(db)
}
