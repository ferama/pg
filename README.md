# pg

`pg` is an alternative CLI tool for PostgreSQL

It aims to be simple to use with less commands to remember

## Features

* very easy to use
* support for multiple database connections
* handy subcommands for database structure exploration
* commands autocomplete ( we all love `tab -> tab` right? :) )
* embedded sql editor for more advanced exploration
* query history navigation

## Handy subcommands

![pg subcommands](https://raw.githubusercontent.com/ferama/pg/main/media/commands.png)


## Sql Editor

![pg sql editor](https://raw.githubusercontent.com/ferama/pg/main/media/editor.png)

## Web UI

Every feature available on the CLI is also reachable from a browser:

```sh
pg web --addr 127.0.0.1 --port 8080
```

Then open http://127.0.0.1:8080. The web UI covers connection browsing,
database/schema creation and chown, table search, a SQL editor with
autocomplete and query history, connection stats and user management.

### Building from source

The web UI is a Vite/React app under `web/` that gets embedded into the `pg`
binary via `go:embed`. `go build .` always succeeds even without Node
installed (a placeholder page is committed at `pkg/webui/dist`), but to get
the real interface you need to build it first:

```sh
cd web && npm ci && npm run build
cd ..
go build .
```

`build.sh` already does this before cross-compiling the release binaries.

## Config file example

Put this in your $HOME/.pg dir

```yaml
connections:
  - name: db1
    url: postgres://user1:pass1@host1:5432/db1
  - name: db2
    url: postgres://user2:pass2@host2:5432/db2
```