package common

import (
	"database/sql"
	"encoding/json"
)

type Config map[string]string

func JsonConfig(data []byte) (cfg Config, err error) {
	err = json.Unmarshal(data, &cfg)
	return
}

func (c Config) Url(path string) (url string) {
	switch {
	case c["SiteRoot"][len(c["SiteRoot"])-1:] == "/" && path[0:1] == "/":
		path = path[1:]
	case c["SiteRoot"][len(c["SiteRoot"])-1:] != "/" && path[0:1] != "/":
		path = "/" + path
	}
	return c["SiteRoot"] + path
}

func (c Config) DB() (db *sql.DB, err error) {
	return sql.Open(c["DBType"], c["DBConStr"])
}
