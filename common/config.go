// This file is part of Darius. See License.txt for license information.

package common

import (
	"database/sql"
	"encoding/json"
)

// Config holds all configuration of darius
type Config map[string]string

// JSONConfig reads json format config
func JSONConfig(data []byte) (cfg Config, err error) {
	err = json.Unmarshal(data, &cfg)
	return
}

// URL returns site root url
func (c Config) URL(path string) (url string) {
	switch {
	case c["SiteRoot"][len(c["SiteRoot"])-1:] == "/" && path[0:1] == "/":
		path = path[1:]
	case c["SiteRoot"][len(c["SiteRoot"])-1:] != "/" && path[0:1] != "/":
		path = "/" + path
	}
	return c["SiteRoot"] + path
}

// DB returns db connection
func (c Config) DB() (db *sql.DB, err error) {
	return sql.Open(c["DBType"], c["DBConStr"])
}
