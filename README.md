# mdpadgo

A simple server to hold your markdown documents.

## Synopsis

* Build from source or down binary.
* Create [config.json](https://github.com/Patrolavia/mdpadgo/blob/master/config.example.json).
* `mdpadgo config.json`

API specification see [API.md](https://github.com/Patrolavia/mdpadgo/blob/master/API.md).

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
