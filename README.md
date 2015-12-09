# Darius - MDPAD backend

Darius is part of MDPAD project.

[![Build Status](https://travis-ci.org/Patrolavia/darius.svg?branch=master)](https://travis-ci.org/Patrolavia/darius)

A simple server to hold markdown documents for your small company or personal use.

## Synopsis

* Download binary or `go get github.com/Patrolavia/darius`.
* Create [config.json](https://github.com/Patrolavia/darius/blob/master/config.example.json).
* `darius config.json`

API specification see [API.md](https://github.com/Patrolavia/darius/blob/master/API.md).

## Configuration

* SiteRoot: Site root url with tailing slash. eg: `http://example.com:12345/my/pad/`.
* FrontEnd: Local path you put frontend files. eg: `../frontend` or `/srv/mdpad/frontend/build`
* Listen: Parameter pass to [http.ListenAndServe](http://golang.org/pkg/net/http/#ListenAndServe). The address and port to bind.
* DBType: Database driver, can be `sqlite3` or `mysql`.
* DBConStr: DB connection string, varies according to `DBType`.
* RedisAddr: Redis server connection string.
* SessSecret: SALT for session, type some random string here.
* SessName: Session name prefix.
* GoogleKeyFile: Where to find your google OAuth credential. This file must be in json format, which can be downloaded from Google Developer Console.
* ValidEditor: Emails of valid pad creators, comma-separated. eg: `a@example.com,b@example.com`. Leave blank if everyone can create pad.

## Darius?

Darius is a blue whale. He is 17 years old, lives in Antarctic Ocean and travels to Southern Atlantic Ocean every winter.

## License

Any version of MIT, GPL or LGPL.
