# Ross 2

[![GoDoc](https://godoc.org/github.com/mbaraa/ross2?status.png)](https://godoc.org/github.com/mbaraa/ross2) [![Go Report Card](https://goreportcard.com/badge/github.com/mbaraa/ross2)](https://goreportcard.com/report/github.com/mbaraa/ross2) ![GitHub](https://img.shields.io/github/license/mbaraa/ross2)

**`Ross 2` is the second, redesigned, refactored, and can make mashed potatoes version of `Ross` (after serious battles with the initial version 🙂👀), `Ross` is a university contest managing platform**

**this is the backend repo of `Ross 2`**

## Project setup

- run the queries inside the file `./db/init_db.sql`
- modify the db's password in the file `config.json`
- set a proper port number and/or set the frontend's address (every address is accepted by default) in the file `config.json`
- compile the server's executable

```
go get
go build -ldflags="-w -s"
```

- run the server
```
./ross2
```
