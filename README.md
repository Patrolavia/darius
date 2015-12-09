# mdpadgo

[![Build Status](https://travis-ci.org/Patrolavia/mdpadgo.svg?branch=master)](https://travis-ci.org/Patrolavia/mdpadgo)

A simple server to hold markdown documents for your small company or personal use.

## Synopsis

* Build from source or down binary.
* Create [config.json](https://github.com/Patrolavia/mdpadgo/blob/master/config.example.json).
* `mdpadgo config.json`

API specification see [API.md](https://github.com/Patrolavia/mdpadgo/blob/master/API.md).

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

## For frontend developers

You can run a local mdpad server in few steps.

1. Copy `config.example.json` to `config.json`.
2. Add a fake record in your hosts file. Normally, it is placed in `/etc/hosts` in linux, `C:\Windows\system32\drivers\etc\hosts` in windows. This record have to have a valid TLD like `.com`, `.org`, `.cc`... etc. Point this record (`example.com` for example) to `127.0.0.1`.
3. Obtain OAuth 2.0 credentials from the Google Developers Console. Specify callback url to someting like `http://example.com:8000/auth/google/oauth2callback`. Save you credential (json format) to `google.json`.
4. Edit `config.json`, change `SiteRoot` to something like `http://example.com:8000/` (remember the tailing slash), and change `FrontEnd` to the path cantains your frontend file.
5. Install ad start a redis server.
6. Build (`go build`) and run the server (`./mdpadgo config.json` in linux, `mdpadgo.exe config.json` in windows).

By these steps, server will put data in memory, so every time you restart it will also cleanup you data. Change `DBConStr` to `data.db` if you want to persist your testing data.

Session data are placed in redis server, and it will not be cleared when restarting mdpadgo server. You may need to delete cookies to reset current session.

## License

Any version of MIT, GPL or LGPL.
